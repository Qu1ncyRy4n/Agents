// Package cli exposes the deliberately small milestone-one command surface.
package cli

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/internal/library"
	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
	"github.com/Qu1ncyRy4n/Agents/internal/navigator"
	"github.com/Qu1ncyRy4n/Agents/internal/render"
	"github.com/Qu1ncyRy4n/Agents/internal/state"
	"github.com/Qu1ncyRy4n/Agents/internal/workspace"
)

func Run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		if _, err := fmt.Fprintln(stdout, "Usage: mogent build [--manifest agents.yaml] [--force]\n       mogent status [--manifest agents.yaml]\n       mogent coverage [--manifest agents.yaml]\n       mogent source list [source] [--manifest agents.yaml]\n       mogent source show <ref> [--manifest agents.yaml]\n       mogent add <ref> [--under heading/path | --append]\n       mogent complete <kind> [--manifest agents.yaml]\n       mogent completion <bash|zsh>\n       mogent tui [--manifest agents.yaml]"); err != nil {
			return fmt.Errorf("write usage: %w", err)
		}
		return nil
	}
	switch args[0] {
	case "build":
		return runBuild(args[1:], stdout, stderr)
	case "status":
		return runStatus(args[1:], stdout, stderr)
	case "coverage":
		return runCoverage(args[1:], stdout, stderr)
	case "source":
		return runSource(args[1:], stdout, stderr)
	case "add":
		return runAdd(args[1:], stdout, stderr)
	case "complete":
		return runComplete(args[1:], stdout, stderr)
	case "completion":
		return runCompletion(args[1:], stdout, stderr)
	case "tui":
		return runTUI(args[1:], stdout, stderr)
	default:
		return unknownCommandError(args[0])
	}
}

func runBuild(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("build", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	force := flags.Bool("force", false, "replace an untracked or directly edited output")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("build accepts no positional arguments")
	}
	value, manifestPath, err := manifest.Load(*manifestFile)
	if err != nil {
		return err
	}
	result, err := render.Build(value, manifestPath)
	if err != nil {
		return err
	}
	outputPath := value.Output
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(filepath.Dir(manifestPath), outputPath)
	}
	statePath := filepath.Join(filepath.Dir(manifestPath), ".mogent", "state.json")
	if err := state.CheckOverwrite(outputPath, statePath, *force); err != nil {
		return err
	}
	if err := render.WriteAtomically(outputPath, result.Content); err != nil {
		return err
	}
	if err := state.Write(statePath, result.Content); err != nil {
		return err
	}
	for _, warning := range result.Warnings {
		if _, err := fmt.Fprintf(stderr, "mogent: warning: %s\n", warning); err != nil {
			return fmt.Errorf("write warning: %w", err)
		}
	}
	if _, err := fmt.Fprintf(stdout, "Wrote %s\n", outputPath); err != nil {
		return fmt.Errorf("write success message: %w", err)
	}
	return nil
}

