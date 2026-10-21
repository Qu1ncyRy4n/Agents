package assembly

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/syntax"
)

// ResolvePlan loads local source libraries and turns plan selections into ordered
// contributions. It intentionally returns partial diagnostics rather than failing
// at the first unresolved choice.
func ResolvePlan(file *syntax.File, filename string) (*Result, map[string]Source, error) {
	if file.Kind != syntax.Plan && file.Kind != syntax.Build {
		return nil, nil, fmt.Errorf("resolve requires a plan or build file")
	}
	root := filepath.Dir(filename)
	bindings := make(map[string]string)
	sources := make(map[string]Source)
	for _, entry := range file.Entries[1:] {
		switch entry.Kind {
		case syntax.Binding:
			value, err := entry.Value.ResolveString(bindings)
			if err != nil {
				return nil, nil, err
			}
			bindings[entry.Name] = value
		case syntax.Block:
			if entry.Name != "source" {
				continue
			}
			alias, err := entry.Label.ResolveString(bindings)
			if err != nil {
				return nil, nil, err
			}
			local, localFound, err := stringField(entry, "local", bindings)
			if err != nil {
				return nil, nil, err
			}
			git, gitFound, err := stringField(entry, "git", bindings)
			if err != nil {
				return nil, nil, err
			}
			if localFound == gitFound {
				return nil, nil, fmt.Errorf("source %q requires exactly one of local or git", alias)
			}
			ref, _, err := stringField(entry, "ref", bindings)
			if err != nil {
				return nil, nil, err
			}
			commit, _, err := stringField(entry, "commit", bindings)
			if err != nil {
				return nil, nil, err
			}
			var libraryRoot string
			if localFound {
				if filepath.IsAbs(local) {
					return nil, nil, fmt.Errorf("source %q local path must be relative", alias)
				}
				libraryRoot = filepath.Join(root, filepath.FromSlash(local))
			} else {
				libraryRoot, commit, err = resolveGitSource(root, alias, git, ref, commit)
				if err != nil {
					return nil, nil, fmt.Errorf("source %q: %w", alias, err)
				}
			}
			library, err := LoadLibrary(libraryRoot)
			if err != nil {
				return nil, nil, fmt.Errorf("source %q: %w", alias, err)
			}
			sources[alias] = Source{Alias: alias, Root: libraryRoot, Local: local, Git: git, Ref: ref, Commit: commit, Lib: library}
		}
	}
	r := &Result{}
	for _, entry := range file.Entries[1:] {
		if entry.Kind != syntax.Block || entry.Name != "include" && entry.Name != "replace" {
			continue
		}
		resolveInclude(r, sources, bindings, entry, false, nil, nil)
	}
	if file.Kind == syntax.Build {
		validateRecordedInputs(r, file)
	}
	validateRelationships(r, sources)
	return r, sources, nil
}

// validateRelationships checks the selected identities after all includes have
// resolved. A replacement does not satisfy a requirement for another identity.
func validateRelationships(result *Result, sources map[string]Source) {
	selected := make(map[string]bool)
	for _, contribution := range result.Contributions {
		if contribution.Node != nil {
			selected[contribution.Source+":"+contribution.Node.Path] = true
		}
	}
	groups := make(map[string][]string)
	for _, contribution := range result.Contributions {
		if contribution.Node == nil {
			continue
		}
		node, source := contribution.Node, contribution.Source
		for _, required := range node.Requires {
			if !selected[source+":"+required] {
				result.Diagnostics = append(result.Diagnostics, diagnostic(contribution.RangeEntry(), Error, fmt.Sprintf("%s:%s requires %s:%s", source, node.Path, source, required)))
			}
		}
		for _, conflict := range node.ConflictsWith {
			if selected[source+":"+conflict] {
				result.Diagnostics = append(result.Diagnostics, diagnostic(contribution.RangeEntry(), Error, fmt.Sprintf("%s:%s conflicts with %s:%s", source, node.Path, source, conflict)))
			}
		}
		if node.ExclusiveGroup != "" {
			library := sources[source].Lib
			if _, found := library.Groups[node.ExclusiveGroup]; !found {
				result.Diagnostics = append(result.Diagnostics, diagnostic(contribution.RangeEntry(), Error, fmt.Sprintf("%s:%s names unknown exclusive group %q", source, node.Path, node.ExclusiveGroup)))
			} else {
				groups[source+":"+node.ExclusiveGroup] = append(groups[source+":"+node.ExclusiveGroup], node.Path)
			}
		}
	}
	for group, members := range groups {
		if len(members) > 1 {
			result.Diagnostics = append(result.Diagnostics, Diagnostic{Severity: Error, Path: group, Message: fmt.Sprintf("exclusive group selects multiple nodes: %s", strings.Join(members, ", "))})
		}
	}
}

