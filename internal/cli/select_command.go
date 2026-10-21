package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/Qu1ncyRy4n/Agents/assembly"
)

func nativeSelectionRequested(args []string) bool {
	if hasFlag(args, "--manifest") || hasFlag(args, "-manifest") {
		return false
	}
	if hasFlag(args, "--file") || hasFlag(args, "-file") {
		return true
	}
	return newSourceRequested(nil)
}

func runNativeSelect(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or authored build")
	dryRun := flags.Bool("dry-run", false, "preview the include without writing")
	all := flags.Bool("all", false, "select all direct children")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--dry-run": false, "-dry-run": false, "--all": false, "-all": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("add requires one module reference, for example qmr:intro/author-ownership")
	}
	entryPath, err := projectEntryFile(*path)
	if err != nil {
		return err
	}
	contents, err := assembly.SelectWithOptions(entryPath, flags.Arg(0), assembly.SelectOptions{All: *all, DryRun: *dryRun})
	if err != nil {
		return err
	}
	if *dryRun {
		_, err = fmt.Fprint(stdout, string(contents))
	} else {
		_, err = fmt.Fprintf(stdout, "Selected %s in %s\nNext: mogent show --build --file %s\n", flags.Arg(0), entryPath, entryPath)
	}
	return err
}

func runSelect(args []string, stdout, stderr io.Writer) error {
	return runNativeSelect(args, stdout, stderr)
}

func runModule(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("module requires add or rm; use mogent module add|rm source:path")
	}
	switch args[0] {
	case "add":
		return runNativeSelect(args[1:], stdout, stderr)
	case "rm":
		return runRemove(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown module command %q; use mogent module add|rm source:path", args[0])
	}
}
