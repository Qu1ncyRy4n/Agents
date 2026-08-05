// Package cli exposes the deliberately small milestone-one command surface.
package cli

import (
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
	"github.com/Qu1ncyRy4n/Agents/internal/navigator"
	"github.com/Qu1ncyRy4n/Agents/internal/render"
	"github.com/Qu1ncyRy4n/Agents/internal/state"
	"github.com/Qu1ncyRy4n/Agents/internal/workspace"
)

func Run(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] == "help" || args[0] == "--help" || args[0] == "-h" {
		if _, err := fmt.Fprintln(stdout, "Usage: mogent build [--manifest agents.yaml] [--force]\n       mogent status [--manifest agents.yaml]\n       mogent coverage [--manifest agents.yaml]\n       mogent source show <ref> [--manifest agents.yaml]\n       mogent add <ref> [--under heading/path | --append]\n       mogent tui [--manifest agents.yaml]"); err != nil {
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
	case "tui":
		return runTUI(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown command %q; use build, status, coverage, source, add, or tui", args[0])
	}
}

func runBuild(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("build", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	force := flags.Bool("force", false, "replace an untracked or directly edited output")
	if err := flags.Parse(args); err != nil {
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
	if err := flags.Parse(args); err != nil {
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

func runCoverage(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("coverage", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	sourceAlias := flags.String("source", "", "show only one source alias")
	contentOnly := flags.Bool("content-only", false, "hide source nodes without body content")
	leavesOnly := flags.Bool("leaves-only", false, "show only terminal source nodes")
	depth := flags.Int("depth", -1, "maximum source-tree depth to show; root headings are depth 0")
	unusedOnly := flags.Bool("unused-only", false, "show only unused source references")
	tree := flags.Bool("tree", false, "show unused source references as an ASCII tree")
	if err := flags.Parse(args); err != nil {
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
		ContentOnly: *contentOnly,
		LeavesOnly:  *leavesOnly,
		LimitDepth:  *depth >= 0,
		MaxDepth:    *depth,
	})
	if *sourceAlias != "" && len(coverage.Sources) == 0 {
		return fmt.Errorf("source %q is not declared in manifest", *sourceAlias)
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
		if _, err := fmt.Fprintf(stdout, "%s  %s  unused %d/%d\n", source.Alias, source.Path, len(source.Unused), source.Total); err != nil {
			return fmt.Errorf("write coverage: %w", err)
		}
		if len(source.Unused) == 0 {
			if _, err := fmt.Fprintln(stdout, "  (none)"); err != nil {
				return fmt.Errorf("write coverage: %w", err)
			}
		}
		for index, node := range source.Unused {
			depth := visibleDepth(source.Unused, index)
			connector := "`-- "
			if hasNextAtDepth(source.Unused, index, depth) {
				connector = "|-- "
			}
			indent := strings.Repeat("|   ", depth)
			if _, err := fmt.Fprintf(stdout, "%s%s%s  %s:%s\n", indent, connector, node.Heading, source.Alias, node.Path); err != nil {
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
		return fmt.Errorf("source requires a subcommand; use source show")
	}
	switch args[0] {
	case "show":
		return runSourceShow(args[1:], stdout, stderr)
	default:
		return fmt.Errorf("unknown source subcommand %q; use source show", args[0])
	}
}

func runSourceShow(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source show", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	showFile := flags.Bool("file", true, "show source file path")
	showLine := flags.Bool("line", true, "show source heading line")
	contentMode := flags.String("content", "full", "content mode: none, snippet, or full")
	snippetLines := flags.Int("lines", 12, "number of content lines when --content=snippet")
	args = reorderArgs(args, map[string]bool{
		"-manifest":  true,
		"--manifest": true,
		"-content":   true,
		"--content":  true,
		"-lines":     true,
		"--lines":    true,
	})
	if err := flags.Parse(args); err != nil {
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
		return err
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
	if err := flags.Parse(args); err != nil {
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
		return err
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
	if _, err := fmt.Fprintln(stdout, "@@"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintf(stdout, "+  - %s: %s\n", result.Heading, result.Reference); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "--- AGENTS.md"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "+++ AGENTS.md"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	if _, err := fmt.Fprintln(stdout, "@@"); err != nil {
		return fmt.Errorf("write add preview: %w", err)
	}
	for _, line := range strings.Split(strings.TrimRight(result.Section, "\n"), "\n") {
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

func runTUI(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("tui", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	if err := flags.Parse(args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("tui accepts no positional arguments")
	}
	return navigator.Run(*manifestFile, stdout)
}