func runStatus(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("status", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("status accepts no positional arguments")
	}
	session, err := workspace.New(*manifestFile)
	if err != nil {
		return err
	}
	status, err := session.Status()
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdout, "Manifest: %s\n", status.ManifestPath); err != nil {
		return fmt.Errorf("write status: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "Output: %s\n", status.OutputPath); err != nil {
		return fmt.Errorf("write status: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "Output status: %s\n", status.Output); err != nil {
		return fmt.Errorf("write status: %w", err)
	}
	if hint := statusHint(status.Output); hint != "" {
		if _, err := fmt.Fprintf(stdout, "Hint: %s\n", hint); err != nil {
			return fmt.Errorf("write status: %w", err)
		}
	}
	if _, err := fmt.Fprintf(stdout, "Sources: %d\n", len(status.Sources)); err != nil {
		return fmt.Errorf("write status: %w", err)
	}
	for _, source := range status.Sources {
		if _, err := fmt.Fprintf(stdout, "- %s: %s (%d nodes)\n", source.Alias, source.Path, source.Nodes); err != nil {
			return fmt.Errorf("write status: %w", err)
		}
	}
	return nil
}

func statusHint(status workspace.OutputStatus) string {
	switch status {
	case workspace.StatusMissing:
		return "run `mogent build` to create the output"
	case workspace.StatusStale:
		return "run `mogent build` to refresh the output"
	case workspace.StatusDirectEdits:
		return "AGENTS.md has direct edits; review them before running `mogent build --force`"
	case workspace.StatusUntracked:
		return "AGENTS.md was not generated by mogent; use `mogent build --force` only after review"
	case workspace.StatusUpToDate:
		return "clean"
	default:
		return ""
	}
}

func runCoverage(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("coverage", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	sourceAlias := flags.String("source", "", "show only one source alias")
	tag := flags.String("tag", "", "show only source nodes with this metadata tag")
	contentOnly := flags.Bool("content-only", false, "hide source nodes without body content")
	leavesOnly := flags.Bool("leaves-only", false, "show only terminal source nodes")
	depth := flags.Int("depth", -1, "maximum source-tree depth to show; root headings are depth 0")
	unusedOnly := flags.Bool("unused-only", false, "show only unused source references")
	tree := flags.Bool("tree", false, "show unused source references as an ASCII tree")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("coverage accepts no positional arguments")
	}
	session, err := workspace.New(*manifestFile)
	if err != nil {
		return err
	}
	coverage := session.CoverageWithOptions(workspace.CoverageOptions{
		SourceAlias: *sourceAlias,
		Tag:         *tag,
		ContentOnly: *contentOnly,
		LeavesOnly:  *leavesOnly,
		LimitDepth:  *depth >= 0,
		MaxDepth:    *depth,
	})
	if *sourceAlias != "" && len(coverage.Sources) == 0 {
		if _, err := session.SourceReferences(*sourceAlias); err != nil {
			return err
		}
		return fmt.Errorf("source %q has no coverage rows after filtering", *sourceAlias)
	}
	if *tree {
		return writeTreeCoverage(stdout, coverage)
	}
	if *unusedOnly {
		return writeUnusedCoverage(stdout, coverage)
	}
	for _, source := range coverage.Sources {
		if _, err := fmt.Fprintf(stdout, "Source %s: %s\n", source.Alias, source.Path); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
		if _, err := fmt.Fprintf(stdout, "Included: %d/%d\n", source.Included, source.Total); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
		if len(source.Unused) == 0 {
			if _, err := fmt.Fprintln(stdout, "Unused: none"); err != nil {
				return fmt.Errorf("write coverage: %w", err)
			}
		} else {
			if _, err := fmt.Fprintln(stdout, "Unused:"); err != nil {
				return fmt.Errorf("write coverage: %w", err)
			}
			for _, node := range source.Unused {
				if _, err := fmt.Fprintf(stdout, "- %s  %s:%s\n", node.Heading, source.Alias, node.Path); err != nil {
					return fmt.Errorf("write coverage: %w", err)
				}
			}
		}
		if _, err := fmt.Fprintln(stdout); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
	}
	return nil
}

func writeTreeCoverage(stdout io.Writer, coverage workspace.Coverage) error {
	for _, source := range coverage.Sources {
		if _, err := fmt.Fprintf(stdout, "%s  %s  included %d/%d, unused %d\n", source.Alias, source.Path, source.Included, source.Total, len(source.Unused)); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
		if len(source.Nodes) == 0 {
			if _, err := fmt.Fprintln(stdout, "  (none)"); err != nil {
				return fmt.Errorf("write coverage: %w", err)
			}
		}
		for index, node := range source.Nodes {
			depth := visibleDepth(source.Nodes, index)
			connector := "`-- "
			if hasNextAtDepth(source.Nodes, index, depth) {
				connector = "|-- "
			}
			indent := strings.Repeat("|   ", depth)
			if _, err := fmt.Fprintf(stdout, "%s%s[%s] %s  %s:%s\n", indent, connector, node.State, node.Heading, source.Alias, node.Path); err != nil {
				return fmt.Errorf("write coverage: %w", err)
			}
		}
		if _, err := fmt.Fprintln(stdout); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
	}
	return nil
}

func visibleDepth(nodes []workspace.CoverageNode, index int) int {
	depth := 0
	path := nodes[index].Path
	for previous := 0; previous < index; previous++ {
		prefix := nodes[previous].Path + "/"
		if strings.HasPrefix(path, prefix) {
			depth++
		}
	}
	return depth
}

func hasNextAtDepth(nodes []workspace.CoverageNode, index int, depth int) bool {
	for next := index + 1; next < len(nodes); next++ {
		if visibleDepth(nodes, next) == depth {
			return true
		}
		if visibleDepth(nodes, next) < depth {
			return false
		}
	}
	return false
}

func writeUnusedCoverage(stdout io.Writer, coverage workspace.Coverage) error {
	for _, source := range coverage.Sources {
		if len(source.Unused) == 0 {
			continue
		}
		if _, err := fmt.Fprintf(stdout, "%s  %s  unused %d/%d\n", source.Alias, source.Path, len(source.Unused), source.Total); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
		for _, node := range source.Unused {
			if _, err := fmt.Fprintf(stdout, "  %s  %s\n", source.Alias+":"+node.Path, node.Heading); err != nil {
				return fmt.Errorf("write coverage: %w", err)
			}
		}
		if _, err := fmt.Fprintln(stdout); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
	}
	return nil
}

func runSource(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("source requires a subcommand; use source list or source show")
	}
	switch args[0] {
	case "list":
		return runSourceList(args[1:], stdout, stderr)
	case "show":
		return runSourceShow(args[1:], stdout, stderr)
	default:
		if suggestion := closestString(args[0], []string{"list", "show"}, 2); suggestion != "" {
			return fmt.Errorf("unknown source subcommand %q; did you mean %q?", args[0], suggestion)
		}
		return fmt.Errorf("unknown source subcommand %q; use source list or source show", args[0])
	}
}

