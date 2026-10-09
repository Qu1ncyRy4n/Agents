package v2

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/state"
)

// TreeChange is one tree output path the plan would replace with a copy of
// the selected library directory. Added, Changed, and Removed list files
// relative to the output path.
type TreeChange struct {
	Output   string
	Path     string
	Absolute string
	Source   string
	Status   FileStatus
	Added    []string
	Changed  []string
	Removed  []string
	Exclude  []string
	hashes   map[string]string
}

// FileCount reports how many files the copied tree contains.
func (t TreeChange) FileCount() int { return len(t.hashes) }

func planTrees(result *Result, root string) error {
	for _, output := range result.Plan.Outputs {
		if output.Kind != "tree" || len(output.Sources) != 1 {
			continue
		}
		planned := output.Sources[0]
		library := result.Libraries[planned.Name]
		tree, found := library.Trees[planned.TreeRoot]
		if !found {
			return fmt.Errorf("output %q: tree %q is not declared by source %q", output.Name, planned.TreeRoot, planned.Name)
		}
		source := filepath.Join(library.Root, filepath.FromSlash(tree.Root))
		if err := verifySourceTree(source); err != nil {
			return fmt.Errorf("output %q: %w", output.Name, err)
		}
		hashes, err := state.DirectoryHashes(source)
		if err != nil {
			return fmt.Errorf("output %q: %w", output.Name, err)
		}
		for _, excluded := range planned.Exclude {
			if _, err := os.Lstat(filepath.Join(source, filepath.FromSlash(excluded))); err != nil {
				return fmt.Errorf("output %q: exclude path %q does not exist in tree source %q", output.Name, excluded, planned.Name)
			}
		}
		hashes = filterTreeHashes(hashes, planned.Exclude)
		for _, relative := range output.Paths {
			change, err := inspectTree(result, output.Name, relative, filepath.Join(root, filepath.FromSlash(relative)), source, hashes)
			if err != nil {
				return err
			}
			change.Exclude = append([]string(nil), planned.Exclude...)
			result.Trees = append(result.Trees, change)
		}
	}
	return nil
}

func inspectTree(result *Result, outputName, relative, absolute, source string, hashes map[string]string) (TreeChange, error) {
	change := TreeChange{Output: outputName, Path: relative, Absolute: absolute, Source: source, hashes: hashes}
	if err := rejectTargetSymlinks(absolute, filepath.Dir(filepath.Dir(result.statePath))); err != nil {
		return change, fmt.Errorf("output %q: %w", outputName, err)
	}
	info, err := os.Lstat(absolute)
	if errors.Is(err, os.ErrNotExist) {
		change.Status = FileNew
		change.Added = sortedPaths(hashes)
		return change, nil
	}
	if err != nil {
		return change, fmt.Errorf("inspect output %q: %w", relative, err)
	}
	if !info.IsDir() {
		return change, fmt.Errorf("output %q exists and is not a directory", relative)
	}
	existing, err := state.DirectoryHashes(absolute)
	if err != nil {
		return change, fmt.Errorf("output %q: %w", relative, err)
	}
	for path, digest := range hashes {
		previous, found := existing[path]
		switch {
		case !found:
			change.Added = append(change.Added, path)
		case previous != digest:
			change.Changed = append(change.Changed, path)
		}
	}
	for path := range existing {
		if _, keep := hashes[path]; !keep {
			change.Removed = append(change.Removed, path)
		}
	}
	sort.Strings(change.Added)
	sort.Strings(change.Changed)
	sort.Strings(change.Removed)
	change.Status = FileChanged
	if len(change.Added)+len(change.Changed)+len(change.Removed) == 0 {
		change.Status = FileUnchanged
	}
	treeState, err := state.InspectDirectory(absolute, result.statePath)
	if err != nil {
		return change, err
	}
	switch treeState {
	case state.OutputModified:
		result.Plan.Diagnostics = append(result.Plan.Diagnostics, Diagnostic{
			Severity: SeverityWarning,
			Code:     "MOGENT208",
			Message:  fmt.Sprintf("tree output %q was edited after the last apply; apply requires --force", relative),
		})
	case state.OutputUntracked:
		result.Plan.Diagnostics = append(result.Plan.Diagnostics, Diagnostic{
			Severity: SeverityWarning,
			Code:     "MOGENT208",
			Message:  fmt.Sprintf("tree output %q exists but is not managed by Mogent; apply requires --force", relative),
		})
	}
	return change, nil
}

func checkTreeOverwrite(change TreeChange, statePath string, force bool) error {
	if force {
		return nil
	}
	treeState, err := state.InspectDirectory(change.Absolute, statePath)
	if err != nil {
		return err
	}
	switch treeState {
	case state.OutputUntracked:
		return fmt.Errorf("refusing to replace untracked tree output %q; rerun with --force", change.Path)
	case state.OutputModified:
		return fmt.Errorf("refusing to replace direct edits in tree output %q; inspect them or rerun with --force", change.Path)
	}
	return nil
}

