package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RenderMarkdown renders every Markdown output in plan. It validates that a
// leaf source's first heading agrees with its server-side section title, then
// emits the source body below the sidecar-owned hierarchy.
func RenderMarkdown(plan *Plan, libraries map[string]*Library) (map[string]string, error) {
	outputs := make(map[string]string)
	for _, output := range plan.Outputs {
		if output.Kind != "markdown" {
			continue
		}
		var rendered strings.Builder
		for _, source := range output.Sources {
			library, found := libraries[source.Name]
			if !found {
				return nil, fmt.Errorf("output %q source %q is not loaded", output.Name, source.Name)
			}
			if err := renderSource(&rendered, library, source.Sections); err != nil {
				return nil, fmt.Errorf("output %q source %q: %w", output.Name, source.Name, err)
			}
		}
		content := strings.TrimSpace(rendered.String())
		if content == "" {
			return nil, fmt.Errorf("output %q renders no Markdown content", output.Name)
		}
		outputs[output.Name] = content + "\n"
	}
	return outputs, nil
}

func renderSource(output *strings.Builder, library *Library, sections []*Section) error {
	previous := make([]*Section, 0)
	for _, section := range sections {
		ancestors, err := sectionAncestors(library, section)
		if err != nil {
			return err
		}
		shared := 0
		for shared < len(previous) && shared < len(ancestors) && previous[shared] == ancestors[shared] {
			shared++
		}
		for index := shared; index < len(ancestors); index++ {
			if index >= 6 {
				return fmt.Errorf("section %q exceeds Markdown heading level 6", section.Path)
			}
			output.WriteString(strings.Repeat("#", index+1))
			output.WriteByte(' ')
			output.WriteString(ancestors[index].Title)
			output.WriteString("\n\n")
		}
		body, err := sectionBody(library, section)
		if err != nil {
			return err
		}
		if body != "" {
			output.WriteString(body)
			output.WriteString("\n\n")
		}
		previous = ancestors
	}
	return nil
}

func sectionAncestors(library *Library, section *Section) ([]*Section, error) {
	parts := strings.Split(section.Path, "/")
	ancestors := make([]*Section, 0, len(parts))
	for end := 1; end <= len(parts); end++ {
		path := strings.Join(parts[:end], "/")
		ancestor, found := library.ByPath[path]
		if !found {
			return nil, fmt.Errorf("section %q has missing ancestor %q", section.Path, path)
		}
		ancestors = append(ancestors, ancestor)
	}
	return ancestors, nil
}

func sectionBody(library *Library, section *Section) (string, error) {
	contents, err := os.ReadFile(filepath.Join(library.Root, filepath.FromSlash(library.MarkdownRoot), filepath.FromSlash(section.Source)))
	if err != nil {
		return "", fmt.Errorf("read section %q: %w", section.Path, err)
	}
	lines := strings.Split(strings.ReplaceAll(string(contents), "\r\n", "\n"), "\n")
	if len(lines) == 0 || lines[0] != "# "+section.Title {
		return "", fmt.Errorf("section %q source heading must be %q", section.Path, "# "+section.Title)
	}
	return strings.TrimSpace(strings.Join(lines[1:], "\n")), nil
}
