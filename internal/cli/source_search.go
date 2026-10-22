package cli

import (
	"flag"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
)

type stringList []string

func (values *stringList) String() string { return strings.Join(*values, ",") }

func (values *stringList) Set(value string) error {
	if strings.TrimSpace(value) == "" {
		return fmt.Errorf("tag must not be empty")
	}
	*values = append(*values, value)
	return nil
}

// runSourceSearch searches catalog metadata only. It intentionally does not
// inspect project selection state or install output, so it remains useful before
// a consumer has chosen any modules.
func runSourceSearch(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source search", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or build")
	var all, any, none stringList
	flags.Var(&all, "tag-all", "require tag; repeat for every required tag")
	flags.Var(&any, "tag-any", "match tag; repeat for any matching tag")
	flags.Var(&none, "tag-none", "exclude tag; repeat for any matching tag")
	addCommands := flags.Bool("add-commands", false, "print explicit module-add commands for matches")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--tag-all": true, "-tag-all": true, "--tag-any": true, "-tag-any": true, "--tag-none": true, "-tag-none": true, "--add-commands": false, "-add-commands": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() > 1 {
		return fmt.Errorf("source search accepts at most one source or source:path scope")
	}
	if len(all) == 0 && len(any) == 0 && len(none) == 0 {
		return fmt.Errorf("source search requires --tag-all, --tag-any, or --tag-none")
	}
	entryPath, err := projectEntryFile(*path)
	if err != nil {
		return err
	}
	file, err := assembly.LoadFile(entryPath)
	if err != nil {
		return err
	}
	_, sources, err := assembly.ResolvePlan(file, entryPath)
	if err != nil {
		return err
	}
	alias, prefix, err := searchScope(flags.Arg(0), sources)
	if err != nil {
		return err
	}
	aliases := make([]string, 0, len(sources))
	for name := range sources {
		if alias == "" || alias == name {
			aliases = append(aliases, name)
		}
	}
	sort.Strings(aliases)
	type match struct {
		reference string
		node      *assembly.Node
	}
	var matches []match
	for _, name := range aliases {
		library := sources[name].Lib
		paths := make([]string, 0, len(library.ByPath))
		for nodePath := range library.ByPath {
			if prefix == "" || nodePath == prefix || strings.HasPrefix(nodePath, prefix+"/") {
				paths = append(paths, nodePath)
			}
		}
		sort.Strings(paths)
		for _, nodePath := range paths {
			node := library.ByPath[nodePath]
			if assembly.MatchesTags(node, all, any, none) {
				matches = append(matches, match{reference: name + ":" + nodePath, node: node})
			}
		}
	}
	if len(matches) == 0 {
		_, err := fmt.Fprintln(stdout, "No source nodes matched the tag query.")
		return err
	}
	rows := [][]string{{"REF", "TYPE", "TAGS"}}
	for _, found := range matches {
		rows = append(rows, []string{found.reference, string(found.node.Type), strings.Join(found.node.Tags, ", ")})
	}
	if err := writeColumns(stdout, rows, true); err != nil {
		return err
	}
	if !*addCommands {
		return nil
	}
	if _, err := fmt.Fprintln(stdout, "---\nExplicit leaf add commands:"); err != nil {
		return err
	}
	for _, found := range matches {
		if len(found.node.Children) != 0 {
			continue
		}
		if _, err := fmt.Fprintf(stdout, "  mogent module add %s\n", found.reference); err != nil {
			return err
		}
	}
	return nil
}

func searchScope(value string, sources map[string]assembly.Source) (string, string, error) {
	if value == "" {
		return "", "", nil
	}
	alias, prefix, found := strings.Cut(value, ":")
	if !found {
		alias, prefix = value, ""
	}
	source, found := sources[alias]
	if !found {
		return "", "", fmt.Errorf("unknown source alias %q", alias)
	}
	if prefix != "" && source.Lib.ByPath[prefix] == nil {
		return "", "", fmt.Errorf("unknown source node %q", value)
	}
	return alias, prefix, nil
}
