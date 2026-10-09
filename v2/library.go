package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

const LibraryFile = "library.mogent.hcl"

// Library is the validated server-side content tree for one v2 library.
type Library struct {
	Path         string
	Root         string
	ID           string
	Name         string
	MarkdownRoot string
	TreeRoot     string
	Trees        map[string]Tree
	Sections     []*Section
	ByPath       map[string]*Section
}

// Tree declares one raw server-side directory root eligible for a tree output.
// Entries are optional discovery metadata for subdirectories of that root;
// they never change what a tree output copies.
type Tree struct {
	Name    string
	Root    string
	Entries []TreeEntry
}

// TreeEntry describes one subdirectory of a tree for browsing.
type TreeEntry struct {
	Name string
	Path string
	TLDR string
	Tags []string
}

// Section is one hierarchy node. Tags are metadata attached to this node; they
// are not an alternative structure.
//
// Offer is the library's promise about a branch's direct children when the
// consumer selects the branch broadly: "foundation" includes them all and a
// drop needs a reason; "choose" needs a decision per child or accept_defaults;
// "optional" applies Default; "opt_in" excludes them unless named. A leaf has
// no Offer.
type Section struct {
	Name     string
	Title    string
	Source   string
	TLDR     string
	Tags     []string
	Offer    string
	Defaults map[string]bool
	Default  *bool
	Children []*Section
	Path     string
}

// Offer values a branch section may declare.
const (
	OfferFoundation = "foundation"
	OfferChoose     = "choose"
	OfferOptional   = "optional"
	OfferOptIn      = "opt_in"
)

type libraryFile struct {
	Library  []libraryBlock  `hcl:"library,block"`
	Content  []contentBlock  `hcl:"content,block"`
	Sections []sectionsBlock `hcl:"sections,block"`
	Trees    []treesBlock    `hcl:"trees,block"`
}

type treesBlock struct {
	Trees []treeBlock `hcl:"tree,block"`
}

type treeBlock struct {
	Name    string           `hcl:"name,label"`
	Root    string           `hcl:"root"`
	Entries []treeEntryBlock `hcl:"entry,block"`
}

type treeEntryBlock struct {
	Name string   `hcl:"name,label"`
	Path string   `hcl:"path"`
	TLDR *string  `hcl:"tldr,optional"`
	Tags []string `hcl:"tags,optional"`
}

type libraryBlock struct {
	Format int    `hcl:"format"`
	ID     string `hcl:"id"`
	Name   string `hcl:"name"`
}

type contentBlock struct {
	MarkdownRoot string  `hcl:"markdown_root"`
	TreeRoot     *string `hcl:"tree_root,optional"`
}

type sectionsBlock struct {
	Sections []sectionBlock `hcl:"section,block"`
}

type sectionBlock struct {
	Name     string                 `hcl:"name,label"`
	Title    string                 `hcl:"title"`
	Source   *string                `hcl:"source,optional"`
	TLDR     *string                `hcl:"tldr,optional"`
	Tags     []string               `hcl:"tags,optional"`
	Curate   *string                `hcl:"curate,optional"`
	Offer    *string                `hcl:"offer,optional"`
	Defaults map[string]bool        `hcl:"defaults,optional"`
	Default  *bool                  `hcl:"default,optional"`
	Legacy   []legacyInclusionBlock `hcl:"inclusion,block"`
	Children []sectionBlock         `hcl:"section,block"`
}

// legacyInclusionBlock recognizes the removed inclusion block so the loader
// can name its replacement instead of reporting an unsupported block type.
type legacyInclusionBlock struct {
	Remain hcl.Body `hcl:",remain"`
}

