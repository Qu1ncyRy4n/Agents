package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/manifest"
	"github.com/Qu1ncyRy4n/Agents/render"
	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"github.com/Qu1ncyRy4n/Agents/state"
)

type configuredOutput struct {
	spec         manifest.Output
	path, source string
}

func configuredOutputs(value *manifest.Manifest, manifestPath string) ([]configuredOutput, error) {
	specs := value.EffectiveOutputs()
	outputs := make([]configuredOutput, 0, len(specs))
	for _, spec := range specs {
		path := filepath.Join(filepath.Dir(manifestPath), filepath.FromSlash(strings.TrimSuffix(spec.Path, "/")))
		if err := rejectTargetSymlinks(path, filepath.Dir(manifestPath)); err != nil {
			return nil, err
		}
		for alias, source := range value.Sources {
			if source.Location == "" || strings.HasPrefix(source.Location, "http") {
				continue
			}
			root, err := render.ResolveSourcePath(value, manifestPath, alias)
			if err != nil {
				return nil, err
			}
			if sameOrWithin(path, filepath.Clean(root)) {
				return nil, fmt.Errorf("output %q is inside source root %q", spec.Path, source.Location)
			}
		}
		output := configuredOutput{spec: spec, path: path}
		if spec.Directory() {
			selection := spec.From
			if selection == "" {
				if len(spec.Include) != 1 || len(spec.Exclude) != 0 || len(spec.Include[0].Tags) != 0 {
					return nil, fmt.Errorf("directory output %q requires one all or source selector without exclusions; tag selection requires indexed metadata and is not a safe tree copy", spec.Path)
				}
				if spec.Include[0].All != "" {
					selection = spec.Include[0].All + ":root"
				} else {
					selection = spec.Include[0].Source
				}
			}
			alias, selected, _ := manifest.SplitReference(selection)
			root, err := render.ResolveSourcePath(value, manifestPath, alias)
			if err != nil {
				return nil, err
			}
			if selected == "root" && spec.Include != nil && spec.Include[0].All != "" {
				output.source = root
			} else {
				output.source = filepath.Join(root, filepath.FromSlash(selected))
			}
			if err := verifySourceDirectory(output.source); err != nil {
				return nil, err
			}
		}
		outputs = append(outputs, output)
	}
	return outputs, nil
}

func sameOrWithin(path, root string) bool {
	return path == root || strings.HasPrefix(path, root+string(filepath.Separator))
}

func rejectTargetSymlinks(path, workspaceRoot string) error {
	for current := filepath.Clean(path); ; current = filepath.Dir(current) {
		info, err := os.Lstat(current)
		if err == nil && info.Mode()&os.ModeSymlink != 0 {
			return fmt.Errorf("output target uses symlink %q", current)
		}
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("inspect output target %q: %w", current, err)
		}
		if current == filepath.Clean(workspaceRoot) || current == filepath.Dir(current) {
			return nil
		}
	}
}

func verifySourceDirectory(root string) error {
	info, err := os.Lstat(root)
	if err != nil {
		return fmt.Errorf("inspect directory output source %q: %w", root, err)
	}
	if info.Mode()&os.ModeSymlink != 0 || !info.IsDir() {
		return fmt.Errorf("directory output source %q must be a non-symlink directory", root)
	}
	return filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("directory output source contains symlink %q", path)
		}
		if !entry.IsDir() && !entry.Type().IsRegular() {
			return fmt.Errorf("directory output source contains non-regular file %q", path)
		}
		return nil
	})
}

// BuildOutputs writes every configured output as one rollback transaction.
func BuildOutputs(value *manifest.Manifest, manifestPath, markdown string, force bool) error {
	return BuildOutputsWithOptions(value, manifestPath, markdown, force, render.Options{})
}

// BuildOutputsWithOptions writes outputs while applying options to Markdown only.
func BuildOutputsWithOptions(value *manifest.Manifest, manifestPath, markdown string, force bool, options render.Options) error {
	outputs, err := configuredOutputs(value, manifestPath)
	if err != nil {
		return err
	}
	statePath := filepath.Join(filepath.Dir(manifestPath), ".mogent", "state.json")
	for _, output := range outputs {
		if output.spec.Directory() {
			status, err := state.InspectDirectory(output.path, statePath)
			if err != nil {
				return err
			}
			if !force && status != state.OutputMissing && status != state.OutputClean {
				return fmt.Errorf("refusing to overwrite %s generated directory %q; rerun with --force", status, output.path)
			}
		} else if err := state.CheckOverwrite(output.path, statePath, force); err != nil {
			return err
		}
	}
	snapshot, err := snapshotOutputTree(outputs, statePath)
	if err != nil {
		return err
	}
	files, dirs := map[string]string{}, map[string]map[string]string{}
	for _, output := range outputs {
		if output.spec.Directory() {
			hashes, err := syncDirectory(output.source, output.path, snapshot.directoryFiles[output.path])
			if err != nil {
				return rollbackOutputTree(err, snapshot)
			}
			dirs[output.path] = hashes
		} else {
			content := markdown
			if len(output.spec.Include) > 0 {
				result, err := render.RenderOutputWithOptions(value, manifestPath, output.spec, options)
				if err != nil {
					return rollbackOutputTree(err, snapshot)
				}
				content = result.Content
			}
			if err := renderfs.WriteAtomically(output.path, []byte(content)); err != nil {
				return rollbackOutputTree(err, snapshot)
			}
			files[output.path] = content
		}
	}
	if err := state.WriteAll(statePath, files, dirs); err != nil {
		return rollbackOutputTree(err, snapshot)
	}
	return nil
}

