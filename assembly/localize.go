package assembly

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

type LocalizeResult struct {
	SourceReference string
	LocalReference  string
	LocalPath       string
	Plan            []byte
}

// Localize copies one complete physical node into the workspace-local library and
// appends a normal local replacement include. It never installs generated output.
func Localize(plan *syntax.File, planPath, reference string, dryRun bool) (*LocalizeResult, error) {
	result, _, err := ResolvePlan(plan, planPath)
	if err != nil {
		return nil, err
	}
	if result.HasErrors() {
		return nil, fmt.Errorf("cannot localize a plan with unresolved errors")
	}
	var wanted *Contribution
	for index := range result.Contributions {
		candidate := &result.Contributions[index]
		if candidate.Node != nil && candidate.Source+":"+candidate.Node.Path == reference {
			wanted = candidate
			break
		}
	}
	if wanted == nil || wanted.Node.Type == Section {
		return nil, fmt.Errorf("localize requires one selected physical source reference %q", reference)
	}
	localRoot := filepath.Join(filepath.Dir(planPath), ".mogent", "local")
	relative := wanted.Node.Path
	localPath := filepath.Join(localRoot, filepath.FromSlash(relative))
	if _, err := os.Lstat(localPath); err == nil {
		return nil, fmt.Errorf("local override %q already exists", localPath)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	copy, err := capturePhysical(wanted.Node)
	if err != nil {
		return nil, err
	}
	updated := cloneFile(plan)
	if !hasSource(updated, "local") {
		updated.Entries = append(updated.Entries, block("source", "local", fields(fieldStringEntry("local", ".mogent/local", false))))
	}
	updated.Entries = append(updated.Entries, block("include", "local:"+relative, []*syntax.Entry{fieldBoolEntry("replace", true)}))
	formatted, err := syntax.Format(updated, syntax.FormatOptions{})
	if err != nil {
		return nil, err
	}
	value := &LocalizeResult{SourceReference: reference, LocalReference: "local:" + relative, LocalPath: localPath, Plan: formatted}
	if dryRun {
		return value, nil
	}
	if err := os.MkdirAll(filepath.Dir(localRoot), 0755); err != nil {
		return nil, err
	}
	stage, err := os.MkdirTemp(filepath.Dir(localRoot), ".mogent-localize-")
	if err != nil {
		return nil, err
	}
	defer func() { _ = os.RemoveAll(stage) }()
	if err := copy.Materialize(stage); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(filepath.Dir(localPath), 0755); err != nil {
		return nil, err
	}
	if err := os.Rename(stage, localPath); err != nil {
		return nil, err
	}
	if err := renderfs.WriteAtomically(planPath, formatted); err != nil {
		_ = os.RemoveAll(localPath)
		return nil, err
	}
	return value, nil
}

func hasSource(file *syntax.File, alias string) bool {
	for _, entry := range file.Entries {
		if entry.Name != "source" || entry.Label == nil {
			continue
		}
		name, _ := entry.Label.ResolveString(nil)
		if name == alias {
			return true
		}
	}
	return false
}

func cloneFile(file *syntax.File) *syntax.File {
	copy := *file
	copy.Entries = make([]*syntax.Entry, 0, len(file.Entries))
	for _, entry := range file.Entries {
		copy.Entries = append(copy.Entries, cloneEntry(entry))
	}
	return &copy
}

func capturePhysical(node *Node) (*content.Manifest, error) {
	inventory, err := physicalInventory(node)
	if err != nil {
		return nil, err
	}
	root := physicalPath(node)
	manifest := &content.Manifest{Files: make(map[string]content.Payload), Directories: map[string]os.FileMode{".": 0755}, Bundles: make(map[string]string)}
	for _, name := range inventory.Paths() {
		if root != "." && name != root && !strings.HasPrefix(name, root+"/") {
			continue
		}
		item := inventory.Nodes[name]
		relative := strings.TrimPrefix(name, root)
		relative = strings.TrimPrefix(relative, "/")
		if relative == "" {
			relative = "."
		}
		if item.Kind == content.Directory {
			manifest.Directories[relative] = item.Mode
			continue
		}
		contents, err := os.ReadFile(item.Source)
		if err != nil {
			return nil, err
		}
		manifest.Files[relative] = content.Payload{Bytes: contents, Mode: item.Mode, Origin: node.Path}
	}
	return manifest, nil
}
