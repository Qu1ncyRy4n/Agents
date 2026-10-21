package assembly

import (
	"strings"

	"github.com/Qu1ncyRy4n/Agents/syntax"
)

// DirectIncludeOverlap identifies a pair of direct includes where one path is
// an ancestor of the other. It is advisory: repeated content and ordering can
// be intentional, so callers must never rewrite either declaration implicitly.
type DirectIncludeOverlap struct {
	Ancestor   string
	Descendant string
}

func DirectIncludeOverlaps(file *syntax.File) []DirectIncludeOverlap {
	type include struct{ alias, path, reference string }
	var includes []include
	for _, entry := range file.Entries[1:] {
		if entry.Name != "include" || entry.Label == nil {
			continue
		}
		reference, err := entry.Label.ResolveString(nil)
		if err != nil {
			continue
		}
		alias, node, found := strings.Cut(reference, ":")
		if !found || alias == "" || node == "" {
			continue
		}
		includes = append(includes, include{alias, node, reference})
	}
	var overlaps []DirectIncludeOverlap
	for earlier := 0; earlier < len(includes); earlier++ {
		for later := earlier + 1; later < len(includes); later++ {
			first, second := includes[earlier], includes[later]
			if first.alias != second.alias {
				continue
			}
			if strings.HasPrefix(second.path, first.path+"/") {
				overlaps = append(overlaps, DirectIncludeOverlap{Ancestor: first.reference, Descendant: second.reference})
			}
			if strings.HasPrefix(first.path, second.path+"/") {
				overlaps = append(overlaps, DirectIncludeOverlap{Ancestor: second.reference, Descendant: first.reference})
			}
		}
	}
	return overlaps
}
