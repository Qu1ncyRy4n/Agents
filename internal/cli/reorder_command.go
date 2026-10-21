package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/Qu1ncyRy4n/Agents/assembly"
)

func runReorder(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("reorder", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or authored build")
	dryRun := flags.Bool("dry-run", false, "preview reordered configuration without writing")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--dry-run": false, "-dry-run": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("reorder accepts no positional arguments")
	}
	entryPath, err := projectEntryFile(*path)
	if err != nil {
		return err
	}
	contents, changed, err := assembly.Reorder(entryPath, *dryRun)
	if err != nil {
		return err
	}
	if *dryRun {
		_, err = fmt.Fprint(stdout, string(contents))
	} else if changed {
		_, err = fmt.Fprintf(stdout, "Reordered %s: bindings, sources, then ordered includes\n", entryPath)
	} else {
		_, err = fmt.Fprintf(stdout, "%s is already ordered\n", entryPath)
	}
	return err
}
