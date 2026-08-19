package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/sourcecache"
)

func TestLocalizeDryRunWritesNothing(t *testing.T) {
	temporary, manifestPath, sourcePath := writeLocalizeFixture(t)
	originalManifest, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	originalSource, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := session.Localize(LocalizeOptions{ManifestHeading: "Instructions/Tests", DryRun: true})
	if err != nil {
		t.Fatal(err)
	}
	if result.LocalReference != "local:instructions/testing" || !strings.Contains(result.Preview, "Use fixtures.") {
		t.Fatalf("result = %#v", result)
	}
	currentManifest, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	currentSource, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(currentManifest) != string(originalManifest) || string(currentSource) != string(originalSource) {
		t.Fatal("dry run changed manifest or shared source")
	}
	if _, err := os.Stat(filepath.Join(temporary, ".mogent")); !os.IsNotExist(err) {
		t.Fatalf("dry run created .mogent: %v", err)
	}
}

func TestLocalizeWritesOverrideProvenanceAndManifest(t *testing.T) {
	temporary, manifestPath, sourcePath := writeLocalizeFixture(t)
	originalSource, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := session.Localize(LocalizeOptions{ManifestHeading: "Instructions/Tests", Rebuild: true})
	if err != nil {
		t.Fatal(err)
	}
	if !result.WroteManifest || !result.Rebuilt {
		t.Fatalf("result = %#v", result)
	}
	currentSource, err := os.ReadFile(sourcePath)
	if err != nil {
		t.Fatal(err)
	}
	if string(currentSource) != string(originalSource) {
		t.Fatal("localization changed shared source")
	}
	local, err := os.ReadFile(filepath.Join(temporary, ".mogent", "library", "instructions", "testing.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(local), "# Testing") || !strings.Contains(string(local), "Use fixtures.") {
		t.Fatalf("local override:\n%s", local)
	}
	provenance, err := os.ReadFile(filepath.Join(temporary, ".mogent", "provenance.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"version: 1", "local:instructions/testing:", "source_ref: shared:instructions/testing", "source_file: core.md", "original_sha256:"} {
		if !strings.Contains(string(provenance), expected) {
			t.Fatalf("provenance missing %q:\n%s", expected, provenance)
		}
	}
	manifestContents, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifestContents), "local: .mogent/library") || !strings.Contains(string(manifestContents), "Tests: local:instructions/testing") {
		t.Fatalf("manifest:\n%s", manifestContents)
	}
	output, err := os.ReadFile(filepath.Join(temporary, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != result.Preview {
		t.Fatalf("output mismatch:\n%s", output)
	}
	if _, err := New(manifestPath); err != nil {
		t.Fatalf("reload localized workspace: %v", err)
	}
}

func TestLocalizeRefusesExistingOverride(t *testing.T) {
	temporary, manifestPath, _ := writeLocalizeFixture(t)
	localPath := filepath.Join(temporary, ".mogent", "library", "instructions", "testing.md")
	if err := os.MkdirAll(filepath.Dir(localPath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(localPath, []byte("# Testing\nExisting.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := session.Localize(LocalizeOptions{ManifestHeading: "Instructions/Tests"}); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("error = %v", err)
	}
}

func TestLocalizePinnedURLSourceRecordsRemoteProvenance(t *testing.T) {
	temporary := t.TempDir()
	commit := strings.Repeat("e", 40)
	cache := filepath.Join(temporary, ".mogent", "sources", "remote", commit)
	if err := os.MkdirAll(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cache, "rules.md"), []byte("# Rules\nRemote guidance.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	hash, err := sourcecache.HashMarkdown(cache)
	if err != nil {
		t.Fatal(err)
	}
	lock := "version: 1\nsources:\n  remote:\n    url: https://example.com/library.git\n    commit: " + commit + "\n    content_sha256: " + hash + "\n"
	if err := os.WriteFile(filepath.Join(temporary, "mogent.lock.yaml"), []byte(lock), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  remote: https://example.com/library.git\noutput: AGENTS.md\ndoc:\n  - Rules: remote:rules\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := session.Localize(LocalizeOptions{ManifestHeading: "Rules"}); err != nil {
		t.Fatal(err)
	}
	provenance, err := os.ReadFile(filepath.Join(temporary, ".mogent", "provenance.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(provenance), "source_location: https://example.com/library.git") || !strings.Contains(string(provenance), "source_file: rules.md") {
		t.Fatalf("provenance:\n%s", provenance)
	}
}

func writeLocalizeFixture(t *testing.T) (string, string, string) {
	t.Helper()
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	sourcePath := filepath.Join(libraryPath, "core.md")
	source := "---\ntldr: Test deliberately.\n---\n# Instructions\n\n## Testing\nUse fixtures.\n\n### Errors\nKeep failures clear.\n"
	if err := os.WriteFile(sourcePath, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Instructions:\n      - Tests: shared:instructions/testing\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return temporary, manifestPath, sourcePath
}
