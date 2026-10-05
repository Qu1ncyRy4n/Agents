package v2

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/zclconf/go-cty/cty"
)

// Plan is a no-write resolution of user-side selections against server-side
// libraries. Errors are diagnostics so callers can present all useful fixes in
// one guided pass.
type Plan struct {
	Outputs     []PlannedOutput
	Diagnostics []Diagnostic
}

type PlannedOutput struct {
	Name    string
	Paths   []string
	Kind    string
	Sources []PlannedSource
}

type PlannedSource struct {
	Name     string
	From     string
	Sections []*Section
	TreeRoot string
}

type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

type Diagnostic struct {
	Severity   Severity
	Code       string
	Message    string
	Suggestion string
}

func (d Diagnostic) String() string {
	value := string(d.Severity) + "[" + d.Code + "]: " + d.Message
	if d.Suggestion != "" {
		return value + "\n\n" + d.Suggestion
	}
	return value
}

// Compile resolves every configured output. libraries maps user-side source
// aliases to their already validated server-side library trees.
func Compile(config *Config, libraries map[string]*Library) *Plan {
	plan := &Plan{}
	for _, output := range config.Outputs {
		plannedOutput := PlannedOutput{Name: output.Name, Paths: append([]string(nil), output.Paths...), Kind: output.Kind}
		for _, configuredSource := range output.Sources {
			library, found := libraries[configuredSource.Name]
			if !found {
				plan.Diagnostics = append(plan.Diagnostics, Diagnostic{
					Severity: SeverityError,
					Code:     "MOGENT101",
					Message:  fmt.Sprintf("output %q source %q is not loaded", output.Name, configuredSource.Name),
				})
				continue
			}
			if output.Kind == "tree" {
				planTreeOutput(plan, output.Name, configuredSource, library)
				plannedOutput.Sources = append(plannedOutput.Sources, PlannedSource{Name: configuredSource.Name, From: configuredSource.From, TreeRoot: strings.TrimPrefix(configuredSource.From, configuredSource.Name+":")})
				continue
			}
			if output.Kind != "markdown" {
				plan.Diagnostics = append(plan.Diagnostics, Diagnostic{
					Severity: SeverityError,
					Code:     "MOGENT103",
					Message:  fmt.Sprintf("tree output %q is not implemented by the v2 planner", output.Name),
				})
				continue
			}
			if configuredSource.From != configuredSource.Name+":"+library.MarkdownRoot {
				plan.Diagnostics = append(plan.Diagnostics, Diagnostic{
					Severity: SeverityError,
					Code:     "MOGENT104",
					Message:  fmt.Sprintf("source %q must select markdown root %q", configuredSource.Name, configuredSource.Name+":"+library.MarkdownRoot),
				})
				continue
			}
			selection, diagnostics := decodeSelection(configuredSource.Select)
			plan.Diagnostics = append(plan.Diagnostics, diagnostics...)
			if selection == nil {
				continue
			}
			resolver := selectionResolver{plan: plan, source: configuredSource.Name}
			resolver.unknownChildren("", selection, library.Sections)
			plannedSource := PlannedSource{Name: configuredSource.Name, From: configuredSource.From}
			for _, section := range library.Sections {
				resolver.resolve(section, selection.children[section.Name], selection, false, &plannedSource.Sections)
			}
			appendTagSelections(library, configuredSource, &plannedSource.Sections)
			if len(plannedSource.Sections) == 0 {
				plan.Diagnostics = append(plan.Diagnostics, Diagnostic{
					Severity: SeverityWarning,
					Code:     "MOGENT202",
					Message:  fmt.Sprintf("output %q source %q selects no content", output.Name, configuredSource.Name),
				})
			}
			plannedOutput.Sources = append(plannedOutput.Sources, plannedSource)
		}
		plan.Outputs = append(plan.Outputs, plannedOutput)
	}
	return plan
}

func planTreeOutput(plan *Plan, outputName string, source OutputSource, library *Library) {
	root := strings.TrimPrefix(source.From, source.Name+":")
	tree, found := library.Trees[root]
	if !found {
		plan.Diagnostics = append(plan.Diagnostics, Diagnostic{Severity: SeverityError, Code: "MOGENT104", Message: fmt.Sprintf("tree output %q source %q must select a declared server-side tree", outputName, source.Name)})
		return
	}
	selection, diagnostics := decodeSelection(source.Select)
	plan.Diagnostics = append(plan.Diagnostics, diagnostics...)
	if selection == nil || !selection.all || len(selection.children) != 0 || len(selection.exclude) != 0 || selection.otherwise != nil {
		plan.Diagnostics = append(plan.Diagnostics, Diagnostic{Severity: SeverityError, Code: "MOGENT105", Message: fmt.Sprintf("tree output %q source %q currently requires select = { all = true }", outputName, source.Name)})
		return
	}
	if tree.Root == "" {
		plan.Diagnostics = append(plan.Diagnostics, Diagnostic{Severity: SeverityError, Code: "MOGENT104", Message: fmt.Sprintf("tree output %q source %q has an empty root", outputName, source.Name)})
	}
}

