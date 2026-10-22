package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/workspace"
)

// runOptimise previews removal of redundant direct descendants by default.
// Repetition can be deliberate, so changing an authored selection requires the
// explicit --apply confirmation.
func runOptimise(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("optimize", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or build")
	apply := flags.Bool("apply", false, "remove redundant direct descendant includes")
	dryRun := flags.Bool("dry-run", false, "preview without writing (the default)")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--apply": false, "-apply": false, "--dry-run": false, "-dry-run": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("optimize accepts no positional arguments")
	}
	if *apply && *dryRun {
		return fmt.Errorf("optimize --apply cannot be combined with --dry-run")
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
	descendants := map[string]bool{}
	for _, overlap := range overlaps {
		descendants[overlap.Descendant] = true
		if _, err := fmt.Fprintf(stdout, "Overlap: %s already includes descendant %s\n", overlap.Ancestor, overlap.Descendant); err != nil {
			return err
		}
	}
	references := make([]string, 0, len(descendants))
	for reference := range descendants {
		references = append(references, reference)
	}
	sort.Strings(references)
	before, err := os.ReadFile(entryPath)
	if err != nil {
		return err
	}
	beforeManifest, err := nativeManifest(entryPath, nil)
	if err != nil {
		return err
	}
	contents, err := assembly.RemoveMany(entryPath, references, true)
	if err != nil {
		return err
	}
	afterManifest, err := nativeManifest(entryPath, contents)
	if err != nil {
		return err
	}
	if *apply {
		if _, err := assembly.RemoveMany(entryPath, references, false); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(stdout, "Applied optimized selection:"); err != nil {
			return err
		}
	} else if _, err := fmt.Fprintln(stdout, "Dry run: no files written. To apply these removals: mogent optimize --apply --file "+entryPath+"\nProposed project change:"); err != nil {
		return err
	}
	if _, err := fmt.Fprint(stdout, workspace.UnifiedDiffNamed(string(before), string(contents), "current "+entryPath, "proposed "+entryPath)); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(stdout, "---\nManaged-output changes:"); err != nil {
		return err
	}
	return writeManifestDiff(stdout, beforeManifest, afterManifest)
}
