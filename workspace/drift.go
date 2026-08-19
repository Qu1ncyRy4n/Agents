package workspace

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/render"
	"github.com/Qu1ncyRy4n/Agents/state"
)

type DriftReport struct {
	Status      OutputStatus
	OutputPath  string
	Expected    string
	Actual      string
	Changed     bool
	DirectEdits bool
}

// Drift inspects generated and on-disk output without changing workspace files.
func (s *Session) Drift() (*DriftReport, error) {
	status, err := s.Status()
	if err != nil {
		return nil, err
	}
	actual, err := os.ReadFile(status.OutputPath)
	if errors.Is(err, os.ErrNotExist) {
		actual = nil
	} else if err != nil {
		return nil, fmt.Errorf("read output for drift: %w", err)
	}
	return &DriftReport{
		Status:      status.Output,
		OutputPath:  status.OutputPath,
		Expected:    s.Output,
		Actual:      string(actual),
		Changed:     string(actual) != s.Output,
		DirectEdits: status.Output == StatusDirectEdits,
	}, nil
}

// RejectDrift explicitly replaces direct edits with the current manifest render.
func (s *Session) RejectDrift(force bool) error {
	if !force {
		return fmt.Errorf("rejecting direct edits requires --force")
	}
	if err := state.CheckOverwrite(s.OutputPath(), s.StatePath(), true); err != nil {
		return err
	}
	if err := render.WriteAtomically(s.OutputPath(), s.Output); err != nil {
		return err
	}
	return state.Write(s.StatePath(), s.Output)
}

// ImportDrift localizes one manifest section only when all direct edits are
// contained within that unambiguously matched section.
func (s *Session) ImportDrift(manifestHeading, from string) (*LocalizeResult, error) {
	report, err := s.Drift()
	if err != nil {
		return nil, err
	}
	if !report.DirectEdits {
		return nil, fmt.Errorf("output has no recorded direct edits to import; current status is %s", report.Status)
	}
	path := splitManifestPath(manifestHeading)
	if len(path) == 0 {
		return nil, fmt.Errorf("drift import requires a manifest heading path")
	}
	expected, err := markdownSection(report.Expected, path)
	if err != nil {
		return nil, fmt.Errorf("find generated section: %w", err)
	}
	actual, err := markdownSection(report.Actual, path)
	if err != nil {
		return nil, fmt.Errorf("find edited section: %w", err)
	}
	if expected.Full == actual.Full {
		return nil, fmt.Errorf("manifest section %q has no direct edits", manifestHeading)
	}
	reconciled := report.Actual[:actual.Start] + expected.Full + report.Actual[actual.End:]
	if reconciled != report.Expected {
		return nil, fmt.Errorf("direct edits exist outside manifest section %q; import refused", manifestHeading)
	}
	content := shiftSectionContent(actual.Content, actual.Level-1)
	return s.Localize(LocalizeOptions{
		ManifestHeading: manifestHeading,
		From:            from,
		Content:         &content,
		Rebuild:         true,
		ForceOutput:     true,
	})
}

type sectionSpan struct {
	Start   int
	End     int
	Level   int
	Full    string
	Content string
}

type markdownHeading struct {
	start        int
	contentStart int
	level        int
	text         string
	path         []string
}

func markdownSection(document string, target []string) (sectionSpan, error) {
	headings := parseMarkdownHeadings(document)
	var matches []int
	for index, heading := range headings {
		if samePath(heading.path, target) {
			matches = append(matches, index)
		}
	}
	if len(matches) == 0 {
		return sectionSpan{}, fmt.Errorf("heading path %q was not found", strings.Join(target, "/"))
	}
	if len(matches) > 1 {
		return sectionSpan{}, fmt.Errorf("heading path %q is ambiguous", strings.Join(target, "/"))
	}
	index := matches[0]
	heading := headings[index]
	end := len(document)
	for _, next := range headings[index+1:] {
		if next.level <= heading.level {
			end = next.start
			break
		}
	}
	return sectionSpan{
		Start:   heading.start,
		End:     end,
		Level:   heading.level,
		Full:    document[heading.start:end],
		Content: document[heading.contentStart:end],
	}, nil
}

func parseMarkdownHeadings(document string) []markdownHeading {
	var headings []markdownHeading
	var stack []markdownHeading
	inFence := false
	offset := 0
	for _, line := range strings.SplitAfter(document, "\n") {
		trimmed := strings.TrimSpace(strings.TrimSuffix(line, "\n"))
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
			offset += len(line)
			continue
		}
		if !inFence {
			level, text, ok := outputHeading(line)
			if ok {
				for len(stack) > 0 && stack[len(stack)-1].level >= level {
					stack = stack[:len(stack)-1]
				}
				path := make([]string, 0, len(stack)+1)
				for _, parent := range stack {
					path = append(path, parent.text)
				}
				path = append(path, text)
				heading := markdownHeading{start: offset, contentStart: offset + len(line), level: level, text: text, path: path}
				headings = append(headings, heading)
				stack = append(stack, heading)
			}
		}
		offset += len(line)
	}
	return headings
}

func outputHeading(line string) (int, string, bool) {
	line = strings.TrimSuffix(line, "\n")
	level := 0
	for level < len(line) && line[level] == '#' {
		level++
	}
	if level == 0 || level > 6 || level >= len(line) || line[level] != ' ' {
		return 0, "", false
	}
	text := strings.TrimSpace(strings.TrimSuffix(strings.TrimSpace(line[level+1:]), "#"))
	return level, text, text != ""
}

func shiftSectionContent(content string, amount int) string {
	if amount <= 0 {
		return strings.TrimSpace(content)
	}
	var output strings.Builder
	inFence := false
	for _, line := range strings.SplitAfter(content, "\n") {
		trimmed := strings.TrimSpace(strings.TrimSuffix(line, "\n"))
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inFence = !inFence
		}
		if !inFence {
			level, _, ok := outputHeading(line)
			if ok && level > amount {
				line = line[amount:]
			}
		}
		output.WriteString(line)
	}
	return strings.TrimSpace(output.String())
}
