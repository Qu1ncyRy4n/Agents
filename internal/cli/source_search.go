package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/syntax"
	"github.com/charmbracelet/x/term"
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
	addConfig := flags.Bool("add-config", false, "print a dynamic tag-query include block")
	searchAll := flags.Bool("all", false, "search path, title, tags, and content")
	searchPath := flags.Bool("path", false, "search module paths")
	searchTitle := flags.Bool("title", false, "search module titles")
	searchTags := flags.Bool("tags", false, "search tags as text")
	searchContent := flags.Bool("content", false, "search rendered text content")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--tag-all": true, "-tag-all": true, "--tag-any": true, "-tag-any": true, "--tag-none": true, "-tag-none": true, "--add-commands": false, "-add-commands": false, "--add-config": false, "-add-config": false, "--all": false, "-all": false, "--path": false, "-path": false, "--title": false, "-title": false, "--tags": false, "-tags": false, "--content": false, "-content": false})
	if err := parseFlags(flags, args); err != nil {
		return err
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
	tagMode := len(all) != 0 || len(any) != 0 || len(none) != 0
	fieldMode := *searchAll || *searchPath || *searchTitle || *searchTags || *searchContent
	if tagMode && fieldMode {
		return fmt.Errorf("tag predicates cannot be combined with text-search fields")
	}
	if tagMode && flags.NArg() > 1 {
		return fmt.Errorf("tag search accepts at most one source or source:path scope")
	}
	if !tagMode && (*addCommands || *addConfig) {
		return fmt.Errorf("--add-commands and --add-config require a tag predicate")
	}
	var scope, query string
	if tagMode {
		scope = flags.Arg(0)
	} else {
		if flags.NArg() == 0 || flags.NArg() > 2 {
			return fmt.Errorf("text search requires a term, optionally after a source or source:path scope")
		}
		if flags.NArg() == 2 {
			scope, query = flags.Arg(0), flags.Arg(1)
		} else {
			query = flags.Arg(0)
		}
	}
	alias, prefix, err := searchScope(scope, sources)
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
	if !tagMode {
		fields := textSearchFields{all: *searchAll || !fieldMode, path: *searchPath, title: *searchTitle, tags: *searchTags, content: *searchContent}
		matches, err := sourceTextMatches(file, entryPath, sources, aliases, prefix, query, fields)
		if err != nil {
			return err
		}
		if len(matches) == 0 {
			_, err := fmt.Fprintln(stdout, "No source nodes matched the text search.")
			return err
		}
		rows := [][]string{{"REF", "TYPE", "TAGS", "MATCH"}}
		styled := supportsTerminalStyle(stdout)
		for _, found := range matches {
			rows = append(rows, []string{highlightText(found.reference, query, styled), string(found.node.Type), highlightText(strings.Join(found.node.Tags, ", "), query, styled), strings.Join(found.fields, ", ")})
		}
		return writeColumns(stdout, rows, true)
	}
	matches := sourceTagMatches(sources, aliases, prefix, all, any, none)
	if len(matches) == 0 {
		_, err := fmt.Fprintln(stdout, "No source nodes matched the tag query.")
		return err
	}
	rows := [][]string{{"REF", "TYPE", "TAGS"}}
	styled := supportsTerminalStyle(stdout)
	for _, found := range matches {
		rows = append(rows, []string{found.reference, string(found.node.Type), formatSearchTags(found.node.Tags, all, any, styled)})
	}
	if err := writeColumns(stdout, rows, true); err != nil {
		return err
	}
	if *addConfig {
		if prefix == "" {
			return fmt.Errorf("--add-config requires a source:path scope so the include has a target")
		}
		if _, err := fmt.Fprintf(stdout, "---\nDynamic tag-query config:\n%s", tagQueryConfig(alias+":"+prefix, all, any, none)); err != nil {
			return err
		}
	}
	if *addCommands {
		return writeExplicitAddCommands(stdout, matches)
	}
	return nil
}

func supportsTerminalStyle(output io.Writer) bool {
	file, terminal := output.(*os.File)
	_, noColor := os.LookupEnv("NO_COLOR")
	return terminal && term.IsTerminal(file.Fd()) && !noColor && os.Getenv("TERM") != "dumb"
}

