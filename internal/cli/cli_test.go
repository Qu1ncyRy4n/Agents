package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/cli"
)

var testConfigRoot string

func TestMain(m *testing.M) {
	configRoot, err := os.MkdirTemp("", "mogent-cli-config-")
	if err != nil {
		panic(err)
	}
	if err := os.Setenv("XDG_CONFIG_HOME", configRoot); err != nil {
		panic(err)
	}
	testConfigRoot = configRoot
	code := m.Run()
	if err := os.RemoveAll(configRoot); err != nil && code == 0 {
		code = 1
	}
	os.Exit(code)
}

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

func TestRunInitGuidesAndWritesStarter(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"init", "--list-templates"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"Starter templates:", "minimal", "personal-go-nix", "research-python"} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("template guide missing %q:\n%s", expected, stdout.String())
		}
	}

	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "# Shared Baseline\n\n## Identity\n\n### Role\nAct.\n\n### Source Of Truth\nRead docs.\n\n## Instructions\n\n### Focused Change Loop\nWork.\n\n## Constraints\n\n### Safe Defaults\nSafe.\n\n## Format\n\n### Clear Handoff\nReport.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "baseline.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"init", "--template", "minimal", "--source", "shared=library", "--manifest", manifestPath, "--dry-run"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Dry run: no files written") || !strings.Contains(stdout.String(), "shared: library") {
		t.Fatalf("init dry run:\n%s", stdout.String())
	}
	if !strings.Contains(stdout.String(), "repo_name:") || !strings.Contains(stdout.String(), filepath.Base(temporary)) || !strings.Contains(stdout.String(), "repo_url: \"\"") {
		t.Fatalf("init did not materialize repository variables:\n%s", stdout.String())
	}
	if _, err := os.Stat(manifestPath); !os.IsNotExist(err) {
		t.Fatalf("init dry run wrote manifest: %v", err)
	}
	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"init", "--template", "minimal", "--source", "shared=library", "--manifest", manifestPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Wrote starter manifest") {
		t.Fatalf("init output:\n%s", stdout.String())
	}
}

func TestRunInitAcceptsExplicitRepositoryVariables(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "# Shared Baseline\n\n## Identity\n\n### Role\nWork on {{ .repo_name }} at {{ .repo_url }}.\n\n### Source Of Truth\nRead docs.\n\n## Instructions\n\n### Focused Change Loop\nWork.\n\n## Constraints\n\n### Safe Defaults\nSafe.\n\n## Format\n\n### Clear Handoff\nReport.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "baseline.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	var stdout, stderr bytes.Buffer
	err := cli.Run([]string{
		"init", "--template", "minimal", "--source", "shared=library",
		"--manifest", manifestPath, "--var", "repo_name=Example",
		"--var", "repo_url=https://example.test/repo", "--dry-run",
	}, &stdout, &stderr)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"repo_name: Example", "repo_url: https://example.test/repo"} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("init output missing %q:\n%s", expected, stdout.String())
		}
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
		"Hint: run `mogent build` to create the output",
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
		"Key: / directory  # heading",
		"# Identity      [included]  shared:identity",
		"# Instructions  [unused]    shared:instructions",
		"# Workflow  [unused]    shared:instructions/workflow",
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
		"shared:instructions",
		"shared:instructions/workflow",
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
	if !strings.Contains(output, "# Child  [unused]  first:parent/child") {
		t.Fatalf("coverage should include content-bearing child:\n%s", output)
	}
}

func TestRunCoverageFiltersByTag(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	tagged := "---\ntags: [go, testing]\n---\n# Go\n\n## Testing\nTest.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "go.md"), []byte(tagged), 0o644); err != nil {
		t.Fatal(err)
	}
	untagged := "# Docs\nDocument.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "docs.md"), []byte(untagged), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Docs: shared:docs\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"coverage", "--manifest", manifestPath, "--tag", "testing", "--unused-only"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	output := stdout.String()
	if !strings.Contains(output, "# Testing  [unused]  shared:go/testing") {
		t.Fatalf("tagged node missing:\n%s", output)
	}
	if strings.Contains(output, "shared:docs") {
		t.Fatalf("untagged node should be hidden:\n%s", output)
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
		"Key: / directory  # heading",
		"# Identity      [included]  shared:identity",
		"# Instructions  [unused]    shared:instructions",
		"# Workflow  [unused]    shared:instructions/workflow",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("tree coverage missing %q:\n%s", expected, stdout.String())
		}
	}
}

