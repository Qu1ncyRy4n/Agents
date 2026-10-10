package v2

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/state"
	"github.com/Qu1ncyRy4n/Agents/textdiff"
)

func compileArtifactOutput(plan *Plan, planned *PlannedOutput, output Output, libraries map[string]*Library) {
	for _, source := range output.Sources {
		library := libraries[source.Name]
		if library == nil {
			plan.Diagnostics = append(plan.Diagnostics, Diagnostic{Severity: SeverityError, Code: "MOGENT101", Message: fmt.Sprintf("source %q is not loaded", source.Name)})
			continue
		}
		root := strings.TrimPrefix(source.From, source.Name+":")
		if source.Select != nil {
			subconfig := &Config{Outputs: []Output{{Name: output.Name, Kind: "markdown", Sources: []OutputSource{source}}}}
			subplan := Compile(subconfig, libraries)
			plan.Diagnostics = append(plan.Diagnostics, subplan.Diagnostics...)
			if len(subplan.Outputs[0].Sources) > 0 {
				planned.Sources = append(planned.Sources, subplan.Outputs[0].Sources[0])
			}
			continue
		}
		if library.Contents[root] == nil {
			plan.Diagnostics = append(plan.Diagnostics, Diagnostic{Severity: SeverityError, Code: "MOGENT104", Message: fmt.Sprintf("source %q has no named content root %q", source.Name, root)})
			continue
		}
		planned.Sources = append(planned.Sources, PlannedSource{Name: source.Name, From: source.From, TreeRoot: root})
	}
}

func composeArtifacts(config *Config, result *Result) (map[string]*content.Manifest, error) {
	manifests := make(map[string]*content.Manifest)
	for index, output := range config.Outputs {
		if output.Kind != "dir-tree" {
			continue
		}
		var contributions []content.Contribution
		for sourceIndex, source := range output.Sources {
			library := result.Libraries[source.Name]
			root := strings.TrimPrefix(source.From, source.Name+":")
			contribution := content.Contribution{Inventory: library.Contents[root], Node: source.Node, Heading: source.Heading, Operation: content.Operation(source.Operation), Into: source.Into, Origin: source.From + "/" + source.Node, Exclude: source.Exclude, Replace: source.Replace, Append: source.Append}
			resolved := result.Sources[source.Name]
			if resolved.Commit != "" {
				contribution.Origin += "@" + resolved.Commit
			}
			if len(source.Heading) > 0 {
				contribution.Origin += " heading " + strings.Join(source.Heading, " / ")
			}
			if source.Select != nil {
				subplan := &Plan{Outputs: []PlannedOutput{{Name: output.Name, Kind: "markdown", Sources: []PlannedSource{result.Plan.Outputs[index].Sources[sourceIndex]}}}}
				rendered, err := RenderMarkdown(subplan, result.Libraries)
				if err != nil {
					return nil, err
				}
				doc, err := content.ParseDocument([]byte(rendered[output.Name]))
				if err != nil {
					return nil, err
				}
				view, err := doc.Render(nil, 1)
				if err != nil {
					return nil, err
				}
				contribution.Rendered = &view
			}
			contributions = append(contributions, contribution)
		}
		manifest, err := content.Compose(contributions)
		if err != nil {
			return nil, fmt.Errorf("output %q: %w", output.Name, err)
		}
		manifests[output.Name] = manifest
	}
	return manifests, nil
}

func manifestHashes(manifest *content.Manifest) map[string]string {
	hashes := make(map[string]string, len(manifest.Files))
	for relative, payload := range manifest.Files {
		hashes[relative] = state.Hash(payload.Bytes)
	}
	return hashes
}

func planArtifacts(result *Result, root string) error {
	for _, output := range result.Plan.Outputs {
		manifest := result.Manifests[output.Name]
		if manifest == nil {
			continue
		}
		for _, relative := range output.Paths {
			change, err := inspectTree(result, output.Name, relative, filepath.Join(root, filepath.FromSlash(relative)), "", manifestHashes(manifest))
			if err != nil {
				return err
			}
			change.Manifest = manifest
			before := make(map[string]content.Payload)
			for member, payload := range manifest.Files {
				if !payload.Rendered {
					continue
				}
				bytes, err := os.ReadFile(filepath.Join(change.Absolute, filepath.FromSlash(member)))
				if os.IsNotExist(err) {
					continue
				}
				if err != nil {
					return err
				}
				before[member] = content.Payload{Bytes: bytes}
			}
			change.Documents = compareDocuments(output.Name, relative, before, manifest.Files)
			result.Trees = append(result.Trees, change)
		}
	}
	return nil
}

func compareDocuments(output, root string, before, after map[string]content.Payload) []FileChange {
	var names []string
	for name, payload := range after {
		if payload.Rendered {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	var changes []FileChange
	for _, name := range names {
		previous, exists := before[name]
		content := string(after[name].Bytes)
		relative := path.Join(root, name)
		change := FileChange{Output: output, Path: relative, Content: content, Status: FileUnchanged}
		oldLabel := "a/" + relative
		if !exists {
			change.Status = FileNew
			oldLabel = "/dev/null"
		}
		change.Diff = textdiff.Unified(string(previous.Bytes), content, oldLabel, "b/"+relative)
		if exists && change.Diff != "" {
			change.Status = FileChanged
		}
		changes = append(changes, change)
	}
	return changes
}

// Update compares expected old/new source manifests, not a potentially edited
// working output. This covers both typed artifacts and legacy raw tree outputs.
func directoryOutputHashes(result *Result) (map[string]map[string]string, error) {
	outputs := make(map[string]map[string]string)
	for _, output := range result.Plan.Outputs {
		if manifest := result.Manifests[output.Name]; manifest != nil {
			outputs[output.Name] = manifestHashes(manifest)
			continue
		}
		if output.Kind != "tree" {
			continue
		}
		source := output.Sources[0]
		library := result.Libraries[source.Name]
		tree := library.Trees[source.TreeRoot]
		directory := filepath.Join(library.Root, filepath.FromSlash(tree.Root))
		if err := verifySourceTree(directory); err != nil {
			return nil, err
		}
		for _, excluded := range source.Exclude {
			if _, err := os.Lstat(filepath.Join(directory, filepath.FromSlash(excluded))); err != nil {
				return nil, fmt.Errorf("excluded path %q in source %q: %w", excluded, source.Name, err)
			}
		}
		hashes, err := state.DirectoryHashes(directory)
		if err != nil {
			return nil, err
		}
		outputs[output.Name] = filterTreeHashes(hashes, source.Exclude)
	}
	return outputs, nil
}
