// Package navigator implements the milestone-two terminal manifest browser.
package navigator

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Qu1ncyRy4n/Agents/internal/library"
	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
	"github.com/Qu1ncyRy4n/Agents/internal/render"
	"github.com/Qu1ncyRy4n/Agents/internal/state"
)

type focus int

const (
	focusTree focus = iota
	focusFinal
	focusSource
)

type detailMode int

const (
	detailFinal detailMode = iota
	detailSource
)

type row struct {
	Entry   *manifest.Entry
	Depth   int
	Index   []int
	From    []string
	Parents []string
}

type candidate struct {
	Heading   string
	Reference string
}

// Model is the Bubble Tea state for browsing and drafting changes to one
// manifest.
type Model struct {
	manifestPath  string
	saved         *manifest.Manifest
	draft         *manifest.Manifest
	sources       map[string]*library.Index
	overrides     map[string]string
	output        string
	rows          []row
	selected      int
	collapsed     map[string]bool
	focus         focus
	detail        detailMode
	width         int
	height        int
	err           string
	message       string
	review        string
	confirmSave   bool
	confirmReject bool
	dirty         bool
	drift         bool
	importedDrift bool
}

// New loads a manifest and prepares an in-memory draft.
func New(manifestPath string) (Model, error) {
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		return Model{}, err
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		return Model{}, err
	}
	sources, err := render.LoadSources(value, loadedPath)
	if err != nil {
		return Model{}, err
	}
	model := Model{
		manifestPath: loadedPath,
		saved:        value,
		draft:        value.Clone(),
		sources:      sources,
		overrides:    make(map[string]string),
		output:       result.Content,
		collapsed:    make(map[string]bool),
		focus:        focusTree,
		detail:       detailFinal,
		width:        100,
		height:       30,
	}
	if changed, err := state.Changed(state.OutputPath(loadedPath, value.Output), state.StatePath(loadedPath)); err != nil {
		model.err = err.Error()
	} else {
		model.drift = changed
	}
	model.refreshRows()
	return model, nil
}

// Run starts the interactive terminal browser.
func Run(manifestPath string, output io.Writer) error {
	model, err := New(manifestPath)
	if err != nil {
		return err
	}
	program := tea.NewProgram(model, tea.WithOutput(output))
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("run tui: %w", err)
	}
	return nil
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := message.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		if m.confirmSave {
			switch msg.String() {
			case "y":
				m.confirmSave = false
				if err := m.saveAndBuild(false); err != nil {
					m.err = err.Error()
					m.message = ""
				} else {
					m.err = ""
					m.message = "saved and built"
					m.saved = m.draft.Clone()
					m.dirty = false
					m.drift = false
					m.importedDrift = false
					m.overrides = make(map[string]string)
					if sources, err := render.LoadSources(m.saved, m.manifestPath); err == nil {
						m.sources = sources
					}
				}
			case "n", "esc":
				m.confirmSave = false
				m.discardDraft("save cancelled; draft discarded")
			case "ctrl+c", "q":
				m.discardDraft("")
				return m, tea.Quit
			}
			return m, nil
		}
		if m.confirmReject {
			switch msg.String() {
			case "y":
				m.confirmReject = false
				if err := m.saveAndBuild(true); err != nil {
					m.err = err.Error()
					m.message = ""
				} else {
					m.err = ""
					m.message = "direct edit rejected; rebuilt from manifest"
					m.drift = false
					m.importedDrift = false
					m.saved = m.draft.Clone()
					m.dirty = false
				}
			case "n", "esc":
				m.confirmReject = false
				m.message = "reject cancelled"
			case "ctrl+c", "q":
				m.discardDraft("")
				return m, tea.Quit
			}
			return m, nil
		}
		switch msg.String() {
		case "ctrl+c", "q":
			m.discardDraft("")
			return m, tea.Quit
		case "up", "k":
			if m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if m.selected < len(m.rows)-1 {
				m.selected++
			}
		case "left", "h":
			m.setCollapsed(true)
		case "right", "l":
			m.setCollapsed(false)
		case "a":
			m.addCandidate()
		case "d":
			m.removeSelected()
		case "e":
			m.localizeSelected()
		case "i":
			m.importDriftIntoSelected()
		case "r":
			if m.drift {
				m.err = ""
				m.review = m.saveReview()
				m.message = "reject direct edit and rebuild AGENTS.md? y/n"
				m.confirmReject = true
			}
		case "K":
			if m.drift {
				m.discardDraft("kept direct edit; no files written")
			}
		case "U":
			m.moveSelected(-1)
		case "D":
			m.moveSelected(1)
		case "tab":
			m.nextFocus()
		case "esc":
			m.focus = focusTree
		case "1":
			m.focus = focusTree
		case "2":
			m.focus = focusFinal
			m.detail = detailFinal
		case "3":
			m.focus = focusSource
			m.detail = detailSource
		case "v":
			if m.detail == detailFinal {
				m.detail = detailSource
				m.focus = focusSource
			} else {
				m.detail = detailFinal
				m.focus = focusFinal
			}
		case "s":
			m.err = ""
			m.review = m.saveReview()
			m.message = "save current manifest and build AGENTS.md? y/n"
			m.confirmSave = true
		}
	}
	return m, nil
}

