package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/internal/presentation"
	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

func hasFlag(args []string, wanted string) bool {
	for _, arg := range args {
		if arg == wanted || strings.HasPrefix(arg, wanted+"=") {
			return true
		}
	}
	return false
}

func runLocalSourceUpdate(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source update local", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "mogent.build", "authored build file")
	dryRun := flags.Bool("dry-run", false, "preview without writing")
	args = reorderArgs(args, map[string]bool{"-file": true, "--file": true, "-dry-run": false, "--dry-run": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("source update local accepts no positional arguments")
	}
	file, err := assembly.LoadFile(*path)
	if err != nil {
		return err
	}
	updated, changes, err := assembly.RefreshLocalInputHashes(file, *path)
	if err != nil {
		return err
	}
	if len(changes) == 0 {
		_, err = fmt.Fprintln(stdout, "Local source inputs are already current")
		return err
	}
	for _, change := range changes {
		if _, err := fmt.Fprintf(stdout, "%s\n  %s -> %s\n", change.Input, change.Old, change.New); err != nil {
			return err
		}
	}
	if *dryRun {
		_, err = fmt.Fprintln(stdout, "Dry run: no files written")
		return err
	}
	if err := renderfs.WriteAtomically(*path, updated); err != nil {
		return err
	}
	_, err = fmt.Fprintln(stdout, "Accepted local source updates. Next: run mogent sequence --file mogent.plan")
	return err
}

func runSequence(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("sequence", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "mogent.plan", "plan file")
	args = reorderArgs(args, map[string]bool{"-file": true, "--file": true})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("sequence accepts no positional arguments")
	}
	return sequenceFile(*path, stdout)
}

func sequenceFile(planPath string, stdout io.Writer) error {
	contents, buildPath, err := sequenceContents(planPath)
	if err != nil {
		return err
	}
	if err := os.WriteFile(buildPath, contents, 0644); err != nil {
		return err
	}
	_, err = fmt.Fprintf(stdout, "Sequenced %s -> %s\n", planPath, buildPath)
	return err
}

func sequenceContents(planPath string) ([]byte, string, error) {
	plan, err := assembly.LoadFile(planPath)
	if err != nil {
		return nil, "", err
	}
	contents, result, err := assembly.Sequence(plan, planPath)
	if err != nil {
		return nil, "", err
	}
	if result.HasErrors() {
		writeAssemblyDiagnostics(stderrDiscard{}, result)
		return nil, "", fmt.Errorf("sequence contains errors")
	}
	buildPath := filepath.Join(filepath.Dir(planPath), "mogent.build")
	if existing, err := os.ReadFile(buildPath); err == nil {
		existingFile, parseErr := syntax.Parse(buildPath, existing)
		if parseErr != nil {
			return nil, "", fmt.Errorf("refusing to replace unreadable %s: %w", buildPath, parseErr)
		}
		origin, found, originErr := buildOrigin(existingFile)
		if originErr != nil || !found || origin == "authored" {
			return nil, "", fmt.Errorf("refusing to replace authored %s; select it explicitly with build --file", buildPath)
		}
		if err := assembly.VerifyBuildChecksum(existing); err != nil {
			return nil, "", fmt.Errorf("refusing to replace %s: %w", buildPath, err)
		}
	} else if !os.IsNotExist(err) {
		return nil, "", err
	}
	return contents, buildPath, nil
}

func buildOrigin(file *syntax.File) (string, bool, error) {
	if file.Kind != syntax.Build || len(file.Entries) == 0 {
		return "", false, nil
	}
	for _, child := range file.Entries[0].Body {
		if child.Kind == syntax.Field && child.Name == "origin" {
			value, err := child.Value.ResolveString(nil)
			return value, true, err
		}
	}
	return "", false, nil
}

// stderrDiscard keeps API diagnostics out of successful command stdout while
// callers that need rich display use show. CLI error messages still name blockers.
type stderrDiscard struct{}

func (stderrDiscard) Write([]byte) (int, error) { return 0, nil }

func writeAssemblyDiagnostics(output io.Writer, result *assembly.Result) {
	for _, diagnostic := range result.Diagnostics {
		_, _ = fmt.Fprintln(output, diagnostic.String())
	}
}