func resolveGitSource(workspace, alias, url, ref, commit string) (string, string, error) {
	if ref == "" {
		ref = "HEAD"
	}
	if commit == "" {
		output, err := exec.Command("git", "ls-remote", url, ref).Output()
		if err != nil {
			return "", "", fmt.Errorf("resolve git ref %q: %w", ref, err)
		}
		fields := strings.Fields(string(output))
		if len(fields) == 0 || len(fields[0]) != 40 {
			return "", "", fmt.Errorf("git ref %q did not resolve to a full commit", ref)
		}
		commit = fields[0]
	}
	for _, r := range commit {
		if !strings.ContainsRune("0123456789abcdef", r) {
			return "", "", fmt.Errorf("commit %q must be lowercase full SHA-1", commit)
		}
	}
	if len(commit) != 40 {
		return "", "", fmt.Errorf("commit %q must be a full SHA-1", commit)
	}
	root := filepath.Join(workspace, ".mogent", "sources", alias, commit)
	if _, err := os.Stat(root); err == nil {
		return root, commit, nil
	} else if !os.IsNotExist(err) {
		return "", "", err
	}
	if err := os.MkdirAll(filepath.Dir(root), 0755); err != nil {
		return "", "", err
	}
	temporary, err := os.MkdirTemp(filepath.Dir(root), ".git-source-")
	if err != nil {
		return "", "", err
	}
	defer func() { _ = os.RemoveAll(temporary) }()
	if output, err := exec.Command("git", "clone", "--no-checkout", "--quiet", url, temporary).CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("clone source: %s: %w", strings.TrimSpace(string(output)), err)
	}
	if output, err := exec.Command("git", "-C", temporary, "checkout", "--quiet", "--detach", commit).CombinedOutput(); err != nil {
		return "", "", fmt.Errorf("checkout %s: %s: %w", commit, strings.TrimSpace(string(output)), err)
	}
	if err := os.RemoveAll(filepath.Join(temporary, ".git")); err != nil {
		return "", "", err
	}
	if err := os.Rename(temporary, root); err != nil {
		return "", "", err
	}
	return root, commit, nil
}

func validateRecordedInputs(result *Result, file *syntax.File) {
	expected := make(map[string]string)
	for _, entry := range file.Entries {
		if entry.Kind != syntax.Block || entry.Name != "resolved" {
			continue
		}
		for _, child := range entry.Body {
			if child.Name != "input" || child.Label == nil {
				continue
			}
			name, err := child.Label.ResolveString(nil)
			if err != nil {
				continue
			}
			digest, found, err := stringField(child, "sha256", nil)
			if err == nil && found {
				expected[name] = digest
			}
		}
	}
	if len(expected) == 0 {
		// Direct authored builds may use live local inputs. build --check can
		// recommend recording hashes, but ordinary assembly stays available.
		return
	}
	seen := make(map[string]bool)
	for _, contribution := range result.Contributions {
		if contribution.Node == nil {
			continue
		}
		name := contribution.Source + ":" + contribution.Node.Path
		if seen[name] {
			continue
		}
		seen[name] = true
		want, found := expected[name]
		if !found {
			result.Diagnostics = append(result.Diagnostics, diagnostic(contribution.RangeEntry(), Error, fmt.Sprintf("build has no recorded hash for %s", name)))
			continue
		}
		got, err := NodeHash(contribution.Node)
		if err != nil {
			result.Diagnostics = append(result.Diagnostics, diagnostic(contribution.RangeEntry(), Error, err.Error()))
			continue
		}
		if want != got {
			result.Diagnostics = append(result.Diagnostics, diagnostic(contribution.RangeEntry(), Error, fmt.Sprintf("recorded input %s changed", name)))
		}
	}
}

