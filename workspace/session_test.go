package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/manifest"
)

func TestSessionMarksAndDiscardsDraftChanges(t *testing.T) {
	_, manifestPath := writeSessionFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	session.Draft.Doc = append(session.Draft.Doc, manifest.Entry{
		Heading: "Instructions",
		From:    []string{"shared:instructions"},
	})
	if err := session.MarkDraftChanged(); err != nil {
		t.Fatal(err)
	}
	if !session.Dirty {
		t.Fatal("expected dirty draft")
	}
	if !strings.Contains(session.Output, "# Instructions") {
		t.Fatalf("preview missing draft content:\n%s", session.Output)
	}
	if err := session.DiscardDraft(); err != nil {
		t.Fatal(err)
	}
	if session.Dirty {
		t.Fatal("expected discard to clear dirty state")
	}
	if strings.Contains(session.Output, "# Instructions") {
		t.Fatalf("discarded preview still contains draft content:\n%s", session.Output)
	}
}

func TestSessionSaveAndBuildWritesManifestAndOutput(t *testing.T) {
	temporary, manifestPath := writeSessionFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	session.Draft.Doc = append(session.Draft.Doc, manifest.Entry{
		Heading: "Instructions",
		From:    []string{"shared:instructions"},
	})
	if err := session.MarkDraftChanged(); err != nil {
		t.Fatal(err)
	}
	if err := session.SaveAndBuild(); err != nil {
		t.Fatal(err)
	}
	if session.Dirty {
		t.Fatal("expected save to clear dirty state")
	}
	manifestContents, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifestContents), "Instructions: shared:instructions") {
		t.Fatalf("saved manifest missing draft entry:\n%s", manifestContents)
	}
	output, err := os.ReadFile(filepath.Join(temporary, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "# Identity\n\nHello.\n\n# Instructions\n\nWork carefully.\n" {
		t.Fatalf("output = %q", output)
	}
}

func TestSessionStatusReportsOutputAndSources(t *testing.T) {
	_, manifestPath := writeSessionFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	status, err := session.Status()
	if err != nil {
		t.Fatal(err)
	}
	if status.Output != StatusMissing {
		t.Fatalf("initial output status = %s", status.Output)
	}
	if len(status.Sources) != 1 || status.Sources[0].Alias != "shared" || status.Sources[0].Nodes != 2 {
		t.Fatalf("sources = %#v", status.Sources)
	}
	if err := session.SaveAndBuild(); err != nil {
		t.Fatal(err)
	}
	status, err = session.Status()
	if err != nil {
		t.Fatal(err)
	}
	if status.Output != StatusUpToDate {
		t.Fatalf("saved output status = %s", status.Output)
	}
	if err := os.WriteFile(status.OutputPath, []byte("direct edit\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	status, err = session.Status()
	if err != nil {
		t.Fatal(err)
	}
	if status.Output != StatusDirectEdits {
		t.Fatalf("edited output status = %s", status.Output)
	}
}

func TestSessionCoverageReportsUnusedSourceNodes(t *testing.T) {
	_, manifestPath := writeSessionFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	coverage := session.Coverage()
	if len(coverage.Sources) != 1 {
		t.Fatalf("sources = %#v", coverage.Sources)
	}
	source := coverage.Sources[0]
	if source.Included != 1 || source.Total != 2 {
		t.Fatalf("coverage = %#v", source)
	}
	if len(source.Unused) != 1 || source.Unused[0].Path != "instructions" {
		t.Fatalf("unused = %#v", source.Unused)
	}
}

func TestSessionCoverageCountsSubtreeAndExclude(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "# Instructions\n\n## Workflow\nWork.\n\n## Testing\nTest.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	config := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - heading: Instructions\n    from: [shared:instructions]\n    exclude: [shared:instructions/testing]\n"
	if err := os.WriteFile(manifestPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	sourceCoverage := session.Coverage().Sources[0]
	if sourceCoverage.Included != 2 || sourceCoverage.Total != 3 {
		t.Fatalf("coverage = %#v", sourceCoverage)
	}
	if len(sourceCoverage.Unused) != 0 {
		t.Fatalf("unused = %#v", sourceCoverage.Unused)
	}
	var excluded bool
	for _, node := range sourceCoverage.Nodes {
		if node.Path == "instructions/testing" && node.State == CoverageExcluded {
			excluded = true
		}
	}
	if !excluded {
		t.Fatalf("coverage nodes missing excluded testing node: %#v", sourceCoverage.Nodes)
	}
}

func TestSessionCoverageAggregatesExcludedDirectorySelectionAsPartial(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library", "engineering")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "errors.md"), []byte("# Errors\nWrap them.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "testing.md"), []byte("# Testing\nTest them.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	config := "sources:\n  shared: library\ndoc:\n  - heading: Engineering\n    from: [shared:engineering]\n    exclude: [shared:engineering/testing]\n"
	if err := os.WriteFile(manifestPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	coverage := session.Coverage()
	states := make(map[string]CoverageState)
	for _, node := range coverage.Sources[0].Nodes {
		states[node.Path] = node.State
	}
	if states["engineering"] != CoveragePartial || states["engineering/errors"] != CoverageInherited || states["engineering/testing"] != CoverageExcluded {
		t.Fatalf("states = %#v", states)
	}
}

func TestSessionSourceNodeReturnsLocationAndSubtreeContent(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "# Go\n\n## Development\nRun gofmt.\n\n### Errors\nHandle errors.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "go.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	config := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Development: shared:go/development\n"
	if err := os.WriteFile(manifestPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	node, err := session.SourceNode("shared:go/development")
	if err != nil {
		t.Fatal(err)
	}
	if node.File != filepath.Join(libraryPath, "go.md") || node.Line != 3 {
		t.Fatalf("location = %s:%d", node.File, node.Line)
	}
	if !strings.Contains(node.Content, "Run gofmt.") || !strings.Contains(node.Content, "## Errors\n\nHandle errors.") {
		t.Fatalf("content = %q", node.Content)
	}
}

func TestSessionAddSourceDryRunDoesNotWriteManifest(t *testing.T) {
	_, manifestPath := writeAddFixture(t)
	original, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := session.AddSource(AddOptions{
		Reference: "shared:instructions/testing",
		Under:     "Instructions",
		Heading:   "Tests",
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	if result.WroteManifest || !strings.Contains(result.Preview, "## Tests") {
		t.Fatalf("result = %#v", result)
	}
	after, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(original) {
		t.Fatalf("dry run changed manifest:\n%s", after)
	}
}

func TestSessionAddSourceWritesManifestOnlyByDefault(t *testing.T) {
	temporary, manifestPath := writeAddFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := session.AddSource(AddOptions{
		Reference: "shared:instructions/testing",
		Under:     "Instructions",
	}, false)
	if err != nil {
		t.Fatal(err)
	}
	if !result.WroteManifest || result.Rebuilt {
		t.Fatalf("result = %#v", result)
	}
	manifestContents, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifestContents), "Testing: shared:instructions/testing") {
		t.Fatalf("manifest missing added entry:\n%s", manifestContents)
	}
	if _, err := os.Stat(filepath.Join(temporary, "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatalf("AGENTS.md should not be rebuilt by default, err=%v", err)
	}
}

func TestSessionAddSourceRebuildUsesOutputProtection(t *testing.T) {
	temporary, manifestPath := writeAddFixture(t)
	if err := os.WriteFile(filepath.Join(temporary, "AGENTS.md"), []byte("hand edit\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = session.AddSource(AddOptions{
		Reference: "shared:instructions/testing",
		Under:     "Instructions",
		Rebuild:   true,
	}, false)
	if err == nil || !strings.Contains(err.Error(), "untracked output") {
		t.Fatalf("expected overwrite protection, got %v", err)
	}
}

func TestSessionAddSourceRejectsDuplicateSiblingHeading(t *testing.T) {
	_, manifestPath := writeAddFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = session.AddSource(AddOptions{
		Reference: "shared:instructions/testing",
		Under:     "Instructions",
		Heading:   "Workflow",
	}, true)
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected duplicate heading error, got %v", err)
	}
}

func TestSessionAddSourceAppendAddsTopLevelEntry(t *testing.T) {
	_, manifestPath := writeAddFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := session.AddSource(AddOptions{
		Reference: "shared:instructions/testing",
		Append:    true,
		Heading:   "Testing",
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	if result.ParentPath != "" || !strings.Contains(result.Preview, "# Testing") {
		t.Fatalf("result = %#v", result)
	}
}

func TestSessionAddSourceSupportsStableRelativePlacement(t *testing.T) {
	_, manifestPath := writeAddFixture(t)
	tests := []struct {
		name    string
		options AddOptions
		want    []string
	}{
		{
			name:    "first child",
			options: AddOptions{Reference: "shared:instructions/testing", Under: "Instructions", First: true, Heading: "Testing"},
			want:    []string{"Testing", "Workflow"},
		},
		{
			name:    "before top-level anchor",
			options: AddOptions{Reference: "shared:instructions/testing", Before: "Identity", Heading: "Testing"},
			want:    []string{"Testing", "Identity", "Instructions"},
		},
		{
			name:    "after nested anchor",
			options: AddOptions{Reference: "shared:instructions/testing", After: "Instructions/Workflow", Heading: "Testing"},
			want:    []string{"Workflow", "Testing"},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			session, err := New(manifestPath)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := session.AddSource(test.options, true); err != nil {
				t.Fatal(err)
			}
			entries := session.Draft.Doc
			if test.name == "first child" || test.name == "after nested anchor" {
				entries = entries[1].Children
			}
			var headings []string
			for _, entry := range entries {
				headings = append(headings, entry.Heading)
			}
			if strings.Join(headings, ",") != strings.Join(test.want, ",") {
				t.Fatalf("headings = %#v, want %#v", headings, test.want)
			}
		})
	}
}

func TestSessionAddDirectoryPreviewExpandsSourceAndReportsRelations(t *testing.T) {
	temporary := t.TempDir()
	sharedPath := filepath.Join(temporary, "shared", "engineering")
	personalPath := filepath.Join(temporary, "personal", "engineering")
	if err := os.MkdirAll(sharedPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(personalPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(sharedPath, "testing.md"), []byte("# Test Strategy\nShared tests.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(personalPath, "migration.md"), []byte("# Staged Migration\nMigrate safely.\n\n## Cutover\nReview cutover.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	config := "sources:\n  shared: shared\n  personal: personal\ndoc:\n  - heading: Instructions\n    children:\n      - Tests: shared:engineering/test-strategy\n"
	if err := os.WriteFile(manifestPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := session.AddSource(AddOptions{Reference: "personal:engineering", Under: "Instructions", Heading: "Engineering"}, true)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"/ engineering/", "# Staged Migration", "# Cutover"} {
		if !strings.Contains(result.SourceTree, expected) {
			t.Fatalf("source tree missing %q:\n%s", expected, result.SourceTree)
		}
	}
	if len(result.Relations) != 1 || result.Relations[0].Reference != "shared:engineering/test-strategy" || result.Relations[0].Exact {
		t.Fatalf("relations = %#v", result.Relations)
	}
	session, err = New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err = session.AddSource(AddOptions{Reference: "shared:engineering", Under: "Instructions", Heading: "Shared Engineering"}, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Relations) != 1 || !result.Relations[0].Exact || result.Relations[0].Reference != "shared:engineering/test-strategy" {
		t.Fatalf("exact relations = %#v", result.Relations)
	}
}

func writeSessionFixture(t *testing.T) (string, string) {
	t.Helper()
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "# Identity\nHello.\n\n# Instructions\nWork carefully.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	config := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n"
	if err := os.WriteFile(manifestPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	return temporary, manifestPath
}

func writeAddFixture(t *testing.T) (string, string) {
	t.Helper()
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "# Identity\nHello.\n\n# Instructions\n\n## Workflow\nWork carefully.\n\n## Testing\nTest deterministically.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	config := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n  - heading: Instructions\n    children:\n      - Workflow: shared:instructions/workflow\n"
	if err := os.WriteFile(manifestPath, []byte(config), 0o644); err != nil {
		t.Fatal(err)
	}
	return temporary, manifestPath
}
