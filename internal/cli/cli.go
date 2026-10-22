// Package cli adapts Mogent's public core operations to command-line parsing
// and text presentation.
package cli

import (
	"io"
)

type commandDefinition struct {
	name  string
	usage string
	run   func([]string, io.Writer, io.Writer) error
}

// commands is the single registry for dispatch, the full command reference,
// command completion, and typo suggestions. Workflow help lives in help.go.
func commands() []commandDefinition {
	return []commandDefinition{
		{name: "init", usage: "init [--plan] [--file path] [--source alias=local-path|git-url] [--dry-run] [--no-interactive]", run: runInit},
		{name: "build", usage: "build [--file mogent.plan|mogent.build] [--check]", run: runBuild},
		{name: "sequence", usage: "sequence [--file mogent.plan]", run: runSequence},
		{name: "show", usage: "show [tree|build|issues|content <alias:path>] [--file path] (catalog: source tree [alias:path])", run: runShow},
		{name: "status", usage: "status [--file path]", run: runStatus},
		{name: "drift", usage: "drift [--file path] [--all] [--diff]", run: runDrift},
		{name: "source", usage: "source add|list|tree|search|content|register|update local ...", run: runSource},
		{name: "add", usage: "add <source:path> [--file path] [--all] [--dry-run|--diff]", run: runAdd},
		{name: "module", usage: "module add|rm|exclude|localize|optimize ...", run: runModule},
		{name: "rm", usage: "rm <source:path> [--file path] [--dry-run|--diff]", run: runRemove},
		{name: "optimize", usage: "optimize [--file path] [--apply|--dry-run]", run: runOptimise},
		{name: "optimise", usage: "optimise [--file path] [--apply|--dry-run] (alias for optimize)", run: runOptimise},
		{name: "except", usage: "except <source:path> [--file path] [--dry-run]", run: runExcept},
		{name: "exclude", usage: "exclude <source:path> [--file path] [--dry-run]", run: runExcept},
		{name: "reorder", usage: "reorder [--file path] [--dry-run]", run: runReorder},
		{name: "completion", usage: "completion <bash|zsh>", run: runCompletion},
		{name: "mint", usage: "mint [-w 1|2] [-r REPO_ROOT] [-n] [-s SEED]", run: runMint},
	}
}

// Run dispatches one Mogent command without owning process exit behavior.
func Run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return writeUsage(stdout)
	}
	if args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		return writeHelp(args[1:], stdout)
	}
	for _, command := range commands() {
		if args[0] == command.name {
			if helpRequested(args[1:]) {
				return writeHelp(commandHelpPath(args), stdout)
			}
			return command.run(args[1:], stdout, stderr)
		}
	}
	return unknownCommandError(args[0])
}
