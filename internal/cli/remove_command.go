package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/syntax"
	"github.com/Qu1ncyRy4n/Agents/workspace"
)

func runRemove(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("module remove", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or authored build")
	dryRun := flags.Bool("dry-run", false, "preview the declaration removal without writing")
	diff := flags.Bool("diff", false, "show resulting project and managed-output changes")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--dry-run": false, "-dry-run": false, "--diff": false, "-diff": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("module remove requires one direct module reference, for example qmr:intro/author-ownership")
	}
	entryPath, err := projectEntryFile(*path)
	if err != nil {
		return err
	}
	beforeBytes, err := os.ReadFile(entryPath)
	if err != nil {
		return err
	}
	before, err := nativeManifest(entryPath, nil)
	if err != nil {
		return err
	}
	contents, err := assembly.Remove(entryPath, flags.Arg(0), *dryRun || *diff)
	if err != nil {
		return err
	}
	if *diff {
		after, err := nativeManifest(entryPath, contents)
		if err != nil {
			return err
		}
		headline := "Applied project change:"
		if *dryRun {
			headline = "Dry run: no files written\nProposed project change:"
		} else if _, err := assembly.Remove(entryPath, flags.Arg(0), false); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(stdout, headline); err != nil {
			return err
		}
		if _, err := fmt.Fprint(stdout, workspace.UnifiedDiffNamed(string(beforeBytes), string(contents), "current "+entryPath, "proposed "+entryPath)); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(stdout, "---\nManaged-output changes:"); err != nil {
			return err
		}
		return writeManifestDiff(stdout, before, after)
	}
	if *dryRun {
		_, err = fmt.Fprint(stdout, string(contents))
		return err
	}
	_, err = fmt.Fprintf(stdout, "Removed %s from %s\nNext: mogent show tree --file %s, then mogent build --check --file %s\n", flags.Arg(0), entryPath, entryPath, entryPath)
	return err
}

func nativeManifest(path string, contents []byte) (*content.Manifest, error) {
	file, err := assembly.LoadFile(path)
	if err != nil {
		return nil, err
	}
	if contents != nil {
		file, err = syntaxParse(path, contents)
		if err != nil {
			return nil, err
		}
	}
	result, _, err := assembly.ResolvePlan(file, path)
	if err != nil {
		return nil, err
	}
	if result.HasErrors() {
		return nil, fmt.Errorf("cannot compose managed outputs with unresolved selection errors")
	}
	return result.Compose()
}

// syntaxParse keeps the output-diff path small without exposing parser details
// from assembly just for a CLI preview.
func syntaxParse(path string, contents []byte) (*syntax.File, error) {
	return syntax.Parse(path, contents)
}

func writeManifestDiff(stdout io.Writer, before, after *content.Manifest) error {
	paths := map[string]bool{}
	for path := range before.Files {
		paths[path] = true
	}
	for path := range after.Files {
		paths[path] = true
	}
	ordered := make([]string, 0, len(paths))
	for path := range paths {
		ordered = append(ordered, path)
	}
	sort.Strings(ordered)
	changed := false
	for _, path := range ordered {
		old, new := before.Files[path], after.Files[path]
		if string(old.Bytes) == string(new.Bytes) && old.Mode == new.Mode {
			continue
		}
		changed = true
		if _, err := fmt.Fprintf(stdout, "Output: %s\n", path); err != nil {
			return err
		}
		if _, err := fmt.Fprint(stdout, workspace.UnifiedDiffNamed(string(old.Bytes), string(new.Bytes), "current "+path, "proposed "+path)); err != nil {
			return err
		}
	}
	if !changed {
		_, err := fmt.Fprintln(stdout, "No managed-output byte changes.")
		return err
	}
	return nil
}