func TestCoverageIsSourceListCoveragePresetAndUnicodeIsPresentationOnly(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library", "engineering")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "testing.md"), []byte("# Testing\nTest behavior.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	if err := os.WriteFile(manifestPath, []byte("sources:\n  shared: library\ndoc:\n  - Engineering: shared:engineering\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var coverageOutput, listOutput, stderr bytes.Buffer
	if err := cli.Run([]string{"coverage", "--manifest", manifestPath}, &coverageOutput, &stderr); err != nil {
		t.Fatal(err)
	}
	if err := cli.Run([]string{"source", "list", "--manifest", manifestPath, "--coverage", "--tree"}, &listOutput, &stderr); err != nil {
		t.Fatal(err)
	}
	if coverageOutput.String() != listOutput.String() {
		t.Fatalf("coverage alias differs:\ncoverage:\n%s\nsource list:\n%s", coverageOutput.String(), listOutput.String())
	}
	listOutput.Reset()
	if err := cli.Run([]string{"source", "list", "--manifest", manifestPath, "--coverage", "--tree", "--chars", "unicode"}, &listOutput, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"▸ directory", "└──", "shared:engineering/testing"} {
		if !strings.Contains(listOutput.String(), expected) {
			t.Fatalf("unicode tree missing %q:\n%s", expected, listOutput.String())
		}
	}
}

