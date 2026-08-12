// Package sourcecache owns immutable URL-source locks and ignored checkouts.
package sourcecache

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/Qu1ncyRy4n/Agents/internal/renderfs"
	"github.com/Qu1ncyRy4n/Agents/internal/sourcepath"
	"gopkg.in/yaml.v3"
)

const lockVersion = 1

var fetchCandidate = fetch

type Lock struct {
	Version int                  `yaml:"version"`
	Sources map[string]LockEntry `yaml:"sources"`
}

type LockEntry struct {
	URL           string `yaml:"url"`
	Subdir        string `yaml:"subdir,omitempty"`
	Commit        string `yaml:"commit"`
	ContentSHA256 string `yaml:"content_sha256"`
}

type Change struct {
	Status string
	Path   string
	Old    string
	New    string
}

type Result struct {
	Alias   string
	URL     string
	Old     *LockEntry
	New     LockEntry
	Changes []Change
	Wrote   bool
}

// Resolve returns a verified local checkout and never performs network access.
func Resolve(manifestPath, alias, sourceURL, subdir string) (string, error) {
	if err := validateRemote(alias, sourceURL); err != nil {
		return "", err
	}
	if err := sourcepath.ValidateSubdir(subdir); err != nil {
		return "", fmt.Errorf("URL source %q: %w", alias, err)
	}
	lock, err := load(lockPath(manifestPath))
	if errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("URL source %q is not pinned; run `mogent source pin %s`", alias, alias)
	}
	if err != nil {
		return "", err
	}
	entry, found := lock.Sources[alias]
	if !found {
		return "", fmt.Errorf("URL source %q is not pinned; run `mogent source pin %s`", alias, alias)
	}
	if entry.URL != sourceURL || entry.Subdir != subdir {
		return "", fmt.Errorf("URL source %q does not match mogent.lock.yaml; pin or update it explicitly", alias)
	}
	if err := validateLockEntry(alias, entry); err != nil {
		return "", err
	}
	path := cachePath(manifestPath, alias, entry.Commit)
	root, err := sourcepath.Resolve(path, entry.Subdir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("pinned cache for URL source %q is missing; run `mogent source pin %s`", alias, alias)
		}
		return "", fmt.Errorf("resolve URL source %q cache root: %w", alias, err)
	}
	hash, err := HashMarkdown(root)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("pinned cache for URL source %q is missing; run `mogent source pin %s`", alias, alias)
		}
		return "", fmt.Errorf("verify URL source %q cache: %w", alias, err)
	}
	if hash != entry.ContentSHA256 {
		return "", fmt.Errorf("pinned cache for URL source %q failed content verification; expected %s, got %s", alias, entry.ContentSHA256, hash)
	}
	return root, nil
}

// Pin installs an initial source or restores the exact existing lock commit.
func Pin(manifestPath, alias, sourceURL, subdir, requestedRef string) (*Result, error) {
	if err := validateRemote(alias, sourceURL); err != nil {
		return nil, err
	}
	lock, err := loadOptional(lockPath(manifestPath))
	if err != nil {
		return nil, err
	}
	old, exists := lock.Sources[alias]
	if err := sourcepath.ValidateSubdir(subdir); err != nil {
		return nil, fmt.Errorf("URL source %q: %w", alias, err)
	}
	if exists && (old.URL != sourceURL || old.Subdir != subdir) {
		return nil, fmt.Errorf("URL source %q identity changed; remove its lock only after reviewing the new location/subdir", alias)
	}
	ref := requestedRef
	if ref == "" && exists {
		if _, resolveErr := Resolve(manifestPath, alias, sourceURL, subdir); resolveErr == nil {
			copy := old
			return &Result{Alias: alias, URL: sourceURL, Old: &copy, New: old, Wrote: false}, nil
		}
		ref = old.Commit
	}
	if ref == "" {
		ref = "HEAD"
	}
	candidate, cleanup, err := fetchCandidate(sourceURL, ref, subdir)
	if err != nil {
		return nil, err
	}
	defer cleanup()
	if exists && requestedRef == "" && (candidate.Commit != old.Commit || candidate.ContentSHA256 != old.ContentSHA256) {
		return nil, fmt.Errorf("fetched content for locked source %q does not match its immutable lock", alias)
	}
	if exists && requestedRef != "" && candidate.Commit != old.Commit {
		return nil, fmt.Errorf("source %q is already pinned; use source update to review a changed revision", alias)
	}
	if err := install(manifestPath, alias, candidate, lock); err != nil {
		return nil, err
	}
	var oldPointer *LockEntry
	if exists {
		copy := old
		oldPointer = &copy
	}
	return &Result{Alias: alias, URL: sourceURL, Old: oldPointer, New: candidate.LockEntry, Wrote: true}, nil
}

