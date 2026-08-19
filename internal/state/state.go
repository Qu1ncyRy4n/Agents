// Package state records the last generated output hash outside version control.
package state

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/Qu1ncyRy4n/Agents/internal/renderfs"
)

type generated struct {
	SHA256 string `json:"sha256"`
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
	if Hash(output) != previous.SHA256 {
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
	if Hash(output) != previous.SHA256 {
		return fmt.Errorf("refusing to overwrite direct edits in %q; inspect the generated and direct changes or rerun with --force", outputPath)
	}
	return nil
}

// Write records exactly the content just written to the generated output.
func Write(statePath, output string) error {
	contents, err := json.MarshalIndent(generated{SHA256: Hash([]byte(output))}, "", "  ")
	if err != nil {
		return fmt.Errorf("encode generated-output state: %w", err)
	}
	contents = append(contents, '\n')
	if err := renderfs.WriteAtomically(statePath, contents); err != nil {
		return fmt.Errorf("write generated-output state: %w", err)
	}
	return nil
}

func Hash(content []byte) string {
	digest := sha256.Sum256(content)
	return hex.EncodeToString(digest[:])
}

func decode(contents []byte) (generated, error) {
	var previous generated
	if err := json.Unmarshal(contents, &previous); err != nil || previous.SHA256 == "" {
		return generated{}, fmt.Errorf("read generated-output state: invalid state file")
	}
	return previous, nil
}
