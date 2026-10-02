// Package v2 loads the additive HCL configuration model described by the v2
// specification. It intentionally does not alter the legacy YAML manifest.
package v2

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
)

const ConfigFile = "mogent.hcl"

// Config is the validated user-side v2 configuration. Selection expressions
// remain as HCL expressions until the planner evaluates them against a library.
type Config struct {
	Path    string
	Sources map[string]Source
	Outputs []Output
}

// LoadLocalLibraries loads every local server-side library declared by config.
// Remote resolution is intentionally deferred to the lock implementation.
func LoadLocalLibraries(config *Config) (map[string]*Library, error) {
	libraries := make(map[string]*Library)
	for alias, source := range config.Sources {
		if source.Git != "" {
			return nil, fmt.Errorf("source %q is remote; run mogent lock after remote v2 resolution is implemented", alias)
		}
		root := source.Local
		if !filepath.IsAbs(root) {
			root = filepath.Join(filepath.Dir(config.Path), root)
		}
		if source.Subdir != "" {
			root = filepath.Join(root, filepath.FromSlash(source.Subdir))
		}
		library, err := LoadLibrary(filepath.Clean(root))
		if err != nil {
			return nil, fmt.Errorf("source %q: %w", alias, err)
		}
		libraries[alias] = library
	}
	return libraries, nil
}

type Source struct {
	Name   string
	Git    string
	Ref    string
	Local  string
	Subdir string
}

type Output struct {
	Name    string
	Paths   []string
	Kind    string
	Sources []OutputSource
}

type OutputSource struct {
	Name    string
	From    string
	Select  hcl.Expression
	TagsAll []string
	TagsAny []string
}

type fileConfig struct {
	Mogent  []mogentBlock  `hcl:"mogent,block"`
	Sources []sourcesBlock `hcl:"sources,block"`
	Outputs []outputsBlock `hcl:"outputs,block"`
}

type mogentBlock struct {
	Format int `hcl:"format"`
}

type sourcesBlock struct {
	Sources []sourceBlock `hcl:"source,block"`
}

type sourceBlock struct {
	Name   string  `hcl:"name,label"`
	Git    *string `hcl:"git,optional"`
	Ref    *string `hcl:"ref,optional"`
	Local  *string `hcl:"local,optional"`
	Subdir *string `hcl:"subdir,optional"`
}

type outputsBlock struct {
	Outputs []outputBlock `hcl:"output,block"`
}

type outputBlock struct {
	Name    string              `hcl:"name,label"`
	Path    *string             `hcl:"path,optional"`
	Paths   []string            `hcl:"paths,optional"`
	Kind    string              `hcl:"kind"`
	Sources []outputSourceBlock `hcl:"source,block"`
}

type outputSourceBlock struct {
	Name    string         `hcl:"name,label"`
	From    string         `hcl:"from"`
	Select  hcl.Expression `hcl:"select,optional"`
	TagsAll []string       `hcl:"tags_all,optional"`
	TagsAny []string       `hcl:"tags_any,optional"`
}

// Load reads and validates one v2 user-side configuration.
func Load(path string) (*Config, error) {
	var raw fileConfig
	if err := hclsimple.DecodeFile(path, nil, &raw); err != nil {
		return nil, err
	}
	if len(raw.Mogent) != 1 || raw.Mogent[0].Format != 2 {
		return nil, fmt.Errorf("%s: require exactly one mogent block with format = 2", path)
	}
	if len(raw.Sources) != 1 {
		return nil, fmt.Errorf("%s: require exactly one sources block", path)
	}
	if len(raw.Outputs) != 1 {
		return nil, fmt.Errorf("%s: require exactly one outputs block", path)
	}

	config := &Config{Path: filepath.Clean(path), Sources: make(map[string]Source)}
	for _, block := range raw.Sources[0].Sources {
		source, err := normalizeSource(block)
		if err != nil {
			return nil, fmt.Errorf("source %q: %w", block.Name, err)
		}
		if _, duplicate := config.Sources[source.Name]; duplicate {
			return nil, fmt.Errorf("duplicate source %q", source.Name)
		}
		config.Sources[source.Name] = source
	}
	if len(config.Sources) == 0 {
		return nil, fmt.Errorf("sources block must declare at least one source")
	}

	outputNames := make(map[string]bool)
	outputPaths := make(map[string]string)
	for _, block := range raw.Outputs[0].Outputs {
		output, err := normalizeOutput(block, config.Sources)
		if err != nil {
			return nil, fmt.Errorf("output %q: %w", block.Name, err)
		}
		if outputNames[output.Name] {
			return nil, fmt.Errorf("duplicate output %q", output.Name)
		}
		outputNames[output.Name] = true
		for _, outputPath := range output.Paths {
			if existing, duplicate := outputPaths[outputPath]; duplicate {
				return nil, fmt.Errorf("output path %q is declared by both %q and %q", outputPath, existing, output.Name)
			}
			outputPaths[outputPath] = output.Name
		}
		config.Outputs = append(config.Outputs, output)
	}
	if len(config.Outputs) == 0 {
		return nil, fmt.Errorf("outputs block must declare at least one output")
	}
	return config, nil
}

