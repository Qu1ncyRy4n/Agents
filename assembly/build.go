package assembly

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

const checksumPlaceholder = "0000000000000000000000000000000000000000000000000000000000000000"

// Sequence resolves a plan into canonical, self-contained build instructions.
// It writes nothing; callers choose when to install the returned build bytes.
func Sequence(plan *syntax.File, planPath string) ([]byte, *Result, error) {
	if plan.Kind != syntax.Plan {
		return nil, nil, fmt.Errorf("sequence requires a plan")
	}
	result, _, err := ResolvePlan(plan, planPath)
	if err != nil || result.HasErrors() {
		return nil, result, err
	}
	build := &syntax.File{Kind: syntax.Build, Entries: []*syntax.Entry{block("build", "", fields(
		fieldStringEntry("format", "1", true),
		fieldStringEntry("origin", "generated", false),
		fieldStringEntry("generated_from", "mogent.plan", false),
		fieldStringEntry("generated_by", "mogent", false),
		fieldStringEntry("plan_sha256", hashBytes([]byte(plan.Original())), false),
		fieldStringEntry("content_sha256", checksumPlaceholder, false),
	))}}
	for _, entry := range plan.Entries[1:] {
		if entry.Kind == syntax.Block && entry.Name == "source" {
			build.Entries = append(build.Entries, cloneEntry(entry))
		}
	}
	for _, contribution := range result.Contributions {
		build.Entries = append(build.Entries, buildContribution(contribution))
	}
	resolved, err := resolvedEntry(result)
	if err != nil {
		return nil, result, err
	}
	build.Entries = append(build.Entries, resolved)
	formatted, err := syntax.Format(build, syntax.FormatOptions{ExpandAliases: true})
	if err != nil {
		return nil, result, err
	}
	digest := hashBytes(formatted)
	formatted = []byte(strings.Replace(string(formatted), checksumPlaceholder, digest, 1))
	return formatted, result, nil
}

func buildContribution(c Contribution) *syntax.Entry {
	body := []*syntax.Entry{fieldStringEntry("into", c.Into, false)}
	if c.Replace {
		body = append(body, fieldBoolEntry("replace", true))
	}
	if c.Inline != "" {
		body = append(body, fieldStringEntry("text", c.Inline, false))
		return block("include", "", body)
	}
	for _, name := range sortedValues(c.Values) {
		if len(body) == 0 || body[len(body)-1].Name != "values" {
			body = append(body, block("values", "", nil))
		}
		body[len(body)-1].Body = append(body[len(body)-1].Body, fieldStringEntry(name, c.Values[name], false))
	}
	return block("include", c.Source+":"+c.Node.Path, body)
}

func resolvedEntry(result *Result) (*syntax.Entry, error) {
	entry := block("resolved", "", nil)
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
		digest, err := NodeHash(contribution.Node)
		if err != nil {
			return nil, err
		}
		entry.Body = append(entry.Body, block("input", name, fields(fieldStringEntry("sha256", digest, false))))
	}
	for _, c := range result.Contributions {
		for _, d := range c.Decision {
			body := []*syntax.Entry{fieldStringEntry("action", d.Action, false)}
			if d.Reason != "" {
				body = append(body, fieldStringEntry("reason", d.Reason, false))
			}
			entry.Body = append(entry.Body, block("decision", d.Path, body))
		}
	}
	entry.Body = append(entry.Body, block("owner", "AGENTS.md", fields(fieldStringEntry("kind", "file", false))))
	entry.Body = append(entry.Body, block("owner", ".agents/skills", fields(fieldStringEntry("kind", "directory", false))))
	return entry, nil
}

func sortedValues(values map[string]string) []string {
	var names []string
	for name := range values {
		names = append(names, name)
	}
	for i := 1; i < len(names); i++ {
		for j := i; j > 0 && names[j] < names[j-1]; j-- {
			names[j], names[j-1] = names[j-1], names[j]
		}
	}
	return names
}