func (c Contribution) RangeEntry() *syntax.Entry {
	return &syntax.Entry{Range: c.Range}
}

func resolveInclude(result *Result, sources map[string]Source, bindings map[string]string, entry *syntax.Entry, replace bool, inherited []Decision, recipes map[string]bool) {
	if entry.Name == "replace" {
		replace = true
	}
	if fieldReplace(entry) {
		replace = true
	}
	if entry.Label == nil {
		inline, found, err := stringField(entry, "text", bindings)
		if err != nil || !found {
			result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, "inline include requires text"))
			return
		}
		into, err := destination(entry, bindings, "AGENTS.md")
		if err != nil {
			result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, err.Error()))
			return
		}
		result.Contributions = append(result.Contributions, Contribution{Inline: inline, Into: into, Replace: replace, Range: entry.Range, Decision: inherited})
		return
	}
	ref, err := entry.Label.ResolveString(bindings)
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, err.Error()))
		return
	}
	alias, nodePath, found := strings.Cut(ref, ":")
	if !found || alias == "" || nodePath == "" {
		result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, "include source must be alias:path"))
		return
	}
	source, found := sources[alias]
	if !found {
		result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, fmt.Sprintf("undefined source alias %q", alias)))
		return
	}
	if strings.HasPrefix(nodePath, "self:") {
		result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, "self: is valid only inside a library preset"))
		return
	}
	if preset, found := source.Lib.Presets[nodePath]; found {
		for _, member := range preset.Entries {
			copied := cloneInclude(member, alias)
			resolveInclude(result, sources, bindings, copied, replace, append(inherited, Decision{Path: alias + ":" + nodePath, Action: "preset", Range: entry.Range}), recipes)
		}
		return
	}
	if recipe, found := source.Lib.Recipes[nodePath]; found {
		key := alias + ":" + nodePath
		if recipes != nil && recipes[key] {
			result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, fmt.Sprintf("recipe include cycle at %s", key)))
			return
		}
		next := make(map[string]bool, len(recipes)+1)
		for name, active := range recipes {
			next[name] = active
		}
		next[key] = true
		for _, member := range recipe.Entries {
			copied := cloneInclude(member, alias)
			resolveInclude(result, sources, bindings, copied, replace, append(inherited, Decision{Path: key, Action: "recipe", Range: entry.Range}), next)
		}
		return
	}
	node, found := source.Lib.Node(nodePath)
	if !found {
		result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, fmt.Sprintf("source %q has no node %q", alias, nodePath)))
		return
	}
	values, err := includeValues(entry, bindings)
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, err.Error()))
		return
	}
	selected, decisions := selectNodes(result, alias, node, entry)
	for _, selectedNode := range selected {
		into, err := destinationForNode(entry, bindings, selectedNode)
		if err != nil {
			result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, err.Error()))
			continue
		}
		if err := validateValues(selectedNode, values, bindings); err != nil {
			result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, err.Error()))
			continue
		}
		merged := append(append([]Decision(nil), inherited...), decisions...)
		result.Contributions = append(result.Contributions, Contribution{Source: alias, Node: selectedNode, Into: into, Replace: replace, Values: values, Range: entry.Range, Decision: merged})
		result.Nodes = append(result.Nodes, ResultNode{Source: alias, Path: selectedNode.Path, Status: Selected, Destination: into, Range: entry.Range})
	}
}

