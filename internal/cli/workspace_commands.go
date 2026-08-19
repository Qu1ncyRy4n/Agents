package cli

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"github.com/Qu1ncyRy4n/Agents/manifest"
	"github.com/Qu1ncyRy4n/Agents/render"
	"github.com/Qu1ncyRy4n/Agents/state"
	"github.com/Qu1ncyRy4n/Agents/workspace"
)

func runDrift(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("drift", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	importHeading := flags.String("import", "", "import direct edits from one manifest heading path")
	from := flags.String("from", "", "source reference to localize for a composed entry")
	reject := flags.Bool("reject", false, "reject direct edits and rebuild from the manifest")
	force := flags.Bool("force", false, "confirm rejection of direct edits")
	args = reorderArgs(args, map[string]bool{
		"-manifest": true, "--manifest": true,
		"-import": true, "--import": true,
		"-from": true, "--from": true,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("drift accepts no positional arguments")
	}
	if *reject && *importHeading != "" {
		return fmt.Errorf("drift accepts only one of --reject or --import")
	}
	session, err := workspace.New(*manifestFile)
	if err != nil {
		return err
	}
	if *reject {
		if err := session.RejectDrift(*force); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(stdout, "Rejected direct edits and rebuilt output"); err != nil {
			return fmt.Errorf("write drift result: %w", err)
		}
		return nil
	}
	if *importHeading != "" {
		result, err := session.ImportDrift(*importHeading, *from)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(stdout, "Imported direct edits from %s into %s\n", result.ManifestHeading, result.LocalReference); err != nil {
			return fmt.Errorf("write drift result: %w", err)
		}
		if _, err := fmt.Fprintf(stdout, "Local file: %s\n", result.LocalPath); err != nil {
			return fmt.Errorf("write drift result: %w", err)
		}
		return nil
	}
	report, err := session.Drift()
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdout, "Output: %s\nDrift status: %s\n", report.OutputPath, report.Status); err != nil {
		return fmt.Errorf("write drift result: %w", err)
	}
	if report.DirectEdits {
		if _, err := fmt.Fprintln(stdout, "Direct edits detected. Use --import <manifest-heading-path> or review and use --reject --force."); err != nil {
			return fmt.Errorf("write drift result: %w", err)
		}
	}
	return nil
}

func runLocalize(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("localize", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	from := flags.String("from", "", "source reference to replace when the entry composes several sources")
	dryRun := flags.Bool("dry-run", false, "preview localization without writing files")
	rebuild := flags.Bool("rebuild", false, "also rebuild generated output")
	args = reorderArgs(args, map[string]bool{
		"-manifest": true, "--manifest": true,
		"-from": true, "--from": true,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("localize requires exactly one manifest heading path")
	}
	session, err := workspace.New(*manifestFile)
	if err != nil {
		return err
	}
	result, err := session.Localize(workspace.LocalizeOptions{
		ManifestHeading: flags.Arg(0),
		From:            *from,
		DryRun:          *dryRun,
		Rebuild:         *rebuild,
	})
	if err != nil {
		return err
	}
	if *dryRun {
		if _, err := fmt.Fprintln(stdout, "Dry run: no files written"); err != nil {
			return fmt.Errorf("write localize result: %w", err)
		}
	}
	for _, line := range []string{
		"Manifest entry: " + result.ManifestHeading,
		"Source: " + result.SourceReference,
		"Local: " + result.LocalReference,
		"Local file: " + result.LocalPath,
		"Provenance: " + result.ProvenancePath,
	} {
		if _, err := fmt.Fprintln(stdout, line); err != nil {
			return fmt.Errorf("write localize result: %w", err)
		}
	}
	if result.Rebuilt {
		if _, err := fmt.Fprintln(stdout, "Rebuilt output"); err != nil {
			return fmt.Errorf("write localize result: %w", err)
		}
	}
	return nil
}

func runBuild(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("build", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	force := flags.Bool("force", false, "replace an untracked or directly edited output")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("build accepts no positional arguments")
	}
	value, manifestPath, err := manifest.Load(*manifestFile)
	if err != nil {
		return err
	}
	result, err := render.Build(value, manifestPath)
	if err != nil {
		return err
	}
	outputPath := value.Output
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(filepath.Dir(manifestPath), outputPath)
	}
	statePath := filepath.Join(filepath.Dir(manifestPath), ".mogent", "state.json")
	if err := state.CheckOverwrite(outputPath, statePath, *force); err != nil {
		return err
	}
	if err := render.WriteAtomically(outputPath, result.Content); err != nil {
		return err
	}
	if err := state.Write(statePath, result.Content); err != nil {
		return err
	}
	for _, warning := range result.Warnings {
		if _, err := fmt.Fprintf(stderr, "mogent: warning: %s\n", warning); err != nil {
			return fmt.Errorf("write warning: %w", err)
		}
	}
	if _, err := fmt.Fprintf(stdout, "Wrote %s\n", outputPath); err != nil {
		return fmt.Errorf("write success message: %w", err)
	}
	return nil
}

func runStatus(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("status", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("status accepts no positional arguments")
	}
	session, err := workspace.New(*manifestFile)
	if err != nil {
		return err
	}
	status, err := session.Status()
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdout, "Manifest: %s\n", status.ManifestPath); err != nil {
		return fmt.Errorf("write status: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "Output: %s\n", status.OutputPath); err != nil {
		return fmt.Errorf("write status: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "Output status: %s\n", status.Output); err != nil {
		return fmt.Errorf("write status: %w", err)
	}
	if hint := statusHint(status.Output); hint != "" {
		if _, err := fmt.Fprintf(stdout, "Hint: %s\n", hint); err != nil {
			return fmt.Errorf("write status: %w", err)
		}
	}
	if _, err := fmt.Fprintf(stdout, "Sources: %d\n", len(status.Sources)); err != nil {
		return fmt.Errorf("write status: %w", err)
	}
	for _, source := range status.Sources {
		if _, err := fmt.Fprintf(stdout, "- %s: %s (%d nodes)\n", source.Alias, source.Path, source.Nodes); err != nil {
			return fmt.Errorf("write status: %w", err)
		}
	}
	return nil
}

func statusHint(status workspace.OutputStatus) string {
	switch status {
	case workspace.StatusMissing:
		return "run `mogent build` to create the output"
	case workspace.StatusStale:
		return "run `mogent build` to refresh the output"
	case workspace.StatusDirectEdits:
		return "AGENTS.md has direct edits; review them before running `mogent build --force`"
	case workspace.StatusUntracked:
		return "AGENTS.md was not generated by mogent; use `mogent build --force` only after review"
	case workspace.StatusUpToDate:
		return "clean"
	default:
		return ""
	}
}
