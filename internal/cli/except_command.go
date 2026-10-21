package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/Qu1ncyRy4n/Agents/assembly"
)

func runExcept(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("except", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or authored build")
	dryRun := flags.Bool("dry-run", false, "preview the exception without writing")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--dry-run": false, "-dry-run": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("except requires one descendant reference, for example qmr:workflow/decision-first")
	}
	entryPath, err := projectEntryFile(*path)
	if err != nil {
		return err
	}
	result, err := assembly.Except(entryPath, flags.Arg(0), *dryRun)
	if err != nil {
		return err
	}
	if *dryRun {
		_, err = fmt.Fprint(stdout, string(result.Content))
	} else {
		_, err = fmt.Fprintf(stdout, "Excepted %s from %s in %s\nNext: mogent show --build --file %s\n", result.Reference, result.From, result.Path, result.Path)
	}
	return err
}
