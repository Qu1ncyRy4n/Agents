package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/internal/presentation"
	"github.com/Qu1ncyRy4n/Agents/library"
	"github.com/Qu1ncyRy4n/Agents/manifest"
	"github.com/Qu1ncyRy4n/Agents/sourcecache"
	"github.com/Qu1ncyRy4n/Agents/workspace"
	"github.com/charmbracelet/x/term"
)

func runCoverage(args []string, stdout, stderr io.Writer) error {
	return runSourceListMode(args, stdout, stderr, true)
}

func runSource(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("source requires a subcommand; use source add, list, show, pin, or update")
	}
	switch args[0] {
	case "add":
		return runSourceAdd(args[1:], stdout, stderr)
	case "list":
		return runSourceList(args[1:], stdout, stderr)
	case "show":
		return runSourceShow(args[1:], stdout, stderr)
	case "pin":
		return runSourcePin(args[1:], stdout, stderr)
	case "update":
		return runSourceUpdate(args[1:], stdout, stderr)
	default:
		if suggestion := closestString(args[0], []string{"add", "list", "show", "pin", "update"}, 2); suggestion != "" {
			return fmt.Errorf("unknown source subcommand %q; did you mean %q?", args[0], suggestion)
		}
		return fmt.Errorf("unknown source subcommand %q; use source add, list, show, pin, or update", args[0])
	}
}

func runSourceAdd(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source add", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	subdir := flags.String("subdir", "", "library root within an HTTP(S) Git repository")
	dryRun := flags.Bool("dry-run", false, "preview without writing")
	args = reorderArgs(args, map[string]bool{
		"-manifest": true, "--manifest": true,
		"-subdir": true, "--subdir": true,
		"-dry-run": false, "--dry-run": false,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 2 {
		return fmt.Errorf("source add requires an alias and path or URL")
	}
	result, err := workspace.AddSourceDeclaration(workspace.SourceAddOptions{
		ManifestPath: *manifestFile,
		Alias:        flags.Arg(0),
		Location:     flags.Arg(1),
		Subdir:       *subdir,
		DryRun:       *dryRun,
	})
	if err != nil {
		return err
	}
	if result.Wrote {
		if _, err := fmt.Fprintf(stdout, "Added source %s to %s\n", result.Alias, result.ManifestPath); err != nil {
			return fmt.Errorf("write source add result: %w", err)
		}
	} else if _, err := fmt.Fprintf(stdout, "Dry run: no files written\n\nProposed manifest:\n%s", result.ManifestYAML); err != nil {
		return fmt.Errorf("write source add preview: %w", err)
	}
	if result.Remote {
		_, err = fmt.Fprintf(stdout, "\nNext: mogent source pin %s\n", result.Alias)
	} else {
		_, err = fmt.Fprintf(stdout, "\nNext: mogent source list %s --tree\n", result.Alias)
	}
	if err != nil {
		return fmt.Errorf("write source add hint: %w", err)
	}
	return nil
}

func runSourcePin(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source pin", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	ref := flags.String("ref", "", "Git branch, tag, or full commit to pin")
	args = reorderArgs(args, map[string]bool{"-manifest": true, "--manifest": true, "-ref": true, "--ref": true})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("source pin requires exactly one source alias")
	}
	value, manifestPath, err := manifest.Load(*manifestFile)
	if err != nil {
		return err
	}
	alias := strings.TrimSuffix(flags.Arg(0), ":")
	source, found := value.Sources[alias]
	if !found {
		return fmt.Errorf("source %q is not declared in manifest", alias)
	}
	result, err := sourcecache.Pin(manifestPath, alias, source.Location, source.Subdir, *ref)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdout, "Pinned %s\nURL: %s\n", alias, result.URL); err != nil {
		return fmt.Errorf("write source pin result: %w", err)
	}
	if result.New.Subdir != "" {
		if _, err := fmt.Fprintf(stdout, "Subdir: %s\n", result.New.Subdir); err != nil {
			return fmt.Errorf("write source pin result: %w", err)
		}
	}
	if _, err := fmt.Fprintf(stdout, "Commit: %s\nContent SHA-256: %s\n", result.New.Commit, result.New.ContentSHA256); err != nil {
		return fmt.Errorf("write source pin result: %w", err)
	}
	if !result.Wrote {
		if _, err := fmt.Fprintln(stdout, "Cache and lock already verified"); err != nil {
			return fmt.Errorf("write source pin result: %w", err)
		}
	}
	return nil
}