func runSourceList(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source list", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	sourceAlias := flags.String("source", "", "show only one source alias")
	tag := flags.String("tag", "", "show only source nodes with this exact metadata tag")
	tagSearch := flags.String("tag-search", "", "show only source nodes whose tags contain this text")
	sortMode := flags.String("sort", "path", "sort mode: path or priority")
	showMetadata := flags.Bool("metadata", false, "show source metadata on separate lines")
	showFile := flags.Bool("file", false, "show source file path")
	showLine := flags.Bool("line", false, "show source heading line")
	search := flags.String("search", "", "search source refs, headings, TLDRs, tags, and direct body text")
	tldrOnly := flags.Bool("tldr", false, "show compact TLDR rows and skip Metadata: none lines")
	args = reorderArgs(args, map[string]bool{
		"-manifest":    true,
		"--manifest":   true,
		"-source":      true,
		"--source":     true,
		"-tag":         true,
		"--tag":        true,
		"-tag-search":  true,
		"--tag-search": true,
		"-sort":        true,
		"--sort":       true,
		"-search":      true,
		"--search":     true,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() > 1 {
		return fmt.Errorf("source list accepts at most one source alias positional argument")
	}
	if flags.NArg() == 1 {
		if *sourceAlias != "" {
			return fmt.Errorf("source list source alias was provided twice")
		}
		*sourceAlias = strings.TrimSuffix(flags.Arg(0), ":")
	}
	session, err := workspace.New(*manifestFile)
	if err != nil {
		return err
	}
	nodes, err := session.SourceNodes(workspace.SourceListOptions{
		SourceAlias: *sourceAlias,
		Tag:         *tag,
		TagSearch:   *tagSearch,
		Search:      *search,
		Sort:        *sortMode,
	})
	if err != nil {
		return err
	}
	if len(nodes) == 0 {
		if _, err := fmt.Fprintln(stdout, "No source nodes matched."); err != nil {
			return fmt.Errorf("write source list: %w", err)
		}
		if err := writeSourceListEmptyHint(stdout, session, *sourceAlias, *tag, *tagSearch, *search); err != nil {
			return err
		}
		return nil
	}
	for _, node := range nodes {
		if err := writeSourceListNode(stdout, node, *showMetadata && !*tldrOnly, *showFile, *showLine); err != nil {
			return err
		}
	}
	return nil
}

func writeSourceListEmptyHint(stdout io.Writer, session *workspace.Session, sourceAlias string, tag string, tagSearch string, search string) error {
	if tagSearch != "" && search == "" {
		candidates, err := session.SourceNodes(workspace.SourceListOptions{
			SourceAlias: sourceAlias,
			Search:      tagSearch,
			Sort:        "path",
		})
		if err == nil && len(candidates) > 0 {
			if _, err := fmt.Fprintf(stdout, "Hint: --tag-search only searches metadata tags. %d source paths/headings/content snippets match %q; try `mogent source list --search %s`.\n", len(candidates), tagSearch, shellToken(tagSearch)); err != nil {
				return fmt.Errorf("write source list: %w", err)
			}
			return nil
		}
	}
	if tag != "" || tagSearch != "" {
		if _, err := fmt.Fprintln(stdout, "Hint: most source files may not have tags yet. Try `mogent source list --search <text>` to search refs, headings, TLDRs, tags, and direct body text."); err != nil {
			return fmt.Errorf("write source list: %w", err)
		}
		return nil
	}
	if _, err := fmt.Fprintln(stdout, "Hint: try `mogent source list --search <text>` or `mogent source list <source-alias>`."); err != nil {
		return fmt.Errorf("write source list: %w", err)
	}
	return nil
}

func writeSourceListNode(stdout io.Writer, node workspace.SourceNode, showMetadata bool, showFile bool, showLine bool) error {
	if _, err := fmt.Fprintf(stdout, "%s  %s", node.Reference, node.Heading); err != nil {
		return fmt.Errorf("write source list: %w", err)
	}
	if len(node.Metadata.Tags) > 0 {
		if _, err := fmt.Fprintf(stdout, "  [%s]", strings.Join(node.Metadata.Tags, ", ")); err != nil {
			return fmt.Errorf("write source list: %w", err)
		}
	}
	if node.Metadata.Priority != nil {
		if _, err := fmt.Fprintf(stdout, "  p=%.2f", *node.Metadata.Priority); err != nil {
			return fmt.Errorf("write source list: %w", err)
		}
	}
	if node.Metadata.TLDR != "" {
		if _, err := fmt.Fprintf(stdout, "  - %s", node.Metadata.TLDR); err != nil {
			return fmt.Errorf("write source list: %w", err)
		}
	}
	if showFile {
		if _, err := fmt.Fprintf(stdout, "  %s", node.File); err != nil {
			return fmt.Errorf("write source list: %w", err)
		}
	}
	if showLine {
		if _, err := fmt.Fprintf(stdout, ":%d", node.Line); err != nil {
			return fmt.Errorf("write source list: %w", err)
		}
	}
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write source list: %w", err)
	}
	if showMetadata {
		if err := writeIndentedSourceMetadata(stdout, node.Metadata); err != nil {
			return err
		}
	}
	return nil
}

func writeIndentedSourceMetadata(stdout io.Writer, metadata library.Metadata) error {
	if metadata.TLDR == "" && len(metadata.Tags) == 0 && metadata.Priority == nil && metadata.Scope == "" && len(metadata.Requires) == 0 && len(metadata.ConflictsWith) == 0 {
		if _, err := fmt.Fprintln(stdout, "  Metadata: none"); err != nil {
			return fmt.Errorf("write source list: %w", err)
		}
		return nil
	}
	if _, err := fmt.Fprintln(stdout, "  Metadata:"); err != nil {
		return fmt.Errorf("write source list: %w", err)
	}
	if err := writeMetadataFields(stdout, metadata, "    "); err != nil {
		return fmt.Errorf("write source list: %w", err)
	}
	return nil
}

func runSourceShow(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source show", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	showFile := flags.Bool("file", true, "show source file path")
	showLine := flags.Bool("line", true, "show source heading line")
	showMetadata := flags.Bool("metadata", false, "show source metadata")
	contentMode := flags.String("content", "full", "content mode: none, snippet, or full")
	snippetLines := flags.Int("lines", 12, "number of content lines when --content=snippet")
	alignSource := flags.Bool("align-source", false, "preview the source as a rendered manifest section")
	under := flags.String("under", "", "manifest heading path to align beneath when using --align-source")
	args = reorderArgs(args, map[string]bool{
		"-manifest":  true,
		"--manifest": true,
		"-content":   true,
		"--content":  true,
		"-lines":     true,
		"--lines":    true,
		"-under":     true,
		"--under":    true,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("source show requires exactly one source reference")
	}
	session, err := workspace.New(*manifestFile)
	if err != nil {
		return err
	}
	node, err := session.SourceNode(flags.Arg(0))
	if err != nil {
		return decorateSourceReferenceError(flags.Arg(0), err)
	}
	if *alignSource {
		return writeAlignedSource(stdout, node, *under, *contentMode, *snippetLines)
	}
	if _, err := fmt.Fprintf(stdout, "Source: %s\n", node.Alias); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "Reference: %s\n", node.Reference); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	if *showFile {
		if _, err := fmt.Fprintf(stdout, "File: %s\n", node.File); err != nil {
			return fmt.Errorf("write source: %w", err)
		}
	}
	if *showLine {
		if _, err := fmt.Fprintf(stdout, "Line: %d\n", node.Line); err != nil {
			return fmt.Errorf("write source: %w", err)
		}
	}
	if _, err := fmt.Fprintf(stdout, "Heading: %s\n", node.Heading); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	if *showMetadata {
		if err := writeSourceMetadata(stdout, node.Metadata); err != nil {
			return err
		}
	}
	content, err := sourceContentForMode(node.Content, *contentMode, *snippetLines)
	if err != nil {
		return err
	}
	if content == "" {
		return nil
	}
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, content); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	return nil
}

