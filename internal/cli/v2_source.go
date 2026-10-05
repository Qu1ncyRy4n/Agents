package cli

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/v2"
)

// v2SourceRequested decides whether a source subcommand addresses a v2
// configuration: an explicit --config, or a mogent.hcl beside no agents.yaml
// when --manifest is absent.
func v2SourceRequested(args []string) bool {
	for _, arg := range args {
		name, _, _ := strings.Cut(arg, "=")
		switch name {
		case "--config", "-config":
			return true
		case "--manifest", "-manifest":
			return false
		}
	}
	if _, err := os.Stat(v2.ConfigFile); err != nil {
		return false
	}
	_, err := os.Stat("agents.yaml")
	return errors.Is(err, os.ErrNotExist)
}

func loadV2Sources(configPath string) (*v2.Config, map[string]v2.ResolvedSource, map[string]*v2.Library, error) {
	config, err := v2.Load(configPath)
	if err != nil {
		return nil, nil, nil, err
	}
	sources, _, err := v2.ResolveSources(config)
	if err != nil {
		return nil, nil, nil, err
	}
	libraries, err := v2.LoadLibraries(sources)
	if err != nil {
		return nil, nil, nil, err
	}
	return config, sources, libraries, nil
}

// printer buffers text output and keeps the first write error so callers
// check once at flush.
type printer struct {
	writer *bufio.Writer
	err    error
}

func newPrinter(output io.Writer) *printer {
	return &printer{writer: bufio.NewWriter(output)}
}

func (p *printer) printf(format string, args ...any) {
	if p.err != nil {
		return
	}
	_, p.err = fmt.Fprintf(p.writer, format, args...)
}

func (p *printer) line(text string) {
	p.printf("%s\n", text)
}

func (p *printer) flush() error {
	if p.err != nil {
		return fmt.Errorf("write output: %w", p.err)
	}
	if err := p.writer.Flush(); err != nil {
		return fmt.Errorf("write output: %w", err)
	}
	return nil
}

