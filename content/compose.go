package content

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"sort"
	"strings"
)

type Operation string

const (
	Copy           Operation = "copy"
	RenderMarkdown Operation = "render-markdown"
)

// Contribution explicitly selects, transforms, and places one physical node or
// document view. Rendered adapts an already authored document contribution.
type Contribution struct {
	Inventory *Inventory
	Node      string
	Heading   []string
	Rendered  *Rendered
	Operation Operation
	Into      string
	Origin    string
	Exclude   []string
	Replace   bool
	Append    bool
	Level     int
}

type Payload struct {
	Bytes      []byte
	Mode       os.FileMode
	Origin     string
	Rendered   bool
	References []string
}

// Event retains contribution history, including equal-byte replacements.
type Event struct {
	Action    string
	Path      string
	Origin    string
	Previous  string
	Operation Operation
}

// Manifest is the complete materialized result in memory. Copy operations
// capture bytes before installation so staging cannot silently reread new data.
type Manifest struct {
	Files       map[string]Payload
	Directories map[string]os.FileMode
	Bundles     map[string]string
	Events      []Event
}

// Compose creates one final manifest without writing anything. Directory union
// is allowed; file collisions and overlapping skill bundles require Replace.
// Append is explicitly restricted to rendered document contributions.
func Compose(contributions []Contribution) (*Manifest, error) {
	m := &Manifest{Files: make(map[string]Payload), Directories: map[string]os.FileMode{".": 0755}, Bundles: make(map[string]string)}
	for index, c := range contributions {
		if err := m.contribute(c); err != nil {
			return nil, fmt.Errorf("contribution %d (%s): %w", index+1, c.Origin, err)
		}
	}
	for destination, payload := range m.Files {
		for _, reference := range payload.References {
			u, err := url.Parse(reference)
			if err != nil {
				return nil, fmt.Errorf("rendered %q has invalid reference %q: %w", destination, reference, err)
			}
			if u.Scheme != "" || u.Host != "" || u.Path == "" || strings.HasPrefix(u.Path, "/") {
				continue
			}
			resolved := path.Clean(path.Join(path.Dir(destination), u.Path))
			if err := ValidatePath(resolved, true); err != nil {
				return nil, fmt.Errorf("rendered %q reference %q escapes the output", destination, reference)
			}
			_, file := m.Files[resolved]
			_, dir := m.Directories[resolved]
			if !file && !dir {
				return nil, fmt.Errorf("rendered %q reference %q has no payload at %q; copy the dependency explicitly", destination, reference, resolved)
			}
		}
	}
	return m, nil
}

func (m *Manifest) contribute(c Contribution) error {
	if err := ValidatePath(c.Into, true); err != nil {
		return err
	}
	if c.Replace && c.Append {
		return fmt.Errorf("replace and append cannot be combined")
	}
	for boundary, origin := range m.Bundles {
		if within(c.Into, boundary) && c.Into != boundary {
			return fmt.Errorf("destination %q is inside skill bundle from %s; replace the complete bundle", c.Into, origin)
		}
	}
	if c.Replace {
		m.remove(c.Into, c.Origin, c.Operation)
	}
	protected := make(map[string]string, len(m.Bundles))
	for boundary, origin := range m.Bundles {
		protected[boundary] = origin
	}
	switch c.Operation {
	case Copy:
		if c.Append || len(c.Heading) > 0 || c.Rendered != nil {
			return fmt.Errorf("copy cannot append or select a heading view")
		}
		if c.Inventory == nil {
			return fmt.Errorf("copy requires a physical inventory")
		}
		node := c.Inventory.Nodes[c.Node]
		if node == nil {
			return fmt.Errorf("node %q does not exist", c.Node)
		}
		if node.Bundle != "" && node.Bundle != node.Path {
			return fmt.Errorf("node %q is a partial skill; select complete bundle %q", node.Path, node.Bundle)
		}
		if node.Kind == File && (c.Into == "." || len(c.Exclude) > 0) {
			return fmt.Errorf("file copy requires a file destination and cannot exclude descendants")
		}
		for _, excluded := range c.Exclude {
			if err := ValidatePath(excluded, false); err != nil {
				return err
			}
			if c.Inventory.Nodes[path.Join(c.Node, excluded)] == nil {
				return fmt.Errorf("excluded path %q does not exist under node %q", excluded, c.Node)
			}
			m.Events = append(m.Events, Event{Action: "exclude-source", Path: path.Join(c.Into, excluded), Origin: c.Origin, Operation: c.Operation})
		}
		for _, p := range c.Inventory.Paths() {
			if !within(p, c.Node) {
				continue
			}
			relative := relativeTo(p, c.Node)
			if excludedPath(relative, c.Exclude) {
				continue
			}
			selected := c.Inventory.Nodes[p]
			if selected.Role == "skill" {
				for _, member := range c.Inventory.Paths() {
					if within(member, p) && excludedPath(relativeTo(member, c.Node), c.Exclude) {
						return fmt.Errorf("exclude would make skill bundle %q incomplete; exclude the entire bundle instead", p)
					}
				}
				bundlePath := path.Join(c.Into, relative)
				if previous, found := m.Bundles[bundlePath]; found {
					return fmt.Errorf("skill bundle collision at %q with %s; request explicit replacement", bundlePath, previous)
				}
				// A new skill cannot absorb pre-existing loose members silently.
				for earlier := range m.Files {
					if within(earlier, bundlePath) {
						return fmt.Errorf("skill bundle %q overlaps earlier payload %q; request explicit replacement", bundlePath, earlier)
					}
				}
				m.Bundles[bundlePath] = c.Origin
			}
			destination := path.Join(c.Into, relative)
			for boundary, origin := range protected {
				if within(destination, boundary) {
					return fmt.Errorf("payload %q would modify skill bundle from %s; replace the complete bundle", destination, origin)
				}
			}
			if selected.Kind == Directory {
				if err := m.directory(destination, selected.Mode); err != nil {
					return err
				}
				continue
			}
			info, err := os.Lstat(selected.Source)
			if err != nil {
				return err
			}
			if !info.Mode().IsRegular() {
				return fmt.Errorf("source %q is no longer a regular file", selected.Source)
			}
			data, err := os.ReadFile(selected.Source)
			if err != nil {
				return err
			}
			if err := m.put(destination, Payload{Bytes: data, Mode: selected.Mode, Origin: c.Origin}, false); err != nil {
				return err
			}
		}
	case RenderMarkdown:
		if c.Into == "." || len(c.Exclude) > 0 {
			return fmt.Errorf("render-markdown requires a file destination and cannot exclude physical descendants")
		}
		var rendered Rendered
		if c.Rendered != nil {
			rendered = *c.Rendered
		} else {
			if c.Inventory == nil {
				return fmt.Errorf("render-markdown requires a document view")
			}
			doc, err := c.Inventory.Document(c.Node)
			if err != nil {
				return err
			}
			level := c.Level
			if level == 0 {
				level = 1
			}
			result, err := doc.Render(c.Heading, level)
			if err != nil {
				return err
			}
			rendered = result
		}
		if err := m.put(c.Into, Payload{Bytes: append([]byte(nil), rendered.Bytes...), Mode: 0644, Origin: c.Origin, Rendered: true, References: append([]string(nil), rendered.References...)}, c.Append); err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown operation %q", c.Operation)
	}
	return nil
}

