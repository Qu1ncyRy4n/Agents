package cli

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/charmbracelet/x/term"
)

const workflowHelp = `Mogent — compose project instructions and skills from reusable libraries.

Usage: mogent <command> [<args>]
       mogent [-h | --help]
       mogent help [<topic> | --all]

Project initialization (advanced: mogent help init)
   init                                Create an authored mogent.build
   source add <name> [<path-or-url>]    Connect a local, Git, or remembered library

Browse content (advanced: mogent help source)
   source list                         List connected libraries
   source show                         Browse their content as a tree

Select modules (advanced: mogent help select)
   add <name:path>                    Select a module in the project file
   select <name:path>                 Alias for add
   module add <name:path>             Explicit module-selection form
   rm <name:path>                     Remove one direct include
   except <name:path>                 Exclude a descendant from a selected ancestor
   optimize                            Inspect direct include overlaps; writes nothing
   reorder                             Group bindings/sources before ordered includes

Preview the selection (advanced: mogent help show)
   show --build                        Show selected modules and destinations
   show --issues                       Explain unresolved choices or errors
   show <name:path> --content          Preview a module's rendered or physical content

Build outputs (advanced: mogent help build)
   build                               Generate the selected outputs
   drift [--diff]                      Compare disk with output state and current build

Further help
   mogent help --all                    Expanded workflow with common flags
   mogent <command> --help              Detailed help for one workflow stage
   mogent help legacy                   HCL v2 and YAML workflows
   mogent help tools                    Completion and library utilities
   mogent help commands                 All command signatures
`

const expandedWorkflowHelp = `Mogent — compose project instructions and skills from reusable libraries.

Usage: mogent <command> [<args>]
       mogent [-h | --help]
       mogent help [<topic> | --all]

Common commands, grouped by workflow:

Project initialization (advanced: mogent help init)
   init                                Create an authored mogent.build
       --plan                          Create mogent.plan instead
       --source <name=path-or-url>      Connect a library; repeat for more sources
       --dry-run                       Preview without writing
       --no-interactive                Skip the optional Gum source picker
   source add <name> [<path-or-url>]    Connect a local, Git, or remembered library
       --dry-run                       Preview the source declaration

Browse content (advanced: mogent help source)
   source list, source ls              List connected libraries and locations
       --known                         List remembered libraries instead
       --tree                          Browse connected libraries' content
   source show [<name:path>]           Browse all content or one subtree
       --tree | --flat                 Hierarchy (default) or complete paths
       --chars ascii|unicode           Choose tree connectors
   source register <name> <path-or-url> Remember a library for future projects

Select modules (advanced: mogent help select)
   add <name:path>                    Select a module in the project file
   select <name:path>                 Alias for add
    module add <name:path>             Explicit module-selection form
        --all                          Select all direct children (for opt-in trees)
        --file <path>                  Select a specific plan/build
        --dry-run                      Preview the include without writing
   rm <name:path>                     Remove one exact direct include
       --diff                         Preview managed-output changes without writing
   optimize                            Explain direct ancestor/descendant overlaps; writes nothing
   except <name:path>                 Exclude a descendant from its selected ancestor
       --file <path>                  Select a specific plan/build
       --dry-run                      Preview the exception without writing
   reorder                            Group bindings/sources before ordered includes
       --dry-run                      Preview the formatting rewrite
   edit include <name:path> --localize Copy an input into an editable local source
       --file <path>                   Select the plan/build to edit
       --dry-run                       Preview localization without writing

Preview the selection (advanced: mogent help show)
   show --build                       Show selected modules and destinations
   show --issues                      Explain unresolved choices or errors
       --file <path>                   Inspect a specific plan/build

Build outputs (advanced: mogent help build)
   build                              Generate the selected outputs
       --check                         Verify existing outputs without writing
       --file <path>                   Build a specific plan/build
    sequence --file <plan>              Generate a build without installing outputs
   status                              Report all managed output baselines and current-build matches
   drift                               Report AGENTS.md baseline and current-build match
       --all                            Include the managed .agents/skills tree
       --diff                           Print current-build versus disk differences

init, source add/list, show and build accept --file <path>.
source add/list, show and build discover mogent.plan or mogent.build by default.
In include examples, replace qmr with your connected library's name.

Further help
   mogent <command> --help             Common usage and advanced options
   mogent help legacy                  HCL v2 and YAML workflows
   mogent help tools                   Completion and library utilities
   mogent help commands                All command signatures
`

