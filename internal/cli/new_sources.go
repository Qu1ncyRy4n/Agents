package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/registry"
)

func projectEntryFile(path string) (string, error) {
	if path != "" {
		return path, nil
	}
	plan, build := false, false
	for name, found := range map[string]*bool{"mogent.plan": &plan, "mogent.build": &build} {
		if _, err := os.Lstat(name); err == nil {
			*found = true
		} else if !os.IsNotExist(err) {
			return "", err
		}
	}
	if plan && build {
		file, err := assembly.LoadFile("mogent.build")
		if err != nil {
			return "", err
		}
		origin, found, err := buildOrigin(file)
		if err != nil || !found || origin != "generated" {
			return "", fmt.Errorf("both a plan and an authored build exist; choose --file explicitly")
		}
	}
	if plan {
		return "mogent.plan", nil
	}
	if build {
		return "mogent.build", nil
	}
	return "", fmt.Errorf("no mogent.plan or mogent.build; run mogent init first")
}

func writeNoSources(stdout io.Writer, path string) error {
	_, err := fmt.Fprintf(stdout, "No sources configured in %s.\nAdd a local library or Git repository:\n  mogent source add <name> <path-or-url> --file %s\nBrowse remembered libraries: mogent source list --known\n", path, path)
	return err
}

func runSourceRegister(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source register", flag.ContinueOnError)
	flags.SetOutput(stderr)
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 2 {
		return fmt.Errorf("source register requires a name and local path or Git URL")
	}
	path, err := registry.DefaultPath()
	if err != nil {
		return err
	}
	entry, err := registry.Register(path, flags.Arg(0), flags.Arg(1))
	if err != nil {
		return err
	}
	_, err = fmt.Fprintf(stdout, "Remembered %s: %s\nAvailable during mogent init, or: mogent source add %s\n", entry.Alias, entry.Location, entry.Alias)
	return err
}

func knownLibraries() ([]registry.Entry, error) {
	path, err := registry.DefaultPath()
	if err != nil {
		return nil, err
	}
	return registry.Load(path)
}

func runNewSourceAdd(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source add", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or authored build")
	dryRun := flags.Bool("dry-run", false, "preview without writing")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--dry-run": false, "-dry-run": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() < 1 || flags.NArg() > 2 {
		return fmt.Errorf("source add requires <name> [path-or-url]; omit location to use a registered library")
	}
	entryPath, err := projectEntryFile(*path)
	if err != nil {
		return err
	}
	alias, location := flags.Arg(0), flags.Arg(1)
	file, err := assembly.LoadFile(entryPath)
	if err != nil {
		return err
	}
	declarations, err := assembly.SourceDeclarations(file)
	if err != nil {
		return err
	}
	for _, source := range declarations {
		if source.Alias == alias {
			return fmt.Errorf("source %q is already connected in %s; browse it with mogent source show %s", alias, entryPath, alias)
		}
	}
	if flags.NArg() == 1 {
		entries, err := knownLibraries()
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if entry.Alias == alias {
				location = entry.Location
				break
			}
		}
		if location == "" {
			return fmt.Errorf("no registered library %q; provide a local path or Git URL", alias)
		}
	}
	content, err := assembly.AddSource(entryPath, alias, location, *dryRun)
	if err != nil {
		return err
	}
	if *dryRun {
		_, err = fmt.Fprint(stdout, string(content))
	} else {
		_, err = fmt.Fprintf(stdout, "Added source %s to %s\nNext: mogent source show --file %s\n", alias, entryPath, entryPath)
	}
	return err
}

func runNewSourceList(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source list", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or build")
	known := flags.Bool("known", false, "list user-registered libraries instead of project sources")
	tree := flags.Bool("tree", false, "browse source content as a tree")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--known": false, "-known": false, "--tree": false, "-tree": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() > 1 || *known && (*path != "" || *tree || flags.NArg() != 0) {
		return fmt.Errorf("source list accepts one project alias, or --known by itself")
	}
	rows := [][]string{{"NAME", "KIND", "LOCATION"}}
	if *known {
		entries, err := knownLibraries()
		if err != nil {
			return err
		}
		if len(entries) == 0 {
			_, err := fmt.Fprintln(stdout, "No remembered libraries.\nRegister a local library or Git repository:\n  mogent source register <name> <path-or-url>")
			return err
		}
		rows = [][]string{{"NAME", "LIBRARY", "LOCATION"}}
		for _, entry := range entries {
			rows = append(rows, []string{entry.Alias, entry.Name, entry.Location})
		}
		return writeColumns(stdout, rows, true)
	}
	entryPath, err := projectEntryFile(*path)
	if err != nil {
		return err
	}
	if *tree {
		show := []string{"--tree", "--file", entryPath}
		if flags.NArg() == 1 {
			show = append(show, flags.Arg(0)+":")
		}
		return runNativeSourceShow(show, stdout, stderr)
	}
	file, err := assembly.LoadFile(entryPath)
	if err != nil {
		return err
	}
	sources, err := assembly.SourceDeclarations(file)
	if err != nil {
		return err
	}
	if len(sources) == 0 {
		return writeNoSources(stdout, entryPath)
	}
	for _, source := range sources {
		if flags.NArg() == 1 && source.Alias != strings.TrimSuffix(flags.Arg(0), ":") {
			continue
		}
		kind, location := "local", source.Local
		if source.Git != "" {
			kind, location = "git", source.Git
		}
		rows = append(rows, []string{source.Alias, kind, location})
	}
	if len(rows) == 1 {
		return fmt.Errorf("unknown source alias %q", flags.Arg(0))
	}
	return writeColumns(stdout, rows, true)
}
