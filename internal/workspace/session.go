// Package workspace owns manifest draft state and durable workspace writes.
package workspace

import (
	"errors"
	"fmt"
	"os"

	"github.com/Qu1ncyRy4n/Agents/internal/library"
	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
	"github.com/Qu1ncyRy4n/Agents/internal/render"
	"github.com/Qu1ncyRy4n/Agents/internal/state"
)

// Session is one loaded mogent workspace. UI and CLI callers mutate Draft,
// rebuild the preview, then explicitly save through the shared safety path.
type Session struct {
	ManifestPath string
	Saved        *manifest.Manifest
	Draft        *manifest.Manifest
	Sources      map[string]*library.Index
	Output       string
	Dirty        bool
}

// New loads an existing manifest and prepares a clean in-memory draft.
func New(manifestPath string) (*Session, error) {
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		return nil, err
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		return nil, err
	}
	sources, err := render.LoadSources(value, loadedPath)
	if err != nil {
		return nil, err
	}
	return &Session{
		ManifestPath: loadedPath,
		Saved:        value,
		Draft:        value.Clone(),
		Sources:      sources,
		Output:       result.Content,
	}, nil
}

// MarkDraftChanged records an in-memory change and rebuilds the preview.
func (s *Session) MarkDraftChanged() error {
	s.Dirty = true
	result, err := render.Build(s.Draft, s.ManifestPath)
	if err != nil {
		return err
	}
	s.Output = result.Content
	return nil
}

// DiscardDraft resets the draft to the last saved manifest.
func (s *Session) DiscardDraft() error {
	s.Draft = s.Saved.Clone()
	s.Dirty = false
	result, err := render.Build(s.Draft, s.ManifestPath)
	if err != nil {
		return err
	}
	s.Output = result.Content
	return nil
}

// SaveAndBuild atomically writes agents.yaml, AGENTS.md, and generated-output
// state. If the output write/state write fails, the manifest and output are
// restored to their pre-save contents.
func (s *Session) SaveAndBuild() error {
	result, err := render.Build(s.Draft, s.ManifestPath)
	if err != nil {
		return err
	}
	outputPath := s.OutputPath()
	statePath := s.StatePath()
	if err := state.CheckOverwrite(outputPath, statePath, false); err != nil {
		return err
	}
	originalManifest, err := os.ReadFile(s.ManifestPath)
	if err != nil {
		return fmt.Errorf("read existing manifest before save: %w", err)
	}
	originalOutput, outputExisted, err := readOptional(outputPath)
	if err != nil {
		return err
	}
	if err := manifest.WriteAtomically(s.ManifestPath, s.Draft); err != nil {
		return err
	}
	if err := render.WriteAtomically(outputPath, result.Content); err != nil {
		return rollbackSave(err, s.ManifestPath, originalManifest, outputPath, originalOutput, outputExisted)
	}
	if err := state.Write(statePath, result.Content); err != nil {
		return rollbackSave(err, s.ManifestPath, originalManifest, outputPath, originalOutput, outputExisted)
	}
	s.Saved = s.Draft.Clone()
	s.Draft = s.Saved.Clone()
	s.Output = result.Content
	s.Dirty = false
	sources, err := render.LoadSources(s.Saved, s.ManifestPath)
	if err != nil {
		return err
	}
	s.Sources = sources
	return nil
}

func readOptional(path string) ([]byte, bool, error) {
	contents, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, fmt.Errorf("read existing output before save: %w", err)
	}
	return contents, true, nil
}

func rollbackSave(saveErr error, manifestPath string, originalManifest []byte, outputPath string, originalOutput []byte, outputExisted bool) error {
	var rollbackErrs []error
	if err := render.WriteAtomically(manifestPath, string(originalManifest)); err != nil {
		rollbackErrs = append(rollbackErrs, fmt.Errorf("restore manifest: %w", err))
	}
	if outputExisted {
		if err := render.WriteAtomically(outputPath, string(originalOutput)); err != nil {
			rollbackErrs = append(rollbackErrs, fmt.Errorf("restore output: %w", err))
		}
	} else if err := os.Remove(outputPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		rollbackErrs = append(rollbackErrs, fmt.Errorf("remove new output after failed save: %w", err))
	}
	return errors.Join(append([]error{saveErr}, rollbackErrs...)...)
}
