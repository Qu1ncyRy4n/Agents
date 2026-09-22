package workspace

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Qu1ncyRy4n/Agents/state"
)

// OutputStatus is the user-facing state of the generated document.
type OutputStatus string

const (
	StatusMissing     OutputStatus = "missing"
	StatusUntracked   OutputStatus = "untracked"
	StatusDirectEdits OutputStatus = "direct edits"
	StatusStale       OutputStatus = "stale"
	StatusUpToDate    OutputStatus = "up to date"
)

type SourceStatus struct {
	Alias string
	Path  string
	Nodes int
}

type Status struct {
	ManifestPath string
	OutputPath   string
	Output       OutputStatus
	Outputs      []NamedOutputStatus
	Sources      []SourceStatus
}

type NamedOutputStatus struct {
	Path   string
	Status OutputStatus
}

// OutputPath returns the absolute generated document path for the current
// draft.
func (s *Session) OutputPath() string {
	outputPath := s.Draft.Output
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(filepath.Dir(s.ManifestPath), outputPath)
	}
	return outputPath
}

// StatePath returns the local generated-output state path for this workspace.
func (s *Session) StatePath() string {
	return filepath.Join(filepath.Dir(s.ManifestPath), ".mogent", "state.json")
}

// Status summarizes the local workspace without writing any files.
func (s *Session) Status() (*Status, error) {
	outputPath := s.OutputPath()
	outputState, err := state.Inspect(outputPath, s.StatePath())
	if err != nil {
		return nil, err
	}
	status := &Status{
		ManifestPath: s.ManifestPath,
		OutputPath:   outputPath,
		Output:       s.outputStatus(outputState, outputPath),
		Sources:      s.sourceStatuses(),
	}
	outputs, err := configuredOutputs(s.Draft, s.ManifestPath)
	if err != nil {
		return nil, err
	}
	for _, output := range outputs[1:] {
		current := state.OutputMissing
		if output.spec.Directory() {
			current, err = state.InspectDirectory(output.path, s.StatePath())
		} else {
			current, err = state.Inspect(output.path, s.StatePath())
		}
		if err != nil {
			return nil, err
		}
		value := s.outputStatus(current, output.path)
		if output.spec.Directory() {
			switch current {
			case state.OutputMissing:
				value = StatusMissing
			case state.OutputUntracked:
				value = StatusUntracked
			case state.OutputModified:
				value = StatusDirectEdits
			case state.OutputClean:
				value = StatusUpToDate
			}
		}
		if output.spec.Directory() && current == state.OutputClean {
			wanted, hashErr := state.DirectoryHashes(output.source)
			actual, actualErr := state.DirectoryHashes(output.path)
			if hashErr != nil || actualErr != nil || !sameHashes(wanted, actual) {
				value = StatusStale
			}
		}
		status.Outputs = append(status.Outputs, NamedOutputStatus{Path: output.path, Status: value})
	}
	return status, nil
}

func sameHashes(first, second map[string]string) bool {
	if len(first) != len(second) {
		return false
	}
	for path, hash := range first {
		if second[path] != hash {
			return false
		}
	}
	return true
}

func (s *Session) outputStatus(outputState state.OutputState, outputPath string) OutputStatus {
	switch outputState {
	case state.OutputMissing:
		return StatusMissing
	case state.OutputUntracked:
		return StatusUntracked
	case state.OutputModified:
		return StatusDirectEdits
	case state.OutputClean:
		output, err := os.ReadFile(outputPath)
		if err != nil {
			return StatusStale
		}
		if string(output) == s.Output {
			return StatusUpToDate
		}
		return StatusStale
	default:
		return OutputStatus(fmt.Sprintf("unknown: %s", outputState))
	}
}

func (s *Session) sourceStatuses() []SourceStatus {
	aliases := make([]string, 0, len(s.Saved.Sources))
	for alias := range s.Saved.Sources {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	sources := make([]SourceStatus, 0, len(aliases))
	for _, alias := range aliases {
		nodes := 0
		if index, found := s.Sources[alias]; found {
			nodes = len(index.ByPath)
		}
		sources = append(sources, SourceStatus{
			Alias: alias,
			Path:  s.Saved.Sources[alias].Display(),
			Nodes: nodes,
		})
	}
	return sources
}
