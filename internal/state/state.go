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
)

type generated struct {
	SHA256 string `json:"sha256"`
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
	var previous generated
	if err := json.Unmarshal(contents, &previous); err != nil || previous.SHA256 == "" {
		return fmt.Errorf("read generated-output state: invalid state file")
	}
	if Hash(output) != previous.SHA256 {
		return fmt.Errorf("refusing to overwrite direct edits in %q; inspect the generated and direct changes or rerun with --force", outputPath)
	}
	return nil
}

// Changed reports whether the current output differs from the last generated
// hash. Missing output or missing state means there is no tracked drift yet.
func Changed(outputPath, statePath string) (bool, error) {
	output, err := os.ReadFile(outputPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("read existing output: %w", err)
	}
	contents, err := os.ReadFile(statePath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("read generated-output state: %w", err)
	}
	var previous generated
	if err := json.Unmarshal(contents, &previous); err != nil || previous.SHA256 == "" {
		return false, fmt.Errorf("read generated-output state: invalid state file")
	}
	return Hash(output) != previous.SHA256, nil
}

// OutputPath resolves a manifest output path beside the manifest when needed.
func OutputPath(manifestPath, output string) string {
	if filepath.IsAbs(output) {
		return output
	}
	return filepath.Join(filepath.Dir(manifestPath), output)
}

// StatePath resolves the local generated-output state path for a manifest.
func StatePath(manifestPath string) string {
	return filepath.Join(filepath.Dir(manifestPath), ".mogent", "state.json")
}

// Write records exactly the content just written to the generated output.
func Write(statePath, output string) error {
	if err := os.MkdirAll(filepath.Dir(statePath), 0o755); err != nil {
		return fmt.Errorf("create generated-output state directory: %w", err)
	}
	contents, err := json.MarshalIndent(generated{SHA256: Hash([]byte(output))}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode generated-output state: %w", err)
	}
	contents = append(contents, '\n')
	return writeAtomically(statePath, contents)
}

func Hash(content []byte) string {
	digest := sha256.Sum256(content)
	return hex.EncodeToString(digest[:])
}

func writeAtomically(path string, contents []byte) error {
	temporary, err := os.CreateTemp(filepath.Dir(path), ".mogent-state-")
	if err != nil {
		return fmt.Errorf("create temporary state: %w", err)
	}
	temporaryPath := temporary.Name()
	if _, err := temporary.Write(contents); err != nil {
		return cleanup(fmt.Errorf("write temporary state: %w", err), temporary, temporaryPath)
	}
	if err := temporary.Close(); err != nil {
		if removeErr := os.Remove(temporaryPath); removeErr != nil {
			return errors.Join(fmt.Errorf("close temporary state: %w", err), removeErr)
		}
		return fmt.Errorf("close temporary state: %w", err)
	}
	if err := os.Rename(temporaryPath, path); err != nil {
		if removeErr := os.Remove(temporaryPath); removeErr != nil {
			return errors.Join(fmt.Errorf("replace state: %w", err), removeErr)
		}
		return fmt.Errorf("replace state: %w", err)
	}
	return nil
}

func cleanup(original error, temporary *os.File, path string) error {
	return errors.Join(original, temporary.Close(), os.Remove(path))
}
