package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/workspace"
)

func runAdd(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	under := flags.String("under", "", "manifest heading path to append under")
	appendRoot := flags.Bool("append", false, "append to the end of the document")
	first := flags.Bool("first", false, "insert first beneath --under")
	last := flags.Bool("last", false, "insert last beneath --under (default)")
	before := flags.String("before", "", "insert before this manifest heading path")
	after := flags.String("after", "", "insert after this manifest heading path")
	heading := flags.String("heading", "", "rendered heading to use in the manifest")
	dryRun := flags.Bool("dry-run", false, "preview the change without writing")
	previewMode := flags.String("preview", "summary", "dry-run preview mode: summary, patch, tree, or full")
	rebuild := flags.Bool("rebuild", false, "also rebuild the generated AGENTS.md")
	args = reorderArgs(args, map[string]bool{
		"-manifest":  true,
		"--manifest": true,
		"-under":     true,
		"--under":    true,
		"-heading":   true,
		"--heading":  true,
		"-before":    true,
		"--before":   true,
		"-after":     true,
		"--after":    true,
		"-preview":   true,
		"--preview":  true,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("add requires exactly one source reference")
	}
	session, err := workspace.New(*manifestFile)
	if err != nil {
		return err
	}
	result, err := session.AddSource(workspace.AddOptions{
		Reference: flags.Arg(0),
		Under:     *under,
		Append:    *appendRoot,
		First:     *first,
		Last:      *last,
		Before:    *before,
		After:     *after,
		Heading:   *heading,
		Rebuild:   *rebuild,
	}, *dryRun)
	if err != nil {
		return decorateSourceReferenceError(flags.Arg(0), err)
	}
	if *dryRun {
		if _, err := fmt.Fprintln(stdout, "Dry run: no files written"); err != nil {
			return fmt.Errorf("write add result: %w", err)
		}
		if err := writeAddPreview(stdout, result, *previewMode); err != nil {
			return err
		}
		return nil
	} else if result.Rebuilt {
		if _, err := fmt.Fprintln(stdout, "Wrote manifest and rebuilt output"); err != nil {
			return fmt.Errorf("write add result: %w", err)
		}
	} else if result.WroteManifest {
		if _, err := fmt.Fprintln(stdout, "Wrote manifest"); err != nil {
			return fmt.Errorf("write add result: %w", err)
		}
	}
	if _, err := fmt.Fprintf(stdout, "Added: %s\n", result.Heading); err != nil {
		return fmt.Errorf("write add result: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "From: %s\n", result.Reference); err != nil {
		return fmt.Errorf("write add result: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "Under: %s\n", parentLabelForOutput(result.ParentPath)); err != nil {
		return fmt.Errorf("write add result: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "Placement: %s\n", result.Placement); err != nil {
		return fmt.Errorf("write add result: %w", err)
	}
	if err := writeAddRelations(stdout, result.Relations); err != nil {
		return err
	}
	return nil
}

func writeAddPreview(stdout io.Writer, result *workspace.AddResult, mode string) error {
	switch strings.ToLower(mode) {
	case "summary":
		return writeAddSummary(stdout, result)
	case "patch":
		return writeAddPatch(stdout, result)
	case "tree":
		return writeAddTree(stdout, result)
	case "full":
		return writeAddFull(stdout, result)
	default:
		return fmt.Errorf("unknown preview mode %q; use summary, patch, tree, or full", mode)
	}
}

func writeAddSummary(stdout io.Writer, result *workspace.AddResult) error {
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "Manifest change:"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "  Add: %s / %s\n", parentLabelForOutput(result.ParentPath), result.Heading); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "  From: %s\n", result.Reference); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "  Placement: %s\n", result.Placement); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if err := writeAddRelations(stdout, result.Relations); err != nil {
		return err
	}
	return nil
}

func writeAddPatch(stdout io.Writer, result *workspace.AddResult) error {
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "--- agents.yaml"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "+++ agents.yaml"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if err := writeAddHunk(stdout, result.ManifestHunk); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(stdout, "--- AGENTS.md"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "+++ AGENTS.md"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if err := writeAddHunk(stdout, result.OutputHunk); err != nil {
		return err
	}
	return nil
}

func writeAddHunk(stdout io.Writer, hunk workspace.DiffHunk) error {
	if _, err := fmt.Fprintf(stdout, "@@ -%d,%d +%d,%d @@\n", hunk.OldStart, hunk.OldCount, hunk.NewStart, hunk.NewCount); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	for _, line := range hunk.Lines {
		if _, err := fmt.Fprintln(stdout, "+"+line); err != nil {
			return fmt.Errorf("write add preview: %w", err)
		}
	}
	return nil
}

func writeAddTree(stdout io.Writer, result *workspace.AddResult) error {
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "Document tree:"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, result.Tree); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if result.SourceTree != "" {
		if _, err := fmt.Fprintln(stdout, "\nInherited source subtree:"); err != nil {
			return fmt.Errorf("write add preview: %w", err)
		}
		if _, err := fmt.Fprintln(stdout, result.SourceTree); err != nil {
			return fmt.Errorf("write add preview: %w", err)
		}
	}
	if err := writeAddRelations(stdout, result.Relations); err != nil {
		return err
	}
	return nil
}

func writeAddRelations(stdout io.Writer, relations []workspace.AddRelation) error {
	if len(relations) == 0 {
		return nil
	}
	if _, err := fmt.Fprintln(stdout, "\nRelated existing selections:"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	for _, relation := range relations {
		label := "review"
		if relation.Exact {
			label = "overlap"
		}
		if _, err := fmt.Fprintf(stdout, "  - [%s] %s - %s\n", label, relation.Reference, relation.Reason); err != nil {
			return fmt.Errorf("write add preview: %w", err)
		}
	}
	return nil
}

func writeAddFull(stdout io.Writer, result *workspace.AddResult) error {
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "Rendered preview:"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, result.Preview); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	return nil
}

func parentLabelForOutput(path string) string {
	if path == "" {
		return "document root"
	}
	return path
}