// LoadLibrary reads and validates library.mogent.hcl from root.
func LoadLibrary(root string) (*Library, error) {
	path := filepath.Join(root, LibraryFile)
	var raw libraryFile
	if err := hclsimple.DecodeFile(path, nil, &raw); err != nil {
		return nil, err
	}
	if len(raw.Library) != 1 || raw.Library[0].Format != 2 {
		return nil, fmt.Errorf("%s: require exactly one library block with format = 2", path)
	}
	if len(raw.Content) != 1 {
		return nil, fmt.Errorf("%s: require exactly one content block", path)
	}
	if len(raw.Sections) != 1 {
		return nil, fmt.Errorf("%s: require exactly one sections block", path)
	}
	if raw.Content[0].MarkdownRoot == "" {
		return nil, fmt.Errorf("%s: markdown_root must not be empty", path)
	}
	if err := validateRelativePath(raw.Content[0].MarkdownRoot, "markdown_root"); err != nil {
		return nil, err
	}
	library := &Library{
		Path:         path,
		Root:         filepath.Clean(root),
		ID:           raw.Library[0].ID,
		Name:         raw.Library[0].Name,
		MarkdownRoot: raw.Content[0].MarkdownRoot,
		ByPath:       make(map[string]*Section),
		Trees:        make(map[string]Tree),
	}
	if raw.Content[0].TreeRoot != nil {
		library.TreeRoot = *raw.Content[0].TreeRoot
		if err := validateRelativePath(library.TreeRoot, "tree_root"); err != nil {
			return nil, err
		}
	}
	if library.ID == "" || library.Name == "" {
		return nil, fmt.Errorf("%s: library id and name must not be empty", path)
	}
	if len(raw.Sections[0].Sections) == 0 {
		return nil, fmt.Errorf("%s: sections block must declare at least one section", path)
	}
	for _, rawSection := range raw.Sections[0].Sections {
		section, err := loadSection(root, library, "", rawSection)
		if err != nil {
			return nil, err
		}
		library.Sections = append(library.Sections, section)
	}
	if len(raw.Trees) > 1 {
		return nil, fmt.Errorf("%s: require at most one trees block", path)
	}
	if len(raw.Trees) == 1 {
		for _, rawTree := range raw.Trees[0].Trees {
			if rawTree.Name == "" || rawTree.Root == "" {
				return nil, fmt.Errorf("%s: tree name and root must not be empty", path)
			}
			if err := validateRelativePath(rawTree.Root, "tree root"); err != nil {
				return nil, err
			}
			if _, duplicate := library.Trees[rawTree.Name]; duplicate {
				return nil, fmt.Errorf("duplicate tree %q", rawTree.Name)
			}
			info, err := os.Stat(filepath.Join(root, filepath.FromSlash(rawTree.Root)))
			if err != nil {
				return nil, fmt.Errorf("tree %q root %q: %w", rawTree.Name, rawTree.Root, err)
			}
			if !info.IsDir() {
				return nil, fmt.Errorf("tree %q root %q is not a directory", rawTree.Name, rawTree.Root)
			}
			tree := Tree{Name: rawTree.Name, Root: rawTree.Root}
			seen := make(map[string]bool, len(rawTree.Entries))
			for _, rawEntry := range rawTree.Entries {
				if rawEntry.Name == "" || seen[rawEntry.Name] {
					return nil, fmt.Errorf("tree %q has an empty or duplicate entry %q", rawTree.Name, rawEntry.Name)
				}
				seen[rawEntry.Name] = true
				if err := validateRelativePath(rawEntry.Path, "entry path"); err != nil {
					return nil, fmt.Errorf("tree %q entry %q: %w", rawTree.Name, rawEntry.Name, err)
				}
				entryInfo, err := os.Stat(filepath.Join(root, filepath.FromSlash(rawTree.Root), filepath.FromSlash(rawEntry.Path)))
				if err != nil {
					return nil, fmt.Errorf("tree %q entry %q path %q: %w", rawTree.Name, rawEntry.Name, rawEntry.Path, err)
				}
				if !entryInfo.IsDir() {
					return nil, fmt.Errorf("tree %q entry %q path %q is not a directory", rawTree.Name, rawEntry.Name, rawEntry.Path)
				}
				for _, tag := range rawEntry.Tags {
					if err := validateTag(tag); err != nil {
						return nil, fmt.Errorf("tree %q entry %q: %w", rawTree.Name, rawEntry.Name, err)
					}
				}
				entry := TreeEntry{Name: rawEntry.Name, Path: rawEntry.Path, Tags: append([]string(nil), rawEntry.Tags...)}
				if rawEntry.TLDR != nil {
					entry.TLDR = *rawEntry.TLDR
				}
				tree.Entries = append(tree.Entries, entry)
			}
			library.Trees[rawTree.Name] = tree
		}
	}
	return library, nil
}

