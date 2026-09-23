package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Qu1ncyRy4n/Agents/library"
	"github.com/Qu1ncyRy4n/Agents/manifest"
	"github.com/Qu1ncyRy4n/Agents/render"
)

func runLib(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("lib requires init, scan, or check")
	}
	switch args[0] {
	case "init":
		return runLibInit(args[1:], stdout, stderr)
	case "scan":
		return runLibScan(args[1:], stdout, stderr)
	case "check":
		return runLibCheck(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown lib command %q", args[0])
	}
}
func runLibInit(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("lib init", flag.ContinueOnError)
	flags.SetOutput(stderr)
	dry := flags.Bool("dry-run", false, "preview without writing")
	args = reorderArgs(args, map[string]bool{"-dry-run": true, "--dry-run": true})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("lib init requires exactly one library path")
	}
	path := flags.Arg(0)
	target := filepath.Join(path, library.SidecarFile)
	if _, err := os.Lstat(target); err == nil {
		return fmt.Errorf("refusing to overwrite existing %s", target)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	sidecar, err := library.Scan(path)
	if err != nil {
		return err
	}
	contents, err := library.MarshalSidecar(sidecar)
	if err != nil {
		return err
	}
	if *dry {
		_, err = fmt.Fprintf(stdout, "Dry run: no files written\n%s", contents)
		return err
	}
	if err := os.WriteFile(target, contents, 0o644); err != nil {
		return fmt.Errorf("write %s: %w", target, err)
	}
	_, err = fmt.Fprintf(stdout, "Created %s\n%s", target, contents)
	return err
}
func runLibScan(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("lib scan", flag.ContinueOnError)
	flags.SetOutput(stderr)
	dry := flags.Bool("dry-run", false, "preview only")
	args = reorderArgs(args, map[string]bool{"-dry-run": true, "--dry-run": true})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("lib scan requires exactly one library path")
	}
	sidecar, err := library.Scan(flags.Arg(0))
	if err != nil {
		return err
	}
	contents, err := library.MarshalSidecar(sidecar)
	if err != nil {
		return err
	}
	if *dry {
		_, err = fmt.Fprintln(stdout, "Dry run: no files written")
		if err != nil {
			return err
		}
	}
	_, err = stdout.Write(contents)
	return err
}
func runLibCheck(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("lib check", flag.ContinueOnError)
	flags.SetOutput(stderr)
	source := flags.String("source", "", "source alias")
	manifestPath := flags.String("manifest", "agents.yaml", "path to manifest")
	args = reorderArgs(args, map[string]bool{"-source": true, "--source": true, "-manifest": true, "--manifest": true})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	var root string
	if *source != "" {
		if flags.NArg() != 0 {
			return fmt.Errorf("lib check with --source accepts no path")
		}
		value, loaded, err := manifest.Load(*manifestPath)
		if err != nil {
			return err
		}
		root, err = render.ResolveSourcePath(value, loaded, *source)
		if err != nil {
			return err
		}
	} else {
		if flags.NArg() != 1 {
			return fmt.Errorf("lib check requires exactly one library path")
		}
		root = flags.Arg(0)
	}
	report, err := library.Check(root)
	if err != nil {
		return err
	}
	for _, warning := range report.Warnings {
		if _, err := fmt.Fprintf(stderr, "mogent: warning: %s\n", warning); err != nil {
			return err
		}
	}
	for _, validation := range report.Errors {
		if _, err := fmt.Fprintf(stderr, "mogent: error: %s\n", validation); err != nil {
			return err
		}
	}
	if len(report.Errors) > 0 {
		return fmt.Errorf("library check failed with %d error(s)", len(report.Errors))
	}
	_, err = fmt.Fprintln(stdout, "Library sidecar check passed")
	return err
}
