// Package state records the last generated output hash outside version control.
package state

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
)

type generated struct {
	Version    int               `json:"version,omitempty"`
	OutputPath string            `json:"output_path,omitempty"`
	SHA256     string            `json:"sha256,omitempty"`
	Outputs    map[string]record `json:"outputs,omitempty"`
}

type record struct {
	Kind   string            `json:"kind"`
	SHA256 string            `json:"sha256,omitempty"`
	Files  map[string]string `json:"files,omitempty"`
}

// OutputState describes the relationship between an output file and mogent's
// last generated-output record.
type OutputState string

const (
	OutputMissing   OutputState = "missing"
	OutputUntracked OutputState = "untracked"
	OutputClean     OutputState = "clean"
	OutputModified  OutputState = "modified"
)

// Inspect reports whether the output is missing, untracked, unchanged from the
// last mogent write, or edited after the last mogent write.
func Inspect(outputPath, statePath string) (OutputState, error) {
	output, err := os.ReadFile(outputPath)
	if errors.Is(err, os.ErrNotExist) {
		return OutputMissing, nil
	}
	if err != nil {
		return "", fmt.Errorf("read existing output: %w", err)
	}
	contents, err := os.ReadFile(statePath)
	if errors.Is(err, os.ErrNotExist) {
		return OutputUntracked, nil
	}
	if err != nil {
		return "", fmt.Errorf("read generated-output state: %w", err)
	}
	previous, err := decode(contents)
	if err != nil {
		return "", err
	}
	record, found, err := previous.output(outputPath, statePath)
	if err != nil {
		return "", err
	}
	if !found || record.Kind != "file" {
		return OutputUntracked, nil
	}
	if Hash(output) != record.SHA256 {
		return OutputModified, nil
	}
	return OutputClean, nil
}

// CheckOverwrite refuses to replace an output that mogent did not generate or
// whose content no longer matches its recorded hash. Source: DI-vukam.
func CheckOverwrite(outputPath, statePath string, force bool) error {
	output, err := os.ReadFile(outputPath)
	if errors.Is(err, os.ErrNotExist) || force {
		return nil
	}
	if err != nil {
		return fmt.Errorf("read existing output: %w", err)
	}
	contents, err := os.ReadFile(statePath)
	if errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("refusing to overwrite untracked output %q; rerun with --force", outputPath)
	}
	if err != nil {
		return fmt.Errorf("read generated-output state: %w", err)
	}
	previous, err := decode(contents)
	if err != nil {
		return err
	}
	record, found, err := previous.output(outputPath, statePath)
	if err != nil {
		return err
	}
	if !found || record.Kind != "file" {
		return fmt.Errorf("refusing to overwrite untracked output %q; output path does not match generated-output state, rerun with --force", outputPath)
	}
	if Hash(output) != record.SHA256 {
		return fmt.Errorf("refusing to overwrite direct edits in %q; inspect the generated and direct changes or rerun with --force", outputPath)
	}
	return nil
}

// Write records exactly the content just written to the generated output.
func Write(statePath, outputPath, output string) error {
	previous, err := readState(statePath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return err
	}
	outputs := make(map[string]record, len(previous.Outputs)+1)
	for path, value := range previous.Outputs {
		key, err := portableStateKey(statePath, path)
		if err != nil {
			return err
		}
		outputs[key] = value
	}
	if previous.OutputPath != "" && previous.SHA256 != "" {
		key, err := portableStateKey(statePath, previous.OutputPath)
		if err != nil {
			return err
		}
		outputs[key] = record{Kind: "file", SHA256: previous.SHA256}
	}
	key, err := outputKey(statePath, outputPath)
	if err != nil {
		return err
	}
	outputs[key] = record{Kind: "file", SHA256: Hash([]byte(output))}
	return writeOutputs(statePath, outputs)
}

// WriteAll records all generated files and directory trees together. It writes
// the portable version-three shape but decode also accepts legacy absolute-path
// state shapes.
func WriteAll(statePath string, files map[string]string, directories map[string]map[string]string) error {
	outputs := make(map[string]record, len(files)+len(directories))
	for path, content := range files {
		key, err := outputKey(statePath, path)
		if err != nil {
			return err
		}
		outputs[key] = record{Kind: "file", SHA256: Hash([]byte(content))}
	}
	for path, entries := range directories {
		key, err := outputKey(statePath, path)
		if err != nil {
			return err
		}
		copy := make(map[string]string, len(entries))
		for name, digest := range entries {
			copy[name] = digest
		}
		outputs[key] = record{Kind: "directory", Files: copy}
	}
	return writeOutputs(statePath, outputs)
}

func portableStateKey(statePath, path string) (string, error) {
	if filepath.IsAbs(path) {
		return outputKey(statePath, path)
	}
	return filepath.ToSlash(path), nil
}