// Stage help deliberately separates the next useful actions from optional
// flags. The registry remains the exhaustive command reference, not onboarding.
var workflowTopics = map[string]string{
	"init": `Project initialization

  mogent init
  mogent source add qmr ../qmr-agents-library

init creates an authored mogent.build for you to edit. Adding a source makes
content available; it does not select modules or install outputs.

With Gum installed, interactive init offers known libraries, a local directory,
a Git repository, or skipping. Otherwise add a source afterward as shown above.

Advanced
  mogent init --source qmr=../qmr-agents-library
  mogent init --plan                    Start with selection intent in mogent.plan
  mogent init --file path               Choose the entry file's location
  mogent init --dry-run                 Preview without writing or prompting
  mogent init --no-interactive          Skip the optional Gum picker

Existing files are never overwritten. Explicit --source bypasses the picker.
For old configurations: mogent help legacy

Next: mogent help source
`,
	"source": `Connect libraries and browse content

  mogent source add qmr ../qmr-agents-library
  mogent source add team https://example.org/team/library.git
  mogent source list
  mogent source show

source list (or source ls) lists project declarations without opening libraries
or fetching Git sources. source show loads the libraries and browses content.
Local paths are relative to the project entry file; absolute paths are converted.

Advanced — browsing
  mogent source show --tree --chars unicode
  mogent source show --flat             Show complete paths
  mogent source show qmr:intro          Limit the tree to a node and its children
  mogent source list --tree             Browse all connected libraries
  mogent source add qmr ../library --dry-run

Advanced — remembered libraries
  mogent source register qmr ~/dev/cdint/qmr-agents-library
  mogent source list --known
  mogent source add qmr                 Use a remembered location in this project
  mogent source update local --file mogent.build --dry-run

Registration remembers a location for discovery and interactive init; builds
use only project declarations. Use --file path to choose another project file.
Legacy source show/content flags: mogent help legacy

Next: mogent help select
`,
	"select": `Select modules

  mogent add qmr:intro/author-ownership
  mogent select qmr:intro/author-ownership
  mogent module add qmr:intro/author-ownership
  mogent rm qmr:intro/author-ownership --diff
  mogent except qmr:workflow/decision-first

Replace qmr and the path with a connected library and a node shown by:
  mogent source show

All three commands append the same validated include to the authored build or
plan. Includes are applied in the order you select them. Use --dry-run to inspect
the exact declaration first; use --file to name a specific plan/build.

For an opt-in tree, select its direct children explicitly:
  mogent add qmr:skills --all

This turns an existing bare include "qmr:skills" into an explicit child list,
or creates that list for a new include. It is intentionally not the default:
opt-in content such as skills should not be installed accidentally.

You can also edit includes directly, for example:
  include "qmr:intro/author-ownership"

mogent add --manifest agents.yaml remains the legacy YAML placement command.

Removal and overlap review
  mogent rm qmr:intro/author-ownership --diff
  mogent module rm qmr:intro/author-ownership
  mogent optimize

rm removes one exact direct include, never an ancestor selection. --diff is a
safe preview of resulting managed-output changes and writes nothing. optimize
only identifies ancestor/descendant direct includes; repeated content can be
intentional, so it never rewrites the file.

Exception behavior
  mogent except qmr:workflow/decision-first
  mogent exclude qmr:workflow/decision-first  # alias

except finds the narrowest direct selected ancestor (here qmr:workflow) and
adds decision-first to that include's exclude list. For a curated choose branch,
it also makes accepting the source defaults explicit. It does not remove another
direct include or rewrite overlapping selections; that needs an explicit remove
or optimization decision.

Authoring order
  mogent reorder --dry-run
  mogent reorder

reorder is explicit formatting only: bindings first, then sources, then includes
in their existing order. Each include's options remain attached. It does not
deduplicate content, remove selections, or change build behavior.

Advanced — intent-first authoring
  mogent init --plan                    Start a new workspace with mogent.plan
  mogent sequence --file mogent.plan    Compile intent into a generated build

Edit includes in the plan when using this workflow; build also sequences it.
Do not edit a generated build as though it were authored.

Advanced — local customization
  mogent edit include qmr:intro/author-ownership --localize --file mogent.build

Localization copies the selected input into an editable project-local source
and rewrites its reference. Add --dry-run to preview. For a plan, pass its path.

Next: mogent help show
`,
	"show": `Preview the selection

  mogent show --build
  mogent show --issues

Preview library content before selecting it:
  mogent show qmr:intro/assistant-agent --content
  mogent show qmr:intro --content --lines 12

The build view lists ordered contributions and their destinations. The issues
view explains unresolved choices and errors. These commands do not install outputs.
Source inspection may load or fetch libraries.

Advanced
  mogent show qmr:intro/author-ownership --build
  mogent show --build --file mogent.build
  mogent show --issues --file mogent.plan
  mogent show --source                  Browse available, not selected, content
  mogent show qmr:skills/review --content
  mogent show qmr:assets/readme.md --content --headings

The build view is a contribution preview, not a rendered-file diff. An empty
build view can mean you have declared sources but have not added includes yet.
Content preview is read-only; --lines 0 (the default) shows all text. Physical
binary payloads report their byte count instead, and --headings applies only to
physical Markdown files.

Next: mogent help build
`,
	"build": `Build outputs

  mogent build
  mogent build --check

build assembles your selected content into project outputs, such as AGENTS.md
and complete skill bundles. With a plan, it sequences the plan before assembly.
--check verifies existing outputs match the expected result without installing
outputs; it is useful after building, not a preview of missing outputs.

Advanced
  mogent build --file mogent.build
  mogent build --file mogent.plan
  mogent sequence --file mogent.plan    Generate the build without installing outputs
  mogent source update local --file mogent.build --dry-run

Builds protect managed outputs against drift. For local-source input changes,
review the update preview; omit --dry-run to accept the new input hashes.
For selection previews: mogent help show
For legacy --dry-run/--force workflows: mogent help legacy

Output inspection
  mogent drift
  mogent drift --diff

drift is read-only. It reports both whether AGENTS.md still matches the last
Mogent output state and whether it matches the current configuration/source
composition. The two can differ after either direct output edits or config/source
changes. Skill-directory detail remains part of build --check in this slice.
`,
	"legacy": `Legacy workflows — explicit configuration formats

HCL v2 (mogent.hcl)
  mogent init --config mogent.hcl --source qmr=../qmr-agents-library
  mogent source list qmr --config mogent.hcl --tldr --tags
  mogent source show qmr:path --config mogent.hcl [--lines n] [--headings]
  mogent plan --config mogent.hcl
  mogent apply --config mogent.hcl
  mogent update qmr --config mogent.hcl [--accept]

YAML (agents.yaml)
  mogent init --list-templates
  mogent init --manifest agents.yaml --template minimal --source shared=../library
  mogent source add shared ../library --manifest agents.yaml
  mogent source list --manifest agents.yaml --tree
  mogent add shared:path --manifest agents.yaml
  mogent build --manifest agents.yaml [--dry-run] [--force]
  mogent status --manifest agents.yaml
  mogent drift --manifest agents.yaml [--diff]
  mogent coverage --manifest agents.yaml
  mogent localize <heading-path> --manifest agents.yaml
  mogent tui --manifest agents.yaml

Advanced
  YAML source pin/update, add placement, drift import/reject, and template flags
  belong to these legacy handlers; they do not imply native plan/build support.
  Use explicit --config or --manifest when mixing formats in one workspace.

All command signatures: mogent help commands
Current workflow: mogent --help
`,
	"tools": `Utilities

  mogent completion bash|zsh            Generate shell completion
  mogent complete commands             List command completion candidates
  mogent complete flags                List flag completion candidates
  mogent mint                          Generate an unused proquint handle

Advanced — legacy library sidecars
  mogent lib init <path>
  mogent lib scan <path>
  mogent lib check <path>

lib manages legacy YAML sidecars, not the native mogent.lib catalog language.
Source-reference completion currently targets the legacy YAML workflow.

All command signatures: mogent help commands
Current workflow: mogent --help
`,
}

