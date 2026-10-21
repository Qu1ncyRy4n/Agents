package assembly

import (
	"fmt"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

type ExceptResult struct {
	Path      string
	Reference string
	From      string
	Relative  string
	Content   []byte
}

// Except adds a descendant exclusion to the narrowest direct selected ancestor.
// It never removes another include or installs outputs; that remains a separate
// decision because a direct include may deliberately overlap its ancestor.
func Except(path, reference string, dryRun bool) (*ExceptResult, error) {
	file, err := LoadFile(path)
	if err != nil {
		return nil, err
	}
	if file.Kind == syntax.Build {
		origin, _, err := stringField(file.Entries[0], "origin", nil)
		if err != nil {
			return nil, err
		}
		if origin != "authored" {
			return nil, fmt.Errorf("cannot edit generated %s; add the exception to its plan", path)
		}
	}
	alias, wanted, found := strings.Cut(reference, ":")
	if !found || alias == "" || wanted == "" || strings.Contains(wanted, ":") {
		return nil, fmt.Errorf("exception reference must use source:path, for example qmr:workflow/decision-first")
	}
	_, sources, err := ResolvePlan(file, path)
	if err != nil {
		return nil, err
	}
	source, found := sources[alias]
	if !found || source.Lib.ByPath[wanted] == nil {
		return nil, fmt.Errorf("unknown module %q", reference)
	}
	var target *syntax.Entry
	ancestor := ""
	for _, entry := range file.Entries[1:] {
		if entry.Name != "include" || entry.Label == nil {
			continue
		}
		label, err := entry.Label.ResolveString(nil)
		if err != nil {
			continue
		}
		candidateAlias, candidatePath, valid := strings.Cut(label, ":")
		if !valid || candidateAlias != alias || candidatePath == "" || candidatePath == wanted || !strings.HasPrefix(wanted, candidatePath+"/") {
			continue
		}
		if len(candidatePath) > len(ancestor) {
			target, ancestor = entry, candidatePath
		}
	}
	if target == nil {
		return nil, fmt.Errorf("%q has no selected ancestor; select %s first, then except its descendant", reference, alias+":"+strings.TrimSuffix(wanted, "/"+strings.Split(wanted, "/")[len(strings.Split(wanted, "/"))-1]))
	}
	relative := strings.TrimPrefix(wanted, ancestor+"/")
	ensureExceptionDefaults(target, source.Lib.ByPath[ancestor])
	for _, existing := range target.Body {
		if existing.Kind != syntax.Field || existing.Name != "exclude" {
			continue
		}
		values, _, err := fieldStrings(target, "exclude")
		if err != nil {
			return nil, err
		}
		for _, value := range values {
			if value == relative || value == wanted {
				return nil, fmt.Errorf("%q is already excepted from %s", reference, labelForEntry(target))
			}
		}
		existing.Value.Items = append(existing.Value.Items, listString(relative))
		return writeExcept(file, path, reference, labelForEntry(target), relative, dryRun)
	}
	target.HasBody = true
	target.Body = append(target.Body, &syntax.Entry{Kind: syntax.Field, Name: "exclude", Value: &syntax.Value{Kind: syntax.List, Items: []syntax.Value{listString(relative)}}})
	return writeExcept(file, path, reference, labelForEntry(target), relative, dryRun)
}

// A choose branch has no implicit content until defaults are accepted. An
// exception against that branch means "use its defaults except this member";
// make that intent explicit so the resulting entry remains buildable.
func ensureExceptionDefaults(entry *syntax.Entry, node *Node) {
	if node == nil || node.Curate != Choose {
		return
	}
	for _, child := range entry.Body {
		if child.Kind == syntax.Field && child.Name == "accept_defaults" {
			child.Value = &syntax.Value{Kind: syntax.Boolean, Bool: true}
			return
		}
	}
	entry.HasBody = true
	entry.Body = append(entry.Body, fieldBoolEntry("accept_defaults", true))
}

func listString(value string) syntax.Value {
	return syntax.Value{Kind: syntax.String, Parts: []syntax.StringPart{{Text: value}}}
}

func labelForEntry(entry *syntax.Entry) string {
	if entry.Label == nil {
		return "include"
	}
	value, _ := entry.Label.ResolveString(nil)
	return value
}

func writeExcept(file *syntax.File, path, reference, from, relative string, dryRun bool) (*ExceptResult, error) {
	contents, err := syntax.Format(file, syntax.FormatOptions{})
	if err != nil {
		return nil, err
	}
	updated, err := syntax.Parse(path, contents)
	if err != nil {
		return nil, err
	}
	result, _, err := ResolvePlan(updated, path)
	if err != nil {
		return nil, err
	}
	if result.HasErrors() {
		return nil, fmt.Errorf("exception %q produces unresolved selection errors; inspect with mogent show --issues", reference)
	}
	value := &ExceptResult{Path: path, Reference: reference, From: from, Relative: relative, Content: contents}
	if !dryRun {
		if err := renderfs.WriteAtomically(path, contents); err != nil {
			return nil, err
		}
	}
	return value, nil
}
