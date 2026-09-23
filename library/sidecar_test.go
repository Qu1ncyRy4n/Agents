package library_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/library"
)

func TestScanProducesExplicitExhaustiveSidecar(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(root, "agents", "rules.md"), "# Rules\n## Tests\nTest.\n")
	sidecar, err := library.Scan(root)
	if err != nil {
		t.Fatal(err)
	}
	contents, err := library.MarshalSidecar(sidecar)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"id: mogent/1", "name: \"" + filepath.Base(root) + "\"", "source: agents", "source: agents/rules", "source: agents/rules/tests"} {
		if !strings.Contains(string(contents), want) {
			t.Fatalf("scan missing %q:\n%s", want, contents)
		}
	}
	writeTestFile(t, filepath.Join(root, library.SidecarFile), string(contents))
	report, err := library.Check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Errors) != 0 || len(report.Warnings) != 0 {
		t.Fatalf("report = %#v", report)
	}
}

func TestLoadWithSidecarLeavesSidecarFreeLibraryUnchanged(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "rules.md"), "# Rules\n## First\nFirst.\n## Second\nSecond.\n")
	raw, err := library.Load(root)
	if err != nil {
		t.Fatal(err)
	}
	loaded, warnings, err := library.LoadWithSidecar(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(warnings) != 0 {
		t.Fatalf("warnings = %#v, want none", warnings)
	}
	if len(loaded.Roots) != len(raw.Roots) || loaded.Roots[0].Path != raw.Roots[0].Path || loaded.Roots[0].Heading != raw.Roots[0].Heading {
		t.Fatalf("roots changed: raw = %#v, loaded = %#v", raw.Roots, loaded.Roots)
	}
	for index, node := range raw.Roots[0].Children {
		got := loaded.Roots[0].Children[index]
		if got.Path != node.Path || got.Heading != node.Heading {
			t.Fatalf("child %d changed: raw = %#v, loaded = %#v", index, node, got)
		}
	}
}

func TestLoadWithSidecarScopesRootContentAndHidesVirtualMarkdownRoot(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(root, "agents", "intro.md"), "# Intro\nWelcome.\n")
	if err := os.MkdirAll(filepath.Join(root, "archive"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(root, "archive", "broken.md"), "---\nnot_a_mogent_field: true\n---\n# Broken\n")
	writeTestFile(t, filepath.Join(root, "library.mogent.yaml"), `schema: {id: mogent/1}
library: {name: Test}
content:
  markdown_roots: [agents]
  directory_roots: [skills]
tree:
  - source: agents/intro
`)
	if err := os.Mkdir(filepath.Join(root, "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	loaded, _, err := library.LoadWithSidecar(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, found := loaded.ByPath["agents"]; found {
		t.Fatal("configured Markdown root was retained as a selectable node")
	}
	if _, found := loaded.ByPath["archive/broken"]; found {
		t.Fatal("unconfigured archive was indexed")
	}
	if len(loaded.Roots) != 1 || loaded.Roots[0].Path != "agents/intro" {
		t.Fatalf("roots = %#v, want agents/intro", loaded.Roots)
	}
}

func TestCheckRejectsInvalidContentRoots(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "agents", "nested"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(root, "agents", "nested", "intro.md"), "# Intro\n")
	writeTestFile(t, filepath.Join(root, "not-directory"), "file\n")
	if err := os.Mkdir(filepath.Join(root, "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(root, library.SidecarFile), `schema: {id: mogent/1}
library: {name: Test}
content:
  markdown_roots: [agents, agents/nested]
  directory_roots: [skills, skills, not-directory, missing, ../outside]
tree: []
`)
	report, err := library.Check(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"overlap", "duplicate content root", "not a directory", "does not exist", "not a normalized"} {
		if !contains(report.Errors, want) {
			t.Fatalf("errors missing %q: %#v", want, report.Errors)
		}
	}
}

func TestCheckRejectsTreeOutsideMarkdownContentRoots(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "agents"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(root, "agents", "intro.md"), "# Intro\n")
	if err := os.Mkdir(filepath.Join(root, "skills"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(root, library.SidecarFile), `schema: {id: mogent/1}
library: {name: Test}
content:
  markdown_roots: [agents]
  directory_roots: [skills]
tree:
  - source: agents/intro
  - source: skills
`)
	report, err := library.Check(root)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(report.Errors, "outside declared Markdown content roots") {
		t.Fatalf("errors = %#v", report.Errors)
	}
}

func TestSidecarWithoutContentStillIndexesEntireLibrary(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "rules.md"), "# Rules\n")
	if err := os.Mkdir(filepath.Join(root, "archive"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeTestFile(t, filepath.Join(root, "archive", "history.md"), "# History\n")
	writeTestFile(t, filepath.Join(root, library.SidecarFile), `schema: {id: mogent/1}
library: {name: Test}
tree:
  - source: archive
    children:
      - source: archive/history
  - source: rules
`)
	loaded, _, err := library.LoadWithSidecar(root)
	if err != nil {
		t.Fatal(err)
	}
	if _, found := loaded.ByPath["archive/history"]; !found {
		t.Fatal("legacy sidecar did not inventory archive content")
	}
}

func TestCheckReportsSidecarContractViolationsAndOwnershipWarnings(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "rules.md"), "---\nrequires: [other]\nconflicts_with: [other]\n---\n# Rules\n## Child\nBody.\n")
	sidecar := `library:
  name: Test
tree:
  - source: rules
    requires: [missing]
    exclusive_group: unused
    tags: [wrong-place]
  - source: rules
  - source: missing-source
    requires: [missing]
  - source: rules/child
    exclusive_group: unused
groups:
  unused:
    mode: exactly-one
`
	writeTestFile(t, filepath.Join(root, library.SidecarFile), sidecar)
	report, err := library.Check(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"duplicate tree source", "unknown relationship target", "not a discovered", "nested under", "invalid mode"} {
		if !contains(report.Errors, want) {
			t.Fatalf("errors missing %q: %#v", want, report.Errors)
		}
	}
	for _, want := range []string{"schema.id is omitted", "tags or tldr", "Markdown frontmatter"} {
		if !contains(report.Warnings, want) {
			t.Fatalf("warnings missing %q: %#v", want, report.Warnings)
		}
	}
}

func TestCheckReportsRequirementCycles(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "rules.md"), "# One\n# Two\n")
	writeTestFile(t, filepath.Join(root, library.SidecarFile), "schema: {id: mogent/1}\nlibrary: {name: Test}\ntree:\n  - source: one\n    requires: [two]\n  - source: two\n    requires: [one]\n")
	report, err := library.Check(root)
	if err != nil {
		t.Fatal(err)
	}
	if !contains(report.Errors, "requires cycle") {
		t.Fatalf("errors = %#v", report.Errors)
	}
}

func contains(values []string, want string) bool {
	for _, value := range values {
		if strings.Contains(value, want) {
			return true
		}
	}
	return false
}
