package assembly

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/syntax"
)

// InitOptions initializes an authored plan/build, never generated outputs.
// Relative source locations are interpreted relative to the entry file.
type InitOptions struct {
	Path    string
	Kind    syntax.FileKind
	Sources map[string]string
	DryRun  bool
}

type InitResult struct {
	Path    string
	Kind    syntax.FileKind
	Content []byte
}

// Initialize creates a minimal new-language entry point with exclusive creation.
// It does not fetch sources, choose content, or create/overwrite AGENTS.md.
func Initialize(options InitOptions) (*InitResult, error) {
	kind := options.Kind
	if kind == "" {
		kind = syntax.Build
	}
	if kind != syntax.Plan && kind != syntax.Build {
		return nil, fmt.Errorf("init kind must be plan or build")
	}
	path := options.Path
	if path == "" {
		path = "mogent." + string(kind)
	}
	if filepath.Base(path) == "mogent.plan" && kind != syntax.Plan || filepath.Base(path) == "mogent.build" && kind != syntax.Build {
		return nil, fmt.Errorf("init filename %q conflicts with %s header", path, kind)
	}
	if _, err := os.Lstat(path); err == nil {
		return nil, fmt.Errorf("refusing to overwrite existing %s", path)
	} else if !os.IsNotExist(err) {
		return nil, err
	}
	// Default discovery must not gain two independent authorities through init.
	if filepath.Base(path) == "mogent.plan" || filepath.Base(path) == "mogent.build" {
		other := "mogent.plan"
		if kind == syntax.Plan {
			other = "mogent.build"
		}
		counterpart := filepath.Join(filepath.Dir(path), other)
		if _, err := os.Lstat(counterpart); err == nil {
			return nil, fmt.Errorf("workspace already has %s; inspect its authority before creating %s", counterpart, path)
		} else if !os.IsNotExist(err) {
			return nil, err
		}
	}
	header := []*syntax.Entry{fieldStringEntry("format", "1", true)}
	if kind == syntax.Build {
		header = append(header, fieldStringEntry("origin", "authored", false))
	}
	file := &syntax.File{Kind: kind, Entries: []*syntax.Entry{block(string(kind), "", header)}}
	aliases := make([]string, 0, len(options.Sources))
	for alias := range options.Sources {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	for _, alias := range aliases {
		entry, err := sourceEntry(path, alias, options.Sources[alias])
		if err != nil {
			return nil, err
		}
		file.Entries = append(file.Entries, entry)
	}
	contents, err := syntax.Format(file, syntax.FormatOptions{})
	if err != nil {
		return nil, err
	}
	if len(aliases) == 0 {
		contents = append(contents, []byte("\n// Add a library (local directory or Git URL):\n//   mogent source add <name> <path-or-url> --file "+path+"\n")...)
	}
	contents = append(contents, []byte("\n// Source declarations do not select modules. Choose includes before building.\n// Inspect available content with: mogent source show --file "+path+"\n")...)
	result := &InitResult{Path: path, Kind: kind, Content: contents}
	if options.DryRun {
		return result, nil
	}
	output, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return nil, fmt.Errorf("create %s: %w", path, err)
	}
	if _, err := output.Write(contents); err != nil {
		_ = output.Close()
		_ = os.Remove(path)
		return nil, err
	}
	if err := output.Close(); err != nil {
		_ = os.Remove(path)
		return nil, err
	}
	return result, nil
}

func sourceEntry(path, alias, location string) (*syntax.Entry, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	if strings.HasPrefix(location, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		location = filepath.Join(home, location[2:])
	}
	if strings.TrimSpace(alias) == "" || strings.TrimSpace(alias) != alias || strings.ContainsAny(alias, ":/\\") || strings.TrimSpace(location) == "" {
		return nil, fmt.Errorf("invalid source binding %q", alias)
	}
	attribute := "local"
	if strings.Contains(location, "://") {
		parsed, err := url.Parse(location)
		if err != nil || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
			return nil, fmt.Errorf("source %q must use a Git URL without credentials/query/fragment", alias)
		}
		if parsed.Scheme != "https" && parsed.Scheme != "http" && parsed.Scheme != "file" {
			return nil, fmt.Errorf("source %q requires an http(s) or file Git URL", alias)
		}
		if parsed.Scheme != "file" && parsed.Host == "" || parsed.Scheme == "file" && parsed.Path == "" {
			return nil, fmt.Errorf("source %q has an incomplete Git URL", alias)
		}
		attribute = "git"
	} else if filepath.IsAbs(location) {
		location, err = filepath.Rel(filepath.Dir(abs), location)
		if err != nil {
			return nil, err
		}
		location = filepath.ToSlash(location)
	}
	return block("source", alias, fields(fieldStringEntry(attribute, location, false))), nil
}
