// Package navigator implements the milestone-two terminal manifest browser.
package navigator

import (
	"fmt"
	"io"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/Qu1ncyRy4n/Agents/manifest"
	"github.com/Qu1ncyRy4n/Agents/workspace"
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
	session     *workspace.Session
	rows        []row
	selected    int
	collapsed   map[string]bool
	focus       focus
	detail      detailMode
	width       int
	height      int
	err         string
	message     string
	confirmSave bool
}

// New loads a manifest and prepares an in-memory draft.
func New(manifestPath string) (Model, error) {
	session, err := workspace.New(manifestPath)
	if err != nil {
		return Model{}, err
	}
	model := Model{
		session:   session,
		collapsed: make(map[string]bool),
		focus:     focusTree,
		detail:    detailFinal,
		width:     100,
		height:    30,
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
				if err := m.session.SaveAndBuild(); err != nil {
					m.err = err.Error()
					m.message = ""
				} else {
					m.err = ""
					m.message = "saved and built"
					m.refreshRows()
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
	if m.session.Dirty {
		stateText = "state: ~ draft"
	}
	parts := []string{
		fmt.Sprintf("manifest: %s", filepath.Base(m.session.ManifestPath)),
		stateText,
		"keys: up/down move, a add, d remove, U/D reorder, tab focus, 1/2/3 views, v detail, s save, q quit",
	}
	if m.confirmSave {
		parts = append(parts, "confirm: y/n")
	}
	if m.message != "" {
		parts = append(parts, m.message)
	}
	if m.err != "" {
		parts = append(parts, "error: "+m.err)
	}
	return truncate(strings.Join(parts, " | "), width)
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
	for index := range m.session.Draft.Doc {
		m.appendEntry(&m.session.Draft.Doc[index], 0, []int{index}, nil)
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
		provenance := "children"
		if len(row.From) > 0 {
			provenance = strings.Join(row.From, ", ")
		}
		line := fmt.Sprintf("%s%s%s[x] %s  <%s>", prefix, strings.Repeat("  ", row.Depth), collapse, row.Entry.Heading, provenance)
		lines = append(lines, truncate(line, width))
	}
	return visibleWindow(lines, m.selected, height)
}

func (m Model) finalView(width, height int) string {
	current := m.rows[m.selected].Entry.Heading
	lines := strings.Split(strings.TrimRight(m.session.Output, "\n"), "\n")
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
		index, found := m.session.Sources[alias]
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
	entry := entryAt(m.session.Draft.Doc, selected.Index)
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

func (m *Model) removeSelected() {
	if len(m.rows) == 0 {
		return
	}
	selected := m.rows[m.selected]
	if len(selected.Index) == 1 {
		m.message = "top-level entries are not removed in this slice"
		return
	}
	parent := entryAt(m.session.Draft.Doc, selected.Index[:len(selected.Index)-1])
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
	siblings := &m.session.Draft.Doc
	if len(selected.Index) > 1 {
		parent := entryAt(m.session.Draft.Doc, selected.Index[:len(selected.Index)-1])
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
	m.message = message
	m.err = ""
	if err := m.session.MarkDraftChanged(); err != nil {
		m.err = err.Error()
	}
	m.refreshRows()
}

func (m *Model) discardDraft(message string) {
	m.err = ""
	m.message = message
	if err := m.session.DiscardDraft(); err != nil {
		m.err = err.Error()
	}
	m.refreshRows()
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
		index, found := m.session.Sources[alias]
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
