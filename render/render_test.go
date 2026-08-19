package render_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/manifest"
	"github.com/Qu1ncyRy4n/Agents/render"
	"github.com/Qu1ncyRy4n/Agents/sourcecache"
)

func TestBuildUsesManifestHeadingsTemplatesAndExclusions(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "rules.md"), "# Rules\nHello {{ .repo_name }}.\n\n## Keep\nKeep this.\n\n## Drop\nDrop this.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\nvars:\n  repo_name: mogent\noutput: AGENTS.md\ndoc:\n  - heading: Local Rules\n    from:\n      - shared:rules\n    exclude:\n      - shared:rules/drop\n")

	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	want := "# Local Rules\n\nHello mogent.\n\n## Keep\n\nKeep this.\n"
	if result.Content != want {
		t.Fatalf("rendered content:\n%s\nwant:\n%s", result.Content, want)
	}
}

func TestBuildSelectsDirectorySubtreeAndNarrowsItWithExclude(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library", "engineering")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "errors.md"), "# Error Handling\nWrap errors.\n")
	writeFile(t, filepath.Join(libraryPath, "tests.md"), "# Testing\nTest behavior.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\ndoc:\n  - heading: Engineering\n    from: [shared:engineering]\n    exclude: [shared:engineering/testing]\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	want := "# Engineering\n\n## Error Handling\n\nWrap errors.\n"
	if result.Content != want {
		t.Fatalf("rendered content:\n%s\nwant:\n%s", result.Content, want)
	}
}

func TestBuildRejectsDirectoryAndDescendantInSameEntry(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library", "engineering")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "tests.md"), "# Testing\nTest behavior.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\ndoc:\n  - heading: Engineering\n    from: [shared:engineering, shared:engineering/testing]\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := render.Build(value, loadedPath); err == nil || !strings.Contains(err.Error(), "overlapping source roots") {
		t.Fatalf("Build error = %v, want overlapping source roots", err)
	}
}

func TestBuildUsesVerifiedPinnedURLCacheOffline(t *testing.T) {
	temporary := t.TempDir()
	commit := strings.Repeat("d", 40)
	cache := filepath.Join(temporary, ".mogent", "sources", "remote", commit)
	if err := os.MkdirAll(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(cache, "rules.md"), "# Rules\nWork carefully.\n")
	hash, err := sourcecache.HashMarkdown(cache)
	if err != nil {
		t.Fatal(err)
	}
	lock := "version: 1\nsources:\n  remote:\n    url: https://example.com/library.git\n    commit: " + commit + "\n    content_sha256: " + hash + "\n"
	writeFile(t, filepath.Join(temporary, "mogent.lock.yaml"), lock)
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  remote: https://example.com/library.git\noutput: AGENTS.md\ndoc:\n  - Rules: remote:rules\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content != "# Rules\n\nWork carefully.\n" {
		t.Fatalf("content = %q", result.Content)
	}
}

func TestBuildUsesOnlyPinnedURLSubdirOffline(t *testing.T) {
	temporary := t.TempDir()
	commit := strings.Repeat("e", 40)
	cache := filepath.Join(temporary, ".mogent", "sources", "remote", commit)
	selected := filepath.Join(cache, "libraries", "cdint")
	if err := os.MkdirAll(selected, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(selected, "rules.md"), "# Rules\nSelected.\n")
	writeFile(t, filepath.Join(cache, "unrelated.md"), "# Rules\nDuplicate outside selected root.\n")
	hash, err := sourcecache.HashMarkdown(selected)
	if err != nil {
		t.Fatal(err)
	}
	lock := "version: 1\nsources:\n  remote:\n    url: https://example.com/repository.git\n    subdir: libraries/cdint\n    commit: " + commit + "\n    content_sha256: " + hash + "\n"
	writeFile(t, filepath.Join(temporary, "mogent.lock.yaml"), lock)
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  remote:\n    location: https://example.com/repository.git\n    subdir: libraries/cdint\ndoc:\n  - Rules: remote:rules\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content != "# Rules\n\nSelected.\n" {
		t.Fatalf("content = %q", result.Content)
	}
}

func TestBuildRequiresExplicitRepositoryTemplateVariables(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "identity.md"), "# Identity\nRepository: {{ .repo_name }}\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\ndoc:\n  - Identity: shared:identity\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = render.Build(value, loadedPath)
	if err == nil || !strings.Contains(err.Error(), "repo_name") {
		t.Fatalf("Build error = %v, want missing repo_name", err)
	}
}

func TestBuildRejectsUnresolvedAndNonDescendantExclusions(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "rules.md"), "# Rules\n## Child\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - heading: Rules\n    from: [shared:rules]\n    exclude: [shared:rules]\n")

	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = render.Build(value, loadedPath)
	if err == nil || !strings.Contains(err.Error(), "cannot exclude an entry source root") {
		t.Fatalf("Build error = %v, want root exclusion error", err)
	}
}

func TestBuildRejectsCollisionsEmptyNodesAndURLs(t *testing.T) {
	temporary := t.TempDir()
	first := filepath.Join(temporary, "first")
	second := filepath.Join(temporary, "second")
	for _, path := range []string{first, second} {
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, filepath.Join(first, "rules.md"), "# Rules\n## Child\nFirst.\n")
	writeFile(t, filepath.Join(second, "rules.md"), "# Rules\n## Child\nSecond.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  one: first\n  two: second\noutput: AGENTS.md\ndoc:\n  - heading: Rules\n    from: [one:rules, two:rules]\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := render.Build(value, loadedPath); err == nil || !strings.Contains(err.Error(), "overlapping") {
		t.Fatalf("collision error = %v", err)
	}

	writeFile(t, filepath.Join(first, "empty.md"), "# Empty\n")
	writeFile(t, manifestPath, "sources:\n  one: first\noutput: AGENTS.md\ndoc:\n  - Empty: one:empty\n")
	value, loadedPath, err = manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := render.Build(value, loadedPath); err == nil || !strings.Contains(err.Error(), "empty source node") {
		t.Fatalf("empty node error = %v", err)
	}

	writeFile(t, manifestPath, "sources:\n  remote: https://example.com/library.git\noutput: AGENTS.md\ndoc:\n  - Remote: remote:rules\n")
	value, loadedPath, err = manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := render.Build(value, loadedPath); err == nil || !strings.Contains(err.Error(), "not pinned") {
		t.Fatalf("URL error = %v", err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
