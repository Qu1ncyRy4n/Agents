package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/cli"
)

func TestRunBuildWritesConfiguredOutput(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "identity.md"), []byte("# Identity\nHello.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  local: library\noutput: generated.md\ndoc:\n  - heading: Agent\n    from: [local:identity]\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"build", "-manifest", manifestPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	output, err := os.ReadFile(filepath.Join(temporary, "generated.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(output) != "# Agent\n\nHello.\n" {
		t.Fatalf("output = %q", output)
	}
}

func TestRunStatusReportsWorkspaceState(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte("# Identity\nHello.\n\n# Instructions\nWork.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"status", "-manifest", manifestPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Manifest: " + manifestPath,
		"Output: " + filepath.Join(temporary, "AGENTS.md"),
		"Output status: missing",
		"Sources: 1",
		"- shared: library (2 nodes)",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("status missing %q:\n%s", expected, stdout.String())
		}
	}
}

func TestRunCoverageReportsUnusedSourceNodes(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte("# Identity\nHello.\n\n# Instructions\n\n## Workflow\nWork.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"coverage", "--manifest", manifestPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Source shared: library",
		"Included: 1/3",
		"Unused:",
		"- Instructions  shared:instructions",
		"- Workflow  shared:instructions/workflow",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("coverage missing %q:\n%s", expected, stdout.String())
		}
	}
}

func TestRunCoverageUnusedOnlyReportsCompactList(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte("# Identity\nHello.\n\n# Instructions\n\n## Workflow\nWork.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"coverage", "--manifest", manifestPath, "--unused-only"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"shared  library  unused 2/3",
		"  shared:instructions  Instructions",
		"  shared:instructions/workflow  Workflow",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("unused coverage missing %q:\n%s", expected, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "Included:") {
		t.Fatalf("unused-only output should not include detailed summary:\n%s", stdout.String())
	}
}

func TestRunCoverageFiltersSourceAndContentOnly(t *testing.T) {
	temporary := t.TempDir()
	firstLibrary := filepath.Join(temporary, "first")
	secondLibrary := filepath.Join(temporary, "second")
	if err := os.Mkdir(firstLibrary, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(secondLibrary, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(firstLibrary, "core.md"), []byte("# Parent\n\n## Child\nBody.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(secondLibrary, "core.md"), []byte("# Other\nBody.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  first: first\n  second: second\noutput: AGENTS.md\ndoc:\n  - Other: second:other\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"coverage", "--manifest", manifestPath, "--source", "first", "--content-only", "--unused-only"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	output := stdout.String()
	if strings.Contains(output, "second") {
		t.Fatalf("coverage should not include filtered source:\n%s", output)
	}
	if strings.Contains(output, "first:parent  Parent") {
		t.Fatalf("content-only should hide empty parent:\n%s", output)
	}
	if !strings.Contains(output, "first:parent/child  Child") {
		t.Fatalf("coverage should include content-bearing child:\n%s", output)
	}
}

func TestRunCoverageTreeOutput(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte("# Identity\nHello.\n\n# Instructions\n\n## Workflow\nWork.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"coverage", "--manifest", manifestPath, "--tree"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"shared  library  unused 2/3",
		"`-- Instructions  shared:instructions",
		"|   `-- Workflow  shared:instructions/workflow",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("tree coverage missing %q:\n%s", expected, stdout.String())
		}
	}
}

func TestRunCoverageLeavesOnlyAndDepth(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "# Root\n\n## Branch\n\n### Leaf\nLeaf body.\n\n# Other\nOther body.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Other: shared:other\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"coverage", "--manifest", manifestPath, "--leaves-only", "--unused-only"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	output := stdout.String()
	if strings.Contains(output, "shared:root  Root") || strings.Contains(output, "shared:root/branch  Branch") {
		t.Fatalf("leaves-only should hide parent nodes:\n%s", output)
	}
	if !strings.Contains(output, "shared:root/branch/leaf  Leaf") {
		t.Fatalf("leaves-only should include terminal leaf:\n%s", output)
	}

	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"coverage", "--manifest", manifestPath, "--depth", "1", "--unused-only"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	output = stdout.String()
	if !strings.Contains(output, "shared:root  Root") || !strings.Contains(output, "shared:root/branch  Branch") {
		t.Fatalf("depth should include root and first child:\n%s", output)
	}
	if strings.Contains(output, "shared:root/branch/leaf  Leaf") {
		t.Fatalf("depth should hide deeper leaf:\n%s", output)
	}
}

func TestRunSourceShowReportsLocationAndContent(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "# Go\n\n## Development\nRun gofmt.\nRun tests.\n\n## Tests\nKeep deterministic.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "go.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Development: shared:go/development\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"source", "show", "shared:go/development", "--manifest", manifestPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Source: shared",
		"Reference: shared:go/development",
		"File: " + filepath.Join(libraryPath, "go.md"),
		"Line: 3",
		"Heading: Development",
		"Run gofmt.\nRun tests.",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("source show missing %q:\n%s", expected, stdout.String())
		}
	}
}

func TestRunSourceShowContentFlags(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "go.md"), []byte("# Go\n\n## Development\nOne.\nTwo.\nThree.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Development: shared:go/development\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"source", "show", "shared:go/development", "--manifest", manifestPath, "--file=false", "--line=false", "--content=snippet", "--lines=2"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	output := stdout.String()
	for _, unexpected := range []string{"File:", "Line:", "Three."} {
		if strings.Contains(output, unexpected) {
			t.Fatalf("source show should not include %q:\n%s", unexpected, output)
		}
	}
	if !strings.Contains(output, "One.\nTwo.") {
		t.Fatalf("snippet missing expected lines:\n%s", output)
	}

	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"source", "show", "shared:go/development", "--manifest", manifestPath, "--content=none"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if strings.Contains(stdout.String(), "One.") {
		t.Fatalf("content=none should hide body:\n%s", stdout.String())
	}
}

func TestRunAddWritesManifestByDefault(t *testing.T) {
	_, manifestPath := writeAddCLIFixture(t)

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"add", "shared:instructions/testing", "--manifest", manifestPath, "--under", "Instructions"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Wrote manifest") || !strings.Contains(stdout.String(), "Added: Testing") {
		t.Fatalf("add output =\n%s", stdout.String())
	}
	manifestContents, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifestContents), "Testing: shared:instructions/testing") {
		t.Fatalf("manifest missing added entry:\n%s", manifestContents)
	}
}