// Update fetches a candidate. It mutates lock/cache only when accept is true.
func Update(manifestPath, alias, sourceURL, subdir, requestedRef string, accept bool) (*Result, error) {
	if err := validateRemote(alias, sourceURL); err != nil {
		return nil, err
	}
	lock, err := load(lockPath(manifestPath))
	if err != nil {
		return nil, err
	}
	old, found := lock.Sources[alias]
	if !found {
		return nil, fmt.Errorf("URL source %q is not pinned; use source pin first", alias)
	}
	if old.URL != sourceURL || old.Subdir != subdir {
		return nil, fmt.Errorf("URL source %q does not match its lock", alias)
	}
	ref := requestedRef
	if accept && !isHex(ref, 40, 64) {
		return nil, fmt.Errorf("accepting an update requires --ref with the full commit shown by the preview")
	}
	if ref == "" {
		ref = "HEAD"
	}
	candidate, cleanup, err := fetchCandidate(sourceURL, ref, subdir)
	if err != nil {
		return nil, err
	}
	defer cleanup()
	if accept && candidate.Commit != requestedRef {
		return nil, fmt.Errorf("fetched commit %s does not match accepted commit %s", candidate.Commit, requestedRef)
	}
	oldPath := cachePath(manifestPath, alias, old.Commit)
	oldRoot, err := sourcepath.Resolve(oldPath, old.Subdir)
	if err != nil {
		return nil, fmt.Errorf("resolve locked source for review: %w", err)
	}
	newRoot, err := sourcepath.Resolve(candidate.Path, candidate.Subdir)
	if err != nil {
		return nil, fmt.Errorf("resolve candidate source for review: %w", err)
	}
	changes, err := markdownChanges(oldRoot, newRoot)
	if err != nil {
		return nil, err
	}
	result := &Result{Alias: alias, URL: sourceURL, Old: &old, New: candidate.LockEntry, Changes: changes}
	if !accept {
		return result, nil
	}
	if err := install(manifestPath, alias, candidate, lock); err != nil {
		return nil, err
	}
	result.Wrote = true
	return result, nil
}

type candidate struct {
	LockEntry
	Path string
}

func fetch(sourceURL, ref, subdir string) (candidate, func(), error) {
	if strings.HasPrefix(ref, "-") || strings.TrimSpace(ref) == "" {
		return candidate{}, func() {}, fmt.Errorf("invalid source ref %q", ref)
	}
	temporary, err := os.MkdirTemp("", "mogent-source-")
	if err != nil {
		return candidate{}, func() {}, fmt.Errorf("create source fetch directory: %w", err)
	}
	cleanup := func() { _ = os.RemoveAll(temporary) }
	commands := [][]string{
		{"init", "--quiet", temporary},
		{"-C", temporary, "remote", "add", "origin", sourceURL},
		{"-C", temporary, "fetch", "--quiet", "--depth=1", "origin", ref},
		{"-C", temporary, "checkout", "--quiet", "--detach", "FETCH_HEAD"},
	}
	for _, arguments := range commands {
		command := exec.Command("git", arguments...)
		if output, err := command.CombinedOutput(); err != nil {
			cleanup()
			return candidate{}, func() {}, fmt.Errorf("fetch URL source: git %s failed: %w: %s", arguments[0], err, strings.TrimSpace(string(output)))
		}
	}
	commitOutput, err := exec.Command("git", "-C", temporary, "rev-parse", "HEAD").Output()
	if err != nil {
		cleanup()
		return candidate{}, func() {}, fmt.Errorf("resolve fetched source commit: %w", err)
	}
	commit := strings.TrimSpace(string(commitOutput))
	root, err := sourcepath.Resolve(temporary, subdir)
	if err != nil {
		cleanup()
		return candidate{}, func() {}, fmt.Errorf("resolve fetched source root: %w", err)
	}
	hash, err := HashMarkdown(root)
	if err != nil {
		cleanup()
		return candidate{}, func() {}, err
	}
	return candidate{LockEntry: LockEntry{URL: sourceURL, Subdir: subdir, Commit: commit, ContentSHA256: hash}, Path: temporary}, cleanup, nil
}

