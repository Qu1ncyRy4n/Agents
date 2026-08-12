package workspace

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/Qu1ncyRy4n/Agents/internal/library"
	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
	"github.com/Qu1ncyRy4n/Agents/internal/render"
	"github.com/Qu1ncyRy4n/Agents/internal/sourcecache"
	"github.com/Qu1ncyRy4n/Agents/internal/state"
	"gopkg.in/yaml.v3"
)

const provenanceVersion = 1

type LocalizeOptions struct {
	ManifestHeading string
	From            string
	Content         *string
	DryRun          bool
	Rebuild         bool
	ForceOutput     bool
}

type LocalizeResult struct {
	ManifestHeading string
	SourceReference string
	LocalReference  string
	LocalPath       string
	ProvenancePath  string
	Preview         string
	WroteManifest   bool
	Rebuilt         bool
}

type provenanceFile struct {
	Version   int                         `yaml:"version"`
	Overrides map[string]provenanceRecord `yaml:"overrides"`
}

type provenanceRecord struct {
	SourceReference string    `yaml:"source_ref"`
	SourceLocation  string    `yaml:"source_location"`
	SourceFile      string    `yaml:"source_file"`
	LocalizedAt     time.Time `yaml:"localized_at"`
	OriginalSHA256  string    `yaml:"original_sha256"`
}