func hashBytes(value []byte) string {
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

// NodeHash describes exactly the selected payload. It is an integrity record,
// not a content-addressed identity or a replacement for retaining source bytes.
func NodeHash(node *Node) (string, error) {
	if node == nil {
		return "", fmt.Errorf("cannot hash nil node")
	}
	if node.Type == Section {
		if node.Source == "" {
			return hashBytes([]byte(node.Path)), nil
		}
		path, err := safePath(filepath.Dir(node.Range.Filename), node.Source)
		if err != nil {
			return "", err
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return "", err
		}
		return hashBytes(contents), nil
	}
	inventory, err := physicalInventory(node)
	if err != nil {
		return "", err
	}
	root := physicalPath(node)
	type record struct {
		Kind string `json:"kind"`
		Path string `json:"path"`
		Mode string `json:"mode"`
		Hash string `json:"sha256"`
	}
	var records []record
	for _, name := range inventory.Paths() {
		if root != "." && name != root && !strings.HasPrefix(name, root+"/") {
			continue
		}
		item := inventory.Nodes[name]
		r := record{Kind: string(item.Kind), Path: name, Mode: fmt.Sprintf("%04o", item.Mode.Perm())}
		if item.Kind == content.File {
			contents, err := os.ReadFile(item.Source)
			if err != nil {
				return "", err
			}
			r.Hash = hashBytes(contents)
		}
		records = append(records, r)
	}
	encoded, err := json.Marshal(records)
	if err != nil {
		return "", err
	}
	return hashBytes(encoded), nil
}

func str(value string) *syntax.Value {
	return &syntax.Value{Kind: syntax.String, Parts: []syntax.StringPart{{Text: value}}}
}

func fieldStringEntry(name, value string, integer bool) *syntax.Entry {
	if integer {
		return &syntax.Entry{Kind: syntax.Field, Name: name, Value: &syntax.Value{Kind: syntax.Integer, Int: 1}}
	}
	return &syntax.Entry{Kind: syntax.Field, Name: name, Value: str(value)}
}

func fieldBoolEntry(name string, value bool) *syntax.Entry {
	return &syntax.Entry{Kind: syntax.Field, Name: name, Value: &syntax.Value{Kind: syntax.Boolean, Bool: value}}
}

func fields(entries ...*syntax.Entry) []*syntax.Entry { return entries }

func block(name, label string, body []*syntax.Entry) *syntax.Entry {
	e := &syntax.Entry{Kind: syntax.Block, Name: name, Body: body, HasBody: true}
	if label != "" {
		e.Label = str(label)
	}
	return e
}

func cloneEntry(entry *syntax.Entry) *syntax.Entry {
	copy := *entry
	copy.Body = make([]*syntax.Entry, 0, len(entry.Body))
	for _, child := range entry.Body {
		copy.Body = append(copy.Body, cloneEntry(child))
	}
	if entry.Label != nil {
		label := *entry.Label
		label.Parts = append([]syntax.StringPart(nil), entry.Label.Parts...)
		copy.Label = &label
	}
	if entry.Value != nil {
		value := *entry.Value
		value.Parts = append([]syntax.StringPart(nil), entry.Value.Parts...)
		value.Items = append([]syntax.Value(nil), entry.Value.Items...)
		copy.Value = &value
	}
	return &copy
}

// VerifyBuildChecksum rejects a generated build modified after sequencing.
func VerifyBuildChecksum(contents []byte) error {
	file, err := syntax.Parse("mogent.build", contents)
	if err != nil {
		return err
	}
	if file.Kind != syntax.Build {
		return fmt.Errorf("not a build file")
	}
	header := file.Entries[0]
	origin, found, err := fieldString(header, "origin")
	if err != nil || !found {
		return fmt.Errorf("build has no origin")
	}
	if origin == "authored" {
		return nil
	}
	want, found, err := fieldString(header, "content_sha256")
	if err != nil || !found || len(want) != 64 {
		return fmt.Errorf("generated build has no valid content checksum")
	}
	actual := hashBytes([]byte(strings.Replace(string(contents), want, checksumPlaceholder, 1)))
	if want != actual {
		return fmt.Errorf("generated build was edited after sequencing")
	}
	return nil
}

// Apply writes the generated root guide and skills subtree. It only owns the two
// declared output roots and preserves unrelated .agents entries.
func Apply(workspace string, manifest *content.Manifest) error {
	if manifest == nil {
		return fmt.Errorf("nil output manifest")
	}
	agents, found := manifest.Files["AGENTS.md"]
	if !found {
		return fmt.Errorf("manifest has no AGENTS.md")
	}
	if err := renderfs.WriteAtomically(filepath.Join(workspace, "AGENTS.md"), agents.Bytes); err != nil {
		return err
	}
	var skills = &content.Manifest{Files: make(map[string]content.Payload), Directories: map[string]os.FileMode{".": 0755}, Bundles: make(map[string]string)}
	for name, payload := range manifest.Files {
		if !strings.HasPrefix(name, ".agents/skills/") {
			continue
		}
		skills.Files[strings.TrimPrefix(name, ".agents/skills/")] = payload
	}
	for name, mode := range manifest.Directories {
		if name == ".agents/skills" || strings.HasPrefix(name, ".agents/skills/") {
			relative := strings.TrimPrefix(name, ".agents/skills")
			relative = strings.TrimPrefix(relative, "/")
			if relative == "" {
				relative = "."
			}
			skills.Directories[relative] = mode
		}
	}
	return replaceDirectory(filepath.Join(workspace, ".agents", "skills"), skills)
}

// CheckWorkspace verifies expected managed outputs without repairing them.
func CheckWorkspace(workspace string, manifest *content.Manifest) error {
	if manifest == nil {
		return fmt.Errorf("nil output manifest")
	}
	for name, payload := range manifest.Files {
		contents, err := os.ReadFile(filepath.Join(workspace, filepath.FromSlash(name)))
		if err != nil {
			return fmt.Errorf("managed output %q: %w", name, err)
		}
		if hashBytes(contents) != hashBytes(payload.Bytes) {
			return fmt.Errorf("managed output %q differs from expected content", name)
		}
	}
	root := filepath.Join(workspace, ".agents", "skills")
	expected := make(map[string]bool)
	for name := range manifest.Files {
		if strings.HasPrefix(name, ".agents/skills/") {
			expected[strings.TrimPrefix(name, ".agents/skills/")] = true
		}
	}
	if err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if !expected[filepath.ToSlash(relative)] {
			return fmt.Errorf("managed skills contains unexpected file %q", relative)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("managed skills tree: %w", err)
	}
	return nil
}

func replaceDirectory(destination string, manifest *content.Manifest) error {
	parent := filepath.Dir(destination)
	if err := os.MkdirAll(parent, 0755); err != nil {
		return err
	}
	stage, err := os.MkdirTemp(parent, ".mogent-stage-")
	if err != nil {
		return err
	}
	defer func() { _ = os.RemoveAll(stage) }()
	if err := manifest.Materialize(stage); err != nil {
		return err
	}
	backup := destination + ".mogent-backup"
	if _, err := os.Lstat(backup); err == nil {
		return fmt.Errorf("refusing to replace existing backup %q", backup)
	} else if !os.IsNotExist(err) {
		return err
	}
	if _, err := os.Lstat(destination); err == nil {
		if err := os.Rename(destination, backup); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.Rename(stage, destination); err != nil {
		if _, backupErr := os.Lstat(backup); backupErr == nil {
			_ = os.Rename(backup, destination)
		}
		return err
	}
	if err := os.RemoveAll(backup); err != nil {
		return err
	}
	return nil
}
