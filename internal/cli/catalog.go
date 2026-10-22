package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
)

func writeCatalog(stdout io.Writer, sources map[string]assembly.Source, requested string, flat, align bool, charSet string, selected map[string]map[string]bool, onlyIncluded bool) error {
	aliases := make([]string, 0, len(sources))
	for alias := range sources {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	requestedAlias, requestedPath, _ := strings.Cut(requested, ":")
	if requested != "" {
		source, found := sources[requestedAlias]
		if !found {
			return fmt.Errorf("unknown source alias %q", requestedAlias)
		}
		if requestedPath != "" && source.Lib.ByPath[requestedPath] == nil {
			return fmt.Errorf("unknown source node %q", requested)
		}
	}
	chars := displayChars(charSet)
	for _, alias := range aliases {
		if requested != "" && alias != requestedAlias {
			continue
		}
		library := sources[alias].Lib
		if _, err := fmt.Fprintf(stdout, "%s: %s\n", alias, library.Name); err != nil {
			return err
		}
		rows := [][]string{{"PATH", "TYPE"}}
		if flat {
			paths := make([]string, 0, len(library.ByPath))
			for path := range library.ByPath {
				if requestedPath == "" || path == requestedPath || strings.HasPrefix(path, requestedPath+"/") {
					paths = append(paths, path)
				}
			}
			sort.Strings(paths)
			for _, path := range paths {
				if selected != nil && catalogIncluded(selected, alias, path) != onlyIncluded {
					continue
				}
				rows = append(rows, []string{path, string(library.ByPath[path].Type)})
			}
		} else {
			var visit func([]*assembly.Node, string)
			visit = func(nodes []*assembly.Node, prefix string) {
				for index, node := range nodes {
					connector, next := chars.branch, chars.vertical
					if index == len(nodes)-1 {
						connector, next = chars.last, chars.blank
					}
					rows = append(rows, []string{prefix + connector + node.Name, string(node.Type)})
					visit(node.Children, prefix+next)
				}
			}
			nodes := catalogRoots(library)
			if requestedPath != "" {
				nodes = []*assembly.Node{library.ByPath[requestedPath]}
			}
			if selected != nil {
				nodes = catalogSelectionFilter(nodes, selected[alias], onlyIncluded)
			}
			visit(nodes, "")
		}
		if len(rows) == 1 {
			if _, err := fmt.Fprintln(stdout, "  No content declared in this library."); err != nil {
				return err
			}
			continue
		}
		if err := writeColumns(stdout, rows, align); err != nil {
			return err
		}
	}
	return nil
}

func catalogIncluded(included map[string]map[string]bool, alias, path string) bool {
	return catalogPathIncluded(included[alias], path)
}

func catalogSelectionFilter(nodes []*assembly.Node, selected map[string]bool, onlyIncluded bool) []*assembly.Node {
	var result []*assembly.Node
	for _, node := range nodes {
		included := catalogPathIncluded(selected, node.Path)
		if included && !onlyIncluded {
			continue
		}
		if included && onlyIncluded {
			result = append(result, node)
			continue
		}
		copy := *node
		copy.Children = catalogSelectionFilter(node.Children, selected, onlyIncluded)
		if len(node.Children) != 0 && len(copy.Children) == 0 {
			continue
		}
		if onlyIncluded && len(node.Children) == 0 {
			continue
		}
		result = append(result, &copy)
	}
	return result
}

func catalogPathIncluded(included map[string]bool, path string) bool {
	for current := path; current != ""; {
		if included[current] {
			return true
		}
		index := strings.LastIndex(current, "/")
		if index < 0 {
			break
		}
		current = current[:index]
	}
	return false
}

// catalogRoots makes a physical directory the visible owner of a top-level
// logical branch only when every payload beneath that branch lives there. This
// keeps the catalog's declared logical structure while making the common
// directory-to-section relationship visible without guessing for mixed trees.
func catalogRoots(library *assembly.Library) []*assembly.Node {
	directories := make(map[string]*assembly.Node)
	for _, node := range library.Nodes {
		if node.Type == assembly.Directory && node.Source != "" {
			directories[node.Source] = node
		}
	}
	attached := make(map[*assembly.Node][]*assembly.Node)
	for _, node := range library.Nodes {
		if node.Type == assembly.Directory {
			continue
		}
		if parent := owningCatalogDirectory(node, directories); parent != nil {
			attached[parent] = append(attached[parent], node)
		}
	}
	var roots []*assembly.Node
	for _, node := range library.Nodes {
		if owner := owningCatalogDirectory(node, directories); owner != nil {
			continue
		}
		if children := attached[node]; len(children) != 0 {
			// The library model remains immutable for resolution; a display-only
			// copy prevents this catalog convenience from changing selection.
			copy := *node
			copy.Children = append(append([]*assembly.Node{}, node.Children...), children...)
			roots = append(roots, &copy)
			continue
		}
		roots = append(roots, node)
	}
	return roots
}

func owningCatalogDirectory(node *assembly.Node, directories map[string]*assembly.Node) *assembly.Node {
	if node.Type == assembly.Directory {
		return nil
	}
	sources := catalogSources(node)
	if len(sources) == 0 {
		return nil
	}
	var owner *assembly.Node
	for source, directory := range directories {
		prefix := source + "/"
		for _, payload := range sources {
			if !strings.HasPrefix(payload, prefix) {
				goto nextDirectory
			}
		}
		if owner == nil || len(source) > len(owner.Source) {
			owner = directory
		}
	nextDirectory:
	}
	return owner
}

func catalogSources(node *assembly.Node) []string {
	var sources []string
	if node.Source != "" {
		sources = append(sources, node.Source)
	}
	for _, child := range node.Children {
		sources = append(sources, catalogSources(child)...)
	}
	return sources
}
