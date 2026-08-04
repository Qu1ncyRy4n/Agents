// Package cli exposes the deliberately small milestone-one command surface.
package cli

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"

	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
	"github.com/Qu1ncyRy4n/Agents/internal/navigator"
	"github.com/Qu1ncyRy4n/Agents/internal/render"
	"github.com/Qu1ncyRy4n/Agents/internal/state"
	"github.com/Qu1ncyRy4n/Agents/internal/workspace"
)

func Run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		if _, err := fmt.Fprintln(stdout, "Usage: mogent build [--manifest agents.yaml] [--force]\n       mogent status [--manifest agents.yaml]\n       mogent coverage [--manifest agents.yaml]\n       mogent tui [--manifest agents.yaml]"); err != nil {
			return fmt.Errorf("write usage: %w", err)
		}
		return nil
	}
	switch args[0] {
	case "build":
		return runBuild(args[1:], stdout, stderr)
	case "status":
		return runStatus(args[1:], stdout, stderr)
	case "coverage":
		return runCoverage(args[1:], stdout, stderr)
	case "tui":
		return runTUI(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q; use build, status, coverage, or tui", args[0])
	}
}

func runBuild(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("build", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	force := flags.Bool("force", false, "replace an untracked or directly edited output")
	if err := flags.Parse(args); err != nil {
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
	if err := flags.Parse(args); err != nil {
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

func runCoverage(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("coverage", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	unusedOnly := flags.Bool("unused-only", false, "show only unused source references")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("coverage accepts no positional arguments")
	}
	session, err := workspace.New(*manifestFile)
	if err != nil {
		return err
	}
	coverage := session.Coverage()
	if *unusedOnly {
		return writeUnusedCoverage(stdout, coverage)
	}
	for _, source := range coverage.Sources {
		if _, err := fmt.Fprintf(stdout, "Source %s: %s\n", source.Alias, source.Path); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
		if _, err := fmt.Fprintf(stdout, "Included: %d/%d\n", source.Included, source.Total); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
		if len(source.Unused) == 0 {
			if _, err := fmt.Fprintln(stdout, "Unused: none"); err != nil {
				return fmt.Errorf("write coverage: %w", err)
			}
		} else {
			if _, err := fmt.Fprintln(stdout, "Unused:"); err != nil {
				return fmt.Errorf("write coverage: %w", err)
			}
			for _, node := range source.Unused {
				if _, err := fmt.Fprintf(stdout, "- %s  %s:%s\n", node.Heading, source.Alias, node.Path); err != nil {
					return fmt.Errorf("write coverage: %w", err)
				}
			}
		}
		if _, err := fmt.Fprintln(stdout); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
	}
	return nil
}

func writeUnusedCoverage(stdout io.Writer, coverage workspace.Coverage) error {
	for _, source := range coverage.Sources {
		if len(source.Unused) == 0 {
			continue
		}
		if _, err := fmt.Fprintf(stdout, "%s  %s  unused %d/%d\n", source.Alias, source.Path, len(source.Unused), source.Total); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
		for _, node := range source.Unused {
			if _, err := fmt.Fprintf(stdout, "  %s  %s\n", source.Alias+":"+node.Path, node.Heading); err != nil {
				return fmt.Errorf("write coverage: %w", err)
			}
		}
		if _, err := fmt.Fprintln(stdout); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
	}
	return nil
}

func runTUI(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("tui", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("tui accepts no positional arguments")
	}
	return navigator.Run(*manifestFile, stdout)
}