func writeAlignedSource(stdout io.Writer, node *workspace.SourceNode, under string, contentMode string, snippetLines int) error {
	level := 1
	if strings.TrimSpace(under) != "" {
		level = len(strings.Split(strings.Trim(strings.TrimSpace(under), "/"), "/")) + 1
	}
	if level > 6 {
		return fmt.Errorf("--under %q would render beyond Markdown heading level 6", under)
	}
	content, err := sourceContentForMode(node.Content, contentMode, snippetLines)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdout, "Aligned preview: %s\n", node.Reference); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	if strings.TrimSpace(under) != "" {
		if _, err := fmt.Fprintf(stdout, "Under: %s\n\n", under); err != nil {
			return fmt.Errorf("write source: %w", err)
		}
	} else if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "%s %s\n\n", strings.Repeat("#", level), node.Heading); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	if strings.TrimSpace(content) != "" {
		if _, err := fmt.Fprintln(stdout, content); err != nil {
			return fmt.Errorf("write source: %w", err)
		}
	}
	return nil
}

func writeSourceMetadata(stdout io.Writer, metadata library.Metadata) error {
	if metadata.TLDR == "" && len(metadata.Tags) == 0 && metadata.Priority == nil && metadata.Scope == "" && len(metadata.Requires) == 0 && len(metadata.ConflictsWith) == 0 {
		if _, err := fmt.Fprintln(stdout, "Metadata: none"); err != nil {
			return fmt.Errorf("write source: %w", err)
		}
		return nil
	}
	if _, err := fmt.Fprintln(stdout, "Metadata:"); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	if err := writeMetadataFields(stdout, metadata, "  "); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	return nil
}