func (m Model) View() string {
	if len(m.rows) == 0 {
		return "mogent tui: manifest has no visible entries\n"
	}
	width := max(m.width, 80)
	height := max(m.height, 12)
	status := m.statusLine(width)
	bodyHeight := max(height-2, 8)
	if width >= 140 {
		treeWidth := 44
		finalWidth := (width - treeWidth - 4) / 2
		sourceWidth := width - treeWidth - finalWidth - 4
		return lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.pane("Manifest tree", m.treeView(treeWidth-2, bodyHeight-2), treeWidth, bodyHeight, m.focus == focusTree),
			" ",
			m.pane("Final AGENTS.md", m.finalView(finalWidth-2, bodyHeight-2), finalWidth, bodyHeight, m.focus == focusFinal),
			" ",
			m.pane("Selected-source context", m.sourceView(sourceWidth-2, bodyHeight-2), sourceWidth, bodyHeight, m.focus == focusSource),
		) + "\n" + status
	}
	treeWidth := width / 2
	detailWidth := width - treeWidth - 1
	title := "Detail: Final"
	content := m.finalView(detailWidth-2, bodyHeight-2)
	active := m.focus == focusFinal
	if m.detail == detailSource {
		title = "Detail: Source"
		content = m.sourceView(detailWidth-2, bodyHeight-2)
		active = m.focus == focusSource
	}
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		m.pane("Manifest tree", m.treeView(treeWidth-2, bodyHeight-2), treeWidth, bodyHeight, m.focus == focusTree),
		" ",
		m.pane(title, content, detailWidth, bodyHeight, active),
	) + "\n" + status
}

func (m Model) statusLine(width int) string {
	stateText := "state: saved"
	if m.dirty {
		stateText = "state: ~ draft"
	}
	parts := []string{
		fmt.Sprintf("manifest: %s", filepath.Base(m.manifestPath)),
		stateText,
		"keys: up/down move, a add, d remove, e localize, i import drift, r reject drift, K keep drift, s save, q quit",
	}
	if m.confirmSave {
		parts = append(parts, "confirm: y/n")
	}
	if m.confirmReject {
		parts = append(parts, "confirm reject: y/n")
	}
	if m.drift {
		parts = append(parts, "drift: direct AGENTS.md edits detected")
	}
	if m.message != "" {
		parts = append(parts, m.message)
	}
	if m.err != "" {
		parts = append(parts, "error: "+m.err)
	}
	return truncate(strings.Join(parts, " | "), width)
}