func normalizeSource(block sourceBlock) (Source, error) {
	source := Source{Name: block.Name}
	if source.Name == "" {
		return Source{}, fmt.Errorf("source name must not be empty")
	}
	if block.Git != nil {
		source.Git = *block.Git
	}
	if block.Ref != nil {
		source.Ref = *block.Ref
	}
	if block.Local != nil {
		source.Local = *block.Local
	}
	if block.Subdir != nil {
		source.Subdir = *block.Subdir
	}
	if (source.Git == "") == (source.Local == "") {
		return Source{}, fmt.Errorf("declare exactly one of git or local")
	}
	if source.Local != "" && source.Ref != "" {
		return Source{}, fmt.Errorf("ref is valid only for git sources")
	}
	if err := validateRelativePath(source.Subdir, "subdir"); err != nil {
		return Source{}, err
	}
	return source, nil
}

func normalizeOutput(block outputBlock, sources map[string]Source) (Output, error) {
	output := Output{Name: block.Name, Kind: block.Kind}
	if output.Name == "" {
		return Output{}, fmt.Errorf("output name must not be empty")
	}
	if block.Path != nil {
		output.Paths = append(output.Paths, *block.Path)
	}
	output.Paths = append(output.Paths, block.Paths...)
	if len(output.Paths) == 0 {
		return Output{}, fmt.Errorf("declare path or paths")
	}
	if block.Path != nil && len(block.Paths) > 0 {
		return Output{}, fmt.Errorf("path and paths cannot be combined")
	}
	if output.Kind != "markdown" && output.Kind != "tree" {
		return Output{}, fmt.Errorf("kind must be markdown or tree")
	}
	seenPaths := make(map[string]bool)
	for _, path := range output.Paths {
		if err := validateRelativePath(path, "path"); err != nil {
			return Output{}, err
		}
		if seenPaths[path] {
			return Output{}, fmt.Errorf("path %q is repeated", path)
		}
		seenPaths[path] = true
	}
	if len(block.Sources) == 0 {
		return Output{}, fmt.Errorf("declare at least one source block")
	}
	for _, blockSource := range block.Sources {
		if _, found := sources[blockSource.Name]; !found {
			return Output{}, fmt.Errorf("source block %q is not declared in sources", blockSource.Name)
		}
		if !strings.HasPrefix(blockSource.From, blockSource.Name+":") {
			return Output{}, fmt.Errorf("from %q must begin with %q", blockSource.From, blockSource.Name+":")
		}
		if blockSource.Select == nil {
			return Output{}, fmt.Errorf("source block %q requires select", blockSource.Name)
		}
		for _, tag := range append(append([]string(nil), blockSource.TagsAll...), blockSource.TagsAny...) {
			if err := validateTag(tag); err != nil {
				return Output{}, err
			}
		}
		output.Sources = append(output.Sources, OutputSource{Name: blockSource.Name, From: blockSource.From, Select: blockSource.Select, TagsAll: blockSource.TagsAll, TagsAny: blockSource.TagsAny})
	}
	return output, nil
}

func validateRelativePath(value, field string) error {
	if value == "" {
		if field == "subdir" {
			return nil
		}
		return fmt.Errorf("%s must not be empty", field)
	}
	if filepath.IsAbs(value) || strings.Contains(value, "\\") || strings.HasPrefix(value, "../") || value == ".." || strings.Contains(value, "/../") || path.Clean(value) != value {
		return fmt.Errorf("%s %q must be a safe relative path", field, value)
	}
	return nil
}