func writeUsage(output io.Writer) error {
	return writeHelpText(output, workflowHelp)
}

var helpFlagPattern = regexp.MustCompile(`--[a-z][a-z-]*|\B-h\b`)

// Styling is presentation-only: pipes, redirected files and NO_COLOR retain
// exactly the same readable plain-text help.
func writeHelpText(output io.Writer, text string) error {
	file, terminal := output.(*os.File)
	_, noColor := os.LookupEnv("NO_COLOR")
	if terminal && term.IsTerminal(file.Fd()) && !noColor && os.Getenv("TERM") != "dumb" {
		text = boldHelp(text)
	}
	_, err := io.WriteString(output, text)
	return err
}

func boldHelp(text string) string {
	bold := func(value string) string { return "\x1b[1m" + value + "\x1b[22m" }
	names := []string{"source add", "source list, source ls", "source list", "source ls", "source register", "source show", "source pin", "source update", "edit include"}
	for _, command := range commands() {
		names = append(names, command.name)
	}
	sort.Slice(names, func(i, j int) bool { return len(names[i]) > len(names[j]) })
	lines := strings.Split(text, "\n")
	for index, line := range lines {
		switch {
		case index == 0, line == "Common commands, grouped by workflow:", line == "Further help", strings.HasPrefix(line, "Advanced"), line == "HCL v2 (mogent.hcl)", line == "YAML (agents.yaml)":
			lines[index] = bold(line)
		case strings.Contains(line, " (advanced: mogent help "):
			title, rest, _ := strings.Cut(line, " (advanced:")
			lines[index] = bold(title) + " (advanced:" + rest
		case strings.HasPrefix(line, "Usage:"):
			lines[index] = bold("Usage:") + strings.TrimPrefix(line, "Usage:")
		case strings.HasPrefix(line, " "):
			prefix := line[:len(line)-len(strings.TrimLeft(line, " "))]
			content := strings.TrimPrefix(line, prefix)
			if strings.HasPrefix(content, "mogent ") {
				prefix += "mogent "
				content = strings.TrimPrefix(content, "mogent ")
			}
			for _, name := range names {
				if content == name || strings.HasPrefix(content, name+" ") {
					content = bold(name) + strings.TrimPrefix(content, name)
					break
				}
			}
			content = helpFlagPattern.ReplaceAllStringFunc(content, bold)
			lines[index] = prefix + content
		}
	}
	return strings.Join(lines, "\n")
}