func (m Model) saveAndBuild(forceOutput bool) error {
	result, err := render.BuildWithSources(m.draft, m.manifestPath, m.sources)
	if err != nil {
		return err
	}
	outputPath := state.OutputPath(m.manifestPath, m.draft.Output)
	statePath := state.StatePath(m.manifestPath)
	if err := state.CheckOverwrite(outputPath, statePath, forceOutput || m.importedDrift); err != nil {
		return err
	}
	originalManifest, err := os.ReadFile(m.manifestPath)
	if err != nil {
		return fmt.Errorf("read existing manifest before save: %w", err)
	}
	originalOutput, outputExisted, err := readOptional(outputPath)
	if err != nil {
		return err
	}
	if err := m.writeOverrides(); err != nil {
		return err
	}
	if err := manifest.WriteAtomically(m.manifestPath, m.draft); err != nil {
		return err
	}
	if err := render.WriteAtomically(outputPath, result.Content); err != nil {
		return rollbackSave(err, m.manifestPath, originalManifest, outputPath, originalOutput, outputExisted)
	}
	if err := state.Write(statePath, result.Content); err != nil {
		return rollbackSave(err, m.manifestPath, originalManifest, outputPath, originalOutput, outputExisted)
	}
	return nil
}

func readOptional(path string) ([]byte, bool, error) {
	contents, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read existing output before save: %w", err)
	}
	return contents, true, nil
}

func rollbackSave(saveErr error, manifestPath string, originalManifest []byte, outputPath string, originalOutput []byte, outputExisted bool) error {
	var rollbackErrs []error
	if err := render.WriteAtomically(manifestPath, string(originalManifest)); err != nil {
		rollbackErrs = append(rollbackErrs, fmt.Errorf("restore manifest: %w", err))
	}
	if outputExisted {
		if err := render.WriteAtomically(outputPath, string(originalOutput)); err != nil {
			rollbackErrs = append(rollbackErrs, fmt.Errorf("restore output: %w", err))
		}
	} else if err := os.Remove(outputPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		rollbackErrs = append(rollbackErrs, fmt.Errorf("remove new output after failed save: %w", err))
	}
	return errors.Join(append([]error{saveErr}, rollbackErrs...)...)
}

func (m Model) pane(title, content string, width, height int, active bool) string {
	borderColor := lipgloss.Color("240")
	if active {
		borderColor = lipgloss.Color("39")
	}
	style := lipgloss.NewStyle().
		Width(width).
		Height(height).
		Border(lipgloss.NormalBorder()).
		BorderForeground(borderColor)
	header := lipgloss.NewStyle().Bold(true).Render(title)
	return style.Render(header + "\n" + limitLines(content, height-3))
}

func (m *Model) refreshRows() {
	m.rows = nil
	for index := range m.draft.Doc {
		m.appendEntry(&m.draft.Doc[index], 0, []int{index}, nil)
	}
	if m.selected >= len(m.rows) {
		m.selected = max(len(m.rows)-1, 0)
	}
}

func (m *Model) appendEntry(entry *manifest.Entry, depth int, index []int, parents []string) {
	m.rows = append(m.rows, row{
		Entry:   entry,
		Depth:   depth,
		Index:   append([]int(nil), index...),
		From:    append([]string(nil), entry.From...),
		Parents: append([]string(nil), parents...),
	})
	if len(entry.Children) == 0 || m.collapsed[rowKey(index)] {
		return
	}
	nextParents := append(append([]string(nil), parents...), entry.Heading)
	for childIndex := range entry.Children {
		childPath := append(append([]int(nil), index...), childIndex)
		m.appendEntry(&entry.Children[childIndex], depth+1, childPath, nextParents)
	}
}

func (m Model) treeView(width, height int) string {
	var lines []string
	for index, row := range m.rows {
		prefix := "  "
		if index == m.selected {
			prefix = "> "
		}
		collapse := " "
		if len(row.Entry.Children) > 0 {
			collapse = "-"
			if m.collapsed[rowKey(row.Index)] {
				collapse = "+"
			}
		}
		marker := "[x]"
		provenance := "children"
		if len(row.From) > 0 {
			provenance = strings.Join(row.From, ", ")
			if hasLocalReference(row.From) {
				marker = "[L]"
			}
		}
		line := fmt.Sprintf("%s%s%s%s %s  <%s>", prefix, strings.Repeat("  ", row.Depth), collapse, marker, row.Entry.Heading, provenance)
		lines = append(lines, truncate(line, width))
	}
	return visibleWindow(lines, m.selected, height)
}