type outputSnapshot struct {
	root           string
	paths          map[string]bool
	backups        map[string]string
	statePath      string
	state          []byte
	stateExisted   bool
	directoryFiles map[string]map[string]string
}

func snapshotOutputTree(outputs []configuredOutput, statePath string) (*outputSnapshot, error) {
	root, err := os.MkdirTemp("", "mogent-output-rollback-")
	if err != nil {
		return nil, err
	}
	s := &outputSnapshot{root: root, paths: map[string]bool{}, backups: map[string]string{}, statePath: statePath, directoryFiles: map[string]map[string]string{}}
	for i, output := range outputs {
		s.paths[output.path] = false
		info, err := os.Lstat(output.path)
		if err == nil {
			s.paths[output.path] = true
			if info.IsDir() {
				hashes, hashErr := state.DirectoryHashes(output.path)
				if hashErr != nil {
					os.RemoveAll(root)
					return nil, hashErr
				}
				s.directoryFiles[output.path] = hashes
			}
			backup := filepath.Join(root, fmt.Sprintf("%d", i))
			s.backups[output.path] = backup
			if err := copyPath(output.path, backup); err != nil {
				os.RemoveAll(root)
				return nil, err
			}
		} else if !errors.Is(err, os.ErrNotExist) {
			os.RemoveAll(root)
			return nil, err
		}
	}
	s.state, s.stateExisted, err = readOptional(statePath)
	if err != nil {
		os.RemoveAll(root)
		return nil, err
	}
	return s, nil
}

func rollbackOutputTree(original error, snapshot *outputSnapshot) error {
	defer os.RemoveAll(snapshot.root)
	var errs []error
	for path, existed := range snapshot.paths {
		if err := os.RemoveAll(path); err != nil {
			errs = append(errs, err)
			continue
		}
		if existed {
			if err := copyPath(snapshot.backups[path], path); err != nil {
				errs = append(errs, err)
			}
		}
	}
	if snapshot.stateExisted {
		if err := renderfs.WriteAtomicallyMode(snapshot.statePath, snapshot.state, 0o600); err != nil {
			errs = append(errs, err)
		}
	} else if err := os.Remove(snapshot.statePath); err != nil && !errors.Is(err, os.ErrNotExist) {
		errs = append(errs, err)
	}
	return errors.Join(append([]error{original}, errs...)...)
}

func syncDirectory(source, target string, old map[string]string) (map[string]string, error) {
	if err := rejectTargetSymlinks(target, filepath.Dir(target)); err != nil {
		return nil, err
	}
	if err := os.MkdirAll(target, 0o755); err != nil {
		return nil, err
	}
	newHashes, err := state.DirectoryHashes(source)
	if err != nil {
		return nil, err
	}
	for relative := range newHashes {
		from, to := filepath.Join(source, filepath.FromSlash(relative)), filepath.Join(target, filepath.FromSlash(relative))
		content, err := os.ReadFile(from)
		if err != nil {
			return nil, err
		}
		if err := renderfs.WriteAtomically(to, content); err != nil {
			return nil, err
		}
		info, err := os.Stat(from)
		if err != nil {
			return nil, err
		}
		if err := os.Chmod(to, info.Mode().Perm()); err != nil {
			return nil, err
		}
	}
	for relative := range old {
		if _, keep := newHashes[relative]; !keep {
			if err := os.Remove(filepath.Join(target, filepath.FromSlash(relative))); err != nil && !errors.Is(err, os.ErrNotExist) {
				return nil, err
			}
		}
	}
	return newHashes, nil
}

func copyPath(source, target string) error {
	info, err := os.Lstat(source)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return filepath.WalkDir(source, func(path string, entry os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			rel, _ := filepath.Rel(source, path)
			destination := filepath.Join(target, rel)
			if entry.IsDir() {
				return os.MkdirAll(destination, 0o755)
			}
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			return renderfs.WriteAtomically(destination, content)
		})
	}
	content, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	return renderfs.WriteAtomically(target, content)
}
