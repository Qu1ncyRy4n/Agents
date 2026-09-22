package workspace

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Qu1ncyRy4n/Agents/manifest"
	"github.com/Qu1ncyRy4n/Agents/render"
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
	outputs := options.Manifest.EffectiveOutputs()
	if len(outputs) == 0 || outputs[0].Directory() {
		return nil, fmt.Errorf("starter manifest requires a Markdown output")
	}
	preview, err := render.RenderOutput(options.Manifest, manifestPath, outputs[0])
	if err != nil {
		return nil, fmt.Errorf("validate starter manifest: %w", err)
	}
	outputPath := filepath.Join(filepath.Dir(manifestPath), outputs[0].Path)
	result := &InitResult{ManifestPath: manifestPath, OutputPath: outputPath, ManifestYAML: string(serialized), Preview: preview.Content}
	if options.DryRun {
		return result, nil
	}
	if err := manifest.WriteAtomically(manifestPath, options.Manifest); err != nil {
		return nil, err
	}
	result.Wrote = true
	if options.Build {
		if err := BuildOutputs(options.Manifest, manifestPath, preview.Content, options.ForceOutput); err != nil {
			_ = os.Remove(manifestPath)
			return nil, err
		}
		result.Built = true
	}
	return result, nil
}