func (m Model) finalView(width, height int) string {
	current := m.rows[m.selected].Entry.Heading
	lines := strings.Split(strings.TrimRight(m.output, "\n"), "\n")
	for index, line := range lines {
		if strings.TrimLeft(line, "# ") == current && strings.HasPrefix(line, "#") {
			lines[index] = "> " + line
		} else {
			lines[index] = "  " + line
		}
		lines[index] = truncate(lines[index], width)
	}
	return visibleWindow(lines, firstHighlighted(lines), height)
}

func (m Model) sourceView(width, height int) string {
	if m.confirmSave || m.confirmReject {
		return truncateLines(strings.Split(m.review, "\n"), width, height)
	}
	selected := m.rows[m.selected]
	if len(selected.From) == 0 {
		var lines []string
		lines = append(lines, "Document section: "+strings.Join(append(selected.Parents, selected.Entry.Heading), " / "))
		lines = append(lines, "Source: child entries only")
		lines = append(lines, "")
		for _, child := range selected.Entry.Children {
			lines = append(lines, "- "+child.Heading)
		}
		candidates := m.candidatesFor(selected)
		if len(candidates) > 0 {
			lines = append(lines, "")
			lines = append(lines, "Available source children:")
			for _, candidate := range candidates {
				lines = append(lines, "+ "+candidate.Heading+"  <"+candidate.Reference+">")
			}
		}
		return truncateLines(lines, width, height)
	}
	var lines []string
	for _, reference := range selected.From {
		alias, path, err := manifest.SplitReference(reference)
		if err != nil {
			lines = append(lines, "Invalid reference: "+reference)
			continue
		}
		index, found := m.sources[alias]
		if !found {
			lines = append(lines, "Missing source: "+alias)
			continue
		}
		node, found := index.ByPath[path]
		if !found {
			lines = append(lines, "Missing path: "+reference)
			continue
		}
		lines = append(lines, "Source: "+reference)
		lines = append(lines, "# "+node.Heading)
		if body := strings.TrimSpace(node.Body); body != "" {
			lines = append(lines, strings.Split(body, "\n")...)
		}
		for _, child := range node.Children {
			lines = append(lines, "## "+child.Heading+"  <"+alias+":"+child.Path+">")
		}
		lines = append(lines, "")
	}
	return truncateLines(lines, width, height)
}

func (m *Model) addCandidate() {
	if len(m.rows) == 0 {
		return
	}
	selected := m.rows[m.selected]
	candidates := m.candidatesFor(selected)
	if len(candidates) == 0 {
		m.message = "no source child available to add"
		return
	}
	entry := entryAt(m.draft.Doc, selected.Index)
	if entry == nil || len(entry.From) > 0 {
		m.message = "select a manifest section with children before adding"
		return
	}
	next := candidates[0]
	entry.Children = append(entry.Children, manifest.Entry{
		Heading: next.Heading,
		From:    []string{next.Reference},
	})
	m.afterDraftChange("added " + next.Heading)
}

func (m *Model) localizeSelected() {
	if len(m.rows) == 0 {
		return
	}
	selected := m.rows[m.selected]
	if len(selected.From) != 1 {
		m.message = "select one shared source node before localizing"
		return
	}
	alias, path, err := manifest.SplitReference(selected.From[0])
	if err != nil {
		m.err = err.Error()
		return
	}
	if alias == "local" {
		m.message = "selected node is already local"
		return
	}
	index, found := m.sources[alias]
	if !found {
		m.err = fmt.Sprintf("source %q is unavailable", alias)
		return
	}
	node, found := index.ByPath[path]
	if !found {
		m.err = fmt.Sprintf("source %q has no heading path %q", alias, path)
		return
	}
	entry := entryAt(m.draft.Doc, selected.Index)
	if entry == nil {
		return
	}
	content, localNode := localOverrideMarkdown(index, node)
	if m.draft.Sources == nil {
		m.draft.Sources = make(map[string]string)
	}
	if _, exists := m.draft.Sources["local"]; !exists {
		m.draft.Sources["local"] = ".mogent/library"
	}
	entry.From = []string{"local:" + path}
	m.overrides[path+".md"] = content
	m.addLocalNode(localNode)
	m.afterDraftChange("localized " + selected.Entry.Heading)
}

