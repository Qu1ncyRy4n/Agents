// Package starter provides small, editable manifest starting points.
package starter

import (
	"fmt"
	"sort"

	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
)

type Template struct {
	Name        string
	Description string
	Required    []string
	Doc         []manifest.Entry
}

func List() []Template {
	values := []Template{
		minimal(),
		goProject(),
		personalGoNix(),
		researchPython(),
	}
	sort.Slice(values, func(i, j int) bool { return values[i].Name < values[j].Name })
	return values
}

func Get(name string) (Template, error) {
	for _, value := range List() {
		if value.Name == name {
			return value, nil
		}
	}
	return Template{}, fmt.Errorf("unknown starter template %q", name)
}

func (t Template) Manifest(sources map[string]string, output string) (*manifest.Manifest, error) {
	normalized := make(map[string]manifest.Source, len(sources))
	for _, alias := range t.Required {
		if sources[alias] == "" {
			return nil, fmt.Errorf("template %q requires --source %s=<path-or-url>", t.Name, alias)
		}
	}
	for alias, location := range sources {
		normalized[alias] = manifest.Source{Location: location}
	}
	value := &manifest.Manifest{Sources: normalized, Output: output, Doc: cloneEntries(t.Doc)}
	if err := value.Validate(); err != nil {
		return nil, err
	}
	return value, nil
}

func minimal() Template {
	return Template{
		Name:        "minimal",
		Description: "shared identity, focused workflow, safe defaults, and clear handoff",
		Required:    []string{"shared"},
		Doc: []manifest.Entry{
			children("Identity", source("Role", "shared:shared-baseline/identity/role"), source("Source Of Truth", "shared:shared-baseline/identity/source-of-truth")),
			children("Instructions", source("Focused Change Loop", "shared:shared-baseline/instructions/focused-change-loop")),
			children("Constraints", source("Safe Defaults", "shared:shared-baseline/constraints/safe-defaults")),
			children("Format", source("Clear Handoff", "shared:shared-baseline/format/clear-handoff")),
		},
	}
}

func goProject() Template {
	value := minimal()
	value.Name = "go"
	value.Description = "minimal shared guidance plus focused Go development and tests"
	value.Required = []string{"shared", "go"}
	value.Doc[1].Children = append(value.Doc[1].Children,
		source("Go Development", "go:go/development"),
		source("Go Tests", "go:go/tests"),
	)
	return value
}

func personalGoNix() Template {
	value := goProject()
	value.Name = "personal-go-nix"
	value.Description = "Go project with independently selected personal Nix safety and scope"
	value.Required = []string{"shared", "go", "personal"}
	value.Doc[2].Children = append(value.Doc[2].Children,
		source("Nix Safety", "personal:lang/nix/safety"),
		source("Nix Scope", "personal:lang/nix/scope"),
	)
	return value
}

func researchPython() Template {
	value := minimal()
	value.Name = "research-python"
	value.Description = "shared baseline plus research protection, validity, and Python/uv workflow"
	value.Required = []string{"shared", "research"}
	value.Doc[1].Children = append(value.Doc[1].Children,
		source("Research Orientation", "research:ucd-research/research-orientation"),
		source("Python Dependencies", "research:python/dependency-management"),
	)
	value.Doc[2].Children = append(value.Doc[2].Children,
		source("Research Data Protection", "research:ucd-research/data-protection"),
		source("Experiment Validity", "research:ucd-research/experiment-validity"),
	)
	return value
}

func source(heading, reference string) manifest.Entry {
	return manifest.Entry{Heading: heading, From: []string{reference}}
}

func children(heading string, entries ...manifest.Entry) manifest.Entry {
	return manifest.Entry{Heading: heading, Children: entries}
}

func cloneEntries(entries []manifest.Entry) []manifest.Entry {
	result := make([]manifest.Entry, len(entries))
	for index, entry := range entries {
		result[index] = manifest.Entry{
			Heading:  entry.Heading,
			From:     append([]string(nil), entry.From...),
			Exclude:  append([]string(nil), entry.Exclude...),
			Children: cloneEntries(entry.Children),
		}
	}
	return result
}
