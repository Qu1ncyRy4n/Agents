package syntax

import "fmt"

func validate(f *File) error {
	if f == nil || len(f.Entries) == 0 {
		return &Diagnostic{Range: Range{Start: Position{Line: 1, Column: 1}}, Message: "require a library, plan, or build header"}
	}
	header := f.Entries[0]
	if header == nil {
		return &Diagnostic{Message: "require a library, plan, or build header"}
	}
	kind := FileKind(header.Name)
	if header.Kind != Block || kind != Library && kind != Plan && kind != Build {
		return located(header, "first declaration must be a library, plan, or build header")
	}
	if err := validateBody(header); err != nil {
		return err
	}
	version := header.field("format")
	if version == nil || version.Value.Kind != Integer || version.Value.Int != 1 {
		return located(header, "header requires format = 1")
	}
	if kind == Build {
		origin := header.field("origin")
		if origin == nil || origin.Value.Kind != String {
			return located(header, "build header requires origin = \"authored\" or \"generated\"")
		}
		value, err := origin.Value.ResolveString(nil)
		if err != nil {
			return err
		}
		if value != "authored" && value != "generated" {
			return located(origin, "build origin must be authored or generated")
		}
	}
	seen := make(map[string]bool)
	for _, e := range f.Entries[1:] {
		if e == nil {
			return located(header, "nil declaration")
		}
		if e.Kind == Binding {
			if kind == Library {
				return located(e, "library declarations do not use consumer bindings")
			}
			if e.Value == nil || e.Value.Kind != String {
				return located(e, "let bindings must be strings")
			}
			if seen["binding:"+e.Name] {
				return located(e, fmt.Sprintf("duplicate binding %q", e.Name))
			}
			seen["binding:"+e.Name] = true
			if err := validateValue(*e.Value); err != nil {
				return err
			}
			continue
		}
		if e.Kind != Block || !allowedRoot(kind, e.Name) {
			return located(e, fmt.Sprintf("unexpected %q in %s file", e.Name, kind))
		}
		if e.Name == "resolved" && kind != Build {
			return located(e, "resolved records belong in a build file")
		}
		if err := validateBody(e); err != nil {
			return err
		}
		if e.Name == "source" || isNode(e.Name) || e.Name == "preset" || e.Name == "group" {
			label, err := staticLabel(e)
			if err != nil {
				return err
			}
			namespace := "library:"
			if e.Name == "source" || e.Name == "group" {
				namespace = e.Name + ":"
			}
			if seen[namespace+label] {
				return located(e, fmt.Sprintf("duplicate %sname %q", namespace, label))
			}
			seen[namespace+label] = true
		}
	}
	return nil
}

func allowedRoot(kind FileKind, name string) bool {
	if kind == Library {
		return isNode(name) || name == "preset" || name == "group"
	}
	return name == "source" || name == "include" || name == "replace" || name == "include-build" || name == "resolved"
}

func isNode(name string) bool {
	return name == "node" || name == "dir" || name == "file" || name == "md-file" || name == "md-section"
}

func allowedField(parent, name string) bool {
	var names []string
	switch parent {
	case "plan":
		names = []string{"format", "author"}
	case "build":
		names = []string{"format", "origin", "author", "generated_from", "generated_by", "plan_sha256", "content_sha256"}
	case "library":
		names = []string{"format", "name"}
	case "source":
		names = []string{"git", "local", "ref", "commit", "subdir"}
	case "include", "replace":
		names = []string{"into", "exclude", "tags_all", "tags_any", "accept_defaults", "replace", "reason", "text"}
	case "exclude":
		names = []string{"reason"}
	case "group":
		names = []string{"mode"}
	case "values":
		return true // Names are data, checked against library vars by the resolver.
	case "input", "decision", "owner":
		return true // Resolved-record field schema belongs to the sequencing layer.
	default:
		if isNode(parent) {
			names = []string{"type", "source", "title", "tldr", "tags", "role", "content", "curate", "default", "reason_required", "requires", "conflicts_with", "exclusive_group", "into", "vars", "heading"}
		}
	}
	for _, candidate := range names {
		if candidate == name {
			return true
		}
	}
	return false
}

func allowedChild(parent, name string) bool {
	if isNode(parent) {
		return isNode(name)
	}
	switch parent {
	case "include", "replace":
		return name == "exclude" || name == "values"
	case "preset":
		return name == "include" || name == "replace"
	case "resolved":
		return name == "input" || name == "decision" || name == "owner"
	}
	return false
}

