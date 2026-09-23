package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/Qu1ncyRy4n/Agents/library"
	"github.com/Qu1ncyRy4n/Agents/render"
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
	Alias          string
	Path           string
	Nodes          int
	DirectoryTree  bool
	Files          int
	DirectoryTrees []DirectoryTreeStatus
}

// DirectoryTreeStatus describes a raw copied subtree declared by a root sidecar.
type DirectoryTreeStatus struct {
	Path  string
	Files int
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
	outputs := s.Draft.EffectiveOutputs()
	if len(outputs) == 0 {
		return ""
	}
	outputPath := outputs[0].Path
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
		Sources:      nil,
	}
	status.Sources, err = s.sourceStatuses()
	if err != nil {
		return nil, err
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
		if !output.spec.Directory() && current == state.OutputClean && len(output.spec.Include) > 0 {
			wanted, renderErr := render.RenderOutput(s.Draft, s.ManifestPath, output.spec)
			actual, readErr := os.ReadFile(output.path)
			if renderErr != nil || readErr != nil || string(actual) != wanted.Content {
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

func (s *Session) sourceStatuses() ([]SourceStatus, error) {
	rawDirectories := s.Saved.RawDirectorySourceAliases()
	aliases := make([]string, 0, len(s.Saved.Sources))
	for alias := range s.Saved.Sources {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	sources := make([]SourceStatus, 0, len(aliases))
	for _, alias := range aliases {
		status := SourceStatus{
			Alias: alias,
			Path:  s.Saved.Sources[alias].Display(),
		}
		path, err := render.ResolveSourcePath(s.Saved, s.ManifestPath, alias)
		if err != nil {
			return nil, fmt.Errorf("resolve source %q: %w", alias, err)
		}
		if rawDirectories[alias] {
			files, err := regularFileCount(path)
			if err != nil {
				return nil, fmt.Errorf("count raw directory source %q: %w", alias, err)
			}
			status.DirectoryTree, status.Files = true, files
		} else if index, found := s.Sources[alias]; found {
			status.Nodes = len(index.ByPath)
		}
		sidecar, _, err := library.LoadSidecar(path)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("load source sidecar %q: %w", alias, err)
		}
		if sidecar != nil && sidecar.Content != nil {
			for _, directory := range sidecar.Content.DirectoryRoots {
				files, err := regularFileCount(filepath.Join(path, filepath.FromSlash(directory)))
				if err != nil {
					return nil, fmt.Errorf("count declared directory root %q for source %q: %w", directory, alias, err)
				}
				status.DirectoryTrees = append(status.DirectoryTrees, DirectoryTreeStatus{Path: directory, Files: files})
			}
		}
		sources = append(sources, status)
	}
	return sources, nil
}

func regularFileCount(root string) (int, error) {
	count := 0
	err := filepath.WalkDir(root, func(_ string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type().IsRegular() {
			count++
		}
		return nil
	})
	return count, err
}