func writeOutputs(statePath string, outputs map[string]record) error {
	contents, err := json.MarshalIndent(generated{Version: 3, Outputs: outputs}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode generated-output state: %w", err)
	}
	contents = append(contents, '\n')
	if err := renderfs.WriteAtomicallyMode(statePath, contents, 0o600); err != nil {
		return fmt.Errorf("write generated-output state: %w", err)
	}
	return nil
}

// InspectDirectory reports tree drift, including untracked files.
func InspectDirectory(outputPath, statePath string) (OutputState, error) {
	info, err := os.Lstat(outputPath)
	if errors.Is(err, os.ErrNotExist) {
		return OutputMissing, nil
	}
	if err != nil {
		return "", fmt.Errorf("inspect existing output: %w", err)
	}
	if !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return OutputUntracked, nil
	}
	previous, err := readState(statePath)
	if errors.Is(err, os.ErrNotExist) {
		return OutputUntracked, nil
	}
	if err != nil {
		return "", err
	}
	record, found, err := previous.output(outputPath, statePath)
	if err != nil {
		return "", err
	}
	if !found || record.Kind != "directory" {
		return OutputUntracked, nil
	}
	files, err := directoryHashes(outputPath)
	if err != nil {
		return "", err
	}
	if len(files) != len(record.Files) {
		return OutputModified, nil
	}
	for path, digest := range files {
		if record.Files[path] != digest {
			return OutputModified, nil
		}
	}
	return OutputClean, nil
}

func DirectoryHashes(path string) (map[string]string, error) { return directoryHashes(path) }

func directoryHashes(root string) (map[string]string, error) {
	files := map[string]string{}
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("symlink in generated directory %q", path)
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return fmt.Errorf("non-regular file in generated directory %q", path)
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files[filepath.ToSlash(rel)] = Hash(content)
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("read generated directory: %w", err)
	}
	return files, nil
}

func normalizedOutputPath(path string) string {
	absolute, err := filepath.Abs(path)
	if err != nil {
		return filepath.Clean(path)
	}
	return filepath.Clean(absolute)
}

// outputKey identifies an output relative to the workspace containing statePath.
// Resolving existing symlinks makes paths invoked through a workspace symlink use
// the same key as paths invoked through the physical checkout location.
func outputKey(statePath, outputPath string) (string, error) {
	root, err := resolvedPath(filepath.Dir(filepath.Dir(statePath)))
	if err != nil {
		return "", fmt.Errorf("resolve workspace root for generated-output state: %w", err)
	}
	output, err := resolvedPath(outputPath)
	if err != nil {
		return "", fmt.Errorf("resolve generated output %q: %w", outputPath, err)
	}
	relative, err := filepath.Rel(root, output)
	if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) || filepath.IsAbs(relative) {
		return "", fmt.Errorf("generated output %q is outside workspace %q", outputPath, root)
	}
	return filepath.ToSlash(relative), nil
}

func resolvedPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}
	var suffix []string
	for {
		if _, err := os.Lstat(abs); err == nil {
			resolved, err := filepath.EvalSymlinks(abs)
			if err != nil {
				return "", err
			}
			return filepath.Join(append([]string{resolved}, suffix...)...), nil
		} else if !errors.Is(err, os.ErrNotExist) {
			return "", err
		}
		parent := filepath.Dir(abs)
		if parent == abs {
			return "", fmt.Errorf("no existing parent")
		}
		suffix = append([]string{filepath.Base(abs)}, suffix...)
		abs = parent
	}
}

func Hash(content []byte) string {
	digest := sha256.Sum256(content)
	return hex.EncodeToString(digest[:])
}

func decode(contents []byte) (generated, error) {
	var previous generated
	if err := json.Unmarshal(contents, &previous); err != nil || (previous.SHA256 == "" && len(previous.Outputs) == 0) {
		return generated{}, fmt.Errorf("read generated-output state: invalid state file")
	}
	return previous, nil
}

func readState(path string) (generated, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return generated{}, err
	}
	value, err := decode(contents)
	if err != nil {
		return generated{}, err
	}
	return value, nil
}

func (g generated) output(path, statePath string) (record, bool, error) {
	key, err := outputKey(statePath, path)
	if err != nil {
		return record{}, false, err
	}
	if value, ok := g.Outputs[key]; ok {
		return value, true, nil
	}
	// Version two keyed outputs by absolute paths. Keep that read path only for
	// migration; all writes emit version three portable keys.
	abs := normalizedOutputPath(path)
	if value, ok := g.Outputs[abs]; ok {
		return value, true, nil
	}
	if g.OutputPath == abs && g.SHA256 != "" {
		return record{Kind: "file", SHA256: g.SHA256}, true, nil
	}
	return record{}, false, nil
}