func writeMetadataFields(stdout io.Writer, metadata library.Metadata, indent string) error {
	if metadata.TLDR != "" {
		if _, err := fmt.Fprintf(stdout, "%sTLDR: %s\n", indent, metadata.TLDR); err != nil {
			return err
		}
	}
	if len(metadata.Tags) > 0 {
		if _, err := fmt.Fprintf(stdout, "%sTags: %s\n", indent, strings.Join(metadata.Tags, ", ")); err != nil {
			return err
		}
	}
	if metadata.Priority != nil {
		if _, err := fmt.Fprintf(stdout, "%sPriority: %.2f\n", indent, *metadata.Priority); err != nil {
			return err
		}
	}
	if metadata.Scope != "" {
		if _, err := fmt.Fprintf(stdout, "%sScope: %s\n", indent, metadata.Scope); err != nil {
			return err
		}
	}
	if len(metadata.Requires) > 0 {
		if _, err := fmt.Fprintf(stdout, "%sRequires: %s\n", indent, strings.Join(metadata.Requires, ", ")); err != nil {
			return err
		}
	}
	if len(metadata.ConflictsWith) > 0 {
		if _, err := fmt.Fprintf(stdout, "%sConflicts with: %s\n", indent, strings.Join(metadata.ConflictsWith, ", ")); err != nil {
			return err
		}
	}
	return nil
}

func reorderArgs(args []string, valueFlags map[string]bool) []string {
	var flagArgs []string
	var positional []string
	for index := 0; index < len(args); index++ {
		arg := args[index]
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			positional = append(positional, arg)
			continue
		}
		flagArgs = append(flagArgs, arg)
		name := arg
		if before, _, found := strings.Cut(arg, "="); found {
			name = before
		}
		if valueFlags[name] && !strings.Contains(arg, "=") && index+1 < len(args) {
			index++
			flagArgs = append(flagArgs, args[index])
		}
	}
	return append(flagArgs, positional...)
}

func parseFlags(flags *flag.FlagSet, args []string) error {
	if err := suggestUnknownFlags(flags, args); err != nil {
		return err
	}
	return flags.Parse(args)
}

func suggestUnknownFlags(flags *flag.FlagSet, args []string) error {
	var known []string
	flags.VisitAll(func(flag *flag.Flag) {
		known = append(known, flag.Name)
	})
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") || arg == "-" || arg == "--" {
			continue
		}
		name := strings.TrimLeft(arg, "-")
		if before, _, found := strings.Cut(name, "="); found {
			name = before
		}
		if name == "h" || name == "help" || flags.Lookup(name) != nil {
			continue
		}
		if suggestion := closestString(name, known, 4); suggestion != "" {
			return fmt.Errorf("unknown flag %q; did you mean --%s?", arg, suggestion)
		}
		return fmt.Errorf("unknown flag %q", arg)
	}
	return nil
}

