package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

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
	input := fs.String("in", "", "Input file path (use - for stdin)")
	output := fs.String("out", "", "Output file path (use - for stdout)")
	to := fs.String("to", "", "Target format (e.g. pdf, docx, png, jpg)")
	from := fs.String("from", "", "Source format (auto-detected if omitted)")
	list := fs.Bool("list", false, "List supported conversions")
	listJSON := fs.Bool("list-json", false, "List supported conversions as JSON")
	version := fs.Bool("version", false, "Show version")
	quiet := fs.Bool("quiet", false, "Quiet mode")
	force := fs.Bool("force", false, "Overwrite existing output file")

	if err := fs.Parse(args[1:]); err != nil {
		return err
	}

	if *version {
		fmt.Println("morph v0.0.1")
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
	var tmpIn string

	if *input == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return err
		}
		tmpIn = "/tmp/morph-stdin-" + randomSuffix()
		if err := os.WriteFile(tmpIn, data, 0644); err != nil {
			return err
		}
		defer os.Remove(tmpIn)
		*input = tmpIn
	} else {
		if _, err := os.Stat(*input); os.IsNotExist(err) {
			return fmt.Errorf("input file does not exist: %s", *input)
		}
	}

	out := *output
	if out == "" {
		out = defaultOutput(*input, dstFmt)
		if *input == tmpIn && strings.Contains(*input, "morph-stdin") {
			// can't infer well, but use stdout if requested
		}
	}
	if out == "-" {
		out = "/tmp/morph-stdout-" + randomSuffix() + "." + dstFmt
		defer func() {
			if data, err := os.ReadFile(out); err == nil {
				os.Stdout.Write(data)
			}
			os.Remove(out)
		}()
	} else if !*force {
		if _, err := os.Stat(out); err == nil {
			return fmt.Errorf("output file already exists: %s (use --force to overwrite)", out)
		}
	}

	srcFmt := strings.ToLower(*from)
	if srcFmt == "" {
		if *input != tmpIn {
			detected, err := detect.Detect(*input)
			if err != nil {
				return err
			}
			srcFmt = detected
		} else {
			srcFmt = dstFmt // fallback
		}
	}

	c := registry.Get(srcFmt, dstFmt)
	if c == nil {
		// try to suggest
		return fmt.Errorf("no converter registered for %s -> %s (use --list to see supported)", srcFmt, dstFmt)
	}

	if !*quiet && out != "/dev/stdout" {
		// don't spam if going to stdout
		if *output != "-" {
			fmt.Printf("Converting %s (%s) -> %s\n", *input, srcFmt, dstFmt)
		}
	}

	err := c.Convert(*input, out)
	if err != nil {
		return err
	}

	if !*quiet && *output != "-" {
		fmt.Printf("Saved to %s\n", out)
	}
	return nil
}

func randomSuffix() string {
	return fmt.Sprintf("%d", os.Getpid())
}

func printUsage() {
	fmt.Print(`morph - privacy-focused offline file format converter

Usage:
  morph --in <file> --to <format> [--out <file>] [--from <format>] [flags]
  morph doctor
  morph completion bash|zsh|fish

Flags:
  --in         Input file (use - for stdin)
  --to         Target format/extension (required)
  --out        Output file path (use - for stdout)
  --from       Source format (optional, auto-detected)
  --list       List supported conversions
  --list-json  List supported conversions as JSON
  --force      Overwrite existing output
  --quiet      Minimal output
  --version    Show version
  -h, --help   Show help

Examples:
  morph --in document.pdf --to docx
  morph --in image.png --to jpg
  morph --in data.csv --to json
  morph --in archive.zip --to tar.gz
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