func (m *Model) importDriftIntoSelected() {
	if !m.drift {
		m.message = "no direct AGENTS.md drift detected"
		return
	}
	if len(m.rows) == 0 {
		return
	}
	selected := m.rows[m.selected]
	if len(selected.From) != 1 {
		m.message = "select one source node before importing drift"
		return
	}
	alias, path, err := manifest.SplitReference(selected.From[0])
	if err != nil {
		m.err = err.Error()
		return
	}
	if alias == "local" {
		m.message = "selected node is already local"
		return
	}
	section, err := extractOutputSection(state.OutputPath(m.manifestPath, m.draft.Output), selected.Entry.Heading, selected.Depth+1)
	if err != nil {
		m.err = err.Error()
		return
	}
	index, found := m.sources[alias]
	if !found {
		m.err = fmt.Sprintf("source %q is unavailable", alias)
		return
	}
	node, found := index.ByPath[path]
	if !found {
		m.err = fmt.Sprintf("source %q has no heading path %q", alias, path)
		return
	}
	entry := entryAt(m.draft.Doc, selected.Index)
	if entry == nil {
		return
	}
	content, localNode := localOverrideMarkdownWithBody(index, node, section)
	if m.draft.Sources == nil {
		m.draft.Sources = make(map[string]string)
	}
	if _, exists := m.draft.Sources["local"]; !exists {
		m.draft.Sources["local"] = ".mogent/library"
	}
	entry.From = []string{"local:" + path}
	m.overrides[path+".md"] = content
	m.addLocalNode(localNode)
	m.importedDrift = true
	m.afterDraftChange("imported direct edit into " + selected.Entry.Heading)
}

func (m *Model) removeSelected() {
	if len(m.rows) == 0 {
		return
	}
	selected := m.rows[m.selected]
	if len(selected.Index) == 1 {
		m.message = "top-level entries are not removed in this slice"
		return
	}
	parent := entryAt(m.draft.Doc, selected.Index[:len(selected.Index)-1])
	if parent == nil || len(parent.Children) <= 1 {
		m.message = "cannot leave a section with no children"
		return
	}
	childIndex := selected.Index[len(selected.Index)-1]
	removed := parent.Children[childIndex].Heading
	parent.Children = append(parent.Children[:childIndex], parent.Children[childIndex+1:]...)
	m.selected = max(m.selected-1, 0)
	m.afterDraftChange("removed " + removed)
}

func (m *Model) moveSelected(delta int) {
	if len(m.rows) == 0 {
		return
	}
	selected := m.rows[m.selected]
	if len(selected.Index) == 0 {
		return
	}
	siblings := &m.draft.Doc
	if len(selected.Index) > 1 {
		parent := entryAt(m.draft.Doc, selected.Index[:len(selected.Index)-1])
		if parent == nil {
			return
		}
		siblings = &parent.Children
	}
	from := selected.Index[len(selected.Index)-1]
	to := from + delta
	if to < 0 || to >= len(*siblings) {
		m.message = "cannot reorder past sibling boundary"
		return
	}
	(*siblings)[from], (*siblings)[to] = (*siblings)[to], (*siblings)[from]
	m.selected += delta
	m.afterDraftChange("reordered " + selected.Entry.Heading)
}

func (m *Model) afterDraftChange(message string) {
	m.dirty = true
	m.message = message
	m.err = ""
	if result, err := render.BuildWithSources(m.draft, m.manifestPath, m.sources); err != nil {
		m.err = err.Error()
	} else {
		m.output = result.Content
	}
	m.refreshRows()
}

func (m *Model) discardDraft(message string) {
	m.draft = m.saved.Clone()
	m.overrides = make(map[string]string)
	m.dirty = false
	m.err = ""
	m.message = message
	m.review = ""
	if sources, err := render.LoadSources(m.draft, m.manifestPath); err == nil {
		m.sources = sources
	}
	if result, err := render.BuildWithSources(m.draft, m.manifestPath, m.sources); err == nil {
		m.output = result.Content
	}
	m.refreshRows()
}