func install(manifestPath, alias string, value candidate, lock *Lock) error {
	if err := validateLockEntry(alias, value.LockEntry); err != nil {
		return err
	}
	destination := cachePath(manifestPath, alias, value.Commit)
	existingRoot, rootErr := sourcepath.Resolve(destination, value.Subdir)
	if rootErr == nil {
		existing, err := HashMarkdown(existingRoot)
		if err != nil {
			return err
		}
		if existing != value.ContentSHA256 {
			return fmt.Errorf("existing cache for %q commit %s has different content", alias, value.Commit)
		}
	} else if !errors.Is(rootErr, os.ErrNotExist) {
		return rootErr
	} else {
		parent := filepath.Dir(destination)
		if err := os.MkdirAll(parent, 0o755); err != nil {
			return fmt.Errorf("create source cache directory: %w", err)
		}
		temporary, err := os.MkdirTemp(parent, ".mogent-source-install-")
		if err != nil {
			return fmt.Errorf("create source cache transaction: %w", err)
		}
		defer func() { _ = os.RemoveAll(temporary) }()
		if err := copyMarkdown(value.Path, temporary); err != nil {
			return err
		}
		if err := os.Rename(temporary, destination); err != nil {
			return fmt.Errorf("install source cache: %w", err)
		}
	}
	lock.Sources[alias] = value.LockEntry
	return writeLock(lockPath(manifestPath), lock)
}

func HashMarkdown(root string) (string, error) {
	files, err := markdownFiles(root)
	if err != nil {
		return "", err
	}
	if len(files) == 0 {
		return "", fmt.Errorf("source %q contains no Markdown files", root)
	}
	digest := sha256.New()
	for _, relative := range files {
		contents, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil {
			return "", err
		}
		if !utf8.Valid(contents) {
			return "", fmt.Errorf("source Markdown %q is not valid UTF-8", relative)
		}
		_, _ = digest.Write([]byte(relative))
		_, _ = digest.Write([]byte{0})
		_, _ = digest.Write(contents)
		_, _ = digest.Write([]byte{0})
	}
	return hex.EncodeToString(digest.Sum(nil)), nil
}

func markdownFiles(root string) ([]string, error) {
	if _, err := os.Stat(root); err != nil {
		return nil, err
	}
	var files []string
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() && entry.Name() == ".git" {
			return filepath.SkipDir
		}
		if entry.Type()&os.ModeSymlink != 0 {
			return fmt.Errorf("source %q contains unsupported symlink %q", root, path)
		}
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(path), ".md") {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, filepath.ToSlash(relative))
		return nil
	})
	sort.Strings(files)
	return files, err
}

func copyMarkdown(source, destination string) error {
	files, err := markdownFiles(source)
	if err != nil {
		return err
	}
	for _, relative := range files {
		contents, err := os.ReadFile(filepath.Join(source, filepath.FromSlash(relative)))
		if err != nil {
			return err
		}
		path := filepath.Join(destination, filepath.FromSlash(relative))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, contents, 0o644); err != nil {
			return err
		}
	}
	return nil
}