func (m *Manifest) directory(destination string, mode os.FileMode) error {
	if _, found := m.Files[destination]; found {
		return fmt.Errorf("file/directory collision at %q", destination)
	}
	if destination != "." {
		if err := m.directory(path.Dir(destination), 0755); err != nil {
			return err
		}
	}
	if _, found := m.Directories[destination]; !found {
		m.Directories[destination] = mode
	}
	return nil
}

func (m *Manifest) put(destination string, payload Payload, appendText bool) error {
	if err := ValidatePath(destination, false); err != nil {
		return err
	}
	if _, found := m.Directories[destination]; found {
		return fmt.Errorf("file/directory collision at %q", destination)
	}
	if err := m.directory(path.Dir(destination), 0755); err != nil {
		return err
	}
	if previous, found := m.Files[destination]; found {
		if !appendText || !previous.Rendered || !payload.Rendered {
			return fmt.Errorf("file collision at %q (%s and %s); request explicit replacement", destination, previous.Origin, payload.Origin)
		}
		joined := strings.TrimRight(string(previous.Bytes), "\r\n") + "\n\n" + string(payload.Bytes)
		payload.Bytes = []byte(joined)
		payload.References = append(append([]string(nil), previous.References...), payload.References...)
		m.Events = append(m.Events, Event{Action: "append", Path: destination, Origin: payload.Origin, Previous: previous.Origin, Operation: RenderMarkdown})
	} else {
		operation := Copy
		if payload.Rendered {
			operation = RenderMarkdown
		}
		m.Events = append(m.Events, Event{Action: "add", Path: destination, Origin: payload.Origin, Operation: operation})
	}
	m.Files[destination] = payload
	return nil
}

func (m *Manifest) remove(destination, origin string, operation Operation) {
	paths := make([]string, 0, len(m.Files))
	for p := range m.Files {
		paths = append(paths, p)
	}
	sort.Strings(paths)
	for _, p := range paths {
		payload := m.Files[p]
		if within(p, destination) {
			m.Events = append(m.Events, Event{Action: "replace-remove", Path: p, Origin: origin, Previous: payload.Origin, Operation: operation})
			delete(m.Files, p)
		}
	}
	for p := range m.Directories {
		if p != "." && within(p, destination) {
			delete(m.Directories, p)
		}
	}
	for p := range m.Bundles {
		if within(p, destination) {
			delete(m.Bundles, p)
		}
	}
}

func within(value, root string) bool {
	return root == "." || value == root || strings.HasPrefix(value, root+"/")
}
func relativeTo(value, root string) string {
	if value == root {
		return "."
	}
	if root == "." {
		return value
	}
	return strings.TrimPrefix(value, root+"/")
}
func excludedPath(value string, excludes []string) bool {
	for _, excluded := range excludes {
		if within(value, excluded) {
			return true
		}
	}
	return false
}
