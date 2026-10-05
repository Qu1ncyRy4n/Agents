package v2

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/state"
	"github.com/Qu1ncyRy4n/Agents/textdiff"
)

// FileStatus classifies one planned output path against the file on disk.
type FileStatus string

const (
	FileNew       FileStatus = "new"
	FileChanged   FileStatus = "changed"
	FileUnchanged FileStatus = "unchanged"
)

// FileChange is one output path the plan would write, with the diff from the
// current file to the rendered content.
type FileChange struct {
	Output   string
	Path     string
	Absolute string
	Status   FileStatus
	Diff     string
	Content  string
}

// Result is a resolved plan plus the file changes it implies. Tree outputs are
// compiled but not yet written, so they produce no FileChange.
type Result struct {
	Plan      *Plan
	Libraries map[string]*Library
	Changes   []FileChange
	statePath string
}

// HasErrors reports whether any diagnostic blocks apply.
func (r *Result) HasErrors() bool {
	for _, diagnostic := range r.Plan.Diagnostics {
		if diagnostic.Severity == SeverityError {
			return true
		}
	}
	return false
}

// PlanLocal loads local libraries, compiles the selection, renders in memory,
// and diffs every Markdown output path against disk. It writes nothing. Drift
// from the recorded generated-output state is reported as MOGENT208.
func PlanLocal(config *Config) (*Result, error) {
	libraries, err := LoadLocalLibraries(config)
	if err != nil {
		return nil, err
	}
	plan := Compile(config, libraries)
	root := filepath.Dir(config.Path)
	result := &Result{Plan: plan, Libraries: libraries, statePath: filepath.Join(root, ".mogent", "state.json")}
	if result.HasErrors() {
		return result, nil
	}
	rendered, err := RenderMarkdown(plan, libraries)
	if err != nil {
		return result, err
	}
	for _, output := range plan.Outputs {
		content, found := rendered[output.Name]
		if !found {
			continue
		}
		for _, relative := range output.Paths {
			change, err := inspectChange(result, output.Name, relative, filepath.Join(root, filepath.FromSlash(relative)), content)
			if err != nil {
				return result, err
			}
			result.Changes = append(result.Changes, change)
		}
	}
	return result, nil
}

func inspectChange(result *Result, outputName, relative, absolute, content string) (FileChange, error) {
	change := FileChange{Output: outputName, Path: relative, Absolute: absolute, Content: content}
	existing, err := os.ReadFile(absolute)
	switch {
	case errors.Is(err, os.ErrNotExist):
		change.Status = FileNew
		change.Diff = textdiff.Unified("", content, "/dev/null", "b/"+relative)
		return change, nil
	case err != nil:
		return change, fmt.Errorf("read output %q: %w", relative, err)
	}
	change.Diff = textdiff.Unified(string(existing), content, "a/"+relative, "b/"+relative)
	change.Status = FileChanged
	if change.Diff == "" {
		change.Status = FileUnchanged
	}
	outputState, err := state.Inspect(absolute, result.statePath)
	if err != nil {
		return change, err
	}
	switch outputState {
	case state.OutputModified:
		result.Plan.Diagnostics = append(result.Plan.Diagnostics, Diagnostic{
			Severity: SeverityWarning,
			Code:     "MOGENT208",
			Message:  fmt.Sprintf("output %q was edited after the last apply; apply requires --force", relative),
		})
	case state.OutputUntracked:
		result.Plan.Diagnostics = append(result.Plan.Diagnostics, Diagnostic{
			Severity: SeverityWarning,
			Code:     "MOGENT208",
			Message:  fmt.Sprintf("output %q exists but is not managed by Mogent; apply requires --force", relative),
		})
	}
	return change, nil
}

// Apply runs PlanLocal and writes every planned Markdown output and its state
// transactionally. An unmanaged or directly edited output is refused unless
// force is set.
func Apply(config *Config, force bool) (*Result, error) {
	result, err := PlanLocal(config)
	if err != nil {
		return result, err
	}
	if result.HasErrors() {
		return result, fmt.Errorf("plan contains errors")
	}
	files := make(map[string]string, len(result.Changes))
	for _, change := range result.Changes {
		if err := state.CheckOverwrite(change.Absolute, result.statePath, force); err != nil {
			return result, err
		}
		files[change.Absolute] = change.Content
	}
	snapshot, err := snapshotFiles(files, result.statePath)
	if err != nil {
		return result, err
	}
	for path, content := range files {
		if err := renderfs.WriteAtomically(path, []byte(content)); err != nil {
			return result, snapshot.restore(err)
		}
	}
	for path, content := range files {
		if err := state.Write(result.statePath, path, content); err != nil {
			return result, snapshot.restore(err)
		}
	}
	return result, nil
}

type fileSnapshot struct {
	files        map[string][]byte
	existed      map[string]bool
	state        []byte
	stateExisted bool
	statePath    string
}

func snapshotFiles(files map[string]string, statePath string) (*fileSnapshot, error) {
	snapshot := &fileSnapshot{files: make(map[string][]byte), existed: make(map[string]bool), statePath: statePath}
	for path := range files {
		snapshot.existed[path] = false
		contents, err := os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("read output %q before apply: %w", path, err)
		}
		snapshot.files[path] = contents
		snapshot.existed[path] = true
	}
	contents, err := os.ReadFile(statePath)
	if errors.Is(err, os.ErrNotExist) {
		return snapshot, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read generated-output state: %w", err)
	}
	snapshot.state, snapshot.stateExisted = contents, true
	return snapshot, nil
}

func (s *fileSnapshot) restore(original error) error {
	var errs []error
	for path, existed := range s.existed {
		if !existed {
			if err := os.Remove(path); err != nil && !errors.Is(err, os.ErrNotExist) {
				errs = append(errs, err)
			}
			continue
		}
		if err := renderfs.WriteAtomically(path, s.files[path]); err != nil {
			errs = append(errs, err)
		}
	}
	if s.stateExisted {
		if err := renderfs.WriteAtomicallyMode(s.statePath, s.state, 0o600); err != nil {
			errs = append(errs, err)
		}
	} else if err := os.Remove(s.statePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		errs = append(errs, err)
	}
	return errors.Join(append([]error{original}, errs...)...)
}