func formatSearchTags(tags, all, any []string, styled bool) string {
	result := append([]string(nil), tags...)
	if !styled {
		return strings.Join(result, ", ")
	}
	queries := append(append([]string{}, all...), any...)
	for index, tag := range result {
		for _, query := range queries {
			if highlighted, found := highlightTagQuery(tag, query); found {
				result[index] = highlighted
				break
			}
		}
	}
	return strings.Join(result, ", ")
}

func highlightTagQuery(tag, query string) (string, bool) {
	if tag == query || strings.HasPrefix(tag, query+"/") {
		return "\x1b[1m" + query + "\x1b[22m" + strings.TrimPrefix(tag, query), true
	}
	if strings.Contains(query, "/") {
		return tag, false
	}
	parts := strings.Split(tag, "/")
	for index, part := range parts {
		if part == query {
			parts[index] = "\x1b[1m" + part + "\x1b[22m"
			return strings.Join(parts, "/"), true
		}
	}
	return tag, false
}

type textSearchFields struct {
	all, path, title, tags, content bool
}

type textSearchMatch struct {
	reference string
	node      *assembly.Node
	fields    []string
}

func sourceTextMatches(file *syntax.File, entryPath string, sources map[string]assembly.Source, aliases []string, prefix, query string, fields textSearchFields) ([]textSearchMatch, error) {
	needle := strings.ToLower(query)
	var matches []textSearchMatch
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
			ref := name + ":" + nodePath
			var found []string
			if (fields.all || fields.path) && containsSearchText(nodePath, needle) {
				found = append(found, "path")
			}
			if (fields.all || fields.title) && containsSearchText(node.Title, needle) {
				found = append(found, "title")
			}
			if (fields.all || fields.tags) && containsSearchText(strings.Join(node.Tags, " "), needle) {
				found = append(found, "tags")
			}
			if fields.all || fields.content {
				preview, err := assembly.PreviewSource(file, entryPath, ref, 0, false)
				if err != nil {
					// A content search should still report searchable metadata when
					// another node cannot render with the project's bindings.
				} else if !preview.Opaque && containsSearchText(string(preview.Content), needle) {
					found = append(found, "content")
				}
			}
			if len(found) != 0 {
				matches = append(matches, textSearchMatch{reference: ref, node: node, fields: found})
			}
		}
	}
	return matches, nil
}

func containsSearchText(value, needle string) bool {
	return needle != "" && strings.Contains(strings.ToLower(value), needle)
}

func highlightText(value, query string, styled bool) string {
	if !styled || query == "" {
		return value
	}
	lower, needle := strings.ToLower(value), strings.ToLower(query)
	index := strings.Index(lower, needle)
	if index < 0 {
		return value
	}
	end := index + len(query)
	return value[:index] + "\x1b[1m" + value[index:end] + "\x1b[22m" + value[end:]
}

type sourceMatch struct {
	reference string
	node      *assembly.Node
}

func sourceTagMatches(sources map[string]assembly.Source, aliases []string, prefix string, all, any, none []string) []sourceMatch {
	var matches []sourceMatch
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
				matches = append(matches, sourceMatch{reference: name + ":" + nodePath, node: node})
			}
		}
	}
	return matches
}

func writeExplicitAddCommands(stdout io.Writer, matches []sourceMatch) error {
	if _, err := fmt.Fprintln(stdout, "---\nExplicit leaf include commands:"); err != nil {
		return err
	}
	for _, found := range matches {
		if len(found.node.Children) != 0 {
			continue
		}
		if _, err := fmt.Fprintf(stdout, "  mogent module include %s\n", found.reference); err != nil {
			return err
		}
	}
	return nil
}

func tagQueryConfig(reference string, all, any, none []string) string {
	var output strings.Builder
	fmt.Fprintf(&output, "include %s {\n", strconv.Quote(reference))
	writeTagConfigField(&output, "tags_all", all)
	writeTagConfigField(&output, "tags_any", any)
	writeTagConfigField(&output, "tags_none", none)
	output.WriteString("}\n")
	return output.String()
}

func writeTagConfigField(output *strings.Builder, name string, tags []string) {
	if len(tags) == 0 {
		return
	}
	quoted := make([]string, 0, len(tags))
	for _, tag := range tags {
		quoted = append(quoted, strconv.Quote(tag))
	}
	fmt.Fprintf(output, "  %s = [%s]\n", name, strings.Join(quoted, ", "))
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
