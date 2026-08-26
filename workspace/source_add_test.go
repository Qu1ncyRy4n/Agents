package workspace_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/workspace"
)

func TestAddSourceDeclarationPreviewsAndWritesLocalSource(t *testing.T) {
	temporary := t.TempDir()
	manifestPath := writeSourceAddFixture(t, temporary)
	personal := filepath.Join(temporary, "personal")
	writeSourceAddLibrary(t, personal, "python.md", "# Python\n\nUse Python.\n")
	original, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := workspace.AddSourceDeclaration(workspace.SourceAddOptions{
		ManifestPath: manifestPath,
		Alias:        "personal",
		Location:     "personal",
		DryRun:       true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.Wrote || !strings.Contains(result.ManifestYAML, "personal: personal") {
		t.Fatalf("unexpected preview: %#v", result)
	}
	afterPreview, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(afterPreview) != string(original) {
		t.Fatal("dry run changed manifest")
	}
	result, err = workspace.AddSourceDeclaration(workspace.SourceAddOptions{
		ManifestPath: manifestPath,
		Alias:        "personal",
		Location:     "personal",
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Wrote {
		t.Fatal("source declaration was not written")
	}
}

func TestAddSourceDeclarationAcceptsUnpinnedURLAndSubdir(t *testing.T) {
	manifestPath := writeSourceAddFixture(t, t.TempDir())
	result, err := workspace.AddSourceDeclaration(workspace.SourceAddOptions{
		ManifestPath: manifestPath,
		Alias:        "remote",
		Location:     "https://example.com/agents.git",
		Subdir:       "libraries/shared",
		DryRun:       true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !result.Remote || !strings.Contains(result.ManifestYAML, "subdir: libraries/shared") {
		t.Fatalf("unexpected URL preview: %#v", result)
	}
}

func TestAddSourceDeclarationRejectsDuplicatesAndInvalidSources(t *testing.T) {
	manifestPath := writeSourceAddFixture(t, t.TempDir())
	tests := []workspace.SourceAddOptions{
		{ManifestPath: manifestPath, Alias: "shared", Location: "missing"},
		{ManifestPath: manifestPath, Alias: "new", Location: "missing"},
		{ManifestPath: manifestPath, Alias: "new", Location: "ssh://example.com/repo.git"},
		{ManifestPath: manifestPath, Alias: "new", Location: "https://example.com/repo.git", Subdir: "../escape"},
	}
	for _, options := range tests {
		if _, err := workspace.AddSourceDeclaration(options); err == nil {
			t.Fatalf("expected error for %#v", options)
		}
	}
}

func writeSourceAddFixture(t *testing.T, temporary string) string {
	t.Helper()
	writeSourceAddLibrary(t, filepath.Join(temporary, "shared"), "rules.md", "# Rules\n\nBe careful.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	contents := "sources:\n  shared: shared\ndoc:\n  - Rules: shared:rules\n"
	if err := os.WriteFile(manifestPath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
	return manifestPath
}

func writeSourceAddLibrary(t *testing.T, directory, name, content string) {
	t.Helper()
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, name), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