func sourceContentForMode(content string, mode string, lines int) (string, error) {
	switch strings.ToLower(mode) {
	case "none":
		return "", nil
	case "snippet":
		return workspace.Snippet(content, lines), nil
	case "full":
		return strings.TrimSpace(content), nil
	default:
		return "", fmt.Errorf("unknown content mode %q; use none, snippet, or full", mode)
	}
}

func runAdd(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("add", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	under := flags.String("under", "", "manifest heading path to append under")
	appendRoot := flags.Bool("append", false, "append to the end of the document")
	heading := flags.String("heading", "", "rendered heading to use in the manifest")
	dryRun := flags.Bool("dry-run", false, "preview the change without writing")
	previewMode := flags.String("preview", "summary", "dry-run preview mode: summary, patch, tree, or full")
	rebuild := flags.Bool("rebuild", false, "also rebuild the generated AGENTS.md")
	args = reorderArgs(args, map[string]bool{
		"-manifest":  true,
		"--manifest": true,
		"-under":     true,
		"--under":    true,
		"-heading":   true,
		"--heading":  true,
		"-preview":   true,
		"--preview":  true,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("add requires exactly one source reference")
	}
	session, err := workspace.New(*manifestFile)
	if err != nil {
		return err
	}
	result, err := session.AddSource(workspace.AddOptions{
		Reference: flags.Arg(0),
		Under:     *under,
		Append:    *appendRoot,
		Heading:   *heading,
		Rebuild:   *rebuild,
	}, *dryRun)
	if err != nil {
		return decorateSourceReferenceError(flags.Arg(0), err)
	}
	if *dryRun {
		if _, err := fmt.Fprintln(stdout, "Dry run: no files written"); err != nil {
			return fmt.Errorf("write add result: %w", err)
		}
		if err := writeAddPreview(stdout, result, *previewMode); err != nil {
			return err
		}
		return nil
	} else if result.Rebuilt {
		if _, err := fmt.Fprintln(stdout, "Wrote manifest and rebuilt output"); err != nil {
			return fmt.Errorf("write add result: %w", err)
		}
	} else if result.WroteManifest {
		if _, err := fmt.Fprintln(stdout, "Wrote manifest"); err != nil {
			return fmt.Errorf("write add result: %w", err)
		}
	}
	if _, err := fmt.Fprintf(stdout, "Added: %s\n", result.Heading); err != nil {
		return fmt.Errorf("write add result: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "From: %s\n", result.Reference); err != nil {
		return fmt.Errorf("write add result: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "Under: %s\n", parentLabelForOutput(result.ParentPath)); err != nil {
		return fmt.Errorf("write add result: %w", err)
	}
	return nil
}

func writeAddPreview(stdout io.Writer, result *workspace.AddResult, mode string) error {
	switch strings.ToLower(mode) {
	case "summary":
		return writeAddSummary(stdout, result)
	case "patch":
		return writeAddPatch(stdout, result)
	case "tree":
		return writeAddTree(stdout, result)
	case "full":
		return writeAddFull(stdout, result)
	default:
		return fmt.Errorf("unknown preview mode %q; use summary, patch, tree, or full", mode)
	}
}

func writeAddSummary(stdout io.Writer, result *workspace.AddResult) error {
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "Manifest change:"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "  Add: %s / %s\n", parentLabelForOutput(result.ParentPath), result.Heading); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "  From: %s\n", result.Reference); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	return nil
}

func writeAddPatch(stdout io.Writer, result *workspace.AddResult) error {
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "--- agents.yaml"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "+++ agents.yaml"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if err := writeAddHunk(stdout, result.ManifestHunk); err != nil {
		return err
	}
	if _, err := fmt.Fprintln(stdout, "--- AGENTS.md"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "+++ AGENTS.md"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if err := writeAddHunk(stdout, result.OutputHunk); err != nil {
		return err
	}
	return nil
}

func writeAddHunk(stdout io.Writer, hunk workspace.DiffHunk) error {
	if _, err := fmt.Fprintf(stdout, "@@ -%d,%d +%d,%d @@\n", hunk.OldStart, hunk.OldCount, hunk.NewStart, hunk.NewCount); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	for _, line := range hunk.Lines {
		if _, err := fmt.Fprintln(stdout, "+"+line); err != nil {
			return fmt.Errorf("write add preview: %w", err)
		}
	}
	return nil
}

func writeAddTree(stdout io.Writer, result *workspace.AddResult) error {
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "Document tree:"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, result.Tree); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	return nil
}

func writeAddFull(stdout io.Writer, result *workspace.AddResult) error {
	if _, err := fmt.Fprintln(stdout); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "Rendered preview:"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, result.Preview); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	return nil
}

func parentLabelForOutput(path string) string {
	if path == "" {
		return "document root"
	}
	return path
}

func runComplete(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("complete", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	sourceAlias := flags.String("source", "", "limit source-ref completion to one source alias")
	prefix := flags.String("prefix", "", "only print candidates with this prefix")
	args = reorderArgs(args, map[string]bool{
		"-manifest": true, "--manifest": true,
		"-source": true, "--source": true,
		"-prefix": true, "--prefix": true,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("complete requires one kind: commands, source-refs, source-aliases, manifest-headings, or flags")
	}
	var candidates []string
	switch flags.Arg(0) {
	case "commands":
		candidates = commandCandidates()
	case "flags":
		candidates = flagCandidates()
	case "source-refs":
		session, err := workspace.New(*manifestFile)
		if err != nil {
			return err
		}
		candidates, err = session.SourceReferences(*sourceAlias)
		if err != nil {
			return err
		}
	case "source-aliases":
		session, err := workspace.New(*manifestFile)
		if err != nil {
			return err
		}
		for _, source := range session.SavedSourceAliases() {
			candidates = append(candidates, source)
		}
	case "manifest-headings":
		session, err := workspace.New(*manifestFile)
		if err != nil {
			return err
		}
		candidates = session.ManifestHeadingPaths()
	default:
		return fmt.Errorf("unknown completion kind %q; use commands, source-refs, source-aliases, manifest-headings, or flags", flags.Arg(0))
	}
	for _, candidate := range candidates {
		if *prefix != "" && !strings.HasPrefix(candidate, *prefix) {
			continue
		}
		if _, err := fmt.Fprintln(stdout, candidate); err != nil {
			return fmt.Errorf("write completions: %w", err)
		}
	}
	return nil
}

func runCompletion(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("completion", flag.ContinueOnError)
	flags.SetOutput(stderr)
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("completion requires one shell: bash or zsh")
	}
	var script string
	switch flags.Arg(0) {
	case "bash":
		script = bashCompletionScript()
	case "zsh":
		script = zshCompletionScript()
	default:
		return fmt.Errorf("unknown shell %q; use bash or zsh", flags.Arg(0))
	}
	if _, err := fmt.Fprint(stdout, script); err != nil {
		return fmt.Errorf("write completion script: %w", err)
	}
	return nil
}

func commandCandidates() []string {
	return []string{"build", "status", "coverage", "source", "add", "complete", "completion", "tui", "help"}
}

func flagCandidates() []string {
	return []string{"--manifest", "--force", "--source", "--tag", "--tag-search", "--search", "--sort", "--metadata", "--file", "--line", "--content", "--lines", "--align-source", "--under", "--append", "--heading", "--dry-run", "--preview", "--rebuild", "--unused-only", "--tree", "--content-only", "--leaves-only", "--depth"}
}

func bashCompletionScript() string {
	return `_mogent_complete_candidates() {
  mogent complete "$1" --prefix "$2" 2>/dev/null
}

_mogent_completion() {
  local cur prev command subcommand
  COMPREPLY=()
  cur="${COMP_WORDS[COMP_CWORD]}"
  prev="${COMP_WORDS[COMP_CWORD-1]}"
  command="${COMP_WORDS[1]}"
  subcommand="${COMP_WORDS[2]}"

  if [[ ${COMP_CWORD} -eq 1 ]]; then
    mapfile -t COMPREPLY < <(_mogent_complete_candidates commands "$cur")
    return 0
  fi

  case "$prev" in
    --manifest|-manifest)
      COMPREPLY=( $(compgen -f -- "$cur") )
      return 0
      ;;
    --source|-source)
      mapfile -t COMPREPLY < <(_mogent_complete_candidates source-aliases "$cur")
      return 0
      ;;
    --under|-under)
      mapfile -t COMPREPLY < <(_mogent_complete_candidates manifest-headings "$cur")
      return 0
      ;;
  esac

  if [[ "$cur" == -* ]]; then
    mapfile -t COMPREPLY < <(_mogent_complete_candidates flags "$cur")
    return 0
  fi

  if [[ "$command" == "source" && "$subcommand" == "show" ]]; then
    mapfile -t COMPREPLY < <(_mogent_complete_candidates source-refs "$cur")
    return 0
  fi
  if [[ "$command" == "source" && "$subcommand" == "list" ]]; then
    mapfile -t COMPREPLY < <(_mogent_complete_candidates source-aliases "$cur")
    return 0
  fi
  if [[ "$command" == "add" ]]; then
    mapfile -t COMPREPLY < <(_mogent_complete_candidates source-refs "$cur")
    return 0
  fi
}

complete -F _mogent_completion mogent
`
}

func zshCompletionScript() string {
	return `#compdef mogent

_mogent_complete_candidates() {
  mogent complete "$1" --prefix "$2" 2>/dev/null
}

_mogent() {
  local -a candidates
  local command subcommand
  command="${words[2]}"
  subcommand="${words[3]}"

  if (( CURRENT == 2 )); then
    candidates=("${(@f)$(_mogent_complete_candidates commands "$PREFIX")}")
    _describe 'command' candidates
    return
  fi

  case "${words[CURRENT-1]}" in
    --manifest|-manifest)
      _files
      return
      ;;
    --source|-source)
      candidates=("${(@f)$(_mogent_complete_candidates source-aliases "$PREFIX")}")
      _describe 'source alias' candidates
      return
      ;;
    --under|-under)
      candidates=("${(@f)$(_mogent_complete_candidates manifest-headings "$PREFIX")}")
      _describe 'manifest heading' candidates
      return
      ;;
  esac

  if [[ "$PREFIX" == -* ]]; then
    candidates=("${(@f)$(_mogent_complete_candidates flags "$PREFIX")}")
    _describe 'flag' candidates
    return
  fi

  if [[ "$command" == "source" && "$subcommand" == "show" ]]; then
    candidates=("${(@f)$(_mogent_complete_candidates source-refs "$PREFIX")}")
    _describe 'source ref' candidates
    return
  fi
  if [[ "$command" == "source" && "$subcommand" == "list" ]]; then
    candidates=("${(@f)$(_mogent_complete_candidates source-aliases "$PREFIX")}")
    _describe 'source alias' candidates
    return
  fi
  if [[ "$command" == "add" ]]; then
    candidates=("${(@f)$(_mogent_complete_candidates source-refs "$PREFIX")}")
    _describe 'source ref' candidates
    return
  fi
}

_mogent "$@"
`
}

func runTUI(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("tui", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("tui accepts no positional arguments")
	}
	return navigator.Run(*manifestFile, stdout)
}

func unknownCommandError(command string) error {
	commands := commandCandidates()
	if suggestion := closestString(command, commands, 3); suggestion != "" {
		return fmt.Errorf("unknown command %q; did you mean %q? Use help for usage", command, suggestion)
	}
	return fmt.Errorf("unknown command %q; use build, status, coverage, source, add, complete, completion, or tui", command)
}

func decorateSourceReferenceError(reference string, err error) error {
	message := err.Error()
	if strings.Contains(message, "invalid source reference") && strings.HasSuffix(strings.TrimSpace(reference), "/") {
		trimmed := strings.TrimSuffix(strings.TrimSpace(reference), "/")
		return fmt.Errorf("%w; remove the trailing slash, or use `mogent source list --search %s` to find descendant heading paths", err, shellToken(trimmed))
	}
	return err
}

func closestString(target string, candidates []string, maxDistance int) string {
	best := ""
	bestDistance := maxDistance + 1
	for _, candidate := range candidates {
		distance := editDistance(strings.ToLower(target), strings.ToLower(candidate))
		if distance < bestDistance {
			best = candidate
			bestDistance = distance
		}
	}
	if bestDistance > maxDistance {
		return ""
	}
	return best
}

func editDistance(first string, second string) int {
	if first == second {
		return 0
	}
	a := []rune(first)
	b := []rune(second)
	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for index := range previous {
		previous[index] = index
	}
	for i, ar := range a {
		current[0] = i + 1
		for j, br := range b {
			cost := 1
			if ar == br {
				cost = 0
			}
			current[j+1] = minimum(current[j]+1, previous[j+1]+1, previous[j]+cost)
		}
		previous, current = current, previous
	}
	return previous[len(b)]
}

func minimum(values ...int) int {
	best := values[0]
	for _, value := range values[1:] {
		if value < best {
			best = value
		}
	}
	return best
}

func shellToken(value string) string {
	if value == "" {
		return "''"
	}
	if strings.ContainsAny(value, " \t\n'\"`$\\") {
		return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
	}
	return value
}