func (m *Model) addLocalNode(node *library.Node) {
	index, found := m.sources["local"]
	if !found {
		index = &library.Index{ByPath: make(map[string]*library.Node)}
		m.sources["local"] = index
	}
	index.ByPath[node.Path] = node
	if !containsRoot(index.Roots, node.Path) {
		index.Roots = append(index.Roots, node)
	}
}

func (m Model) writeOverrides() error {
	if len(m.overrides) == 0 {
		return nil
	}
	root := m.draft.Sources["local"]
	if root == "" {
		return fmt.Errorf("local source is missing")
	}
	if !filepath.IsAbs(root) {
		root = filepath.Join(filepath.Dir(m.manifestPath), root)
	}
	for relative, content := range m.overrides {
		path := filepath.Join(root, filepath.FromSlash(relative))
		if !strings.HasPrefix(filepath.Clean(path), filepath.Clean(root)+string(os.PathSeparator)) {
			return fmt.Errorf("local override path escapes local library: %s", relative)
		}
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return fmt.Errorf("create local override directory: %w", err)
		}
		if err := render.WriteAtomically(path, content); err != nil {
			return fmt.Errorf("write local override %q: %w", relative, err)
		}
	}
	return nil
}

func (m Model) saveReview() string {
	result, err := render.BuildWithSources(m.draft, m.manifestPath, m.sources)
	if err != nil {
		return "Save review unavailable: " + err.Error()
	}
	var lines []string
	lines = append(lines, "Save review")
	lines = append(lines, "")
	lines = append(lines, "agents.yaml:")
	currentManifest, readErr := os.ReadFile(m.manifestPath)
	nextManifest, marshalErr := manifest.Marshal(m.draft)
	if readErr != nil {
		lines = append(lines, "  read current manifest: "+readErr.Error())
	} else if marshalErr != nil {
		lines = append(lines, "  render draft manifest: "+marshalErr.Error())
	} else {
		lines = append(lines, simpleDiff(string(currentManifest), string(nextManifest), 8)...)
	}
	lines = append(lines, "")
	lines = append(lines, m.draft.Output+":")
	outputPath := state.OutputPath(m.manifestPath, m.draft.Output)
	currentOutput, readErr := os.ReadFile(outputPath)
	if errors.Is(readErr, os.ErrNotExist) {
		currentOutput = nil
		readErr = nil
	}
	if readErr != nil {
		lines = append(lines, "  read current output: "+readErr.Error())
	} else {
		lines = append(lines, simpleDiff(string(currentOutput), result.Content, 10)...)
	}
	lines = append(lines, "")
	lines = append(lines, "local override files:")
	if len(m.overrides) == 0 {
		lines = append(lines, "  none")
	} else {
		for _, path := range sortedOverridePaths(m.overrides) {
			lines = append(lines, "  .mogent/library/"+path)
		}
	}
	return strings.Join(lines, "\n")
}

func simpleDiff(before, after string, limit int) []string {
	if before == after {
		return []string{"  no changes"}
	}
	beforeLines := strings.Split(strings.TrimRight(before, "\n"), "\n")
	afterLines := strings.Split(strings.TrimRight(after, "\n"), "\n")
	var lines []string
	maxLines := max(len(beforeLines), len(afterLines))
	for index := 0; index < maxLines && len(lines) < limit; index++ {
		var oldLine, newLine string
		if index < len(beforeLines) {
			oldLine = beforeLines[index]
		}
		if index < len(afterLines) {
			newLine = afterLines[index]
		}
		if oldLine == newLine {
			continue
		}
		if index < len(beforeLines) {
			lines = append(lines, "- "+oldLine)
		}
		if len(lines) >= limit {
			break
		}
		if index < len(afterLines) {
			lines = append(lines, "+ "+newLine)
		}
	}
	if len(lines) == 0 {
		return []string{"  changed"}
	}
	if maxLines > len(lines) {
		lines = append(lines, "  ...")
	}
	return lines
}