func appendTagSelections(library *Library, source OutputSource, selected *[]*Section) {
	if len(source.TagsAll) == 0 && len(source.TagsAny) == 0 {
		return
	}
	seen := make(map[string]bool, len(*selected))
	for _, section := range *selected {
		seen[section.Path] = true
	}
	for _, section := range libraryLeaves(library.Sections) {
		if seen[section.Path] || !matchesTags(library.EffectiveTags(section), source.TagsAll, source.TagsAny) {
			continue
		}
		*selected = append(*selected, section)
		seen[section.Path] = true
	}
}

func libraryLeaves(sections []*Section) []*Section {
	var leaves []*Section
	for _, section := range sections {
		if len(section.Children) == 0 {
			leaves = append(leaves, section)
			continue
		}
		leaves = append(leaves, libraryLeaves(section.Children)...)
	}
	return leaves
}

func matchesTags(tags, all, any []string) bool {
	for _, query := range all {
		if !containsTag(tags, query) {
			return false
		}
	}
	if len(any) == 0 {
		return true
	}
	for _, query := range any {
		if containsTag(tags, query) {
			return true
		}
	}
	return false
}

func containsTag(tags []string, query string) bool {
	for _, tag := range tags {
		if tag == query || strings.HasPrefix(tag, query+"/") {
			return true
		}
	}
	return false
}

type selectionNode struct {
	include        *bool
	all            bool
	otherwise      *bool
	acceptDefaults bool
	forceExclude   bool
	reason         string
	children       map[string]*selectionNode
	exclude        map[string]*selectionNode
}

func decodeSelection(expression interface {
	Value(*hcl.EvalContext) (cty.Value, hcl.Diagnostics)
}) (*selectionNode, []Diagnostic) {
	value, diagnostics := expression.Value(nil)
	if diagnostics.HasErrors() {
		return nil, []Diagnostic{{Severity: SeverityError, Code: "MOGENT100", Message: diagnostics.Error()}}
	}
	node, diagnostic := decodeSelectionValue(value)
	if diagnostic != nil {
		return nil, []Diagnostic{*diagnostic}
	}
	return node, nil
}

func decodeSelectionValue(value cty.Value) (*selectionNode, *Diagnostic) {
	if !value.IsKnown() || value.IsNull() || !value.Type().IsObjectType() {
		return nil, &Diagnostic{Severity: SeverityError, Code: "MOGENT100", Message: "select must be a known object"}
	}
	node := &selectionNode{children: make(map[string]*selectionNode), exclude: make(map[string]*selectionNode)}
	for key, child := range value.AsValueMap() {
		switch key {
		case "all":
			boolean, diagnostic := selectionBool(key, child)
			if diagnostic != nil {
				return nil, diagnostic
			}
			node.all = boolean
		case "else":
			if child.Type() != cty.String {
				return nil, &Diagnostic{Severity: SeverityError, Code: "MOGENT100", Message: "select else must be include or exclude"}
			}
			switch child.AsString() {
			case "include":
				include := true
				node.otherwise = &include
			case "exclude":
				exclude := false
				node.otherwise = &exclude
			default:
				return nil, &Diagnostic{Severity: SeverityError, Code: "MOGENT100", Message: "select else must be include or exclude"}
			}
		case "accept_defaults":
			boolean, diagnostic := selectionBool(key, child)
			if diagnostic != nil {
				return nil, diagnostic
			}
			node.acceptDefaults = boolean
		case "force_exclude":
			boolean, diagnostic := selectionBool(key, child)
			if diagnostic != nil {
				return nil, diagnostic
			}
			node.forceExclude = boolean
		case "reason":
			if child.Type() != cty.String {
				return nil, &Diagnostic{Severity: SeverityError, Code: "MOGENT100", Message: "select reason must be a string"}
			}
			node.reason = child.AsString()
		case "exclude":
			excluded, diagnostic := decodeSelectionValue(child)
			if diagnostic != nil {
				return nil, diagnostic
			}
			node.exclude = excluded.children
		default:
			parsed, diagnostic := decodeNodeValue(child)
			if diagnostic != nil {
				return nil, diagnostic
			}
			node.children[key] = parsed
		}
	}
	return node, nil
}

