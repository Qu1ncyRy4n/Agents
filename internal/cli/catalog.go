package cli

import (
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
)

func writeCatalog(stdout io.Writer, sources map[string]assembly.Source, requested string, flat, align bool, charSet string) error {
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
			nodes := library.Nodes
			if requestedPath != "" {
				nodes = []*assembly.Node{library.ByPath[requestedPath]}
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
