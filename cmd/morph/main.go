package main

import (
	"fmt"
	"os"
	"runtime/debug"

	"github.com/AliFnieer/morph/internal/cli"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Printf("morph %s\n", version)
		fmt.Printf("commit: %s\n", commit)
		fmt.Printf("built: %s\n", date)
		if info, ok := debug.ReadBuildInfo(); ok {
			fmt.Printf("go: %s\n", info.GoVersion)
		}
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "doctor" {
		cli.RunDoctor()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "completion" {
		if len(os.Args) > 2 {
			switch os.Args[2] {
			case "bash":
				fmt.Print(`_morph() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="--in --out --to --from --list --list-json --force --quiet --verbose --version --help doctor completion"
    if [[ ${cur} == -* ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi
}
complete -F _morph morph
`)
				return
			case "zsh":
				fmt.Print("#compdef morph\n_morph() {\n  _arguments '--in[Input file]:file:_files' '--to[Target format]:format:' '--out[Output file]:file:_files' '--from[Source format]:format:' '--list[List conversions]' '--list-json[List as JSON]' '--force[Overwrite]' '--quiet[Quiet]' '--verbose[Verbose]' '--version[Version]' '--help[Help]' 'doctor[Run doctor]' 'completion[Generate completion]:shell:(bash zsh fish)'\n}\n")
				return
			case "fish":
				fmt.Print("complete -c morph -l in -d 'Input file' -r -F\ncomplete -c morph -l to -d 'Target format'\ncomplete -c morph -l out -d 'Output file' -r -F\ncomplete -c morph -l from -d 'Source format'\ncomplete -c morph -l list -d 'List conversions'\ncomplete -c morph -l list-json -d 'List as JSON'\ncomplete -c morph -l force -d 'Overwrite'\ncomplete -c morph -l quiet -d 'Quiet'\ncomplete -c morph -l verbose -d 'Verbose'\ncomplete -c morph -l version -d 'Version'\ncomplete -c morph -s h -l help -d 'Help'\ncomplete -c morph -a doctor -d 'Run doctor'\ncomplete -c morph -a version -d 'Show version'\ncomplete -c morph -a completion -d 'Generate completion'\n")
				return
			}
		}
		fmt.Fprintln(os.Stderr, "Usage: morph completion bash|zsh|fish")
		os.Exit(1)
	}
	if err := cli.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