func decodeNodeValue(value cty.Value) (*selectionNode, *Diagnostic) {
	if value.Type() == cty.Bool {
		include := value.True()
		return &selectionNode{include: &include, children: make(map[string]*selectionNode), exclude: make(map[string]*selectionNode)}, nil
	}
	return decodeSelectionValue(value)
}

func selectionBool(key string, value cty.Value) (bool, *Diagnostic) {
	if value.Type() != cty.Bool {
		return false, &Diagnostic{Severity: SeverityError, Code: "MOGENT100", Message: fmt.Sprintf("select %s must be a boolean", key)}
	}
	return value.True(), nil
}

type selectionResolver struct {
	plan   *Plan
	source string
}

func (r selectionResolver) resolve(section *Section, selection *selectionNode, parent *selectionNode, inherited bool, selected *[]*Section) {
	if selection == nil {
		selection = &selectionNode{children: make(map[string]*selectionNode), exclude: make(map[string]*selectionNode)}
	}
	if excluded, found := parent.exclude[section.Name]; found {
		selection = excluded
		value := false
		selection.include = &value
	}

	include := inherited
	if parent.otherwise != nil {
		include = *parent.otherwise
	}
	if parent.all {
		include = true
	}
	if selection.include != nil {
		include = *selection.include
	}
	if len(selection.children) > 0 || selection.all || selection.acceptDefaults {
		include = true
	}
	if !include {
		return
	}
	if selection.acceptDefaults && section.Offer != OfferChoose {
		r.plan.Diagnostics = append(r.plan.Diagnostics, Diagnostic{
			Severity: SeverityError,
			Code:     "MOGENT108",
			Message:  fmt.Sprintf("accept_defaults is valid only under a choose offer; %s:%s offers %s", r.source, section.Path, offerName(section)),
		})
		return
	}
	if len(section.Children) == 0 {
		*selected = append(*selected, section)
		return
	}
	r.unknownChildren(section.Path, selection, section.Children)
	if section.Offer == OfferFoundation && selection.otherwise != nil && !*selection.otherwise {
		r.plan.Diagnostics = append(r.plan.Diagnostics, Diagnostic{
			Severity: SeverityError,
			Code:     "MOGENT107",
			Message:  fmt.Sprintf("else = \"exclude\" cannot drop foundation children of %s:%s; use force_exclude = true with a reason per child", r.source, section.Path),
		})
		return
	}

	broad := inherited || parent.all || selection.all || selection.include != nil && *selection.include || selection.acceptDefaults
	if section.Offer == OfferChoose && broad && !selection.acceptDefaults {
		missing := make([]string, 0)
		for _, child := range section.Children {
			if _, found := selection.children[child.Name]; !found {
				missing = append(missing, child.Name)
			}
		}
		if len(missing) > 0 {
			sort.Strings(missing)
			r.plan.Diagnostics = append(r.plan.Diagnostics, Diagnostic{
				Severity:   SeverityError,
				Code:       "MOGENT102",
				Message:    fmt.Sprintf("incomplete selection for %s:%s", r.source, section.Path),
				Suggestion: suggestedChoices(section, missing),
			})
			return
		}
	}

	for _, child := range section.Children {
		childSelection := selection.children[child.Name]
		childInherited := false
		if childSelection != nil && !r.allowChildSelection(section, child, childSelection) {
			continue
		}
		if selection.acceptDefaults && section.Offer == OfferChoose {
			defaultValue := section.Defaults[child.Name]
			childSelection = &selectionNode{include: &defaultValue, children: make(map[string]*selectionNode), exclude: make(map[string]*selectionNode)}
			r.plan.Diagnostics = append(r.plan.Diagnostics, Diagnostic{Severity: SeverityWarning, Code: "MOGENT203", Message: fmt.Sprintf("accepted default %t for %s:%s", defaultValue, r.source, child.Path)})
		} else if childSelection == nil {
			childInherited = broad
			switch section.Offer {
			case OfferFoundation:
				childInherited = true
			case OfferOptional:
				childInherited = *section.Default
			case OfferOptIn, OfferChoose:
				childInherited = false
			}
		}
		r.resolve(child, childSelection, selection, childInherited, selected)
	}
}

