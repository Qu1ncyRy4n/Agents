package v2

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

// ResolvedSource is one declared source after its library root on disk is
// known. For git sources Commit is the exact revision read; Pinned reports
// whether that commit came from mogent.hcl or was resolved from ref just now.
type ResolvedSource struct {
	Alias      string
	Root       string
	Commit     string
	Pinned     bool
	LocalHead  string
	LocalDirty bool
}

// ResolveSources locates every declared source. A git source is read from
// .mogent/sources/<alias>/<commit>/ beside mogent.hcl and fetched once when
// that checkout is absent. A git source without commit is resolved over the
// network at its ref and reported as MOGENT206. A local source inside a Git
// work tree reports its HEAD as MOGENT207.
func ResolveSources(config *Config) (map[string]ResolvedSource, []Diagnostic, error) {
	resolved := make(map[string]ResolvedSource, len(config.Sources))
	var diagnostics []Diagnostic
	base := filepath.Dir(config.Path)
	for _, alias := range sortedSourceNames(config) {
		source := config.Sources[alias]
		entry := ResolvedSource{Alias: alias}
		if source.Local != "" {
			local := source.Local
			if !filepath.IsAbs(local) {
				local = filepath.Join(base, local)
			}
			local = filepath.Clean(local)
			entry.Root = local
			if head, dirty, ok := localGitInfo(local); ok {
				entry.LocalHead, entry.LocalDirty = head, dirty
				state := "clean"
				if dirty {
					state = "dirty"
				}
				diagnostics = append(diagnostics, Diagnostic{
					Severity: SeverityInfo,
					Code:     "MOGENT207",
					Message:  fmt.Sprintf("local source %q is at %s (%s); local sources are not pinned", alias, head, state),
				})
			}
		} else {
			cacheRoot := filepath.Join(base, ".mogent", "sources", alias)
			commit := source.Commit
			entry.Pinned = commit != ""
			if commit != "" {
				checkout := filepath.Join(cacheRoot, commit)
				if _, err := os.Stat(checkout); errors.Is(err, os.ErrNotExist) {
					if _, err := fetchGit(source.Git, commit, cacheRoot); err != nil {
						return nil, nil, fmt.Errorf("source %q: %w", alias, err)
					}
				} else if err != nil {
					return nil, nil, fmt.Errorf("source %q: inspect cache: %w", alias, err)
				}
			} else {
				ref := source.Ref
				if ref == "" {
					ref = "HEAD"
				}
				fetched, err := fetchGit(source.Git, ref, cacheRoot)
				if err != nil {
					return nil, nil, fmt.Errorf("source %q: %w", alias, err)
				}
				commit = fetched
				diagnostics = append(diagnostics, Diagnostic{
					Severity: SeverityWarning,
					Code:     "MOGENT206",
					Message:  fmt.Sprintf("source %q has no commit; resolved %s to %s, apply will pin it in %s", alias, ref, commit, filepath.Base(config.Path)),
				})
			}
			entry.Commit = commit
			entry.Root = filepath.Join(cacheRoot, commit)
		}
		if source.Subdir != "" {
			entry.Root = filepath.Join(entry.Root, filepath.FromSlash(source.Subdir))
		}
		resolved[alias] = entry
	}
	return resolved, diagnostics, nil
}

// LoadLibraries loads the sidecar library at every resolved source root.
func LoadLibraries(sources map[string]ResolvedSource) (map[string]*Library, error) {
	libraries := make(map[string]*Library, len(sources))
	for alias, source := range sources {
		library, err := LoadLibrary(source.Root)
		if err != nil {
			return nil, fmt.Errorf("source %q: %w", alias, err)
		}
		libraries[alias] = library
	}
	return libraries, nil
}

