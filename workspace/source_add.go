package workspace

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/library"
	"github.com/Qu1ncyRy4n/Agents/manifest"
)

// SourceAddOptions describes one explicit manifest source declaration.
type SourceAddOptions struct {
	ManifestPath string
	Alias        string
	Location     string
	Subdir       string
	DryRun       bool
}

// SourceAddResult reports the normalized declaration and proposed manifest.
type SourceAddResult struct {
	ManifestPath string
	Alias        string
	Source       manifest.Source
	ManifestYAML string
	Remote       bool
	Wrote        bool
}

// AddSourceDeclaration validates and adds a source without selecting content.
func AddSourceDeclaration(options SourceAddOptions) (*SourceAddResult, error) {
	value, manifestPath, err := manifest.Load(options.ManifestPath)
	if err != nil {
		return nil, err
	}
	alias := strings.TrimSpace(strings.TrimSuffix(options.Alias, ":"))
	location := strings.TrimSpace(options.Location)
	if alias == "" {
		return nil, fmt.Errorf("source alias must not be empty")
	}
	if _, exists := value.Sources[alias]; exists {
		return nil, fmt.Errorf("source %q is already declared in manifest", alias)
	}
	remote, err := validateSourceLocation(location)
	if err != nil {
		return nil, err
	}
	if !remote {
		resolved := location
		if !filepath.IsAbs(resolved) {
			resolved = filepath.Join(filepath.Dir(manifestPath), resolved)
		}
		if _, err := library.Load(resolved); err != nil {
			return nil, err
		}
	}
	value.Sources[alias] = manifest.Source{Location: location, Subdir: strings.TrimSpace(options.Subdir)}
	serialized, err := manifest.Marshal(value)
	if err != nil {
		return nil, err
	}
	result := &SourceAddResult{
		ManifestPath: manifestPath,
		Alias:        alias,
		Source:       value.Sources[alias],
		ManifestYAML: string(serialized),
		Remote:       remote,
	}
	if options.DryRun {
		return result, nil
	}
	if err := manifest.WriteAtomically(manifestPath, value); err != nil {
		return nil, err
	}
	result.Wrote = true
	return result, nil
}

func validateSourceLocation(location string) (bool, error) {
	if location == "" {
		return false, fmt.Errorf("source location must not be empty")
	}
	parsed, err := url.Parse(location)
	if err != nil {
		return false, fmt.Errorf("parse source URL: %w", err)
	}
	if parsed.Scheme == "" {
		return false, nil
	}
	if (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" {
		return false, fmt.Errorf("source location must be a local directory or HTTP(S) Git URL")
	}
	return true, nil
}