func cloneInclude(e *syntax.Entry, alias string) *syntax.Entry {
	copy := *e
	if copy.Label != nil {
		label, _ := copy.Label.ResolveString(nil)
		if strings.HasPrefix(label, "self:") {
			copy.Label = &syntax.Value{Kind: syntax.String, Parts: []syntax.StringPart{{Text: alias + ":" + strings.TrimPrefix(label, "self:")}}, Range: copy.Label.Range}
		}
	}
	return &copy
}

func selectNodes(result *Result, alias string, root *Node, entry *syntax.Entry) ([]*Node, []Decision) {
	tagsAll, _, err := fieldStrings(entry, "tags_all")
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, err.Error()))
		return nil, nil
	}
	tagsAny, _, err := fieldStrings(entry, "tags_any")
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, err.Error()))
		return nil, nil
	}
	excluded, exclusionReasons, err := exclusions(entry)
	if err != nil {
		result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, err.Error()))
		return nil, nil
	}
	filter := func(nodes []*Node) ([]*Node, []Decision) {
		var kept []*Node
		var decisions []Decision
		for _, node := range nodes {
			relative := strings.TrimPrefix(node.Path, root.Path+"/")
			if excluded[node.Path] || excluded[relative] {
				reason := exclusionReasons[node.Path]
				if reason == "" {
					reason = exclusionReasons[relative]
				}
				if root.Curate == Foundation && strings.TrimSpace(reason) == "" {
					result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, fmt.Sprintf("foundation %s:%s requires a reason to exclude", alias, node.Path)))
				}
				decisions = append(decisions, Decision{Path: alias + ":" + node.Path, Action: "exclude", Reason: reason, Range: entry.Range})
				continue
			}
			if node.ReasonRequired {
				reason, _, _ := stringField(entry, "reason", nil)
				if strings.TrimSpace(reason) == "" {
					result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, fmt.Sprintf("%s:%s requires a nonblank reason", alias, node.Path)))
					result.Nodes = append(result.Nodes, ResultNode{Source: alias, Path: node.Path, Status: Invalid, Reason: "missing required inclusion reason", Range: entry.Range})
					continue
				}
				decisions = append(decisions, Decision{Path: alias + ":" + node.Path, Action: "include", Reason: reason, Range: entry.Range})
			}
			kept = append(kept, node)
		}
		return kept, decisions
	}
	if len(tagsAll) > 0 || len(tagsAny) > 0 {
		var selected []*Node
		seen := make(map[string]bool)
		for _, child := range root.Children {
			if matchesTags(child, tagsAll, tagsAny) {
				selected = append(selected, child)
				seen[child.Path] = true
			}
		}
		if entry.Selection != nil {
			for _, value := range entry.Selection.Items {
				name, err := value.ResolveString(nil)
				if err != nil {
					result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, err.Error()))
					continue
				}
				child := childByPath(root, name)
				if child == nil {
					result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, fmt.Sprintf("%s:%s has no child %q", alias, root.Path, name)))
					continue
				}
				if !seen[child.Path] {
					selected, seen[child.Path] = append(selected, child), true
				}
			}
		}
		if len(selected) == 0 {
			result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Warning, fmt.Sprintf("tag query under %s:%s matched nothing", alias, root.Path)))
		}
		selected, decisions := filter(selected)
		return selected, decisions
	}
	if root.Curate == Choose && entry.Selection == nil {
		if accept, _ := boolField(entry, "accept_defaults"); !accept {
			result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, fmt.Sprintf("UNDECIDED: %s:%s requires explicit choices or accept_defaults", alias, root.Path)))
			result.Nodes = append(result.Nodes, ResultNode{Source: alias, Path: root.Path, Status: Undecided, Reason: "requires explicit choices or accept_defaults", Range: entry.Range})
			return nil, nil
		}
		var selected []*Node
		var decisions []Decision
		for _, child := range root.Children {
			if child.Default == nil {
				result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, fmt.Sprintf("UNDECIDED: %s:%s has no default", alias, child.Path)))
				result.Nodes = append(result.Nodes, ResultNode{Source: alias, Path: child.Path, Status: Undecided, Reason: "choose child has no default", Range: entry.Range})
				continue
			}
			if *child.Default {
				selected = append(selected, child)
				decisions = append(decisions, Decision{Path: alias + ":" + child.Path, Action: "accepted-default", Range: entry.Range})
			}
		}
		selected, filtered := filter(selected)
		return selected, append(decisions, filtered...)
	}
	if entry.Selection == nil && root.Curate == Foundation {
		candidates := root.Children
		if len(candidates) == 0 {
			candidates = []*Node{root}
		}
		selected, decisions := filter(candidates)
		return selected, decisions
	}
	if entry.Selection == nil && root.Curate == Optional {
		var selected []*Node
		for _, child := range root.Children {
			if child.Default != nil && *child.Default {
				selected = append(selected, child)
			}
		}
		selected, decisions := filter(selected)
		return selected, decisions
	}
	if entry.Selection == nil && root.Curate == OptIn {
		result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Warning, fmt.Sprintf("OPT_IN: %s:%s selects no content; choose explicit children or use mogent add %s:%s --all", alias, root.Path, alias, root.Path)))
		result.Nodes = append(result.Nodes, ResultNode{Source: alias, Path: root.Path, Status: Undecided, Reason: "opt-in content requires explicit child selection", Range: entry.Range})
		return nil, nil
	}
	if entry.Selection != nil {
		selected := make([]*Node, 0, len(entry.Selection.Items))
		seen := make(map[string]bool)
		for _, value := range entry.Selection.Items {
			name, err := value.ResolveString(nil)
			if err != nil {
				result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, err.Error()))
				continue
			}
			child := childByPath(root, name)
			if child == nil {
				result.Diagnostics = append(result.Diagnostics, diagnostic(entry, Error, fmt.Sprintf("%s:%s has no child %q", alias, root.Path, name)))
				continue
			}
			if !seen[child.Path] {
				selected = append(selected, child)
				seen[child.Path] = true
			}
		}
		selected, decisions := filter(selected)
		return selected, decisions
	}
	if root.Type == Section && root.Source == "" {
		selected, decisions := filter(sectionLeaves(root))
		return selected, decisions
	}
	selected, decisions := filter([]*Node{root})
	return selected, decisions
}