// treeBackup remembers how to undo one tree replacement.
type treeBackup struct {
	target string
	backup string
}

// stageTree replaces target with a fresh copy of source. The previous target,
// when present, is moved aside and returned so a failed transaction can put it
// back; the caller removes it after the transaction commits.
func stageTree(source, target string, excludes []string) (treeBackup, error) {
	backup := treeBackup{target: target}
	parent := filepath.Dir(target)
	if err := os.MkdirAll(parent, 0o755); err != nil {
		return backup, fmt.Errorf("create tree output parent: %w", err)
	}
	staging, err := os.MkdirTemp(parent, ".mogent-tree-")
	if err != nil {
		return backup, fmt.Errorf("create tree staging directory: %w", err)
	}
	if err := copyTree(source, staging, excludes); err != nil {
		_ = os.RemoveAll(staging)
		return backup, err
	}
	if _, err := os.Lstat(target); err == nil {
		aside, err := os.MkdirTemp(parent, ".mogent-tree-old-")
		if err != nil {
			_ = os.RemoveAll(staging)
			return backup, fmt.Errorf("create tree backup directory: %w", err)
		}
		if err := os.Remove(aside); err != nil {
			_ = os.RemoveAll(staging)
			return backup, fmt.Errorf("prepare tree backup: %w", err)
		}
		if err := os.Rename(target, aside); err != nil {
			_ = os.RemoveAll(staging)
			return backup, fmt.Errorf("move previous tree output aside: %w", err)
		}
		backup.backup = aside
	} else if !errors.Is(err, os.ErrNotExist) {
		_ = os.RemoveAll(staging)
		return backup, fmt.Errorf("inspect tree output: %w", err)
	}
	if err := os.Rename(staging, target); err != nil {
		_ = os.RemoveAll(staging)
		return backup, errors.Join(fmt.Errorf("install tree output: %w", err), backup.restore())
	}
	return backup, nil
}

func (b treeBackup) restore() error {
	var errs []error
	if err := os.RemoveAll(b.target); err != nil {
		errs = append(errs, err)
	}
	if b.backup != "" {
		if err := os.Rename(b.backup, b.target); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (b treeBackup) discard() error {
	if b.backup == "" {
		return nil
	}
	return os.RemoveAll(b.backup)
}

func copyTree(source, target string, excludes []string) error {
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if relative != "." && isExcludedTreePath(filepath.ToSlash(relative), excludes) {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		destination := filepath.Join(target, relative)
		if entry.IsDir() {
			return os.MkdirAll(destination, 0o755)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		return renderfs.WriteAtomicallyMode(destination, content, info.Mode().Perm())
	})
}

func verifySourceTree(root string) error {
	info, err := os.Lstat(root)
	if err != nil {
		return fmt.Errorf("inspect tree source %q: %w", root, err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return fmt.Errorf("tree source %q must be a non-symlink directory", root)
	}
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("tree source contains symlink %q", path)
		}
		if !entry.IsDir() && !entry.Type().IsRegular() {
			return fmt.Errorf("tree source contains non-regular file %q", path)
		}
		return nil
	})
}

func rejectTargetSymlinks(path, workspaceRoot string) error {
	canonicalRoot, err := filepath.EvalSymlinks(workspaceRoot)
	if err != nil {
		return fmt.Errorf("resolve workspace root %q: %w", workspaceRoot, err)
	}
	relative, err := filepath.Rel(workspaceRoot, path)
	if err != nil {
		return fmt.Errorf("relativize output target %q: %w", path, err)
	}
	path = filepath.Join(canonicalRoot, relative)
	for current := filepath.Clean(path); ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("output target uses symlink %q", current)
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("inspect output target %q: %w", current, err)
		}
		if current == canonicalRoot || current == filepath.Dir(current) {
			return nil
		}
	}
}

func filterTreeHashes(hashes map[string]string, excludes []string) map[string]string {
	filtered := make(map[string]string, len(hashes))
	for relative, digest := range hashes {
		if !isExcludedTreePath(relative, excludes) {
			filtered[relative] = digest
		}
	}
	return filtered
}

func isExcludedTreePath(relative string, excludes []string) bool {
	for _, excluded := range excludes {
		if relative == excluded || strings.HasPrefix(relative, excluded+"/") {
			return true
		}
	}
	return false
}

func sortedPaths(hashes map[string]string) []string {
	paths := make([]string, 0, len(hashes))
	for path := range hashes {
		paths = append(paths, path)
	}
	sort.Strings(paths)
	return paths
}