// Localize creates an ordinary local Markdown source and changes exactly one
// manifest entry to use it. Shared source files are read-only inputs.
func (s *Session) Localize(options LocalizeOptions) (*LocalizeResult, error) {
	heading := strings.TrimSpace(options.ManifestHeading)
	if heading == "" {
		return nil, fmt.Errorf("localize requires a manifest heading path")
	}
	matches := findEntryPaths(s.Draft.Doc, splitManifestPath(heading), nil)
	if len(matches) == 0 {
		return nil, fmt.Errorf("manifest heading path %q was not found", heading)
	}
	if len(matches) > 1 {
		return nil, fmt.Errorf("manifest heading path %q is ambiguous", heading)
	}
	entry := matches[0].Entry
	if len(entry.From) == 0 {
		return nil, fmt.Errorf("manifest heading path %q contains children and has no source to localize", heading)
	}
	reference, index, err := localizationReference(entry.From, options.From)
	if err != nil {
		return nil, err
	}
	alias, sourcePath, _ := manifest.SplitReference(reference)
	if alias == "local" {
		return nil, fmt.Errorf("manifest heading path %q already uses local source %q", heading, reference)
	}
	node, err := s.SourceNode(reference)
	if err != nil {
		return nil, err
	}
	localRoot := filepath.Join(filepath.Dir(s.ManifestPath), ".mogent", "library")
	if err := validateLocalSource(s.Draft, s.ManifestPath, localRoot); err != nil {
		return nil, err
	}
	localPath := filepath.Join(localRoot, filepath.FromSlash(sourcePath)+".md")
	if _, err := os.Stat(localPath); err == nil {
		return nil, fmt.Errorf("local override %q already exists; edit it explicitly instead of overwriting it", localPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("inspect local override %q: %w", localPath, err)
	}
	originalMarkdown, err := localizedMarkdown(node, nil)
	if err != nil {
		return nil, err
	}
	markdown, err := localizedMarkdown(node, options.Content)
	if err != nil {
		return nil, err
	}
	localReference := "local:" + sourcePath
	entry.From[index] = localReference
	entry.Exclude = localizeExclusions(entry.Exclude, alias, sourcePath)
	s.Draft.Sources["local"] = manifest.Source{Location: ".mogent/library"}

	temporaryRoot, err := os.MkdirTemp("", "mogent-localize-")
	if err != nil {
		return nil, fmt.Errorf("create localization preview directory: %w", err)
	}
	defer func() { _ = os.RemoveAll(temporaryRoot) }()
	if err := copyMarkdownTree(localRoot, temporaryRoot); err != nil {
		return nil, err
	}
	temporaryPath := filepath.Join(temporaryRoot, filepath.FromSlash(sourcePath)+".md")
	if err := writeNewFile(temporaryPath, markdown); err != nil {
		return nil, err
	}
	previewManifest := s.Draft.Clone()
	previewManifest.Sources["local"] = manifest.Source{Location: temporaryRoot}
	preview, err := render.Build(previewManifest, s.ManifestPath)
	if err != nil {
		return nil, fmt.Errorf("validate localized render: %w", err)
	}

	provenancePath := filepath.Join(filepath.Dir(s.ManifestPath), ".mogent", "provenance.yaml")
	result := &LocalizeResult{
		ManifestHeading: heading,
		SourceReference: reference,
		LocalReference:  localReference,
		LocalPath:       localPath,
		ProvenancePath:  provenancePath,
		Preview:         preview.Content,
	}
	if options.DryRun {
		return result, nil
	}
	if options.Rebuild {
		if err := state.CheckOverwrite(s.OutputPath(), s.StatePath(), options.ForceOutput); err != nil {
			return nil, err
		}
	}
	record, err := s.provenanceRecord(reference, node, originalMarkdown)
	if err != nil {
		return nil, err
	}
	provenance, err := loadProvenance(provenancePath)
	if err != nil {
		return nil, err
	}
	if _, exists := provenance.Overrides[localReference]; exists {
		return nil, fmt.Errorf("provenance for %q already exists", localReference)
	}
	provenance.Overrides[localReference] = record
	provenanceBytes, err := marshalProvenance(provenance)
	if err != nil {
		return nil, err
	}
	manifestBytes, err := manifest.Marshal(s.Draft)
	if err != nil {
		return nil, err
	}
	paths := []string{localPath, provenancePath, s.ManifestPath}
	if options.Rebuild {
		paths = append(paths, s.OutputPath(), s.StatePath())
	}
	snapshots, err := snapshotPaths(paths)
	if err != nil {
		return nil, err
	}
	writeErr := writeLocalization(localPath, markdown, provenancePath, provenanceBytes, s.ManifestPath, manifestBytes)
	if writeErr == nil && options.Rebuild {
		writeErr = render.WriteAtomically(s.OutputPath(), preview.Content)
		if writeErr == nil {
			writeErr = state.Write(s.StatePath(), preview.Content)
		}
	}
	if writeErr != nil {
		return nil, rollbackPaths(writeErr, snapshots)
	}
	s.Saved = s.Draft.Clone()
	s.Draft = s.Saved.Clone()
	s.Output = preview.Content
	s.Dirty = false
	s.Sources, err = render.LoadSources(s.Saved, s.ManifestPath)
	if err != nil {
		return nil, err
	}
	result.WroteManifest = true
	result.Rebuilt = options.Rebuild
	return result, nil
}

func localizationReference(references []string, requested string) (string, int, error) {
	if requested == "" {
		if len(references) != 1 {
			return "", 0, fmt.Errorf("entry composes %d sources; use --from to choose one", len(references))
		}
		return references[0], 0, nil
	}
	for index, reference := range references {
		if reference == requested {
			return reference, index, nil
		}
	}
	return "", 0, fmt.Errorf("entry does not use source reference %q", requested)
}

func localizedMarkdown(node *SourceNode, contentOverride *string) ([]byte, error) {
	var output bytes.Buffer
	if hasMetadata(node.Metadata) {
		metadata, err := yaml.Marshal(node.Metadata)
		if err != nil {
			return nil, fmt.Errorf("encode local metadata: %w", err)
		}
		output.WriteString("---\n")
		output.Write(metadata)
		output.WriteString("---\n")
	}
	output.WriteString("# ")
	output.WriteString(node.Heading)
	output.WriteString("\n\n")
	content := node.Content
	if contentOverride != nil {
		content = *contentOverride
	}
	if content := strings.TrimSpace(content); content != "" {
		output.WriteString(content)
		output.WriteByte('\n')
	}
	return output.Bytes(), nil
}

func hasMetadata(metadata library.Metadata) bool {
	return len(metadata.Tags) > 0 || metadata.TLDR != "" || metadata.Priority != nil || metadata.Scope != "" || len(metadata.Requires) > 0 || len(metadata.ConflictsWith) > 0
}

func localizeExclusions(exclusions []string, alias, root string) []string {
	localized := append([]string(nil), exclusions...)
	prefix := alias + ":" + root + "/"
	for index, exclusion := range localized {
		if strings.HasPrefix(exclusion, prefix) {
			localized[index] = "local:" + strings.TrimPrefix(exclusion, alias+":")
		}
	}
	return localized
}

func validateLocalSource(value *manifest.Manifest, manifestPath, expected string) error {
	declared, exists := value.Sources["local"]
	if !exists {
		return nil
	}
	if declared.Subdir != "" {
		return fmt.Errorf("source alias %q already uses subdir %q, not .mogent/library", "local", declared.Subdir)
	}
	resolved := declared.Location
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(filepath.Dir(manifestPath), resolved)
	}
	if filepath.Clean(resolved) != filepath.Clean(expected) {
		return fmt.Errorf("source alias %q already points to %q, not .mogent/library", "local", declared.Location)
	}
	return nil
}

func copyMarkdownTree(source, target string) error {
	if _, err := os.Stat(source); errors.Is(err, os.ErrNotExist) {
		return nil
	} else if err != nil {
		return fmt.Errorf("inspect local library: %w", err)
	}
	return filepath.WalkDir(source, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() || !strings.EqualFold(filepath.Ext(path), ".md") {
			return nil
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		contents, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		return writeNewFile(filepath.Join(target, relative), contents)
	})
}