func sortedSourceNames(config *Config) []string {
	names := make([]string, 0, len(config.Sources))
	for name := range config.Sources {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// fetchGit checks out ref from remote into cacheRoot/<commit>/ without its
// .git directory and returns the resolved commit. An existing checkout for
// that commit is kept.
func fetchGit(remote, ref, cacheRoot string) (string, error) {
	if strings.TrimSpace(ref) == "" || strings.HasPrefix(ref, "-") {
		return "", fmt.Errorf("invalid git ref %q", ref)
	}
	if _, err := exec.LookPath("git"); err != nil {
		return "", fmt.Errorf("git is required to fetch %s: %w", remote, err)
	}
	if err := os.MkdirAll(cacheRoot, 0o755); err != nil {
		return "", fmt.Errorf("create source cache: %w", err)
	}
	temporary, err := os.MkdirTemp(cacheRoot, ".fetch-")
	if err != nil {
		return "", fmt.Errorf("create fetch directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(temporary) }()
	steps := []struct {
		name string
		args []string
	}{
		{"init", []string{"init", "--quiet", temporary}},
		{"remote add", []string{"-C", temporary, "remote", "add", "origin", remote}},
		{"fetch", []string{"-c", "uploadpack.allowAnySHA1InWant=true", "-C", temporary, "fetch", "--quiet", "--depth=1", "origin", ref}},
		{"checkout", []string{"-C", temporary, "checkout", "--quiet", "--detach", "FETCH_HEAD"}},
	}
	for _, step := range steps {
		command := exec.Command("git", step.args...)
		command.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
		if output, err := command.CombinedOutput(); err != nil {
			return "", fmt.Errorf("fetch %s at %s: git %s failed: %w: %s", remote, ref, step.name, err, strings.TrimSpace(string(output)))
		}
	}
	output, err := exec.Command("git", "-C", temporary, "rev-parse", "HEAD").Output()
	if err != nil {
		return "", fmt.Errorf("resolve fetched commit: %w", err)
	}
	commit := strings.TrimSpace(string(output))
	if isFullCommit(ref) && commit != ref {
		return "", fmt.Errorf("fetch %s at %s resolved to a different commit %s", remote, ref, commit)
	}
	if err := os.RemoveAll(filepath.Join(temporary, ".git")); err != nil {
		return "", fmt.Errorf("strip fetched .git directory: %w", err)
	}
	destination := filepath.Join(cacheRoot, commit)
	if _, err := os.Stat(destination); err == nil {
		return commit, nil
	}
	if err := os.Rename(temporary, destination); err != nil {
		return "", fmt.Errorf("install source checkout: %w", err)
	}
	return commit, nil
}

func localGitInfo(path string) (string, bool, bool) {
	if _, err := exec.LookPath("git"); err != nil {
		return "", false, false
	}
	head, err := exec.Command("git", "-C", path, "rev-parse", "--short", "HEAD").Output()
	if err != nil {
		return "", false, false
	}
	status, err := exec.Command("git", "-C", path, "status", "--porcelain").Output()
	if err != nil {
		return "", false, false
	}
	return strings.TrimSpace(string(head)), strings.TrimSpace(string(status)) != "", true
}

// SetCommits returns configPath's bytes with commit set on each named source
// block. Every other byte of the file, including comments, is preserved.
func SetCommits(configPath string, commits map[string]string) ([]byte, error) {
	contents, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", configPath, err)
	}
	file, diagnostics := hclwrite.ParseConfig(contents, configPath, hcl.Pos{Line: 1, Column: 1})
	if diagnostics.HasErrors() {
		return nil, fmt.Errorf("parse %s for pinning: %s", configPath, diagnostics.Error())
	}
	sources := file.Body().FirstMatchingBlock("sources", nil)
	if sources == nil {
		return nil, fmt.Errorf("%s has no sources block", configPath)
	}
	aliases := make([]string, 0, len(commits))
	for alias := range commits {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	for _, alias := range aliases {
		block := sources.Body().FirstMatchingBlock("source", []string{alias})
		if block == nil {
			return nil, fmt.Errorf("%s has no source %q to pin", configPath, alias)
		}
		block.Body().SetAttributeValue("commit", cty.StringVal(commits[alias]))
	}
	return file.Bytes(), nil
}