func runNewBuild(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("build", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "plan or build file")
	check := flags.Bool("check", false, "verify without writing workspace outputs")
	args = reorderArgs(args, map[string]bool{"-file": true, "--file": true, "-check": false, "--check": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 || *path == "" {
		return fmt.Errorf("new build requires --file mogent.plan or mogent.build")
	}
	file, err := assembly.LoadFile(*path)
	if err != nil {
		return err
	}
	buildPath := *path
	var generatedBuild []byte
	if file.Kind == syntax.Plan {
		if *check {
			result, _, err := assembly.ResolvePlan(file, *path)
			if err != nil {
				return err
			}
			if result.HasErrors() {
				writeAssemblyDiagnostics(stderr, result)
				return fmt.Errorf("build check contains errors")
			}
			manifest, err := result.Compose()
			if err != nil {
				return err
			}
			if err := assembly.CheckWorkspace(filepath.Dir(*path), manifest); err != nil {
				return err
			}
			_, err = fmt.Fprintln(stdout, "Build check passed")
			return err
		}
		generatedBuild, buildPath, err = sequenceContents(*path)
		if err != nil {
			return err
		}
		file, err = syntax.Parse(buildPath, generatedBuild)
		if err != nil {
			return err
		}
	}
	if file.Kind != syntax.Build {
		return fmt.Errorf("build file must contain a build header")
	}
	contents := generatedBuild
	if contents == nil {
		contents, err = os.ReadFile(buildPath)
		if err != nil {
			return err
		}
	}
	if err := assembly.VerifyBuildChecksum(contents); err != nil {
		return err
	}
	result, _, err := assembly.ResolvePlan(file, buildPath)
	if err != nil {
		return err
	}
	if result.HasErrors() {
		writeAssemblyDiagnostics(stderr, result)
		return fmt.Errorf("build contains errors")
	}
	manifest, err := result.Compose()
	if err != nil {
		return err
	}
	if *check {
		if err := assembly.CheckWorkspace(filepath.Dir(buildPath), manifest); err != nil {
			return err
		}
		_, err := fmt.Fprintln(stdout, "Build check passed")
		return err
	}
	workspace := filepath.Dir(buildPath)
	if generatedBuild != nil {
		if err := assembly.ApplyGenerated(workspace, manifest, generatedBuild); err != nil {
			return err
		}
	} else if err := assembly.Apply(workspace, manifest); err != nil {
		return err
	}
	_, err = fmt.Fprintf(stdout, "Built %s\n", workspace)
	return err
}

func runShow(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("show", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "plan or build file (discovered by default)")
	sourceView := flags.Bool("source", false, "show source catalog")
	contentView := flags.Bool("content", false, "show one source module's rendered or physical content")
	lines := flags.Int("lines", 0, "content preview lines; 0 shows all")
	headings := flags.Bool("headings", false, "show heading addresses for a physical Markdown file")
	tree := flags.Bool("tree", false, "show source hierarchy (default source view)")
	flat := flags.Bool("flat", false, "show full source paths without tree connectors")
	chars := flags.String("chars", "", "tree characters: ascii or unicode")
	buildView := flags.Bool("build", false, "show ordered contributions")
	issues := flags.Bool("issues", false, "show only diagnostics")
	args = reorderArgs(args, map[string]bool{"-file": true, "--file": true, "-source": false, "--source": false, "-content": false, "--content": false, "-lines": true, "--lines": true, "-headings": false, "--headings": false, "-build": false, "--build": false, "-issues": false, "--issues": false, "-tree": false, "--tree": false, "-flat": false, "--flat": false, "-chars": true, "--chars": true})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() > 1 || *sourceView && *buildView || *contentView && *buildView || *contentView && *issues {
		return fmt.Errorf("show accepts at most one node and only one of --source, --content, or --build")
	}
	if *lines < 0 {
		return fmt.Errorf("--lines must be zero or greater")
	}
	if *tree && *flat || (*flat || *chars != "") && !*sourceView {
		return fmt.Errorf("--flat and --chars apply to source show; choose only one of --tree or --flat")
	}
	requested := ""
	if flags.NArg() == 1 {
		requested = flags.Arg(0)
		if _, _, found := strings.Cut(requested, ":"); !found {
			if !*sourceView {
				return fmt.Errorf("show node must be alias:path")
			}
			requested += ":"
		}
	}
	if (*contentView || *headings) && requested == "" {
		return fmt.Errorf("--content and --headings require one source:path reference")
	}
	if *headings && !*contentView {
		return fmt.Errorf("--headings requires --content")
	}
	entryPath, err := projectEntryFile(*path)
	if err != nil {
		return err
	}
	file, err := assembly.LoadFile(entryPath)
	if err != nil {
		return err
	}
	if *contentView {
		preview, err := assembly.PreviewSource(file, entryPath, requested, *lines, *headings)
		if err != nil {
			return err
		}
		if _, err := fmt.Fprintf(stdout, "%s\nType: %s\n", preview.Reference, preview.Type); err != nil {
			return err
		}
		if preview.Source != "" {
			if _, err := fmt.Fprintf(stdout, "Source: %s\n", preview.Source); err != nil {
				return err
			}
		}
		if preview.Opaque {
			_, err := fmt.Fprintf(stdout, "Payload: opaque (%d bytes)\n", preview.Bytes)
			return err
		}
		if len(preview.Headings) > 0 {
			if _, err := fmt.Fprintln(stdout, "Heading addresses:"); err != nil {
				return err
			}
			for _, heading := range preview.Headings {
				quoted := make([]string, len(heading))
				for index, title := range heading {
					quoted[index] = strconv.Quote(title)
				}
				if _, err := fmt.Fprintf(stdout, "  heading = [%s]\n", strings.Join(quoted, ", ")); err != nil {
					return err
				}
			}
			return nil
		}
		if preview.Children != nil {
			for _, child := range preview.Children {
				if _, err := fmt.Fprintf(stdout, "  %s\n", child); err != nil {
					return err
				}
			}
			return nil
		}
		if len(preview.Content) > 0 {
			_, err = fmt.Fprintf(stdout, "---\n%s", preview.Content)
		}
		return err
	}
	result, sources, err := assembly.ResolvePlan(file, entryPath)
	if err != nil {
		return err
	}
	if *issues {
		writeAssemblyDiagnostics(stdout, result)
		for _, node := range result.Nodes {
			if node.Status != assembly.Selected && (requested == "" || requested == node.Source+":"+node.Path) {
				_, _ = fmt.Fprintf(stdout, "%s:%s [%s] %s\n", node.Source, node.Path, node.Status, node.Reason)
			}
		}
		if result.HasErrors() {
			return fmt.Errorf("show found errors")
		}
		return nil
	}
	if *sourceView {
		if len(sources) == 0 {
			if requested != "" {
				return fmt.Errorf("unknown source node %q", requested)
			}
			return writeNoSources(stdout, entryPath)
		}
		config, err := presentation.Load()
		if err != nil {
			return err
		}
		if *chars != "" {
			config.Chars = *chars
		}
		if config.Chars != "ascii" && config.Chars != "unicode" {
			return fmt.Errorf("--chars must be ascii or unicode")
		}
		return writeCatalog(stdout, sources, requested, *flat, config.Align, config.Chars)
	}
	if *tree {
		return writeContributionTree(stdout, result.Contributions, requested)
	}
	for _, contribution := range result.Contributions {
		name := contribution.Source
		if contribution.Node != nil {
			name += ":" + contribution.Node.Path
		} else {
			name = "inline"
		}
		if requested != "" && name != requested {
			continue
		}
		if _, err := fmt.Fprintf(stdout, "%s -> %s\n", name, contribution.Into); err != nil {
			return err
		}
	}
	if *buildView {
		for _, overlap := range assembly.DirectIncludeOverlaps(file) {
			if _, err := fmt.Fprintf(stdout, "warning: %s also directly includes %s; deliberate repetition is allowed. Inspect with mogent optimize --file %s\n", overlap.Ancestor, overlap.Descendant, entryPath); err != nil {
				return err
			}
		}
	}
	if requested != "" {
		for _, node := range result.Nodes {
			if requested == node.Source+":"+node.Path && node.Status != assembly.Selected {
				if _, err := fmt.Fprintf(stdout, "%s [%s] %s\n", requested, node.Status, node.Reason); err != nil {
					return err
				}
			}
		}
	}
	writeAssemblyDiagnostics(stdout, result)
	if result.HasErrors() {
		return fmt.Errorf("show found errors")
	}
	return nil
}

type contributionTreeNode struct {
	key      string
	label    string
	source   string
	count    int
	children []*contributionTreeNode
	byKey    map[string]*contributionTreeNode
}

func (n *contributionTreeNode) child(key, label, source string) *contributionTreeNode {
	if found := n.byKey[key]; found != nil {
		return found
	}
	child := &contributionTreeNode{key: key, label: label, source: source, byKey: map[string]*contributionTreeNode{}}
	n.byKey[key] = child
	n.children = append(n.children, child)
	return child
}

// writeContributionTree shows the composed result rather than the source
// catalog: Markdown contributions are nested by their heading ancestry, and
// source aliases occupy a narrow column instead of expanding every tree label.
func writeContributionTree(stdout io.Writer, contributions []assembly.Contribution, requested string) error {
	roots := map[string]*contributionTreeNode{}
	var outputs []*contributionTreeNode
	for _, contribution := range contributions {
		name := contribution.Source
		if contribution.Node != nil {
			name += ":" + contribution.Node.Path
		} else {
			name = "inline"
		}
		if requested != "" && requested != name {
			continue
		}
		output := contribution.Into
		suffix := ""
		if strings.HasPrefix(output, ".agents/skills/") {
			suffix = strings.TrimPrefix(output, ".agents/skills/")
			output = ".agents/skills"
		}
		root := roots[output]
		if root == nil {
			root = &contributionTreeNode{key: output, label: output, byKey: map[string]*contributionTreeNode{}}
			roots[output] = root
			outputs = append(outputs, root)
		}
		current := root
		if contribution.Node != nil && contribution.Node.Type == assembly.Section {
			var ancestors []*assembly.Node
			for node := contribution.Node; node != nil; node = node.Parent {
				ancestors = append(ancestors, node)
			}
			for index := len(ancestors) - 1; index >= 0; index-- {
				node := ancestors[index]
				label := node.Title
				if label == "" {
					label = node.Name
				}
				current = current.child(contribution.Source+":"+node.Path, label, contribution.Source)
			}
		} else {
			label := suffix
			if label == "" && contribution.Node != nil {
				label = contribution.Node.Name
			}
			if label == "" {
				label = "inline"
			}
			current = current.child(name, label, contribution.Source)
		}
		current.count++
	}
	if len(outputs) == 0 {
		_, err := fmt.Fprintln(stdout, "No selected contributions.")
		return err
	}
	rows := [][]string{{"SOURCE", "RESULT"}}
	for outputIndex, output := range outputs {
		connector, prefix := "├── ", "│   "
		if outputIndex == len(outputs)-1 {
			connector, prefix = "└── ", "    "
		}
		rows = append(rows, []string{"", connector + output.label})
		appendContributionTreeRows(&rows, output.children, prefix)
	}
	return writeColumns(stdout, rows, true)
}

func appendContributionTreeRows(rows *[][]string, nodes []*contributionTreeNode, prefix string) {
	for index, node := range nodes {
		connector, next := "├── ", prefix+"│   "
		if index == len(nodes)-1 {
			connector, next = "└── ", prefix+"    "
		}
		label := node.label
		if node.count > 1 {
			label += fmt.Sprintf(" ×%d", node.count)
		}
		*rows = append(*rows, []string{node.source, prefix + connector + label})
		appendContributionTreeRows(rows, node.children, next)
	}
}

func runEdit(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "include" {
		return fmt.Errorf("edit requires include; use edit include alias:path --localize")
	}
	flags := flag.NewFlagSet("edit include", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "mogent.plan", "plan file")
	localize := flags.Bool("localize", false, "create editable local source")
	dryRun := flags.Bool("dry-run", false, "preview without writing")
	args = reorderArgs(args[1:], map[string]bool{"-file": true, "--file": true, "-localize": false, "--localize": false, "-dry-run": false, "--dry-run": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 || !*localize {
		return fmt.Errorf("edit include currently requires one reference and --localize")
	}
	plan, err := assembly.LoadFile(*path)
	if err != nil {
		return err
	}
	result, err := assembly.Localize(plan, *path, flags.Arg(0), *dryRun)
	if err != nil {
		return err
	}
	if *dryRun {
		_, err = fmt.Fprintf(stdout, "Dry run: local source %s\n\n%s", result.LocalPath, result.Plan)
		return err
	}
	_, err = fmt.Fprintf(stdout, "Localized %s as %s\n", result.SourceReference, result.LocalReference)
	return err
}