func validateBody(e *Entry) error {
	if e.Label != nil {
		if e.Label.Kind != String {
			return located(e, "declaration name must be a string")
		}
		if err := validateValue(*e.Label); err != nil {
			return err
		}
	}
	if e.Selection != nil {
		if e.Selection.Kind != List {
			return located(e, "child selection must be a string array")
		}
		if err := validateValue(*e.Selection); err != nil {
			return err
		}
	}
	if e.Selection != nil && e.Name != "include" && e.Name != "replace" {
		return located(e, "child selections are only valid on include or replace")
	}
	if e.Name == "include" || e.Name == "replace" {
		if e.Name == "replace" && e.Label == nil {
			return located(e, "replace requires a source reference")
		}
		if e.Label == nil && e.field("text") == nil {
			return located(e, "include requires a source reference or inline text")
		}
		if e.Label != nil && e.field("text") != nil {
			return located(e, "include cannot combine a source reference and inline text")
		}
	} else if e.Name == "replace" || e.Name == "include-build" || e.Name == "exclude" || e.Name == "source" || e.Name == "preset" || e.Name == "group" || isNode(e.Name) || e.Name == "input" || e.Name == "decision" || e.Name == "owner" {
		if e.Label == nil {
			return located(e, fmt.Sprintf("%s requires a quoted name", e.Name))
		}
	} else if e.Label != nil {
		return located(e, fmt.Sprintf("%s does not take a name", e.Name))
	}
	if e.Name == "include-build" {
		if e.HasBody {
			return located(e, "include-build does not take an options block")
		}
		return nil
	}
	if e.Name != "include" && e.Name != "replace" && !e.HasBody {
		return located(e, fmt.Sprintf("%s requires a block", e.Name))
	}
	fields, nodes := make(map[string]bool), make(map[string]bool)
	valuesSeen := false
	for _, child := range e.Body {
		if child == nil {
			return located(e, "nil declaration")
		}
		if child.Kind == Field {
			if !allowedField(e.Name, child.Name) {
				return located(child, fmt.Sprintf("unknown field %q in %s", child.Name, e.Name))
			}
			if fields[child.Name] {
				return located(child, fmt.Sprintf("duplicate field %q", child.Name))
			}
			fields[child.Name] = true
			if child.Value == nil {
				return located(child, "field requires a value")
			}
			if err := validateValue(*child.Value); err != nil {
				return err
			}
			if e.Name == "values" && child.Value.Kind != String {
				return located(child, "guidance values must be strings")
			}
			continue
		}
		if child.Kind != Block || !allowedChild(e.Name, child.Name) {
			return located(child, fmt.Sprintf("unexpected %q inside %s", child.Name, e.Name))
		}
		if child.Name == "values" {
			if valuesSeen {
				return located(child, "include accepts only one values block")
			}
			valuesSeen = true
		}
		if err := validateBody(child); err != nil {
			return err
		}
		if isNode(child.Name) {
			label, err := staticLabel(child)
			if err != nil {
				return err
			}
			if nodes[label] {
				return located(child, fmt.Sprintf("duplicate sibling node %q", label))
			}
			nodes[label] = true
		}
	}
	if e.Name == "replace" {
		if flag := e.field("replace"); flag != nil && (flag.Value.Kind != Boolean || !flag.Value.Bool) {
			return located(flag, "replace alias cannot disable replacement")
		}
	}
	if isNode(e.Name) && e.Name != "node" {
		if kind := e.field("type"); kind != nil {
			value, err := kind.Value.ResolveString(nil)
			if err != nil {
				return err
			}
			if value != e.Name {
				return located(kind, "typed node alias contradicts explicit type")
			}
		}
	}
	return nil
}

func validateValue(v Value) error {
	switch v.Kind {
	case String:
		for _, part := range v.Parts {
			if part.Reference != "" && (!validIdentifier(part.Reference) || part.Text != "") {
				return &Diagnostic{Range: v.Range, Message: "string part must be literal text or a named reference"}
			}
		}
	case Boolean, Integer:
	case List:
		for _, item := range v.Items {
			if item.Kind != String {
				return &Diagnostic{Range: item.Range, Message: "array items must be strings"}
			}
			if err := validateValue(item); err != nil {
				return err
			}
		}
	default:
		return &Diagnostic{Range: v.Range, Message: fmt.Sprintf("unsupported value kind %q", v.Kind)}
	}
	return nil
}
