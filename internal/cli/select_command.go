package cli

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/workspace"
)

func runAdd(args []string, stdout, stderr io.Writer) error {
	return runNativeSelect(args, stdout, stderr)
}

func runNativeSelect(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or authored build")
	dryRun := flags.Bool("dry-run", false, "preview the include without writing")
	diff := flags.Bool("diff", false, "show project and managed-output changes")
	all := flags.Bool("all", false, "select all direct children")
	var tagsAll, tagsAny, tagsNone stringList
	flags.Var(&tagsAll, "tag-all", "require tag; repeat for every required tag")
	flags.Var(&tagsAny, "tag-any", "match tag; repeat for any matching tag")
	flags.Var(&tagsNone, "tag-none", "exclude tag; repeat for any matching tag")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--dry-run": false, "-dry-run": false, "--diff": false, "-diff": false, "--all": false, "-all": false, "--tag-all": true, "-tag-all": true, "--tag-any": true, "-tag-any": true, "--tag-none": true, "-tag-none": true})
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
	var before []byte
	var beforeManifest *content.Manifest
	if *diff {
		before, err = os.ReadFile(entryPath)
		if err != nil {
			return err
		}
		beforeManifest, err = nativeManifest(entryPath, nil)
		if err != nil {
			return err
		}
	}
	selection := assembly.SelectOptions{All: *all, TagsAll: tagsAll, TagsAny: tagsAny, TagsNone: tagsNone, DryRun: *dryRun || *diff}
	contents, err := assembly.SelectWithOptions(entryPath, flags.Arg(0), selection)
	if err != nil {
		return err
	}
	if *diff {
		afterManifest, err := nativeManifest(entryPath, contents)
		if err != nil {
			return err
		}
		headline := "Applied project change:"
		if *dryRun {
			headline = "Dry run: no files written\nProposed project change:"
		} else if _, err := assembly.SelectWithOptions(entryPath, flags.Arg(0), assembly.SelectOptions{All: *all, TagsAll: tagsAll, TagsAny: tagsAny, TagsNone: tagsNone}); err != nil {
			return err
		}
		if _, err := fmt.Fprintln(stdout, headline); err != nil {
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
	if *dryRun {
		_, err = fmt.Fprint(stdout, string(contents))
	} else {
		_, err = fmt.Fprintf(stdout, "Added %s in %s\nNext: mogent show --tree --file %s\n", flags.Arg(0), entryPath, entryPath)
	}
	return err
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
