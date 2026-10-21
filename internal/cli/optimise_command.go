package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/Qu1ncyRy4n/Agents/assembly"
)

// runOptimise is intentionally an inspection command. Overlap can be a useful
// deliberate repetition, so it names candidates without changing selection
// order, exceptions, or generated outputs.
func runOptimise(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("optimize", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or build")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("optimize accepts no positional arguments")
	}
	entryPath, err := projectEntryFile(*path)
	if err != nil {
		return err
	}
	file, err := assembly.LoadFile(entryPath)
	if err != nil {
		return err
	}
	overlaps := assembly.DirectIncludeOverlaps(file)
	if len(overlaps) == 0 {
		_, err = fmt.Fprintln(stdout, "No direct ancestor/descendant include overlaps found. No files written.")
		return err
	}
	for _, overlap := range overlaps {
		if _, err := fmt.Fprintf(stdout, "Overlap: %s already includes descendant %s\n  Deliberate repetition is allowed. To preview removing the direct descendant: mogent rm %s --diff --file %s\n", overlap.Ancestor, overlap.Descendant, overlap.Descendant, entryPath); err != nil {
			return err
		}
	}
	_, err = fmt.Fprintln(stdout, "No files written; optimize never rewrites selections automatically.")
	return err
}
