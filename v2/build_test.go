package v2

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const applyFixtureRendered = "# Intro\n\n## Role\n\nBe helpful.\n"

func applyFixture(t *testing.T) *Config {
	t.Helper()
	root := t.TempDir()
	writeLibraryFile(t, root, "library/agents/intro/role.md", "# Role\n\nBe helpful.\n")
	writeLibraryFile(t, root, "library/"+LibraryFile, `
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
`)
	writeLibraryFile(t, root, ConfigFile, `
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
`)
	config, err := Load(filepath.Join(root, ConfigFile))
	if err != nil {
		t.Fatal(err)
	}
	return config
}

func TestPlanLocalReportsNewFileWithoutWriting(t *testing.T) {
	config := applyFixture(t)
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if len(result.Changes) != 1 || result.Changes[0].Status != FileNew {
		t.Fatalf("changes = %#v", result.Changes)
	}
	if !strings.Contains(result.Changes[0].Diff, "+## Role") {
		t.Fatalf("diff = %q", result.Changes[0].Diff)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(config.Path), "AGENTS.md")); !os.IsNotExist(err) {
		t.Fatalf("plan wrote output: %v", err)
	}
	if codes := diagnosticCodes(result.Plan); len(codes) != 0 {
		t.Fatalf("diagnostics = %v", codes)
	}
}

func TestApplyWritesOutputAndStateThenPlanIsUnchanged(t *testing.T) {
	config := applyFixture(t)
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	root := filepath.Dir(config.Path)
	written, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(written) != applyFixtureRendered {
		t.Fatalf("output = %q", written)
	}
	if _, err := os.Stat(filepath.Join(root, ".mogent", "state.json")); err != nil {
		t.Fatalf("state not written: %v", err)
	}
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changes[0].Status != FileUnchanged || result.Changes[0].Diff != "" {
		t.Fatalf("second plan = %#v", result.Changes[0])
	}
	if codes := diagnosticCodes(result.Plan); len(codes) != 0 {
		t.Fatalf("diagnostics = %v", codes)
	}
}

func TestPlanWarnsAndApplyRefusesEditedOutputUntilForced(t *testing.T) {
	config := applyFixture(t)
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	output := filepath.Join(filepath.Dir(config.Path), "AGENTS.md")
	if err := os.WriteFile(output, []byte("# Intro\n\nhand edit\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if result.Changes[0].Status != FileChanged {
		t.Fatalf("status = %q", result.Changes[0].Status)
	}
	if codes := diagnosticCodes(result.Plan); len(codes) != 1 || codes[0] != "MOGENT208" {
		t.Fatalf("diagnostics = %v", codes)
	}
	if _, err := Apply(config, false); err == nil || !strings.Contains(err.Error(), "direct edits") {
		t.Fatalf("apply without force error = %v", err)
	}
	if _, err := Apply(config, true); err != nil {
		t.Fatal(err)
	}
	restored, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	if string(restored) != applyFixtureRendered {
		t.Fatalf("forced apply output = %q", restored)
	}
}

func TestPlanWarnsAndApplyRefusesUnmanagedOutputUntilForced(t *testing.T) {
	config := applyFixture(t)
	output := filepath.Join(filepath.Dir(config.Path), "AGENTS.md")
	if err := os.WriteFile(output, []byte("# Hand written\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if codes := diagnosticCodes(result.Plan); len(codes) != 1 || codes[0] != "MOGENT208" {
		t.Fatalf("diagnostics = %v", codes)
	}
	if !strings.Contains(result.Plan.Diagnostics[0].Message, "not managed") {
		t.Fatalf("message = %q", result.Plan.Diagnostics[0].Message)
	}
	if _, err := Apply(config, false); err == nil || !strings.Contains(err.Error(), "untracked") {
		t.Fatalf("apply without force error = %v", err)
	}
	if _, err := Apply(config, true); err != nil {
		t.Fatal(err)
	}
}

func diagnosticCodes(plan *Plan) []string {
	codes := make([]string, 0, len(plan.Diagnostics))
	for _, diagnostic := range plan.Diagnostics {
		codes = append(codes, diagnostic.Code)
	}
	return codes
}