func matchesTags(node *Node, all, any []string) bool {
	matches := func(query string) bool {
		for _, tag := range node.Tags {
			if tag == query || strings.HasPrefix(tag, query+"/") {
				return true
			}
		}
		return false
	}
	for _, query := range all {
		if !matches(query) {
			return false
		}
	}
	if len(any) == 0 {
		return true
	}
	for _, query := range any {
		if matches(query) {
			return true
		}
	}
	return false
}

func exclusions(entry *syntax.Entry) (map[string]bool, map[string]string, error) {
	excluded := make(map[string]bool)
	reasons := make(map[string]string)
	values, found, err := fieldStrings(entry, "exclude")
	if err != nil {
		return nil, nil, err
	}
	if found {
		for _, value := range values {
			excluded[value] = true
		}
	}
	for _, child := range entry.Body {
		if child.Kind != syntax.Block || child.Name != "exclude" || child.Label == nil {
			continue
		}
		name, err := child.Label.ResolveString(nil)
		if err != nil {
			return nil, nil, err
		}
		reason, _, err := stringField(child, "reason", nil)
		if err != nil {
			return nil, nil, err
		}
		excluded[name], reasons[name] = true, reason
	}
	return excluded, reasons, nil
}

func sectionLeaves(node *Node) []*Node {
	if node.Source != "" || len(node.Children) == 0 {
		return []*Node{node}
	}
	var result []*Node
	for _, child := range node.Children {
		result = append(result, sectionLeaves(child)...)
	}
	return result
}