// allowChildSelection validates a consumer's explicit entry for one child
// against the parent's offer. It returns false when the child must be skipped,
// after recording the diagnostic that explains why.
func (r selectionResolver) allowChildSelection(section, child *Section, selection *selectionNode) bool {
	location := r.source + ":" + child.Path
	if selection.forceExclude {
		if section.Offer != OfferFoundation {
			r.plan.Diagnostics = append(r.plan.Diagnostics, Diagnostic{
				Severity: SeverityError,
				Code:     "MOGENT107",
				Message:  fmt.Sprintf("force_exclude is valid only for a foundation child; %s is offered as %s, use %s = false", location, offerName(section), child.Name),
			})
			return false
		}
		if strings.TrimSpace(selection.reason) == "" {
			r.plan.Diagnostics = append(r.plan.Diagnostics, Diagnostic{
				Severity: SeverityError,
				Code:     "MOGENT107",
				Message:  fmt.Sprintf("force_exclude of foundation section %s requires a non-empty reason", location),
			})
			return false
		}
		r.plan.Diagnostics = append(r.plan.Diagnostics, Diagnostic{
			Severity: SeverityWarning,
			Code:     "MOGENT205",
			Message:  fmt.Sprintf("foundation section %s force-excluded: %s", location, strings.TrimSpace(selection.reason)),
		})
		return false
	}
	if selection.reason != "" {
		r.plan.Diagnostics = append(r.plan.Diagnostics, Diagnostic{
			Severity: SeverityError,
			Code:     "MOGENT107",
			Message:  fmt.Sprintf("reason is valid only with force_exclude = true at %s", location),
		})
		return false
	}
	if section.Offer == OfferFoundation && selection.include != nil && !*selection.include {
		r.plan.Diagnostics = append(r.plan.Diagnostics, Diagnostic{
			Severity: SeverityError,
			Code:     "MOGENT107",
			Message:  fmt.Sprintf("dropping foundation section %s requires force_exclude = true and a reason", location),
		})
		return false
	}
	if section.Offer == OfferOptIn && (selection.all || len(selection.children) > 0 || selection.acceptDefaults || selection.include != nil && *selection.include) {
		r.plan.Diagnostics = append(r.plan.Diagnostics, Diagnostic{
			Severity: SeverityWarning,
			Code:     "MOGENT204",
			Message:  fmt.Sprintf("selected opt-in section %s", location),
		})
	}
	return true
}

func offerName(section *Section) string {
	if section.Offer == "" {
		return "nothing"
	}
	return section.Offer
}

func (r selectionResolver) unknownChildren(parentPath string, selection *selectionNode, children []*Section) {
	known := make(map[string]bool, len(children))
	for _, child := range children {
		known[child.Name] = true
	}
	for name := range selection.children {
		if known[name] {
			continue
		}
		location := r.source
		if parentPath != "" {
			location += ":" + parentPath
		}
		r.plan.Diagnostics = append(r.plan.Diagnostics, Diagnostic{
			Severity:   SeverityError,
			Code:       "MOGENT106",
			Message:    fmt.Sprintf("selection %q has no child %q", location, name),
			Suggestion: "Available children: " + childNames(children) + ". Correct the selection path or spelling.",
		})
	}
}

func childNames(children []*Section) string {
	names := make([]string, 0, len(children))
	for _, child := range children {
		names = append(names, child.Name)
	}
	sort.Strings(names)
	return strings.Join(names, ", ")
}

func suggestedChoices(section *Section, missing []string) string {
	var output strings.Builder
	output.WriteString("Add decisions for every direct child:\n\n")
	parts := strings.Split(section.Path, "/")
	for index, part := range parts {
		output.WriteString(strings.Repeat("  ", index+1))
		output.WriteString(part)
		output.WriteString(" = {\n")
	}
	for _, name := range missing {
		output.WriteString(strings.Repeat("  ", len(parts)+1))
		if strings.Contains(name, "-") {
			output.WriteString("\"")
			output.WriteString(name)
			output.WriteString("\"")
		} else {
			output.WriteString(name)
		}
		output.WriteString(" = ")
		output.WriteString(fmt.Sprintf("%t", section.Defaults[name]))
		output.WriteString("\n")
	}
	for index := len(parts); index > 0; index-- {
		output.WriteString(strings.Repeat("  ", index))
		output.WriteString("}\n")
	}
	output.WriteString(fmt.Sprintf("\nOr use accept_defaults = true inside the %s block.", section.Name))
	return output.String()
}
