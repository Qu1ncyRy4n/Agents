package cli

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

type sourceBindings map[string]string

func (s sourceBindings) String() string {
	values := make([]string, 0, len(s))
	for alias, path := range s {
		values = append(values, alias+"="+path)
	}
	sort.Strings(values)
	return strings.Join(values, ",")
}

func (s sourceBindings) Set(value string) error {
	alias, path, found := strings.Cut(value, "=")
	alias, path = strings.TrimSpace(alias), strings.TrimSpace(path)
	if !found || alias == "" || path == "" {
		return fmt.Errorf("source binding %q must use alias=path-or-url", value)
	}
	if _, exists := s[alias]; exists {
		return fmt.Errorf("source alias %q was provided twice", alias)
	}
	s[alias] = path
	return nil
}

func runInit(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	flags.SetOutput(stderr)
	entryFile := flags.String("file", "", "plan/build path (default mogent.build, or mogent.plan with --plan)")
	plan := flags.Bool("plan", false, "create selection intent instead of an authored build")
	noInteractive := flags.Bool("no-interactive", false, "skip the optional Gum source picker")
	dryRun := flags.Bool("dry-run", false, "preview initialization without writing")
	sources := make(sourceBindings)
	flags.Var(sources, "source", "source declaration alias=local-path|git-url; repeat as needed (does not select content)")
	args = reorderArgs(args, map[string]bool{
		"-file": true, "--file": true, "-plan": false, "--plan": false,
		"-no-interactive": false, "--no-interactive": false, "-source": true, "--source": true,
		"-dry-run": false, "--dry-run": false,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("init accepts no positional arguments")
	}
	kind := syntax.Build
	if *plan {
		kind = syntax.Plan
	}
	var picker sourcePicker
	if !*dryRun && !*noInteractive && len(sources) == 0 {
		picker = interactiveInitPicker(stdout, stderr)
	}
	result, err := initializeProject(assembly.InitOptions{Path: *entryFile, Kind: kind, Sources: sources, DryRun: *dryRun}, picker)
	if err != nil {
		return err
	}
	if *dryRun {
		_, err = fmt.Fprint(stdout, string(result.Content))
		return err
	}
	if _, err := fmt.Fprintf(stdout, "Created %s\nSources are declarations only; choose includes before building.\n", result.Path); err != nil {
		return err
	}
	parsed, err := syntax.Parse(result.Path, result.Content)
	if err != nil {
		return err
	}
	declarations, err := assembly.SourceDeclarations(parsed)
	if err != nil {
		return err
	}
	if len(declarations) == 0 {
		return writeNoSources(stdout, result.Path)
	}
	_, err = fmt.Fprintf(stdout, "Next: mogent source tree --file %s\n", result.Path)
	return err
}
