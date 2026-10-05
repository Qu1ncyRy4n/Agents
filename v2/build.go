package v2

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"

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
// compiled but not yet written, so they produce no FileChange. Pins lists the
// commits apply writes into mogent.hcl for sources that had none.
type Result struct {
	Plan      *Plan
	Libraries map[string]*Library
	Sources   map[string]ResolvedSource
	Rendered  map[string]string
	Changes   []FileChange
	Pins      map[string]string
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

// PlanConfig resolves sources, compiles the selection, renders in memory, and
// diffs every Markdown output path against disk. It writes nothing except a
// source checkout fetched into .mogent/sources/. Drift from the recorded
// generated-output state is reported as MOGENT208.
func PlanConfig(config *Config) (*Result, error) {
	result, err := render(config)
	if err != nil || result.HasErrors() {
		return result, err
	}
	root := filepath.Dir(config.Path)
	for _, output := range result.Plan.Outputs {
		content, found := result.Rendered[output.Name]
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
	result.Pins = make(map[string]string)
	for alias, source := range result.Sources {
		if source.Commit != "" && !source.Pinned {
			result.Pins[alias] = source.Commit
		}
	}
	return result, nil
}

func render(config *Config) (*Result, error) {
	sources, diagnostics, err := ResolveSources(config)
	if err != nil {
		return nil, err
	}
	libraries, err := LoadLibraries(sources)
	if err != nil {
		return nil, err
	}
	plan := Compile(config, libraries)
	plan.Diagnostics = append(diagnostics, plan.Diagnostics...)
	result := &Result{
		Plan:      plan,
		Libraries: libraries,
		Sources:   sources,
		statePath: filepath.Join(filepath.Dir(config.Path), ".mogent", "state.json"),
	}
	if result.HasErrors() {
		return result, nil
	}
	rendered, err := RenderMarkdown(plan, libraries)
	if err != nil {
		return result, err
	}
	result.Rendered = rendered
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

// Apply runs PlanConfig and writes every planned Markdown output, its state,
// and any newly resolved commit pins transactionally. An unmanaged or directly
// edited output is refused unless force is set.
func Apply(config *Config, force bool) (*Result, error) {
	result, err := PlanConfig(config)
	if err != nil {
		return result, err
	}
	if result.HasErrors() {
		return result, fmt.Errorf("plan contains errors")
	}
	writes := make(map[string]string, len(result.Changes)+1)
	for _, change := range result.Changes {
		if err := state.CheckOverwrite(change.Absolute, result.statePath, force); err != nil {
			return result, err
		}
		writes[change.Absolute] = change.Content
	}
	if len(result.Pins) > 0 {
		pinned, err := SetCommits(config.Path, result.Pins)
		if err != nil {
			return result, err
		}
		writes[config.Path] = string(pinned)
	}
	snapshot, err := snapshotFiles(writes, result.statePath)
	if err != nil {
		return result, err
	}
	for path, content := range writes {
		if err := renderfs.WriteAtomically(path, []byte(content)); err != nil {
			return result, snapshot.restore(err)
		}
	}
	for _, change := range result.Changes {
		if err := state.Write(result.statePath, change.Absolute, change.Content); err != nil {
			return result, snapshot.restore(err)
		}
	}
	return result, nil
}

// UpdateResult describes one git source after re-resolving its ref.
type UpdateResult struct {
	Alias       string
	OldCommit   string
	NewCommit   string
	Changes     []FileChange
	Diagnostics []Diagnostic
}

// Changed reports whether the ref now resolves to a different commit.
func (u UpdateResult) Changed() bool { return u.OldCommit != u.NewCommit }

// Update re-resolves ref for one git source, or every git source when alias
// is empty, and renders the configured outputs under the new commit so the
// caller can show the content diff. Each source is evaluated against the
// current configuration independently. With accept, every changed commit is
// written into mogent.hcl; nothing else is written.
func Update(config *Config, alias string, accept bool) ([]UpdateResult, error) {
	var aliases []string
	for _, name := range sortedSourceNames(config) {
		if config.Sources[name].Git == "" {
			continue
		}
		if alias == "" || name == alias {
			aliases = append(aliases, name)
		}
	}
	if alias != "" && len(aliases) == 0 {
		return nil, fmt.Errorf("source %q is not a declared git source", alias)
	}
	if len(aliases) == 0 {
		return nil, fmt.Errorf("no git sources to update")
	}
	for _, name := range aliases {
		if config.Sources[name].Commit == "" {
			return nil, fmt.Errorf("source %q is unpinned; run mogent apply to pin it before update", name)
		}
	}
	current, err := render(config)
	if err != nil {
		return nil, err
	}
	if current.HasErrors() {
		return nil, fmt.Errorf("current configuration has plan errors; fix them before update")
	}
	base := filepath.Dir(config.Path)
	results := make([]UpdateResult, 0, len(aliases))
	pins := make(map[string]string)
	for _, name := range aliases {
		source := config.Sources[name]
		ref := source.Ref
		if ref == "" {
			ref = "HEAD"
		}
		commit, err := fetchGit(source.Git, ref, filepath.Join(base, ".mogent", "sources", name))
		if err != nil {
			return results, fmt.Errorf("source %q: %w", name, err)
		}
		update := UpdateResult{Alias: name, OldCommit: source.Commit, NewCommit: commit}
		if !update.Changed() {
			results = append(results, update)
			continue
		}
		candidate := *config
		candidate.Sources = make(map[string]Source, len(config.Sources))
		for key, value := range config.Sources {
			candidate.Sources[key] = value
		}
		source.Commit = commit
		candidate.Sources[name] = source
		next, err := render(&candidate)
		if err != nil {
			return results, fmt.Errorf("source %q at %s: %w", name, commit, err)
		}
		update.Diagnostics = next.Plan.Diagnostics
		if !next.HasErrors() {
			for _, output := range current.Plan.Outputs {
				before, after := current.Rendered[output.Name], next.Rendered[output.Name]
				for _, relative := range output.Paths {
					change := FileChange{Output: output.Name, Path: relative, Status: FileUnchanged, Content: after}
					if diff := textdiff.Unified(before, after, "a/"+relative, "b/"+relative); diff != "" {
						change.Status, change.Diff = FileChanged, diff
					}
					update.Changes = append(update.Changes, change)
				}
			}
			pins[name] = commit
		}
		results = append(results, update)
	}
	if accept && len(pins) > 0 {
		pinned, err := SetCommits(config.Path, pins)
		if err != nil {
			return results, err
		}
		if err := renderfs.WriteAtomically(config.Path, pinned); err != nil {
			return results, err
		}
	}
	sort.Slice(results, func(i, j int) bool { return results[i].Alias < results[j].Alias })
	return results, nil
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
			return nil, fmt.Errorf("read %q before apply: %w", path, err)
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