// runV2SourceList prints every declared library as a tree of section paths
// with titles and offers, followed by its raw trees.
func runV2SourceList(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source list", flag.ContinueOnError)
	flags.SetOutput(stderr)
	configPath := flags.String("config", v2.ConfigFile, "path to v2 HCL configuration")
	showTLDR := flags.Bool("tldr", false, "show each section's tldr")
	showTags := flags.Bool("tags", false, "show each section's direct tags")
	args = reorderArgs(args, map[string]bool{"-config": true, "--config": true, "-tldr": false, "--tldr": false, "-tags": false, "--tags": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() > 1 {
		return fmt.Errorf("source list accepts at most one source alias")
	}
	config, sources, libraries, err := loadV2Sources(*configPath)
	if err != nil {
		return err
	}
	aliases := make([]string, 0, len(config.Sources))
	for alias := range config.Sources {
		if flags.Arg(0) == "" || alias == flags.Arg(0) {
			aliases = append(aliases, alias)
		}
	}
	if len(aliases) == 0 {
		return fmt.Errorf("source %q is not declared in %s", flags.Arg(0), *configPath)
	}
	sort.Strings(aliases)
	out := newPrinter(stdout)
	for index, alias := range aliases {
		if index > 0 {
			out.line("")
		}
		library := libraries[alias]
		out.printf("%s: %s (%s)\n", alias, library.Name, describeSource(config.Sources[alias], sources[alias]))
		out.printf("%s:%s\n", alias, library.MarkdownRoot)
		for _, section := range library.Sections {
			writeSectionRows(out, section, nil, 1, *showTLDR, *showTags)
		}
		for _, name := range sortedTreeNames(library) {
			tree := library.Trees[name]
			out.printf("%s:%s  (tree, root %s)\n", alias, name, tree.Root)
			for _, entry := range tree.Entries {
				row := "  " + entry.Name
				if *showTags && len(entry.Tags) > 0 {
					row += "  #" + strings.Join(entry.Tags, " #")
				}
				if *showTLDR && entry.TLDR != "" {
					row += "  - " + entry.TLDR
				}
				out.line(row)
			}
		}
	}
	return out.flush()
}

func describeSource(source v2.Source, resolved v2.ResolvedSource) string {
	if source.Local != "" {
		description := "local " + source.Local
		if resolved.LocalHead != "" {
			description += " at " + resolved.LocalHead
			if resolved.LocalDirty {
				description += ", dirty"
			}
		}
		return description
	}
	description := "git " + source.Git
	if resolved.Pinned {
		return description + " at " + resolved.Commit
	}
	return description + " unpinned, resolved " + resolved.Commit
}

func writeSectionRows(out *printer, section, parent *v2.Section, depth int, showTLDR, showTags bool) {
	row := strings.Repeat("  ", depth) + section.Name + "  " + section.Title
	if label := offerLabel(section); label != "" {
		row += "  [" + label + "]"
	}
	if parent != nil && parent.Offer == v2.OfferChoose {
		row += fmt.Sprintf("  (default %t)", parent.Defaults[section.Name])
	}
	if showTags && len(section.Tags) > 0 {
		row += "  #" + strings.Join(section.Tags, " #")
	}
	if showTLDR && section.TLDR != "" {
		row += "  - " + section.TLDR
	}
	out.line(row)
	for _, child := range section.Children {
		writeSectionRows(out, child, section, depth+1, showTLDR, showTags)
	}
}

func offerLabel(section *v2.Section) string {
	switch section.Offer {
	case "":
		return ""
	case v2.OfferOptional:
		return fmt.Sprintf("optional, default %t", *section.Default)
	default:
		return section.Offer
	}
}

func sortedTreeNames(library *v2.Library) []string {
	names := make([]string, 0, len(library.Trees))
	for name := range library.Trees {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// runV2SourceShow prints one section or tree in full: identity, offer,
// effective tags, tldr, source file, and the first lines of its body.
func runV2SourceShow(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("source show", flag.ContinueOnError)
	flags.SetOutput(stderr)
	configPath := flags.String("config", v2.ConfigFile, "path to v2 HCL configuration")
	lines := flags.Int("lines", 12, "number of body lines to show; 0 shows the whole body")
	args = reorderArgs(args, map[string]bool{"-config": true, "--config": true, "-lines": true, "--lines": true})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 {
		return fmt.Errorf("source show requires exactly one alias:path reference")
	}
	alias, path, found := strings.Cut(flags.Arg(0), ":")
	if !found || alias == "" || path == "" {
		return fmt.Errorf("reference %q must be alias:section/path or alias:tree", flags.Arg(0))
	}
	_, _, libraries, err := loadV2Sources(*configPath)
	if err != nil {
		return err
	}
	library, ok := libraries[alias]
	if !ok {
		return fmt.Errorf("source %q is not declared in %s", alias, *configPath)
	}
	out := newPrinter(stdout)
	if tree, ok := library.Trees[path]; ok {
		out.printf("%s:%s\nKind:   tree\nRoot:   %s\n", alias, path, filepath.Join(library.Root, filepath.FromSlash(tree.Root)))
		for _, entry := range tree.Entries {
			row := fmt.Sprintf("Entry:  %s (%s)", entry.Name, entry.Path)
			if len(entry.Tags) > 0 {
				row += "  #" + strings.Join(entry.Tags, " #")
			}
			if entry.TLDR != "" {
				row += "  - " + entry.TLDR
			}
			out.line(row)
		}
		return out.flush()
	}
	section, ok := library.ByPath[path]
	if !ok {
		return fmt.Errorf("source %q has no section or tree %q; use `mogent source list %s`", alias, path, alias)
	}
	out.printf("%s:%s\nTitle:  %s\n", alias, path, section.Title)
	if label := offerLabel(section); label != "" {
		out.printf("Offer:  %s\n", label)
		if section.Offer == v2.OfferChoose {
			for _, child := range section.Children {
				out.printf("        %s = %t\n", child.Name, section.Defaults[child.Name])
			}
		}
	}
	if parentPath := filepath.ToSlash(filepath.Dir(path)); parentPath != "." {
		if parent := library.ByPath[parentPath]; parent != nil && parent.Offer != "" {
			row := "Parent: " + parentPath + " offers " + parent.Offer
			if parent.Offer == v2.OfferChoose {
				row += fmt.Sprintf(", default %t", parent.Defaults[section.Name])
			}
			out.line(row)
		}
	}
	direct := make(map[string]bool, len(section.Tags))
	for _, tag := range section.Tags {
		direct[tag] = true
	}
	var inherited []string
	for _, tag := range library.EffectiveTags(section) {
		if !direct[tag] {
			inherited = append(inherited, tag)
		}
	}
	if len(section.Tags) > 0 || len(inherited) > 0 {
		row := "Tags:   " + strings.Join(section.Tags, ", ")
		if len(inherited) > 0 {
			row += "  (inherited: " + strings.Join(inherited, ", ") + ")"
		}
		out.line(strings.TrimRight(row, " "))
	}
	if section.TLDR != "" {
		out.printf("TLDR:   %s\n", section.TLDR)
	}
	if len(section.Children) > 0 {
		names := make([]string, 0, len(section.Children))
		for _, child := range section.Children {
			names = append(names, child.Name)
		}
		out.printf("Children: %s\n", strings.Join(names, ", "))
	}
	if section.Source != "" {
		file := filepath.Join(library.Root, filepath.FromSlash(library.MarkdownRoot), filepath.FromSlash(section.Source))
		out.printf("Source: %s\n", file)
		contents, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		body := strings.Split(strings.TrimRight(string(contents), "\n"), "\n")
		if *lines > 0 && len(body) > *lines {
			body = append(body[:*lines], fmt.Sprintf("... (%d more lines)", len(body)-*lines))
		}
		out.line("---")
		for _, row := range body {
			out.line(row)
		}
	}
	return out.flush()
}
