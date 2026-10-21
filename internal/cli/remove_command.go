package cli

import (
	"flag"
	"fmt"
	"io"
	"sort"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/syntax"
	"github.com/Qu1ncyRy4n/Agents/workspace"
)

func runRemove(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("rm", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or authored build")
	dryRun := flags.Bool("dry-run", false, "preview the declaration removal without writing")
	diff := flags.Bool("diff", false, "preview resulting managed-output changes without writing")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--dry-run": false, "-dry-run": false, "--diff": false, "-diff": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("rm requires one direct module reference, for example qmr:intro/author-ownership")
	}
	entryPath, err := projectEntryFile(*path)
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
		if _, err := fmt.Fprintln(stdout, "Dry run: no files written\nExpected managed-output changes:"); err != nil {
			return err
		}
		return writeManifestDiff(stdout, before, after)
	}
	if *dryRun {
		_, err = fmt.Fprint(stdout, string(contents))
		return err
	}
	_, err = fmt.Fprintf(stdout, "Removed %s from %s\nNext: mogent rm %s --diff --file %s, then mogent build --check --file %s\n", flags.Arg(0), entryPath, flags.Arg(0), entryPath, entryPath)
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
		if _, err := fmt.Fprint(stdout, workspace.UnifiedDiff(string(old.Bytes), string(new.Bytes))); err != nil {
			return err
		}
	}
	if !changed {
		_, err := fmt.Fprintln(stdout, "No managed-output byte changes.")
		return err
	}
	return nil
}
