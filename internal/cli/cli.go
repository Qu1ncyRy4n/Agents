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
		{name: "init", usage: "init [--plan] [--file path] [--source alias=local-path|git-url] [--dry-run] [--no-interactive] (legacy: --config mogent.hcl | --manifest agents.yaml)", run: runInit},
		{name: "build", usage: "build [--file mogent.plan|mogent.build] [--check] (legacy: --manifest agents.yaml [--dry-run] [--force] [--preserve-html-comments])", run: runBuild},
		{name: "plan", usage: "plan [--config mogent.hcl]", run: runPlan},
		{name: "sequence", usage: "sequence [--file mogent.plan]", run: runSequence},
		{name: "show", usage: "show [alias:path] [--file path] [--content [--lines n] [--headings] | --build | --issues] (catalog: source show [alias:path])", run: runShow},
		{name: "apply", usage: "apply [--config mogent.hcl] [--force]", run: runApply},
		{name: "update", usage: "update [alias] [--config mogent.hcl] [--accept]", run: runUpdate},
		{name: "lib", usage: "lib <init|scan|check> ...", run: runLib},
		{name: "status", usage: "status [--file path] (legacy: --manifest agents.yaml)", run: runStatus},
		{name: "drift", usage: "drift [--file path] [--all] [--diff] (legacy: --manifest agents.yaml [--import path | --reject --force])", run: runDrift},
		{name: "coverage", usage: "coverage [source] [--manifest agents.yaml]", run: runCoverage},
		{name: "source", usage: "source add <name> [path-or-url] [--file path] [--dry-run] | source list|ls [alias] [--file path] [--tree] | source list --known | source register <name> <path-or-url> | source show [alias:path] [--tree|--flat] [--chars ascii|unicode] | source <pin|update> ... (legacy: --config | --manifest)", run: runSource},
		{name: "edit", usage: "edit include <alias:path> --localize [--file mogent.plan] [--dry-run]", run: runEdit},
		{name: "add", usage: "add <source:path> [--file path] [--dry-run] (legacy: --manifest agents.yaml [--under path | --before path | --after path | --append])", run: runAdd},
		{name: "select", usage: "select <source:path> [--file path] [--dry-run]", run: runSelect},
		{name: "module", usage: "module add|rm <source:path> [--file path] [--dry-run]", run: runModule},
		{name: "rm", usage: "rm <source:path> [--file path] [--dry-run|--diff]", run: runRemove},
		{name: "optimize", usage: "optimize [--file path]", run: runOptimise},
		{name: "optimise", usage: "optimise [--file path] (alias for optimize)", run: runOptimise},
		{name: "except", usage: "except <source:path> [--file path] [--dry-run]", run: runExcept},
		{name: "exclude", usage: "exclude <source:path> [--file path] [--dry-run]", run: runExcept},
		{name: "reorder", usage: "reorder [--file path] [--dry-run]", run: runReorder},
		{name: "localize", usage: "localize <manifest-heading-path> [--from source-ref]", run: runLocalize},
		{name: "complete", usage: "complete <kind> [--manifest agents.yaml]", run: runComplete},
		{name: "completion", usage: "completion <bash|zsh>", run: runCompletion},
		{name: "tui", usage: "tui [--manifest agents.yaml]", run: runTUI},
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