func childByPath(root *Node, value string) *Node {
	parts := strings.Split(value, "/")
	current := root
	for _, part := range parts {
		var next *Node
		for _, child := range current.Children {
			if child.Name == part {
				next = child
				break
			}
		}
		if next == nil {
			return nil
		}
		current = next
	}
	return current
}

func includeValues(entry *syntax.Entry, bindings map[string]string) (map[string]string, error) {
	values := make(map[string]string)
	for _, child := range entry.Body {
		if child.Name != "values" {
			continue
		}
		for _, value := range child.Body {
			if value.Kind != syntax.Field {
				return nil, fmt.Errorf("values accepts only fields")
			}
			resolved, err := value.Value.ResolveString(bindings)
			if err != nil {
				return nil, err
			}
			values[value.Name] = resolved
		}
	}
	return values, nil
}

func validateValues(node *Node, values map[string]string, bindings map[string]string) error {
	if node.Type != Section {
		if len(values) > 0 {
			return fmt.Errorf("%s is not generated guidance and cannot receive values", node.Path)
		}
		return nil
	}
	for _, variable := range node.Variables {
		if _, found := values[variable]; found {
			continue
		}
		if value, found := bindings[variable]; found {
			values[variable] = value
			continue
		}
		return fmt.Errorf("%s requires value %q", node.Path, variable)
	}
	for value := range values {
		if !contains(node.Variables, value) {
			return fmt.Errorf("%s does not declare value %q", node.Path, value)
		}
	}
	return nil
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func destinationForNode(entry *syntax.Entry, bindings map[string]string, node *Node) (string, error) {
	defaultPath := ""
	switch node.Type {
	case Section:
		defaultPath = "AGENTS.md"
	case Directory, File, Markdown:
		if strings.HasPrefix(node.Path, "skills/") || node.Parent != nil && strings.HasPrefix(node.Parent.Path, "skills/") {
			defaultPath = ".agents/" + node.Path
		} else if node.Into != "" {
			defaultPath = node.Into
		} else {
			return "", fmt.Errorf("UNDEFINED: %s has no destination", node.Path)
		}
	}
	return destination(entry, bindings, defaultPath)
}

func destination(entry *syntax.Entry, bindings map[string]string, fallback string) (string, error) {
	value, found, err := stringField(entry, "into", bindings)
	if err != nil {
		return "", err
	}
	if !found {
		value = fallback
	}
	if value == "" || filepath.IsAbs(value) || strings.HasPrefix(value, "../") || value == ".." || strings.Contains(value, "/../") || strings.Contains(value, "\\") || value == ".mogent" || strings.HasPrefix(value, ".mogent/") {
		return "", fmt.Errorf("destination %q is not permitted", value)
	}
	if value != "AGENTS.md" && !strings.HasPrefix(value, ".agents/") {
		return "", fmt.Errorf("destination %q must be AGENTS.md or under .agents", value)
	}
	return value, nil
}

func stringField(entry *syntax.Entry, name string, bindings map[string]string) (string, bool, error) {
	v := field(entry, name)
	if v == nil {
		return "", false, nil
	}
	value, err := v.ResolveString(bindings)
	return value, true, err
}

func boolField(entry *syntax.Entry, name string) (bool, bool) {
	v := field(entry, name)
	if v == nil || v.Kind != syntax.Boolean {
		return false, false
	}
	return v.Bool, true
}

func fieldReplace(entry *syntax.Entry) bool {
	value, found := boolField(entry, "replace")
	return found && value
}

func diagnostic(entry *syntax.Entry, severity Severity, message string) Diagnostic {
	return Diagnostic{Severity: severity, Range: entry.Range, Message: message}
}

// LoadFile is a convenience for callers that own a plan/build path.
func LoadFile(path string) (*syntax.File, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return syntax.Parse(path, bytes)
}
