package render_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/manifest"
	"github.com/Qu1ncyRy4n/Agents/internal/render"
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

func TestBuildProvidesRepositoryNameTemplateVariable(t *testing.T) {
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
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Content, "Repository: "+filepath.Base(temporary)) {
		t.Fatalf("missing repo_name template value: %q", result.Content)
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
	if _, err := render.Build(value, loadedPath); err == nil || !strings.Contains(err.Error(), "URL sources") {
		t.Fatalf("URL error = %v", err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