func sortedOverridePaths(overrides map[string]string) []string {
	paths := make([]string, 0, len(overrides))
	for path := range overrides {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}

func localOverrideMarkdown(index *library.Index, node *library.Node) (string, *library.Node) {
	return localOverrideMarkdownWithBody(index, node, strings.TrimSpace(node.Body))
}

func localOverrideMarkdownWithBody(index *library.Index, node *library.Node, body string) (string, *library.Node) {
	ancestors := ancestorsFor(index, node.Path)
	var output strings.Builder
	for level, ancestor := range ancestors {
		output.WriteString(strings.Repeat("#", level+1))
		output.WriteByte(' ')
		output.WriteString(ancestor.Heading)
		output.WriteString("\n\n")
		if ancestor.Path == node.Path {
			if body := strings.TrimSpace(body); body != "" {
				output.WriteString(body)
				output.WriteString("\n\n")
			}
			if strings.TrimSpace(body) == strings.TrimSpace(node.Body) {
				for _, child := range ancestor.Children {
					output.WriteString(strings.Repeat("#", level+2))
					output.WriteByte(' ')
					output.WriteString(child.Heading)
					output.WriteString("\n\n")
					writeNodeMarkdown(&output, child, level+2)
				}
			}
		}
	}
	localNode := cloneNode(node)
	localNode.Body = strings.TrimSpace(body) + "\n"
	if strings.TrimSpace(body) != strings.TrimSpace(node.Body) {
		localNode.Children = nil
	}
	return strings.TrimSpace(output.String()) + "\n", localNode
}

func writeNodeMarkdown(output *strings.Builder, node *library.Node, level int) {
	if body := strings.TrimSpace(node.Body); body != "" {
		output.WriteString(body)
		output.WriteString("\n\n")
	}
	for _, child := range node.Children {
		output.WriteString(strings.Repeat("#", level+1))
		output.WriteByte(' ')
		output.WriteString(child.Heading)
		output.WriteString("\n\n")
		writeNodeMarkdown(output, child, level+1)
	}
}

func ancestorsFor(index *library.Index, path string) []*library.Node {
	parts := strings.Split(path, "/")
	ancestors := make([]*library.Node, 0, len(parts))
	for i := range parts {
		ancestorPath := strings.Join(parts[:i+1], "/")
		if node, found := index.ByPath[ancestorPath]; found {
			ancestors = append(ancestors, node)
		}
	}
	return ancestors
}

func cloneNode(node *library.Node) *library.Node {
	clone := &library.Node{
		Path:    node.Path,
		Heading: node.Heading,
		Body:    node.Body,
	}
	for _, child := range node.Children {
		clone.Children = append(clone.Children, cloneNode(child))
	}
	return clone
}

func containsRoot(nodes []*library.Node, path string) bool {
	for _, node := range nodes {
		if node.Path == path {
			return true
		}
	}
	return false
}

func hasLocalReference(references []string) bool {
	for _, reference := range references {
		alias, _, err := manifest.SplitReference(reference)
		if err == nil && alias == "local" {
			return true
		}
	}
	return false
}

func extractOutputSection(path, heading string, level int) (string, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read edited output: %w", err)
	}
	lines := strings.Split(string(contents), "\n")
	start := -1
	for index, line := range lines {
		lineLevel, lineHeading, ok := markdownHeading(line)
		if ok && lineLevel == level && lineHeading == heading {
			start = index + 1
			break
		}
	}
	if start < 0 {
		return "", fmt.Errorf("direct edit import could not find heading %q in AGENTS.md", heading)
	}
	end := len(lines)
	for index := start; index < len(lines); index++ {
		lineLevel, _, ok := markdownHeading(lines[index])
		if ok && lineLevel <= level {
			end = index
			break
		}
	}
	return strings.TrimSpace(strings.Join(lines[start:end], "\n")) + "\n", nil
}

func markdownHeading(line string) (int, string, bool) {
	if !strings.HasPrefix(line, "#") {
		return 0, "", false
	}
	level := 0
	for level < len(line) && line[level] == '#' {
		level++
	}
	if level == 0 || level > 6 || level == len(line) || line[level] != ' ' {
		return 0, "", false
	}
	return level, strings.TrimSpace(line[level:]), true
}

