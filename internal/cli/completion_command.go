package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/workspace"
)

func runComplete(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("complete", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	sourceAlias := flags.String("source", "", "limit source-ref completion to one source alias")
	prefix := flags.String("prefix", "", "only print candidates with this prefix")
	args = reorderArgs(args, map[string]bool{
		"-manifest": true, "--manifest": true,
		"-source": true, "--source": true,
		"-prefix": true, "--prefix": true,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("complete requires one kind: commands, source-refs, source-aliases, manifest-headings, or flags")
	}
	var candidates []string
	switch flags.Arg(0) {
	case "commands":
		candidates = commandCandidates()
	case "flags":
		candidates = flagCandidates()
	case "source-refs":
		session, err := workspace.New(*manifestFile)
		if err != nil {
			return err
		}
		candidates, err = session.SourceReferences(*sourceAlias)
		if err != nil {
			return err
		}
	case "source-aliases":
		session, err := workspace.New(*manifestFile)
		if err != nil {
			return err
		}
		for _, source := range session.SavedSourceAliases() {
			candidates = append(candidates, source)
		}
	case "manifest-headings":
		session, err := workspace.New(*manifestFile)
		if err != nil {
			return err
		}
		candidates = session.ManifestHeadingPaths()
	default:
		return fmt.Errorf("unknown completion kind %q; use commands, source-refs, source-aliases, manifest-headings, or flags", flags.Arg(0))
	}
	for _, candidate := range candidates {
		if *prefix != "" && !strings.HasPrefix(candidate, *prefix) {
			continue
		}
		if _, err := fmt.Fprintln(stdout, candidate); err != nil {
			return fmt.Errorf("write completions: %w", err)
		}
	}
	return nil
}

func runCompletion(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("completion", flag.ContinueOnError)
	flags.SetOutput(stderr)
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("completion requires one shell: bash or zsh")
	}
	var script string
	switch flags.Arg(0) {
	case "bash":
		script = bashCompletionScript()
	case "zsh":
		script = zshCompletionScript()
	default:
		return fmt.Errorf("unknown shell %q; use bash or zsh", flags.Arg(0))
	}
	if _, err := fmt.Fprint(stdout, script); err != nil {
		return fmt.Errorf("write completion script: %w", err)
	}
	return nil
}

func commandCandidates() []string {
	definitions := commands()
	candidates := make([]string, 0, len(definitions)+1)
	for _, command := range definitions {
		candidates = append(candidates, command.name)
	}
	return append(candidates, "help")
}

func flagCandidates() []string {
	return []string{"--manifest", "--template", "--output", "--list-templates", "--build", "--force", "--source", "--subdir", "--from", "--import", "--reject", "--ref", "--accept", "--tag", "--tag-search", "--search", "--sort", "--metadata", "--tldr", "--file", "--line", "--content", "--lines", "--align-source", "--align", "--chars", "--fit", "--width", "--coverage", "--under", "--append", "--first", "--last", "--before", "--after", "--heading", "--dry-run", "--preview", "--rebuild", "--unused-only", "--tree", "--content-only", "--leaves-only", "--depth"}
}

func bashCompletionScript() string {
	return `_mogent_complete_candidates() {
  mogent complete "$1" --prefix "$2" 2>/dev/null
}

_mogent_completion() {
  local cur prev command subcommand
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
  command="${COMP_WORDS[1]}"
  subcommand="${COMP_WORDS[2]}"

  if [[ ${COMP_CWORD} -eq 1 ]]; then
    mapfile -t COMPREPLY < <(_mogent_complete_candidates commands "$cur")
    return 0
  fi

	case "$prev" in
    --manifest|-manifest)
      COMPREPLY=( $(compgen -f -- "$cur") )
      return 0
      ;;
    --source|-source)
      mapfile -t COMPREPLY < <(_mogent_complete_candidates source-aliases "$cur")
      return 0
      ;;
		--under|-under|--before|-before|--after|-after)
      mapfile -t COMPREPLY < <(_mogent_complete_candidates manifest-headings "$cur")
      return 0
			;;
		--import|-import)
			mapfile -t COMPREPLY < <(_mogent_complete_candidates manifest-headings "$cur")
			return 0
			;;
  esac

  if [[ "$cur" == -* ]]; then
    mapfile -t COMPREPLY < <(_mogent_complete_candidates flags "$cur")
    return 0
  fi

  if [[ "$command" == "source" && "$subcommand" == "show" ]]; then
    mapfile -t COMPREPLY < <(_mogent_complete_candidates source-refs "$cur")
    return 0
  fi
  if [[ "$command" == "source" && "$subcommand" == "list" ]]; then
    mapfile -t COMPREPLY < <(_mogent_complete_candidates source-aliases "$cur")
    return 0
  fi
	if [[ "$command" == "add" ]]; then
    mapfile -t COMPREPLY < <(_mogent_complete_candidates source-refs "$cur")
    return 0
	fi
	if [[ "$command" == "localize" ]]; then
		mapfile -t COMPREPLY < <(_mogent_complete_candidates manifest-headings "$cur")
		return 0
	fi
}

complete -F _mogent_completion mogent
`
}

func zshCompletionScript() string {
	return `#compdef mogent

_mogent_complete_candidates() {
  mogent complete "$1" --prefix "$2" 2>/dev/null
}

_mogent() {
  local -a candidates
  local command subcommand
  command="${words[2]}"
  subcommand="${words[3]}"

  if (( CURRENT == 2 )); then
    candidates=("${(@f)$(_mogent_complete_candidates commands "$PREFIX")}")
    _describe 'command' candidates
    return
  fi

	case "${words[CURRENT-1]}" in
    --manifest|-manifest)
      _files
      return
      ;;
    --source|-source)
      candidates=("${(@f)$(_mogent_complete_candidates source-aliases "$PREFIX")}")
      _describe 'source alias' candidates
      return
      ;;
		--under|-under|--before|-before|--after|-after)
      candidates=("${(@f)$(_mogent_complete_candidates manifest-headings "$PREFIX")}")
      _describe 'manifest heading' candidates
      return
			;;
		--import|-import)
			candidates=("${(@f)$(_mogent_complete_candidates manifest-headings "$PREFIX")}")
			_describe 'manifest heading' candidates
			return
			;;
  esac

  if [[ "$PREFIX" == -* ]]; then
    candidates=("${(@f)$(_mogent_complete_candidates flags "$PREFIX")}")
    _describe 'flag' candidates
    return
  fi

  if [[ "$command" == "source" && "$subcommand" == "show" ]]; then
    candidates=("${(@f)$(_mogent_complete_candidates source-refs "$PREFIX")}")
    _describe 'source ref' candidates
    return
  fi
  if [[ "$command" == "source" && "$subcommand" == "list" ]]; then
    candidates=("${(@f)$(_mogent_complete_candidates source-aliases "$PREFIX")}")
    _describe 'source alias' candidates
    return
  fi
	if [[ "$command" == "add" ]]; then
    candidates=("${(@f)$(_mogent_complete_candidates source-refs "$PREFIX")}")
    _describe 'source ref' candidates
    return
	fi
	if [[ "$command" == "localize" ]]; then
		candidates=("${(@f)$(_mogent_complete_candidates manifest-headings "$PREFIX")}")
		_describe 'manifest heading' candidates
		return
	fi
}

_mogent "$@"
`
}