func TestSourceListCLICharactersOverrideXDGConfig(t *testing.T) {
	configDir := filepath.Join(testConfigRoot, "mogent")
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(configDir, "config.yaml")
	if err := os.WriteFile(configPath, []byte("display:\n  chars: unicode\n  align: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(configPath) })

	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library", "lang")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "go.md"), []byte("# Go\nUse Go.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	if err := os.WriteFile(manifestPath, []byte("sources:\n  shared: library\ndoc:\n  - Go: shared:lang/go\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"source", "list", "--manifest", manifestPath, "--tree"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "▸ directory") {
		t.Fatalf("XDG character setting not applied:\n%s", stdout.String())
	}
	stdout.Reset()
	if err := cli.Run([]string{"source", "list", "--manifest", manifestPath, "--tree", "--chars", "ascii"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "/ directory") || strings.Contains(stdout.String(), "▸ directory") {
		t.Fatalf("CLI character override not applied:\n%s", stdout.String())
	}
}

func TestRunCoverageTLDRShowsFileSummaryOnce(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "---\ntldr: Prefer deliberate tests.\n---\n# Testing\n\n## Unit\nUse fixtures.\n\n## Integration\nTest boundaries.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "testing.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Unit: shared:testing/unit\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"coverage", "--manifest", manifestPath, "--tree", "--tldr"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if count := strings.Count(stdout.String(), "Prefer deliberate tests."); count != 1 {
		t.Fatalf("TLDR count = %d, want 1:\n%s", count, stdout.String())
	}
}

func TestRunSourceTreeFitsTLDRToExplicitWidth(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library", "communication")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "---\ntldr: Choose an optional communication voice without changing engineering policy or repository safety rules.\n---\n# Personas\nChoose deliberately.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "personas.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	if err := os.WriteFile(manifestPath, []byte("sources:\n  personal: library\ndoc:\n  - Personas: personal:communication/personas\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"source", "list", "personal", "--manifest", manifestPath, "--tree", "--tldr", "--width", "64"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	output := stdout.String()
	if !strings.Contains(output, "TLDR: Choose an optional communication voice") {
		t.Fatalf("fitted TLDR missing continuation:\n%s", output)
	}
	for _, line := range strings.Split(strings.TrimSpace(output), "\n") {
		if len(line) > 64 {
			t.Fatalf("line exceeds width (%d): %q\n%s", len(line), line, output)
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
	if strings.Contains(output, "shared:root\n") || strings.Contains(output, "shared:root/branch\n") {
		t.Fatalf("leaves-only should hide parent nodes:\n%s", output)
	}
	if !strings.Contains(output, "# Leaf  [unused]  shared:root/branch/leaf") {
		t.Fatalf("leaves-only should include terminal leaf:\n%s", output)
	}

	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"coverage", "--manifest", manifestPath, "--depth", "1", "--unused-only"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	output = stdout.String()
	if !strings.Contains(output, "shared:root") || !strings.Contains(output, "shared:root/branch") {
		t.Fatalf("depth should include root and first child:\n%s", output)
	}
	if strings.Contains(output, "shared:root/branch/leaf") {
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

func TestRunSourceListShowsCompactRowsAndFilters(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	goSource := "---\ntags: [lang/go, testing/unit]\ntldr: Prefer table tests.\npriority: 0.9\n---\n# Go\n\n## Testing\nTest.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "go.md"), []byte(goSource), 0o644); err != nil {
		t.Fatal(err)
	}
	docsSource := "---\ntags: [docs/readme]\ntldr: Keep docs current.\npriority: 0.3\n---\n# Docs\nDocument.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "docs.md"), []byte(docsSource), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Docs: shared:docs\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"source", "list", "--manifest", manifestPath, "--tag-search", "go", "--sort", "priority"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	output := stdout.String()
	if !strings.Contains(output, "#  shared:go/testing  Testing  lang/go,testing/unit  p=0.90  Prefer table tests.") {
		t.Fatalf("source list missing tagged row:\n%s", output)
	}
	if strings.Contains(output, "shared:docs") {
		t.Fatalf("tag-search should hide docs node:\n%s", output)
	}
}

func TestRunSourceListShowsMetadataAndLocation(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "---\ntags: [risk/security]\ntldr: Lock down dangerous changes.\npriority: 1.0\nscope: org\nrequires: [shared:workflow]\nconflicts_with: [shared:security/loose]\n---\n# Security\nBody.\n"
	sourcePath := filepath.Join(libraryPath, "security.md")
	if err := os.WriteFile(sourcePath, []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Security: shared:security\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"source", "list", "--manifest", manifestPath, "--metadata", "--file", "--line"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"#  shared:security  Security  risk/security  p=1.00  Lock down dangerous changes.  " + sourcePath + ":9",
		"  Metadata:",
		"    TLDR: Lock down dangerous changes.",
		"    Tags: risk/security",
		"    Priority: 1.00",
		"    Scope: org",
		"    Requires: shared:workflow",
		"    Conflicts with: shared:security/loose",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("source list metadata missing %q:\n%s", expected, stdout.String())
		}
	}
}

func TestRunSourceListSearchesRefsHeadingsAndBodyWithHelpfulEmptyTagSearch(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "python.md"), []byte("# Python\n\n## uv Dependency Management\nUse uv.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Python: shared:python\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"source", "list", "shared", "--manifest", manifestPath, "--search", "uv"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "shared:python/uv-dependency-management") {
		t.Fatalf("source search missing uv child:\n%s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"source", "list", "--manifest", manifestPath, "--tag-search", "python"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"No source nodes matched.",
		"--tag-search only searches metadata tags",
		"try `mogent source list --search python`",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("empty tag-search hint missing %q:\n%s", expected, stdout.String())
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

func TestRunSourceShowAlignsSourceAndSuggestsDescendants(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	engineeringPath := filepath.Join(libraryPath, "engineering")
	if err := os.Mkdir(engineeringPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(engineeringPath, "runtime.md"), []byte("# Runtime Artifacts\nKeep caches out.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Constraints: shared:engineering/runtime-artifacts\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"source", "show", "shared:engineering/runtime-artifacts", "--manifest", manifestPath, "--align-source", "--under", "Constraints"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Aligned preview: shared:engineering/runtime-artifacts",
		"Under: Constraints",
		"## Runtime Artifacts",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("aligned source missing %q:\n%s", expected, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"source", "show", "shared:engineering", "--manifest", manifestPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{"Kind: directory", "## Runtime Artifacts", "Keep caches out."} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("directory inspection missing %q:\n%s", expected, stdout.String())
		}
	}
}

func TestRunSourceShowMetadata(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	source := "---\ntags: [go, testing]\ntldr: Prefer deterministic Go tests.\npriority: 0.75\nscope: team\nrequires: [shared:workflow]\nconflicts_with: [shared:testing/fast-only]\n---\n# Go\n\n## Development\nRun gofmt.\n"
	if err := os.WriteFile(filepath.Join(libraryPath, "go.md"), []byte(source), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Development: shared:go/development\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"source", "show", "shared:go/development", "--manifest", manifestPath, "--metadata", "--content=none"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Metadata:",
		"TLDR: Prefer deterministic Go tests.",
		"Tags: go, testing",
		"Priority: 0.75",
		"Scope: team",
		"Requires: shared:workflow",
		"Conflicts with: shared:testing/fast-only",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("metadata output missing %q:\n%s", expected, stdout.String())
		}
	}
}

func TestRunCompletePrintsCandidates(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte("# Identity\nHello.\n\n# Instructions\n\n## Workflow\nWork.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - heading: Instructions\n    children:\n      - Workflow: shared:instructions/workflow\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"complete", "source-refs", "--manifest", manifestPath, "--prefix", "shared:inst"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "shared:instructions/workflow") {
		t.Fatalf("source-ref completion missing workflow:\n%s", stdout.String())
	}

	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"complete", "manifest-headings", "--manifest", manifestPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Instructions/Workflow") {
		t.Fatalf("manifest heading completion missing nested path:\n%s", stdout.String())
	}
}

func TestRunCompletionPrintsShellScripts(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"completion", "bash"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"_mogent_completion()",
		"complete -F _mogent_completion mogent",
		"mogent complete",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("bash completion script missing %q:\n%s", expected, stdout.String())
		}
	}

	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"completion", "zsh"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"#compdef mogent",
		"_mogent()",
		"mogent complete",
	} {
		if !strings.Contains(stdout.String(), expected) {
			t.Fatalf("zsh completion script missing %q:\n%s", expected, stdout.String())
		}
	}
}

func TestRunSuggestsUnknownFlag(t *testing.T) {
	var stdout, stderr bytes.Buffer
	err := cli.Run([]string{"status", "--manfiest", "agents.yaml"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected unknown flag to fail")
	}
	if !strings.Contains(err.Error(), "did you mean --manifest?") {
		t.Fatalf("missing flag suggestion:\n%v", err)
	}
}

func TestRunSourceAddPreviewsAndWrites(t *testing.T) {
	temporary, manifestPath := writeAddCLIFixture(t)
	personal := filepath.Join(temporary, "personal")
	if err := os.Mkdir(personal, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(personal, "python.md"), []byte("# Python\n\nUse Python.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"source", "add", "personal", "personal", "--manifest", manifestPath, "--dry-run"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Dry run: no files written") || !strings.Contains(stdout.String(), "Next: mogent source list personal --tree") {
		t.Fatalf("source add preview =\n%s", stdout.String())
	}
	contents, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(contents), "personal") {
		t.Fatalf("dry run changed manifest:\n%s", contents)
	}

	stdout.Reset()
	if err := cli.Run([]string{"source", "add", "personal", "personal", "--manifest", manifestPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Added source personal") {
		t.Fatalf("source add output =\n%s", stdout.String())
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

func TestRunLocalizeDryRunAndWrite(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "core.md"), []byte("# Instructions\n\n## Testing\nUse fixtures.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - Tests: shared:instructions/testing\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"localize", "Tests", "--manifest", manifestPath, "--dry-run"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Dry run: no files written") || !strings.Contains(stdout.String(), "Local: local:instructions/testing") {
		t.Fatalf("dry-run output:\n%s", stdout.String())
	}
	if _, err := os.Stat(filepath.Join(temporary, ".mogent")); !os.IsNotExist(err) {
		t.Fatalf("dry run created .mogent: %v", err)
	}

	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"localize", "Tests", "--manifest", manifestPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Manifest entry: Tests") {
		t.Fatalf("localize output:\n%s", stdout.String())
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

func TestRunAddSuggestsManifestHeadingPath(t *testing.T) {
	_, manifestPath := writeAddCLIFixture(t)

	var stdout, stderr bytes.Buffer
	err := cli.Run([]string{"add", "shared:instructions/testing", "--manifest", manifestPath, "--under", "Instrctions"}, &stdout, &stderr)
	if err == nil {
		t.Fatal("expected bad manifest heading path to fail")
	}
	if !strings.Contains(err.Error(), `did you mean "Instructions"`) {
		t.Fatalf("missing manifest heading suggestion:\n%v", err)
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
		"Inherited source subtree:",
		"`-- # Testing",
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

func TestRunAddSupportsBeforeAndFirstPlacement(t *testing.T) {
	_, manifestPath := writeAddCLIFixture(t)
	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"add", "shared:instructions/testing", "--manifest", manifestPath, "--before", "Identity", "--heading", "Tests", "--dry-run", "--preview=tree"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "+ |- Tests  <shared:instructions/testing>\n  |- Identity") {
		t.Fatalf("before preview has wrong order:\n%s", stdout.String())
	}
	stdout.Reset()
	if err := cli.Run([]string{"add", "shared:instructions/testing", "--manifest", manifestPath, "--under", "Instructions", "--first", "--heading", "Tests", "--dry-run", "--preview=tree"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "+    |- Tests  <shared:instructions/testing>\n     `- Workflow") {
		t.Fatalf("first preview has wrong order:\n%s", stdout.String())
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