func helpRequested(args []string) bool {
	for _, arg := range args {
		if arg == "--" {
			break
		}
		if arg == "--help" || arg == "-h" {
			return true
		}
	}
	return false
}

func commandHelpPath(args []string) []string {
	if hasFlag(args[1:], "--manifest") || hasFlag(args[1:], "-manifest") || hasFlag(args[1:], "--config") || hasFlag(args[1:], "-config") {
		return []string{"legacy"}
	}
	path := []string{args[0]}
	if args[0] == "source" && len(args) > 1 && !strings.HasPrefix(args[1], "-") {
		path = append(path, args[1])
	}
	return path
}

func writeHelp(path []string, output io.Writer) error {
	if len(path) == 0 {
		return writeUsage(output)
	}
	if len(path) > 2 || len(path) == 2 && path[0] != "source" {
		return fmt.Errorf("help accepts one topic, or help source <subcommand>; use mogent --help")
	}
	topic := path[0]
	if len(path) == 2 {
		switch path[1] {
		case "add", "list", "ls", "register", "update":
		case "show", "pin":
			topic = "legacy"
		default:
			return fmt.Errorf("unknown source help topic %q; use mogent help source", path[1])
		}
	}
	switch topic {
	case "--all":
		return writeHelpText(output, expandedWorkflowHelp)
	case "commands", "all":
		var reference strings.Builder
		reference.WriteString("All command signatures (including legacy flags)\nFor the guided workflow: mogent --help\n\n")
		for _, command := range commands() {
			reference.WriteString("  mogent " + command.usage + "\n")
		}
		return writeHelpText(output, reference.String())
	case "sequence", "edit", "select", "module", "rm", "except", "exclude", "reorder", "optimise", "optimize":
		topic = "select"
	case "plan", "apply", "update", "status", "drift", "coverage", "add", "localize", "tui":
		topic = "legacy"
	case "lib", "mint", "complete", "completion":
		topic = "tools"
	}
	text, found := workflowTopics[topic]
	if !found {
		return fmt.Errorf("unknown help topic %q; use --all, init, source, select, show, build, legacy, tools, or commands", path[0])
	}
	return writeHelpText(output, text)
}
