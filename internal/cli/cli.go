// Package cli adapts Mogent's public core operations to command-line parsing
// and text presentation.
package cli

import (
	"fmt"
	"io"
)

type commandDefinition struct {
	name  string
	usage string
	run   func([]string, io.Writer, io.Writer) error
}

// commands is the single registry for top-level dispatch, human usage, command
// completion, and typo suggestions.
func commands() []commandDefinition {
	return []commandDefinition{
		{name: "init", usage: "init [--template name] [--source alias=path] [--var name=value]", run: runInit},
		{name: "build", usage: "build [--manifest agents.yaml] [--force]", run: runBuild},
		{name: "status", usage: "status [--manifest agents.yaml]", run: runStatus},
		{name: "drift", usage: "drift [--manifest agents.yaml]", run: runDrift},
		{name: "coverage", usage: "coverage [source] [--manifest agents.yaml]", run: runCoverage},
		{name: "source", usage: "source <add|list|show|pin|update> ...", run: runSource},
		{name: "add", usage: "add <ref> [--under path [--first|--last] | --before path | --after path | --append]", run: runAdd},
		{name: "localize", usage: "localize <manifest-heading-path> [--from source-ref]", run: runLocalize},
		{name: "complete", usage: "complete <kind> [--manifest agents.yaml]", run: runComplete},
		{name: "completion", usage: "completion <bash|zsh>", run: runCompletion},
		{name: "tui", usage: "tui [--manifest agents.yaml]", run: runTUI},
	}
}

// Run dispatches one Mogent command without owning process exit behavior.
func Run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		return writeUsage(stdout)
	}
	for _, command := range commands() {
		if args[0] == command.name {
			return command.run(args[1:], stdout, stderr)
		}
	}
	return unknownCommandError(args[0])
}

func writeUsage(output io.Writer) error {
	if _, err := fmt.Fprintln(output, "Usage:"); err != nil {
		return fmt.Errorf("write usage: %w", err)
	}
	for _, command := range commands() {
		if _, err := fmt.Fprintf(output, "  mogent %s\n", command.usage); err != nil {
			return fmt.Errorf("write usage: %w", err)
		}
	}
	return nil
}