func loadSection(root string, library *Library, parent string, raw sectionBlock) (*Section, error) {
	if raw.Name == "" || raw.Title == "" {
		return nil, fmt.Errorf("section name and title must not be empty")
	}
	if err := validateSectionName(raw.Name); err != nil {
		return nil, err
	}
	path := raw.Name
	if parent != "" {
		path = parent + "/" + raw.Name
	}
	if _, duplicate := library.ByPath[path]; duplicate {
		return nil, fmt.Errorf("duplicate section path %q", path)
	}
	section := &Section{Name: raw.Name, Title: raw.Title, Tags: append([]string(nil), raw.Tags...), Path: path}
	if raw.Source != nil {
		section.Source = *raw.Source
		if err := validateRelativePath(section.Source, "section source"); err != nil {
			return nil, err
		}
		file := filepath.Join(root, filepath.FromSlash(library.MarkdownRoot), filepath.FromSlash(section.Source))
		info, err := os.Stat(file)
		if err != nil {
			return nil, fmt.Errorf("section %q source %q: %w", path, section.Source, err)
		}
		if info.IsDir() {
			return nil, fmt.Errorf("section %q source %q is a directory", path, section.Source)
		}
	}
	if raw.TLDR != nil {
		section.TLDR = *raw.TLDR
	}
	for _, tag := range section.Tags {
		if err := validateTag(tag); err != nil {
			return nil, fmt.Errorf("section %q: %w", path, err)
		}
	}
	if len(raw.Legacy) > 0 {
		return nil, fmt.Errorf("section %q uses the removed inclusion block; set curate = \"foundation|choose|optional|opt_in\" on the section, with defaults or default beside it", path)
	}
	if err := validateOffer(path, raw, section); err != nil {
		return nil, err
	}
	if section.Source != "" && len(raw.Children) > 0 {
		return nil, fmt.Errorf("section %q cannot declare both source and child sections", path)
	}
	if section.Source == "" && len(raw.Children) == 0 {
		return nil, fmt.Errorf("section %q must declare source or child sections", path)
	}
	library.ByPath[path] = section
	for _, rawChild := range raw.Children {
		child, err := loadSection(root, library, path, rawChild)
		if err != nil {
			return nil, err
		}
		section.Children = append(section.Children, child)
	}
	return section, nil
}

func validateOffer(path string, raw sectionBlock, section *Section) error {
	if raw.Offer != nil {
		return fmt.Errorf("section %q uses removed offer; rename it to curate", path)
	}
	if raw.Curate == nil {
		if len(raw.Defaults) != 0 || raw.Default != nil {
			return fmt.Errorf("section %q declares defaults or default without curate", path)
		}
		return nil
	}
	offer := *raw.Curate
	switch offer {
	case OfferFoundation, OfferChoose, OfferOptional, OfferOptIn:
	default:
		return fmt.Errorf("section %q has unknown curate value %q; use foundation, choose, optional, or opt_in", path, offer)
	}
	if len(raw.Children) == 0 {
		return fmt.Errorf("section %q curate requires child sections", path)
	}
	section.Offer = offer
	section.Defaults = raw.Defaults
	section.Default = raw.Default
	if offer == OfferChoose {
		if raw.Default != nil {
			return fmt.Errorf("section %q choose curate takes defaults, not default", path)
		}
		if len(raw.Defaults) != len(raw.Children) {
			return fmt.Errorf("section %q choose curate requires defaults for every direct child", path)
		}
		for _, child := range raw.Children {
			if _, found := raw.Defaults[child.Name]; !found {
				return fmt.Errorf("section %q choose curate has no default for child %q", path, child.Name)
			}
		}
		return nil
	}
	if len(raw.Defaults) != 0 {
		return fmt.Errorf("section %q only a choose curate may declare defaults", path)
	}
	if offer == OfferOptional && raw.Default == nil {
		return fmt.Errorf("section %q optional curate requires default", path)
	}
	if offer != OfferOptional && raw.Default != nil {
		return fmt.Errorf("section %q only an optional curate may declare default", path)
	}
	return nil
}

func validateSectionName(value string) error {
	if value == "" || value == "." || value == ".." || value == "all" || value == "else" || value == "exclude" || value == "accept_defaults" || value == "force_exclude" || value == "reason" || filepath.Base(value) != value {
		return fmt.Errorf("section name %q is invalid or reserved", value)
	}
	return nil
}

func validateTag(value string) error {
	if err := validateRelativePath(value, "tag"); err != nil || value == "." {
		return fmt.Errorf("tag %q must be a slash-separated name", value)
	}
	return nil
}

// EffectiveTags returns ancestor and direct tags for a library section.
func (l *Library) EffectiveTags(section *Section) []string {
	seen := make(map[string]bool)
	var collect func(*Section)
	collect = func(current *Section) {
		if slash := filepath.ToSlash(filepath.Dir(current.Path)); slash != "." {
			collect(l.ByPath[slash])
		}
		for _, tag := range current.Tags {
			seen[tag] = true
		}
	}
	collect(section)
	tags := make([]string, 0, len(seen))
	for tag := range seen {
		tags = append(tags, tag)
	}
	sort.Strings(tags)
	return tags
}
