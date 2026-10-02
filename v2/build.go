package v2

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/state"
)

// BuildLocal compiles, renders, and transactionally writes local-library v2
// Markdown outputs. Remote sources and tree outputs remain unsupported here.
func BuildLocal(config *Config, force, dryRun bool) (*Plan, error) {
	libraries, err := LoadLocalLibraries(config)
	if err != nil {
		return nil, err
	}
	plan := Compile(config, libraries)
	for _, diagnostic := range plan.Diagnostics {
		if diagnostic.Severity == SeverityError {
			return plan, fmt.Errorf("plan contains errors")
		}
	}
	rendered, err := RenderMarkdown(plan, libraries)
	if err != nil {
		return plan, err
	}
	if dryRun {
		return plan, nil
	}
	root := filepath.Dir(config.Path)
	statePath := filepath.Join(root, ".mogent", "state.json")
	files := make(map[string]string)
	for _, output := range plan.Outputs {
		content, found := rendered[output.Name]
		if !found {
			continue
		}
		for _, relative := range output.Paths {
			path := filepath.Join(root, filepath.FromSlash(relative))
			if err := state.CheckOverwrite(path, statePath, force); err != nil {
				return plan, err
			}
			files[path] = content
		}
	}
	snapshot, err := snapshotFiles(files, statePath)
	if err != nil {
		return plan, err
	}
	for path, content := range files {
		if err := renderfs.WriteAtomically(path, []byte(content)); err != nil {
			return plan, snapshot.restore(err)
		}
	}
	for path, content := range files {
		if err := state.Write(statePath, path, content); err != nil {
			return plan, snapshot.restore(err)
		}
	}
	return plan, nil
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
			return nil, fmt.Errorf("read output %q before build: %w", path, err)
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