func (m Model) candidatesFor(selected row) []candidate {
	if len(selected.Entry.Children) == 0 {
		return nil
	}
	parentRefs := sourceParents(selected.Entry.Children)
	if len(parentRefs) == 0 {
		return nil
	}
	selectedRefs := make(map[string]bool)
	for _, child := range selected.Entry.Children {
		for _, reference := range child.From {
			selectedRefs[reference] = true
		}
	}
	var candidates []candidate
	for _, parentRef := range parentRefs {
		alias, path, err := manifest.SplitReference(parentRef)
		if err != nil {
			continue
		}
		index, found := m.sources[alias]
		if !found {
			continue
		}
		node, found := index.ByPath[path]
		if !found {
			continue
		}
		for _, child := range node.Children {
			reference := alias + ":" + child.Path
			if selectedRefs[reference] {
				continue
			}
			candidates = append(candidates, candidate{Heading: child.Heading, Reference: reference})
		}
	}
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].Reference < candidates[j].Reference
	})
	return candidates
}

func sourceParents(entries []manifest.Entry) []string {
	seen := make(map[string]bool)
	var parents []string
	for _, entry := range entries {
		for _, reference := range entry.From {
			alias, path, err := manifest.SplitReference(reference)
			if err != nil {
				continue
			}
			parts := strings.Split(path, "/")
			if len(parts) < 2 {
				continue
			}
			parentRef := alias + ":" + strings.Join(parts[:len(parts)-1], "/")
			if !seen[parentRef] {
				parents = append(parents, parentRef)
				seen[parentRef] = true
			}
		}
	}
	return parents
}

func entryAt(entries []manifest.Entry, index []int) *manifest.Entry {
	if len(index) == 0 {
		return nil
	}
	currentEntries := entries
	var current *manifest.Entry
	for _, position := range index {
		if position < 0 || position >= len(currentEntries) {
			return nil
		}
		current = &currentEntries[position]
		currentEntries = current.Children
	}
	return current
}

func (m *Model) setCollapsed(collapsed bool) {
	if len(m.rows) == 0 {
		return
	}
	current := m.rows[m.selected]
	if len(current.Entry.Children) == 0 {
		return
	}
	m.collapsed[rowKey(current.Index)] = collapsed
	m.refreshRows()
}

func (m *Model) nextFocus() {
	if m.width >= 140 {
		m.focus = (m.focus + 1) % 3
		return
	}
	if m.focus == focusTree {
		if m.detail == detailFinal {
			m.focus = focusFinal
		} else {
			m.focus = focusSource
		}
		return
	}
	m.focus = focusTree
}

func rowKey(index []int) string {
	parts := make([]string, 0, len(index))
	for _, value := range index {
		parts = append(parts, fmt.Sprint(value))
	}
	return strings.Join(parts, ".")
}

func firstHighlighted(lines []string) int {
	for index, line := range lines {
		if strings.HasPrefix(line, "> ") {
			return index
		}
	}
	return 0
}

func visibleWindow(lines []string, selected, height int) string {
	if height <= 0 {
		return ""
	}
	if len(lines) <= height {
		return strings.Join(lines, "\n")
	}
	start := selected - height/2
	if start < 0 {
		start = 0
	}
	if start+height > len(lines) {
		start = len(lines) - height
	}
	return strings.Join(lines[start:start+height], "\n")
}

func truncateLines(lines []string, width, height int) string {
	for index := range lines {
		lines[index] = truncate(lines[index], width)
	}
	return visibleWindow(lines, 0, height)
}

func limitLines(value string, height int) string {
	lines := strings.Split(value, "\n")
	if len(lines) <= height {
		return value
	}
	return strings.Join(lines[:height], "\n")
}

func truncate(value string, width int) string {
	if width <= 0 || lipgloss.Width(value) <= width {
		return value
	}
	if width <= 3 {
		return strings.Repeat(".", max(width, 0))
	}
	runes := []rune(value)
	for len(runes) > 0 && lipgloss.Width(string(runes)+"...") > width {
		runes = runes[:len(runes)-1]
	}
	return string(runes) + "..."
}