func writeNewFile(path string, contents []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return fmt.Errorf("create local override directory: %w", err)
	}
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return fmt.Errorf("create local override %q: %w", path, err)
	}
	if _, err := file.Write(contents); err != nil {
		return errors.Join(fmt.Errorf("write local override %q: %w", path, err), file.Close())
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close local override %q: %w", path, err)
	}
	return nil
}

func (s *Session) provenanceRecord(reference string, node *SourceNode, markdown []byte) (provenanceRecord, error) {
	alias, _, _ := manifest.SplitReference(reference)
	source := s.Saved.Sources[alias]
	root := source.Location
	resolved := root
	if strings.HasPrefix(resolved, "http://") || strings.HasPrefix(resolved, "https://") {
		var err error
		resolved, err = sourcecache.Resolve(s.ManifestPath, alias, resolved, source.Subdir)
		if err != nil {
			return provenanceRecord{}, err
		}
	} else if resolved == "~" || strings.HasPrefix(resolved, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return provenanceRecord{}, fmt.Errorf("resolve source home directory: %w", err)
		}
		resolved = filepath.Join(home, strings.TrimPrefix(strings.TrimPrefix(resolved, "~"), "/"))
	}
	if !filepath.IsAbs(resolved) {
		resolved = filepath.Join(filepath.Dir(s.ManifestPath), resolved)
	}
	relativeFile, err := filepath.Rel(filepath.Clean(resolved), node.File)
	if err != nil || strings.HasPrefix(relativeFile, "..") {
		return provenanceRecord{}, fmt.Errorf("resolve source file provenance for %q", reference)
	}
	digest := sha256.Sum256(markdown)
	return provenanceRecord{
		SourceReference: reference,
		SourceLocation:  root,
		SourceFile:      filepath.ToSlash(relativeFile),
		LocalizedAt:     time.Now().UTC(),
		OriginalSHA256:  hex.EncodeToString(digest[:]),
	}, nil
}

func loadProvenance(path string) (*provenanceFile, error) {
	contents, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &provenanceFile{Version: provenanceVersion, Overrides: make(map[string]provenanceRecord)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read localization provenance: %w", err)
	}
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	decoder.KnownFields(true)
	var value provenanceFile
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("parse localization provenance: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("parse localization provenance: expected one YAML document")
	}
	if value.Version != provenanceVersion || value.Overrides == nil {
		return nil, fmt.Errorf("parse localization provenance: unsupported or incomplete version")
	}
	return &value, nil
}

func marshalProvenance(value *provenanceFile) ([]byte, error) {
	// yaml.v3 orders string map keys, keeping this file deterministic.
	var output bytes.Buffer
	encoder := yaml.NewEncoder(&output)
	encoder.SetIndent(2)
	if err := encoder.Encode(value); err != nil {
		return nil, fmt.Errorf("encode localization provenance: %w", err)
	}
	if err := encoder.Close(); err != nil {
		return nil, fmt.Errorf("finish localization provenance: %w", err)
	}
	return output.Bytes(), nil
}

func writeLocalization(localPath string, markdown []byte, provenancePath string, provenance []byte, manifestPath string, manifestBytes []byte) error {
	if err := writeNewFile(localPath, markdown); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(provenancePath), 0o755); err != nil {
		return fmt.Errorf("create provenance directory: %w", err)
	}
	if err := render.WriteAtomically(provenancePath, string(provenance)); err != nil {
		return err
	}
	return render.WriteAtomically(manifestPath, string(manifestBytes))
}

type pathSnapshot struct {
	path    string
	content []byte
	existed bool
}

func snapshotPaths(paths []string) ([]pathSnapshot, error) {
	snapshots := make([]pathSnapshot, 0, len(paths))
	for _, path := range paths {
		content, existed, err := readOptional(path)
		if err != nil {
			return nil, err
		}
		snapshots = append(snapshots, pathSnapshot{path: path, content: content, existed: existed})
	}
	return snapshots, nil
}

func rollbackPaths(original error, snapshots []pathSnapshot) error {
	var rollbackErrors []error
	for _, snapshot := range snapshots {
		if snapshot.existed {
			restoreErr := os.MkdirAll(filepath.Dir(snapshot.path), 0o755)
			if restoreErr == nil {
				restoreErr = render.WriteAtomically(snapshot.path, string(snapshot.content))
			}
			if restoreErr != nil {
				rollbackErrors = append(rollbackErrors, fmt.Errorf("restore %q: %w", snapshot.path, restoreErr))
			}
		} else if err := os.Remove(snapshot.path); err != nil && !errors.Is(err, os.ErrNotExist) {
			rollbackErrors = append(rollbackErrors, fmt.Errorf("remove new %q: %w", snapshot.path, err))
		}
	}
	sort.Slice(rollbackErrors, func(i, j int) bool { return rollbackErrors[i].Error() < rollbackErrors[j].Error() })
	return errors.Join(append([]error{original}, rollbackErrors...)...)
}
