package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/cli"
)

func writeV2Fixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"library/agents/intro/role.md": "# Role\n\nBe helpful.\n",
		"library/library.mogent.hcl": `
library {
  format = 2
  id = "test"
  name = "Test"
}
content {
  markdown_root = "agents"
}
sections {
  section "intro" {
    title = "Intro"
    section "role" {
      title = "Role"
      source = "intro/role.md"
    }
  }
}
`,
		"mogent.hcl": `
mogent { format = 2 }
sources {
  source "shared" { local = "library" }
}
outputs {
  output "agents" {
    path = "AGENTS.md"
    kind = "markdown"
    source "shared" {
      from = "shared:agents"
      select = { all = true }
    }
  }
}
`,
	}
	for name, contents := range files {
		path := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(contents), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	return root, filepath.Join(root, "mogent.hcl")
}

func TestRunPlanPrintsDiffAndSummaryWithoutWriting(t *testing.T) {
	root, configPath := writeV2Fixture(t)
	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"plan", "--config", configPath}, &stdout, &stderr); err != nil {
		t.Fatalf("plan error = %v\nstderr:\n%s", err, stderr.String())
	}
	for _, want := range []string{"AGENTS.md: new\n", "+++ b/AGENTS.md\n", "+## Role\n", "Plan: 1 to add, 0 to change, 0 unchanged\n"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("plan output missing %q:\n%s", want, stdout.String())
		}
	}
	if _, err := os.Stat(filepath.Join(root, "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatalf("plan wrote output: %v", err)
	}
}

func TestRunApplyWritesOutputsThenPlanIsUnchanged(t *testing.T) {
	root, configPath := writeV2Fixture(t)
	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"apply", "--config", configPath}, &stdout, &stderr); err != nil {
		t.Fatalf("apply error = %v\nstderr:\n%s", err, stderr.String())
	}
	if !strings.Contains(stdout.String(), "Applied 1 output(s): AGENTS.md") {
		t.Fatalf("apply output = %q", stdout.String())
	}
	written, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(written) != "# Intro\n\n## Role\n\nBe helpful.\n" {
		t.Fatalf("output = %q", written)
	}
	stdout.Reset()
	stderr.Reset()
	if err := cli.Run([]string{"plan", "--config", configPath}, &stdout, &stderr); err != nil {
		t.Fatalf("plan error = %v", err)
	}
	if !strings.Contains(stdout.String(), "AGENTS.md: unchanged\n") || !strings.Contains(stdout.String(), "Plan: 0 to add, 0 to change, 1 unchanged\n") {
		t.Fatalf("second plan output = %q", stdout.String())
	}
}

func TestRunApplyRefusesEditedOutputUntilForced(t *testing.T) {
	root, configPath := writeV2Fixture(t)
	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"apply", "--config", configPath}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# edited\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	stderr.Reset()
	err := cli.Run([]string{"apply", "--config", configPath}, &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "direct edits") {
		t.Fatalf("apply error = %v", err)
	}
	if !strings.Contains(stderr.String(), "MOGENT208") {
		t.Fatalf("apply stderr = %q", stderr.String())
	}
	if err := cli.Run([]string{"apply", "--config", configPath, "--force"}, &stdout, &stderr); err != nil {
		t.Fatalf("forced apply error = %v", err)
	}
}

func TestRunV2SourceListAndShow(t *testing.T) {
	root, configPath := writeV2Fixture(t)
	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"source", "list", "--config", configPath, "--tldr"}, &stdout, &stderr); err != nil {
		t.Fatalf("source list error = %v\n%s", err, stderr.String())
	}
	for _, want := range []string{"shared: Test (local library)\n", "shared:agents\n", "  intro  Intro\n", "    role  Role\n"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("source list missing %q:\n%s", want, stdout.String())
		}
	}
	stdout.Reset()
	if err := cli.Run([]string{"source", "show", "shared:intro/role", "--config", configPath, "--lines", "1"}, &stdout, &stderr); err != nil {
		t.Fatalf("source show error = %v\n%s", err, stderr.String())
	}
	for _, want := range []string{"shared:intro/role\n", "Title:  Role\n", "Source: " + filepath.Join(root, "library", "agents", "intro", "role.md") + "\n", "---\n# Role\n... (2 more lines)\n"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("source show missing %q:\n%s", want, stdout.String())
		}
	}
	err := cli.Run([]string{"source", "show", "shared:intro/ghost", "--config", configPath}, &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "no section or tree") {
		t.Fatalf("missing section error = %v", err)
	}
	t.Chdir(root)
	stdout.Reset()
	if err := cli.Run([]string{"source", "list"}, &stdout, &stderr); err != nil {
		t.Fatalf("implicit v2 source list error = %v", err)
	}
	if !strings.Contains(stdout.String(), "shared:agents\n") {
		t.Fatalf("implicit v2 source list output:\n%s", stdout.String())
	}
}

func TestRunUpdateRequiresGitSource(t *testing.T) {
	_, configPath := writeV2Fixture(t)
	var stdout, stderr bytes.Buffer
	err := cli.Run([]string{"update", "--config", configPath}, &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "no git sources") {
		t.Fatalf("update error = %v", err)
	}
	err = cli.Run([]string{"update", "shared", "--config", configPath}, &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "not a declared git source") {
		t.Fatalf("update alias error = %v", err)
	}
}

func TestRunBuildRejectsHCLConfig(t *testing.T) {
	_, configPath := writeV2Fixture(t)
	var stdout, stderr bytes.Buffer
	err := cli.Run([]string{"build", "--manifest", configPath}, &stdout, &stderr)
	if err == nil || !strings.Contains(err.Error(), "mogent apply") {
		t.Fatalf("build error = %v", err)
	}
}