func runSourceUpdate(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source update", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	ref := flags.String("ref", "", "Git branch, tag, or full commit to review")
	accept := flags.Bool("accept", false, "accept the reviewed candidate and update the lock")
	args = reorderArgs(args, map[string]bool{"-manifest": true, "--manifest": true, "-ref": true, "--ref": true, "-accept": false, "--accept": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("source update requires exactly one source alias")
	}
	value, manifestPath, err := manifest.Load(*manifestFile)
	if err != nil {
		return err
	}
	alias := strings.TrimSuffix(flags.Arg(0), ":")
	source, found := value.Sources[alias]
	if !found {
		return fmt.Errorf("source %q is not declared in manifest", alias)
	}
	result, err := sourcecache.Update(manifestPath, alias, source.Location, source.Subdir, *ref, *accept)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdout, "Source: %s\nOld commit: %s\nNew commit: %s\n", alias, result.Old.Commit, result.New.Commit); err != nil {
		return fmt.Errorf("write source update result: %w", err)
	}
	if len(result.Changes) == 0 {
		if _, err := fmt.Fprintln(stdout, "Markdown changes: none"); err != nil {
			return fmt.Errorf("write source update result: %w", err)
		}
	} else {
		if _, err := fmt.Fprintln(stdout, "Markdown changes:"); err != nil {
			return fmt.Errorf("write source update result: %w", err)
		}
		for _, change := range result.Changes {
			if _, err := fmt.Fprintf(stdout, "- %s  %s\n", change.Status, change.Path); err != nil {
				return fmt.Errorf("write source update result: %w", err)
			}
			if err := writeSourceChange(stdout, change); err != nil {
				return err
			}
		}
	}
	if result.Wrote {
		if _, err := fmt.Fprintln(stdout, "Accepted update and wrote mogent.lock.yaml"); err != nil {
			return fmt.Errorf("write source update result: %w", err)
		}
	} else if _, err := fmt.Fprintf(stdout, "Preview only: rerun with --ref %s --accept to install this exact candidate\n", result.New.Commit); err != nil {
		return fmt.Errorf("write source update result: %w", err)
	}
	return nil
}

func writeSourceChange(stdout io.Writer, change sourcecache.Change) error {
	if _, err := fmt.Fprintf(stdout, "  --- %s (old)\n  +++ %s (new)\n", change.Path, change.Path); err != nil {
		return fmt.Errorf("write source update diff: %w", err)
	}
	if change.Old == "" {
		if _, err := fmt.Fprintln(stdout, "  -(absent)"); err != nil {
			return fmt.Errorf("write source update diff: %w", err)
		}
	} else {
		for _, line := range strings.Split(strings.TrimSuffix(change.Old, "\n"), "\n") {
			if _, err := fmt.Fprintln(stdout, "  -"+line); err != nil {
				return fmt.Errorf("write source update diff: %w", err)
			}
		}
	}
	if change.New == "" {
		if _, err := fmt.Fprintln(stdout, "  +(absent)"); err != nil {
			return fmt.Errorf("write source update diff: %w", err)
		}
	} else {
		for _, line := range strings.Split(strings.TrimSuffix(change.New, "\n"), "\n") {
			if _, err := fmt.Fprintln(stdout, "  +"+line); err != nil {
				return fmt.Errorf("write source update diff: %w", err)
			}
		}
	}
	return nil
}

func runSourceList(args []string, stdout, stderr io.Writer) error {
	return runSourceListMode(args, stdout, stderr, false)
}

