package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Qu1ncyRy4n/Agents/manifest"
	"github.com/Qu1ncyRy4n/Agents/render"
	"github.com/Qu1ncyRy4n/Agents/state"
)

type InitOptions struct {
	ManifestPath string
	Manifest     *manifest.Manifest
	DryRun       bool
	Build        bool
	ForceOutput  bool
}

type InitResult struct {
	ManifestPath string
	OutputPath   string
	ManifestYAML string
	Preview      string
	Wrote        bool
	Built        bool
}

// Initialize validates a complete starter manifest before writing anything.
func Initialize(options InitOptions) (*InitResult, error) {
	manifestPath, err := filepath.Abs(options.ManifestPath)
	if err != nil {
		return nil, fmt.Errorf("resolve manifest path: %w", err)
	}
	if _, err := os.Stat(manifestPath); err == nil {
		return nil, fmt.Errorf("manifest %q already exists", manifestPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("inspect manifest path: %w", err)
	}
	serialized, err := manifest.Marshal(options.Manifest)
	if err != nil {
		return nil, err
	}
	preview, err := render.Build(options.Manifest, manifestPath)
	if err != nil {
		return nil, fmt.Errorf("validate starter manifest: %w", err)
	}
	outputPath := options.Manifest.Output
	if !filepath.IsAbs(outputPath) {
		outputPath = filepath.Join(filepath.Dir(manifestPath), outputPath)
	}
	result := &InitResult{ManifestPath: manifestPath, OutputPath: outputPath, ManifestYAML: string(serialized), Preview: preview.Content}
	if options.DryRun {
		return result, nil
	}
	if options.Build {
		statePath := filepath.Join(filepath.Dir(manifestPath), ".mogent", "state.json")
		if err := state.CheckOverwrite(outputPath, statePath, options.ForceOutput); err != nil {
			return nil, err
		}
	}
	paths := []string{manifestPath}
	if options.Build {
		paths = append(paths, outputPath, filepath.Join(filepath.Dir(manifestPath), ".mogent", "state.json"))
	}
	snapshots, err := snapshotPaths(paths)
	if err != nil {
		return nil, err
	}
	if err := manifest.WriteAtomically(manifestPath, options.Manifest); err != nil {
		return nil, err
	}
	result.Wrote = true
	if options.Build {
		if err := render.WriteAtomically(outputPath, preview.Content); err != nil {
			return nil, rollbackPaths(err, snapshots)
		}
		if err := state.Write(filepath.Join(filepath.Dir(manifestPath), ".mogent", "state.json"), preview.Content); err != nil {
			return nil, rollbackPaths(err, snapshots)
		}
		result.Built = true
	}
	return result, nil
}
