package cli_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/cli"
)

func writeArtifactCLI(t *testing.T, extra string) (string, string) {
	t.Helper()
	root := t.TempDir()
	files := map[string]string{
		"library/library.mogent.hcl": `library {
  format = 2
  id = "typed"
  name = "Typed library"
}
content "physical" { directory = "payloads" }
`,
		"library/payloads/review/SKILL.md":           "---\nname: review\ndescription: Review a change.\n---\n# Review\n\n## Procedure\nInspect changes.\n\n### Verify\nRun checks.\n",
		"library/payloads/review/scripts/check.sh":   "#!/bin/sh\nexit 0\n",
		"library/payloads/overrides/review/SKILL.md": "---\nname: review\ndescription: Repo review.\n---\n# Repo review\n",
		"library/payloads/settings.toml":             "# retain this comment\n[tool]\nx = 1\n",
		"mogent.hcl": `mogent { format = 2 }
sources {
  source "shared" { local = "library" }
}
outputs {
  output "package" {
    path = "agent-export"
    kind = "dir-tree"
    source "shared" {
      from = "shared:physical"
      node = "review"
      operation = "copy"
      into = ".agents/skills/review"
    }
    source "shared" {
      from = "shared:physical"
      node = "review/SKILL.md"
      heading = ["Review", "Procedure"]
      operation = "render-markdown"
      into = "AGENTS.md"
    }
    source "shared" {
      from = "shared:physical"
      node = "settings.toml"
      operation = "copy"
      into = "config/tool.toml"
    }
` + extra + `
  }
}
`,
	}
	for name, body := range files {
		file := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
			t.Fatal(err)
		}
		mode := os.FileMode(0644)
		if strings.HasSuffix(name, ".sh") {
			mode = 0755
		}
		if err := os.WriteFile(file, []byte(body), mode); err != nil {
			t.Fatal(err)
		}
	}
	return root, filepath.Join(root, "mogent.hcl")
}

func TestTypedCLIUsageDiscoversNodesAndHeadingAddresses(t *testing.T) {
	_, config := writeArtifactCLI(t, "")
	var out, errOut bytes.Buffer
	if err := cli.Run([]string{"source", "list", "shared", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"shared:physical  (content, dir)", "review [dir] (skill)", "review/SKILL.md [file]"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("missing %q:\n%s", want, out.String())
		}
	}
	out.Reset()
	if err := cli.Run([]string{"source", "show", "shared:physical/review/SKILL.md", "--headings", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Kind:   file", "Bundle: review", `heading = ["Review", "Procedure"]`, `heading = ["Review", "Procedure", "Verify"]`} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("missing %q:\n%s", want, out.String())
		}
	}
	out.Reset()
	if err := cli.Run([]string{"help"}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "[--headings]") {
		t.Fatalf("help=%s", out.String())
	}
}

