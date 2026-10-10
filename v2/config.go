// Package v2 loads the additive HCL configuration model described by the v2
// specification. It intentionally does not alter the legacy YAML manifest.
package v2

import (
	"fmt"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/content"
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

// Source is one declared library location. Git is pinned by Commit, which
// Mogent writes into the configuration; Ref records the branch or tag the
// consumer intends to follow. Local sources read whatever is on disk.
type Source struct {
	Name   string
	Git    string
	Ref    string
	Commit string
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
	Name      string
	From      string
	Select    hcl.Expression
	TagsAll   []string
	TagsAny   []string
	Exclude   []string
	Node      string
	Operation string
	Into      string
	Heading   []string
	Replace   bool
	Append    bool
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
	Commit *string `hcl:"commit,optional"`
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
	Name      string         `hcl:"name,label"`
	From      string         `hcl:"from"`
	Select    hcl.Expression `hcl:"select,optional"`
	TagsAll   []string       `hcl:"tags_all,optional"`
	TagsAny   []string       `hcl:"tags_any,optional"`
	Exclude   []string       `hcl:"exclude,optional"`
	Node      *string        `hcl:"node,optional"`
	Operation *string        `hcl:"operation,optional"`
	Into      *string        `hcl:"into,optional"`
	Heading   []string       `hcl:"heading,optional"`
	Replace   *bool          `hcl:"replace,optional"`
	Append    *bool          `hcl:"append,optional"`
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
			if outputPath == "." || outputPath == filepath.Base(config.Path) || outputPath == ".mogent" || strings.HasPrefix(outputPath, ".mogent/") {
				return nil, fmt.Errorf("output path %q overlaps configuration or workspace state", outputPath)
			}
			for earlier, name := range outputPaths {
				if strings.HasPrefix(earlier, outputPath+"/") || strings.HasPrefix(outputPath, earlier+"/") {
					return nil, fmt.Errorf("output paths %q (%s) and %q (%s) overlap", earlier, name, outputPath, output.Name)
				}
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
	if block.Commit != nil {
		source.Commit = *block.Commit
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
	if source.Local != "" && source.Commit != "" {
		return Source{}, fmt.Errorf("commit is valid only for git sources")
	}
	if source.Git != "" {
		if err := validateGitURL(source.Git); err != nil {
			return Source{}, err
		}
		if source.Commit != "" && !isFullCommit(source.Commit) {
			return Source{}, fmt.Errorf("commit %q must be a full 40-character lowercase hex SHA", source.Commit)
		}
		if strings.HasPrefix(source.Ref, "-") {
			return Source{}, fmt.Errorf("ref %q is invalid", source.Ref)
		}
	}
	if err := validateRelativePath(source.Subdir, "subdir"); err != nil {
		return Source{}, err
	}
	return source, nil
}

func validateGitURL(value string) error {
	parsed, err := url.Parse(value)
	if err != nil || parsed.User != nil {
		return fmt.Errorf("git %q must be an http(s) or file URL without embedded credentials", value)
	}
	switch parsed.Scheme {
	case "http", "https":
		if parsed.Host == "" {
			return fmt.Errorf("git %q must name a host", value)
		}
	case "file":
		if parsed.Path == "" {
			return fmt.Errorf("git %q must name a path", value)
		}
	default:
		return fmt.Errorf("git %q must be an http(s) or file URL", value)
	}
	return nil
}

func isFullCommit(value string) bool {
	if len(value) != 40 {
		return false
	}
	for _, character := range value {
		if !strings.ContainsRune("0123456789abcdef", character) {
			return false
		}
	}
	return true
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
	if output.Kind != "markdown" && output.Kind != "tree" && output.Kind != "dir-tree" {
		return Output{}, fmt.Errorf("kind must be markdown, tree, or dir-tree")
	}
	seenPaths := make(map[string]bool)
	for index, path := range output.Paths {
		if (output.Kind == "tree" || output.Kind == "dir-tree") && len(path) > 1 && strings.HasSuffix(path, "/") {
			path = strings.TrimSuffix(path, "/")
			output.Paths[index] = path
		}
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
		// gohcl represents an absent optional expression as a static null
		// expression, rather than a nil interface.
		if blockSource.Select != nil {
			value, diagnostics := blockSource.Select.Value(nil)
			if !diagnostics.HasErrors() && value.IsNull() {
				blockSource.Select = nil
			}
		}
		if _, found := sources[blockSource.Name]; !found {
			return Output{}, fmt.Errorf("source block %q is not declared in sources", blockSource.Name)
		}
		if !strings.HasPrefix(blockSource.From, blockSource.Name+":") {
			return Output{}, fmt.Errorf("from %q must begin with %q", blockSource.From, blockSource.Name+":")
		}
		if blockSource.Select == nil && output.Kind != "dir-tree" {
			return Output{}, fmt.Errorf("source block %q requires select", blockSource.Name)
		}
		for _, tag := range append(append([]string(nil), blockSource.TagsAll...), blockSource.TagsAny...) {
			if err := validateTag(tag); err != nil {
				return Output{}, err
			}
		}
		if output.Kind != "tree" && output.Kind != "dir-tree" && len(blockSource.Exclude) > 0 {
			return Output{}, fmt.Errorf("source block %q: exclude is valid only for tree outputs", blockSource.Name)
		}
		seenExcludes := make(map[string]bool, len(blockSource.Exclude))
		for _, excluded := range blockSource.Exclude {
			if err := validateRelativePath(excluded, "exclude"); err != nil || excluded == "." {
				return Output{}, fmt.Errorf("source block %q: exclude %q must be a safe non-root relative path", blockSource.Name, excluded)
			}
			if seenExcludes[excluded] {
				return Output{}, fmt.Errorf("source block %q: exclude path %q is repeated", blockSource.Name, excluded)
			}
			seenExcludes[excluded] = true
		}
		source := OutputSource{Name: blockSource.Name, From: blockSource.From, Select: blockSource.Select, TagsAll: blockSource.TagsAll, TagsAny: blockSource.TagsAny, Exclude: append([]string(nil), blockSource.Exclude...)}
		if output.Kind == "dir-tree" {
			if blockSource.Operation == nil || blockSource.Into == nil {
				return Output{}, fmt.Errorf("source %q requires operation and into", blockSource.Name)
			}
			source.Operation, source.Into = *blockSource.Operation, *blockSource.Into
			if source.Operation != "copy" && source.Operation != "render-markdown" {
				return Output{}, fmt.Errorf("operation must be copy or render-markdown")
			}
			if err := content.ValidatePath(source.Into, true); err != nil {
				return Output{}, err
			}
			if blockSource.Node != nil {
				source.Node = *blockSource.Node
				if err := content.ValidatePath(source.Node, true); err != nil {
					return Output{}, err
				}
			}
			if blockSource.Replace != nil {
				source.Replace = *blockSource.Replace
			}
			if blockSource.Append != nil {
				source.Append = *blockSource.Append
			}
			source.Heading = append([]string(nil), blockSource.Heading...)
			if source.Replace && source.Append {
				return Output{}, fmt.Errorf("replace and append cannot be combined")
			}
			if source.Node != "" && source.Select != nil {
				return Output{}, fmt.Errorf("node and authored select cannot be combined")
			}
			if source.Operation == "copy" {
				if source.Select != nil || len(source.Heading) > 0 || source.Append {
					return Output{}, fmt.Errorf("copy uses node, not select, heading, or append")
				}
				if source.Node == "" {
					source.Node = "."
				}
			} else {
				if source.Node == "" && source.Select == nil {
					return Output{}, fmt.Errorf("render-markdown requires node or authored select")
				}
				if source.Into == "." || len(source.Exclude) > 0 {
					return Output{}, fmt.Errorf("render-markdown requires a file destination and cannot exclude physical paths")
				}
			}
			if source.Select != nil && len(source.Heading) > 0 {
				return Output{}, fmt.Errorf("heading applies only to a physical file node")
			}
			if len(source.TagsAll)+len(source.TagsAny) > 0 && source.Select == nil {
				return Output{}, fmt.Errorf("tag queries require authored select")
			}
			for _, title := range source.Heading {
				if strings.TrimSpace(title) == "" {
					return Output{}, fmt.Errorf("heading titles must not be empty")
				}
			}
		} else if blockSource.Node != nil || blockSource.Operation != nil || blockSource.Into != nil || len(blockSource.Heading) > 0 || blockSource.Replace != nil || blockSource.Append != nil {
			return Output{}, fmt.Errorf("typed contribution attributes require kind = dir-tree")
		}
		output.Sources = append(output.Sources, source)
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
