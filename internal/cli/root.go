package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/AliFnieer/morph/internal/detect"
	"github.com/AliFnieer/morph/internal/registry"
)

func Run(args []string) error {
	if len(args) < 2 {
		printUsage()
		return nil
	}

	if args[1] == "-h" || args[1] == "--help" || args[1] == "help" {
		printUsage()
		return nil
	}

	fs := flag.NewFlagSet("morph", flag.ContinueOnError)
	input := fs.String("in", "", "Input file path (use - for stdin, or glob pattern)")
	output := fs.String("out", "", "Output file path (use - for stdout)")
	to := fs.String("to", "", "Target format (e.g. pdf, docx, png, jpg)")
	from := fs.String("from", "", "Source format (auto-detected if omitted)")
	list := fs.Bool("list", false, "List supported conversions")
	listJSON := fs.Bool("list-json", false, "List supported conversions as JSON")
	version := fs.Bool("version", false, "Show version")
	quiet := fs.Bool("quiet", false, "Quiet mode")
	force := fs.Bool("force", false, "Overwrite existing output file")
	verbose := fs.Bool("verbose", false, "Verbose output")
	recursive := fs.Bool("recursive", false, "Recursive for directories")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	if *version {
		fmt.Println("morph v0.2.0")
		return nil
	}

	if *listJSON {
		registry.ListJSON()
		return nil
	}

	if *list {
		registry.List()
		return nil
	}

	if *input == "" || *to == "" {
		printUsage()
		return errors.New("missing required flags: --in and --to")
	}

	dstFmt := strings.ToLower(*to)

	// handle glob
	if strings.ContainsAny(*input, "*?[]") && *input != "-" {
		matches, err := filepath.Glob(*input)
		if err != nil {
			return err
		}
		if len(matches) == 0 {
			return fmt.Errorf("no files matched pattern: %s", *input)
		}
		for _, m := range matches {
			if err := convertOne(m, *to, "", *from, *force, *quiet, *verbose); err != nil {
				return err
			}
		}
		return nil
	}

	// handle directory
	if info, err := os.Stat(*input); err == nil && info.IsDir() {
		var files []string
		if *recursive {
			if err := filepath.Walk(*input, func(path string, fi os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if !fi.IsDir() {
					files = append(files, path)
				}
				return nil
			}); err != nil {
				return err
			}
		} else {
			entries, _ := os.ReadDir(*input)
			for _, e := range entries {
				if !e.IsDir() {
					files = append(files, filepath.Join(*input, e.Name()))
				}
			}
		}
		if len(files) == 0 {
			return errors.New("no files found in directory")
		}
		for _, f := range files {
			if err := convertOne(f, *to, "", *from, *force, *quiet, *verbose); err != nil {
				fmt.Fprintf(os.Stderr, "Error converting %s: %v\n", f, err)
				continue
			}
		}
		return nil
	}

	return convertOne(*input, dstFmt, *output, *from, *force, *quiet, *verbose)
}

func convertOne(input, dstFmt, output, from string, force, quiet, verbose bool) error {
	start := time.Now()
	var tmpIn string

	if input == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		tmpIn = "/tmp/morph-stdin-" + fmt.Sprintf("%d", os.Getpid()) + "-" + fmt.Sprintf("%d", time.Now().UnixNano())
		if err := os.WriteFile(tmpIn, data, 0644); err != nil {
			return err
		}
		defer func() { _ = os.Remove(tmpIn) }()
		input = tmpIn
	} else {
		if _, err := os.Stat(input); os.IsNotExist(err) {
			return fmt.Errorf("input file does not exist: %s", input)
		}
	}

	out := output
	if out == "" {
		out = defaultOutput(input, dstFmt)
	}
	if out == "-" {
		out = "/tmp/morph-stdout-" + fmt.Sprintf("%d", os.Getpid()) + "-" + fmt.Sprintf("%d", time.Now().UnixNano()) + "." + dstFmt
		defer func() {
			if data, err := os.ReadFile(out); err == nil {
				_, _ = os.Stdout.Write(data)
			}
			_ = os.Remove(out)
		}()
	} else if !force {
		if _, err := os.Stat(out); err == nil {
			return fmt.Errorf("output file already exists: %s (use --force to overwrite)", out)
		}
	}

	srcFmt := strings.ToLower(from)
	if srcFmt == "" {
		if !strings.Contains(input, "morph-stdin") {
			detected, err := detect.Detect(input)
			if err != nil {
				return err
			}
			srcFmt = detected
		} else {
			srcFmt = dstFmt
		}
	}

	c := registry.Get(srcFmt, dstFmt)
	if c == nil {
		return suggest(srcFmt, dstFmt)
	}

	if verbose && output != "-" {
		fmt.Printf("[verbose] detected: %s, converter: %s, in: %s\n", srcFmt, c.Name(), input)
	}

	if !quiet && output != "-" {
		fmt.Printf("Converting %s (%s) -> %s\n", input, srcFmt, dstFmt)
	}

	err := c.Convert(input, out)
	if err != nil {
		return err
	}

	if verbose && output != "-" {
		fmt.Printf("[verbose] done in %v\n", time.Since(start))
	}

	if !quiet && output != "-" {
		fmt.Printf("Saved to %s\n", out)
	}
	return nil
}

func suggest(from, to string) error {
	// try to find similar
	all := registry.AllPairs()
	var suggestions []string
	for _, p := range all {
		if strings.Contains(p.From, from[:min(len(from), 2)]) || strings.Contains(p.To, to[:min(len(to), 2)]) {
			suggestions = append(suggestions, fmt.Sprintf("%s -> %s", p.From, p.To))
		}
		if len(suggestions) >= 3 {
			break
		}
	}
	if len(suggestions) > 0 {
		return fmt.Errorf("no converter registered for %s -> %s (use --list to see supported). Did you mean: %s", from, to, strings.Join(suggestions, ", "))
	}
	return fmt.Errorf("no converter registered for %s -> %s (use --list to see supported)", from, to)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func printUsage() {
	fmt.Print(`morph - privacy-focused offline file format converter

Usage:
  morph --in <file> --to <format> [--out <file>] [--from <format>] [flags]
  morph doctor
  morph completion bash|zsh|fish
  morph version

Flags:
  --in         Input file (use - for stdin, glob, or directory)
  --to         Target format/extension (required)
  --out        Output file path (use - for stdout)
  --from       Source format (optional, auto-detected)
  --list       List supported conversions
  --list-json  List supported conversions as JSON
  --force      Overwrite existing output
  --quiet      Minimal output
  --verbose    Verbose output
  --recursive  Process directories recursively
  --version    Show version
  -h, --help   Show help

Examples:
  morph --in document.pdf --to docx
  morph --in "*.png" --to jpg
  morph --in ./images/ --to webp
  morph --in data.csv --to json
  cat data.json | morph --in - --to yaml --out -
`)
}

func defaultOutput(in, ext string) string {
	dot := strings.LastIndex(in, ".")
	if dot > 0 {
		return in[:dot] + "." + ext
	}
	return in + "." + ext
}