func TestTypedCLIPlanApplyAndDriftWorkflow(t *testing.T) {
	root, config := writeArtifactCLI(t, "")
	var out, errOut bytes.Buffer
	if err := cli.Run([]string{"plan", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"agent-export/: new (4 files)", "+# Procedure", "add AGENTS.md <- shared:physical/review/SKILL.md", "[render-markdown]", "[copy]", "(local " + filepath.Join(root, "library") + ")", "Plan: 1 to add"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("missing %q:\n%s", want, out.String())
		}
	}
	if _, err := os.Stat(filepath.Join(root, "agent-export")); !os.IsNotExist(err) {
		t.Fatalf("plan wrote target: %v", err)
	}
	out.Reset()
	if err := cli.Run([]string{"apply", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Applied 1 output(s): agent-export/") {
		t.Fatalf("apply=%s", out.String())
	}
	original, err := os.ReadFile(filepath.Join(root, "library/payloads/review/SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	copied, err := os.ReadFile(filepath.Join(root, "agent-export/.agents/skills/review/SKILL.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(original, copied) {
		t.Fatal("skill was transformed during copy")
	}
	info, err := os.Stat(filepath.Join(root, "agent-export/.agents/skills/review/scripts/check.sh"))
	if err != nil || info.Mode().Perm() != 0755 {
		t.Fatalf("script mode=%v %v", info, err)
	}
	out.Reset()
	if err := cli.Run([]string{"plan", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Plan: 0 to add, 0 to change, 1 unchanged") {
		t.Fatalf("second plan=%s", out.String())
	}
	if err := os.WriteFile(filepath.Join(root, "agent-export/AGENTS.md"), []byte("hand edit\n"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := cli.Run([]string{"apply", "--config", config}, &out, &errOut); err == nil || !strings.Contains(err.Error(), "direct edits") {
		t.Fatalf("drift=%v", err)
	}
	if err := cli.Run([]string{"apply", "--config", config, "--force"}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
}

func TestTypedCLICollisionThenExplicitReplacement(t *testing.T) {
	extra := `    source "shared" {
      from = "shared:physical"
      node = "overrides/review"
      operation = "copy"
      into = ".agents/skills/review"
    }
`
	root, config := writeArtifactCLI(t, extra)
	var out, errOut bytes.Buffer
	err := cli.Run([]string{"plan", "--config", config}, &out, &errOut)
	if err == nil || !strings.Contains(err.Error(), "collision") || strings.Contains(out.String(), "Plan:") {
		t.Fatalf("collision=%v output=%s", err, out.String())
	}
	data, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	updated := strings.Replace(string(data), `node = "overrides/review"`, `node = "overrides/review"`+"\n      replace = true", 1)
	if err := os.WriteFile(config, []byte(updated), 0644); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	if err := cli.Run([]string{"plan", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "replace-remove .agents/skills/review/scripts/check.sh") || !strings.Contains(out.String(), "previous: shared:physical/review") {
		t.Fatalf("replacement plan=%s", out.String())
	}
	if err := cli.Run([]string{"apply", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "agent-export/.agents/skills/review/scripts/check.sh")); !os.IsNotExist(err) {
		t.Fatalf("stale resource survived: %v", err)
	}
}

func artifactGit(t *testing.T, root string, args ...string) {
	t.Helper()
	full := append([]string{"-C", root, "-c", "user.name=test", "-c", "user.email=test@example.test", "-c", "init.defaultBranch=main"}, args...)
	if output, err := exec.Command("git", full...).CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
}

func TestTypedCLIUpdateShowsDirectoryAndRenderedChanges(t *testing.T) {
	root, config := writeArtifactCLI(t, "")
	upstream := filepath.Join(root, "library")
	artifactGit(t, upstream, "init", "--quiet")
	artifactGit(t, upstream, "add", ".")
	artifactGit(t, upstream, "commit", "--quiet", "-m", "Initial library")
	data, err := os.ReadFile(config)
	if err != nil {
		t.Fatal(err)
	}
	remote := strings.Replace(string(data), `local = "library"`, `git = "file://`+upstream+`"`, 1)
	if err := os.WriteFile(config, []byte(remote), 0644); err != nil {
		t.Fatal(err)
	}
	var out, errOut bytes.Buffer
	if err := cli.Run([]string{"apply", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(upstream, "payloads/review/SKILL.md"), []byte("---\nname: review\ndescription: Review a change.\n---\n# Review\n## Procedure\nUpdated procedure.\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(upstream, "payloads/review/new.txt"), []byte("new resource"), 0644); err != nil {
		t.Fatal(err)
	}
	artifactGit(t, upstream, "add", ".")
	artifactGit(t, upstream, "commit", "--quiet", "-m", "Revise library")
	out.Reset()
	if err := cli.Run([]string{"update", "shared", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"agent-export/: changed", "+ .agents/skills/review/new.txt", "+Updated procedure.", "Run `mogent update shared --accept`"} {
		if !strings.Contains(out.String(), want) {
			t.Fatalf("missing %q:\n%s", want, out.String())
		}
	}
	out.Reset()
	if err := cli.Run([]string{"update", "shared", "--config", config, "--accept"}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "Pinned shared to") {
		t.Fatalf("accepted update=%s", out.String())
	}
	guide, err := os.ReadFile(filepath.Join(root, "agent-export/AGENTS.md"))
	if err != nil || strings.Contains(string(guide), "Updated procedure") {
		t.Fatalf("update changed output=%q %v", guide, err)
	}
	if err := cli.Run([]string{"apply", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
}

func TestTypedCLIReportsModeAndDirectoryDrift(t *testing.T) {
	root, config := writeArtifactCLI(t, "")
	var out, errOut bytes.Buffer
	if err := cli.Run([]string{"apply", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	target := filepath.Join(root, "agent-export")
	if err := os.Chmod(filepath.Join(target, ".agents/skills/review/scripts/check.sh"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(target, "untracked-empty"), 0755); err != nil {
		t.Fatal(err)
	}
	out.Reset()
	errOut.Reset()
	if err := cli.Run([]string{"plan", "--config", config}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "~ .agents/skills/review/scripts/check.sh") || !strings.Contains(out.String(), "- untracked-empty/") || !strings.Contains(errOut.String(), "MOGENT208") {
		t.Fatalf("drift plan=%s\n%s", out.String(), errOut.String())
	}
	if err := cli.Run([]string{"apply", "--config", config}, &out, &errOut); err == nil {
		t.Fatal("metadata drift silently overwritten")
	}
	if err := cli.Run([]string{"apply", "--config", config, "--force"}, &out, &errOut); err != nil {
		t.Fatal(err)
	}
}