func TestRunAddDryRunDoesNotWrite(t *testing.T) {
	_, manifestPath := writeAddCLIFixture(t)
	original, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"add", "shared:instructions/testing", "--manifest", manifestPath, "--under", "Instructions", "--heading", "Tests", "--dry-run"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Dry run: no files written") || !strings.Contains(stdout.String(), "Add: Instructions / Tests") {
		t.Fatalf("dry-run output =\n%s", stdout.String())
	}
	if strings.Contains(stdout.String(), "Rendered preview:") {
		t.Fatalf("default dry-run should not show full preview:\n%s", stdout.String())
	}
	after, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(after) != string(original) {
		t.Fatalf("dry run changed manifest:\n%s", after)
	}
}

func TestRunAddDryRunPreviewModes(t *testing.T) {
	_, manifestPath := writeAddCLIFixture(t)

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"add", "shared:instructions/testing", "--manifest", manifestPath, "--under", "Instructions", "--heading", "Tests", "--dry-run", "--preview=patch"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"--- agents.yaml",
		"+++ agents.yaml",
		"@@ -8,0 +9,1 @@",
		"+      - Tests: shared:instructions/testing",
		"--- AGENTS.md",
		"@@ -9,0 +10,4 @@",
		"+## Tests",
		"+Test deterministically.",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("patch preview missing %q:\n%s", expected, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"add", "shared:instructions/testing", "--manifest", manifestPath, "--under", "Instructions", "--dry-run", "--preview=tree"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Document tree:",
		"  `- Instructions",
		"     |- Workflow  <shared:instructions/workflow>",
		"+    `- Testing  <shared:instructions/testing>",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("tree preview missing %q:\n%s", expected, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"add", "shared:instructions/testing", "--manifest", manifestPath, "--under", "Instructions", "--heading", "Tests", "--dry-run", "--preview=full"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Rendered preview:") || !strings.Contains(stdout.String(), "## Tests") {
		t.Fatalf("full preview missing rendered content:\n%s", stdout.String())
	}
}

func writeAddCLIFixture(t *testing.T) (string, string) {
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
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Identity: shared:identity\n  - heading: Instructions\n    children:\n      - Workflow: shared:instructions/workflow\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return temporary, manifestPath
}
