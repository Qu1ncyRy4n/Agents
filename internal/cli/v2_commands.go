package cli

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/v2"
)

// runPlan shows what apply would write: one unified diff per Markdown output
// path, then diagnostics. It never writes files.
func runPlan(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("plan", flag.ContinueOnError)
	flags.SetOutput(stderr)
	configPath := flags.String("config", v2.ConfigFile, "path to v2 HCL configuration")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("plan accepts no positional arguments")
	}
	config, err := v2.Load(*configPath)
	if err != nil {
		return err
	}
	result, err := v2.PlanConfig(config)
	if result != nil {
		if writeErr := writePlanChanges(stdout, result); writeErr != nil {
			return writeErr
		}
		if writeErr := writePlanDiagnostics(stderr, result.Plan); writeErr != nil {
			return writeErr
		}
	}
	if err != nil {
		return err
	}
	if result.HasErrors() {
		return fmt.Errorf("plan contains errors")
	}
	return nil
}

// runApply writes every planned output transactionally after the same checks
// plan performs. Unmanaged or hand-edited outputs require --force.
func runApply(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("apply", flag.ContinueOnError)
	flags.SetOutput(stderr)
	configPath := flags.String("config", v2.ConfigFile, "path to v2 HCL configuration")
	force := flags.Bool("force", false, "replace an unmanaged or directly edited output")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("apply accepts no positional arguments")
	}
	config, err := v2.Load(*configPath)
	if err != nil {
		return err
	}
	result, err := v2.Apply(config, *force)
	if result != nil {
		if writeErr := writePlanDiagnostics(stderr, result.Plan); writeErr != nil {
			return writeErr
		}
	}
	if err != nil {
		return err
	}
	paths := make([]string, 0, len(result.Changes))
	for _, change := range result.Changes {
		paths = append(paths, change.Path)
	}
	if _, err := fmt.Fprintf(stdout, "Applied %d file(s): %s\n", len(paths), strings.Join(paths, ", ")); err != nil {
		return fmt.Errorf("write apply result: %w", err)
	}
	for _, alias := range sortedKeys(result.Pins) {
		if _, err := fmt.Fprintf(stdout, "Pinned %s to %s in %s\n", alias, result.Pins[alias], *configPath); err != nil {
			return fmt.Errorf("write apply result: %w", err)
		}
	}
	return nil
}

// runUpdate re-resolves git sources at their ref and shows the output diff
// the new commit would produce. --accept rewrites the commit in mogent.hcl.
func runUpdate(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("update", flag.ContinueOnError)
	flags.SetOutput(stderr)
	configPath := flags.String("config", v2.ConfigFile, "path to v2 HCL configuration")
	accept := flags.Bool("accept", false, "write the new commit into the configuration")
	args = reorderArgs(args, map[string]bool{"-config": true, "--config": true, "-accept": false, "--accept": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() > 1 {
		return fmt.Errorf("update accepts at most one source alias")
	}
	config, err := v2.Load(*configPath)
	if err != nil {
		return err
	}
	results, err := v2.Update(config, flags.Arg(0), *accept)
	hasErrors := false
	for _, update := range results {
		if writeErr := writeUpdate(stdout, stderr, update, *accept, *configPath); writeErr != nil {
			return writeErr
		}
		for _, diagnostic := range update.Diagnostics {
			if diagnostic.Severity == v2.SeverityError {
				hasErrors = true
			}
		}
	}
	if err != nil {
		return err
	}
	if hasErrors {
		return fmt.Errorf("update produces plan errors; the new commit was not accepted")
	}
	return nil
}

func writeUpdate(stdout, stderr io.Writer, update v2.UpdateResult, accepted bool, configPath string) error {
	if !update.Changed() {
		_, err := fmt.Fprintf(stdout, "%s: up to date at %s\n", update.Alias, update.OldCommit)
		return err
	}
	if _, err := fmt.Fprintf(stdout, "%s: %s -> %s\n", update.Alias, update.OldCommit, update.NewCommit); err != nil {
		return err
	}
	for _, change := range update.Changes {
		if _, err := fmt.Fprintf(stdout, "%s: %s\n", change.Path, change.Status); err != nil {
			return err
		}
		if change.Diff != "" {
			if _, err := fmt.Fprint(stdout, change.Diff); err != nil {
				return err
			}
		}
	}
	for _, diagnostic := range update.Diagnostics {
		if _, err := fmt.Fprintln(stderr, diagnostic.String()); err != nil {
			return err
		}
	}
	if accepted && len(update.Changes) > 0 {
		_, err := fmt.Fprintf(stdout, "Pinned %s to %s in %s\n", update.Alias, update.NewCommit, configPath)
		return err
	}
	if !accepted {
		_, err := fmt.Fprintf(stdout, "Run `mogent update %s --accept` to pin %s\n", update.Alias, update.NewCommit)
		return err
	}
	return nil
}

func sortedKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

func writePlanChanges(stdout io.Writer, result *v2.Result) error {
	counts := map[v2.FileStatus]int{}
	for _, change := range result.Changes {
		counts[change.Status]++
		if _, err := fmt.Fprintf(stdout, "%s: %s\n", change.Path, change.Status); err != nil {
			return fmt.Errorf("write plan change: %w", err)
		}
		if change.Diff == "" {
			continue
		}
		if _, err := fmt.Fprint(stdout, change.Diff); err != nil {
			return fmt.Errorf("write plan diff: %w", err)
		}
	}
	for _, output := range result.Plan.Outputs {
		if output.Kind != "tree" {
			continue
		}
		for _, path := range output.Paths {
			if _, err := fmt.Fprintf(stdout, "%s: skipped (tree outputs are not yet written by apply)\n", path); err != nil {
				return fmt.Errorf("write plan change: %w", err)
			}
		}
	}
	if result.HasErrors() {
		return nil
	}
	if _, err := fmt.Fprintf(stdout, "Plan: %d to add, %d to change, %d unchanged\n", counts[v2.FileNew], counts[v2.FileChanged], counts[v2.FileUnchanged]); err != nil {
		return fmt.Errorf("write plan summary: %w", err)
	}
	return nil
}

func writePlanDiagnostics(stderr io.Writer, plan *v2.Plan) error {
	for _, diagnostic := range plan.Diagnostics {
		if _, err := fmt.Fprintln(stderr, diagnostic.String()); err != nil {
			return fmt.Errorf("write diagnostic: %w", err)
		}
	}
	return nil
}
