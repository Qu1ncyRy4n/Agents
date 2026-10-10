package v2

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/state"
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
			result.Trees = append(result.Trees, change)
		}
	}
	return nil
}
