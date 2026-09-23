package main

import (
	"fmt"
	"os"

	"github.com/AliFnieer/morph/internal/cli"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "doctor" {
		cli.RunDoctor()
		return
	}
	if len(os.Args) > 1 && os.Args[1] == "completion" {
		if len(os.Args) > 2 {
			shell := os.Args[2]
			if shell == "bash" {
				fmt.Print(`_morph() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    opts="--in --out --to --from --list --list-json --force --quiet --version --help doctor completion"
    if [[ ${cur} == -* ]] ; then
        COMPREPLY=( $(compgen -W "${opts}" -- ${cur}) )
        return 0
    fi
}
complete -F _morph morph
`)
				return
			} else if shell == "zsh" {
				fmt.Print("#compdef morph\n_morph() {\n  _arguments '--in[Input file]:file:_files' '--to[Target format]:format:' '--out[Output file]:file:_files' '--from[Source format]:format:' '--list[List conversions]' '--list-json[List as JSON]' '--force[Overwrite]' '--quiet[Quiet]' '--version[Version]' '--help[Help]' 'doctor[Run doctor]' 'completion[Generate completion]:shell:(bash zsh fish)'\n}\n")
				return
			} else if shell == "fish" {
				fmt.Print("complete -c morph -l in -d 'Input file' -r -F\ncomplete -c morph -l to -d 'Target format'\ncomplete -c morph -l out -d 'Output file' -r -F\ncomplete -c morph -l from -d 'Source format'\ncomplete -c morph -l list -d 'List conversions'\ncomplete -c morph -l list-json -d 'List as JSON'\ncomplete -c morph -l force -d 'Overwrite'\ncomplete -c morph -l quiet -d 'Quiet'\ncomplete -c morph -l version -d 'Version'\ncomplete -c morph -s h -l help -d 'Help'\ncomplete -c morph -a doctor -d 'Run doctor'\ncomplete -c morph -a completion -d 'Generate completion'\n")
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