func runSourceListMode(args []string, stdout, stderr io.Writer, coveragePreset bool) error {
	presentationConfig, err := presentation.Load()
	if err != nil {
		return err
	}
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
	showCoverage := flags.Bool("coverage", coveragePreset, "overlay manifest selection state")
	tree := flags.Bool("tree", coveragePreset, "show source nodes as a hierarchy")
	unusedOnly := flags.Bool("unused-only", false, "show only unused source headings")
	contentOnly := flags.Bool("content-only", false, "hide organizational directories and headings without body content")
	leavesOnly := flags.Bool("leaves-only", false, "show only terminal source nodes")
	depth := flags.Int("depth", -1, "maximum source-tree depth to show; root nodes are depth 0")
	chars := flags.String("chars", presentationConfig.Chars, "tree character set: ascii or unicode")
	align := flags.Bool("align", presentationConfig.Align, "align fields into readable columns")
	fit := flags.String("fit", presentationConfig.Fit, "fit long rows to the terminal: term or none")
	width := flags.Int("width", presentationConfig.Width, "wrap output at this width; zero uses --fit")
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
		"-depth":       true,
		"--depth":      true,
		"-chars":       true,
		"--chars":      true,
		"-fit":         true,
		"--fit":        true,
		"-width":       true,
		"--width":      true,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() > 1 {
		if coveragePreset {
			return fmt.Errorf("coverage accepts at most one source alias positional argument")
		}
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
	if *chars != "ascii" && *chars != "unicode" {
		return fmt.Errorf("unknown display character set %q; use ascii or unicode", *chars)
	}
	if *fit != "term" && *fit != "none" {
		return fmt.Errorf("unknown display fit %q; use term or none", *fit)
	}
	if *width < 0 {
		return fmt.Errorf("display width must be zero or greater")
	}
	if *tree && *sortMode != "path" {
		return fmt.Errorf("tree display requires --sort path")
	}
	if *tree && (*showMetadata || *showFile || *showLine) {
		return fmt.Errorf("tree display does not support --metadata, --file, or --line; use --tldr or omit --tree")
	}
	if *showCoverage || *tree || *unusedOnly || *contentOnly || *leavesOnly || *depth >= 0 {
		coverage := session.CoverageWithOptions(workspace.CoverageOptions{
			SourceAlias: *sourceAlias,
			Tag:         *tag,
			ContentOnly: *contentOnly,
			LeavesOnly:  *leavesOnly,
			LimitDepth:  *depth >= 0,
			MaxDepth:    *depth,
		})
		return writeInventory(stdout, nodes, coverage, inventoryOptions{
			Coverage:   *showCoverage,
			Tree:       *tree,
			UnusedOnly: *unusedOnly,
			ShowTLDR:   *tldrOnly,
			Align:      *align,
			Chars:      *chars,
			Width:      outputWidth(stdout, *fit, *width),
		})
	}
	if err := writeAlignedSourceList(stdout, nodes, *align, *showFile, *showLine, *chars); err != nil {
		return err
	}
	if *showMetadata && !*tldrOnly {
		for _, node := range nodes {
			if _, err := fmt.Fprintf(stdout, "\n%s\n", node.Reference); err != nil {
				return fmt.Errorf("write source list: %w", err)
			}
			if err := writeIndentedSourceMetadata(stdout, node.Metadata); err != nil {
				return err
			}
		}
	}
	return nil
}

func outputWidth(stdout io.Writer, fit string, configured int) int {
	if configured > 0 {
		return configured
	}
	if fit != "term" {
		return 0
	}
	file, ok := stdout.(interface{ Fd() uintptr })
	if !ok || !term.IsTerminal(file.Fd()) {
		return 0
	}
	width, _, err := term.GetSize(file.Fd())
	if err != nil || width <= 0 {
		return 0
	}
	return width
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
	if _, err := fmt.Fprintf(stdout, "Kind: %s\n", node.Kind); err != nil {
		return fmt.Errorf("write source: %w", err)
	}
	if *showFile && node.File != "" {
		if _, err := fmt.Fprintf(stdout, "File: %s\n", node.File); err != nil {
			return fmt.Errorf("write source: %w", err)
		}
	}
	if *showLine && node.Line > 0 {
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