func markdownChanges(oldRoot, newRoot string) ([]Change, error) {
	oldFiles, err := markdownMap(oldRoot)
	if err != nil {
		return nil, fmt.Errorf("read locked source for review: %w", err)
	}
	newFiles, err := markdownMap(newRoot)
	if err != nil {
		return nil, err
	}
	paths := make(map[string]bool, len(oldFiles)+len(newFiles))
	for path := range oldFiles {
		paths[path] = true
	}
	for path := range newFiles {
		paths[path] = true
	}
	ordered := make([]string, 0, len(paths))
	for path := range paths {
		ordered = append(ordered, path)
	}
	sort.Strings(ordered)
	var changes []Change
	for _, path := range ordered {
		oldContent, oldFound := oldFiles[path]
		newContent, newFound := newFiles[path]
		status := ""
		switch {
		case !oldFound:
			status = "added"
		case !newFound:
			status = "removed"
		case oldContent != newContent:
			status = "modified"
		}
		if status != "" {
			changes = append(changes, Change{Status: status, Path: path, Old: oldContent, New: newContent})
		}
	}
	return changes, nil
}

func markdownMap(root string) (map[string]string, error) {
	files, err := markdownFiles(root)
	if err != nil {
		return nil, err
	}
	result := make(map[string]string, len(files))
	for _, relative := range files {
		contents, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(relative)))
		if err != nil {
			return nil, err
		}
		if !utf8.Valid(contents) {
			return nil, fmt.Errorf("source Markdown %q is not valid UTF-8", relative)
		}
		result[relative] = string(contents)
	}
	return result, nil
}

func validateRemote(alias, sourceURL string) error {
	if alias == "" {
		return fmt.Errorf("URL source alias is empty")
	}
	for _, character := range alias {
		if !unicode.IsLetter(character) && !unicode.IsDigit(character) && !strings.ContainsRune("._-", character) {
			return fmt.Errorf("URL source alias %q is unsafe for cache paths", alias)
		}
	}
	parsed, err := url.Parse(sourceURL)
	if err != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" || parsed.User != nil {
		return fmt.Errorf("URL source %q must be an HTTP(S) Git URL without embedded credentials", alias)
	}
	return nil
}

func validateLockEntry(alias string, entry LockEntry) error {
	if err := validateRemote(alias, entry.URL); err != nil {
		return err
	}
	if !isHex(entry.Commit, 40, 64) {
		return fmt.Errorf("URL source %q lock has invalid full commit %q", alias, entry.Commit)
	}
	if !isHex(entry.ContentSHA256, 64, 64) {
		return fmt.Errorf("URL source %q lock has invalid content hash", alias)
	}
	if err := sourcepath.ValidateSubdir(entry.Subdir); err != nil {
		return fmt.Errorf("URL source %q lock: %w", alias, err)
	}
	return nil
}

func isHex(value string, minimum, maximum int) bool {
	if len(value) < minimum || len(value) > maximum {
		return false
	}
	for _, character := range value {
		if !strings.ContainsRune("0123456789abcdef", character) {
			return false
		}
	}
	return true
}

func lockPath(manifestPath string) string {
	return filepath.Join(filepath.Dir(manifestPath), "mogent.lock.yaml")
}
func cachePath(manifestPath, alias, commit string) string {
	return filepath.Join(filepath.Dir(manifestPath), ".mogent", "sources", alias, commit)
}

func loadOptional(path string) (*Lock, error) {
	value, err := load(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Lock{Version: lockVersion, Sources: make(map[string]LockEntry)}, nil
	}
	return value, err
}

func load(path string) (*Lock, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	decoder.KnownFields(true)
	var value Lock
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("parse mogent lock: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("parse mogent lock: expected one YAML document")
	}
	if value.Version != lockVersion || value.Sources == nil {
		return nil, fmt.Errorf("parse mogent lock: unsupported or incomplete version")
	}
	for alias, entry := range value.Sources {
		if err := validateLockEntry(alias, entry); err != nil {
			return nil, err
		}
	}
	return &value, nil
}

func writeLock(path string, value *Lock) error {
	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(value); err != nil {
		return fmt.Errorf("encode mogent lock: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return fmt.Errorf("finish mogent lock: %w", err)
	}
	return renderfs.WriteAtomically(path, output.Bytes())
}
