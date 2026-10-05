package v2

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

const upstreamSidecar = `
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
`

func gitRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	full := append([]string{"-c", "user.name=test", "-c", "user.email=test@example.test", "-c", "init.defaultBranch=main", "-C", dir}, args...)
	output, err := exec.Command("git", full...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, output)
	}
	return strings.TrimSpace(string(output))
}

// newUpstream creates a committed library repository and returns its path.
func newUpstream(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git is not installed")
	}
	dir := t.TempDir()
	writeLibraryFile(t, dir, "agents/intro/role.md", "# Role\n\nBe helpful.\n")
	writeLibraryFile(t, dir, LibraryFile, upstreamSidecar)
	gitRun(t, dir, "init", "--quiet")
	gitRun(t, dir, "add", ".")
	gitRun(t, dir, "commit", "--quiet", "-m", "Add library")
	return dir
}

func gitConsumer(t *testing.T, upstream, commit string) *Config {
	t.Helper()
	root := t.TempDir()
	pin := ""
	if commit != "" {
		pin = "\n    commit = \"" + commit + "\""
	}
	writeLibraryFile(t, root, ConfigFile, `# consumer comment
mogent { format = 2 }

sources {
  source "shared" {
    git = "file://`+upstream+`"
    ref = "main"`+pin+`
  }
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

func TestUnpinnedGitSourceResolvesAndApplyPinsIt(t *testing.T) {
	upstream := newUpstream(t)
	head := gitRun(t, upstream, "rev-parse", "HEAD")
	config := gitConsumer(t, upstream, "")
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if codes := diagnosticCodes(result.Plan); len(codes) != 1 || codes[0] != "MOGENT206" {
		t.Fatalf("diagnostics = %v", codes)
	}
	if source := result.Sources["shared"]; source.Commit != head || source.Pinned {
		t.Fatalf("resolved = %#v, want commit %s", source, head)
	}
	if result.Pins["shared"] != head {
		t.Fatalf("pins = %v", result.Pins)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(config.Path), ".mogent", "sources", "shared", head, LibraryFile)); err != nil {
		t.Fatalf("checkout not cached: %v", err)
	}
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(config.Path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(contents), "commit = \""+head+"\"") || !strings.HasPrefix(string(contents), "# consumer comment\n") {
		t.Fatalf("config after apply:\n%s", contents)
	}
	reloaded, err := Load(config.Path)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Sources["shared"].Commit != head {
		t.Fatalf("reloaded commit = %q", reloaded.Sources["shared"].Commit)
	}
}

func TestPinnedGitSourceReadsCacheWithoutRemote(t *testing.T) {
	upstream := newUpstream(t)
	config := gitConsumer(t, upstream, "")
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	if err := os.RemoveAll(upstream); err != nil {
		t.Fatal(err)
	}
	pinned, err := Load(config.Path)
	if err != nil {
		t.Fatal(err)
	}
	result, err := PlanConfig(pinned)
	if err != nil {
		t.Fatalf("offline plan: %v", err)
	}
	if codes := diagnosticCodes(result.Plan); len(codes) != 0 {
		t.Fatalf("diagnostics = %v", codes)
	}
	if !result.Sources["shared"].Pinned || result.Changes[0].Status != FileUnchanged {
		t.Fatalf("offline result = %#v", result.Changes)
	}
}

func TestPinnedGitSourceFetchesMissingCacheAndRejectsUnknownCommit(t *testing.T) {
	upstream := newUpstream(t)
	head := gitRun(t, upstream, "rev-parse", "HEAD")
	config := gitConsumer(t, upstream, head)
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if codes := diagnosticCodes(result.Plan); len(codes) != 0 {
		t.Fatalf("diagnostics = %v", codes)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(config.Path), ".mogent", "sources", "shared", head)); err != nil {
		t.Fatalf("checkout not cached: %v", err)
	}
	missing := gitConsumer(t, upstream, strings.Repeat("0", 40))
	if _, err := PlanConfig(missing); err == nil || !strings.Contains(err.Error(), "fetch") {
		t.Fatalf("unknown commit error = %v", err)
	}
}

func TestUpdateShowsOutputDiffAndAcceptsNewCommit(t *testing.T) {
	upstream := newUpstream(t)
	old := gitRun(t, upstream, "rev-parse", "HEAD")
	config := gitConsumer(t, upstream, "")
	if _, err := Apply(config, false); err != nil {
		t.Fatal(err)
	}
	config, err := Load(config.Path)
	if err != nil {
		t.Fatal(err)
	}
	results, err := Update(config, "", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(results) != 1 || results[0].Changed() {
		t.Fatalf("up-to-date results = %#v", results)
	}
	writeLibraryFile(t, upstream, "agents/intro/role.md", "# Role\n\nBe very helpful.\n")
	gitRun(t, upstream, "commit", "--quiet", "-am", "Revise role")
	next := gitRun(t, upstream, "rev-parse", "HEAD")
	results, err = Update(config, "shared", false)
	if err != nil {
		t.Fatal(err)
	}
	update := results[0]
	if update.OldCommit != old || update.NewCommit != next || len(update.Changes) != 1 || update.Changes[0].Status != FileChanged {
		t.Fatalf("update = %#v", update)
	}
	if !strings.Contains(update.Changes[0].Diff, "+Be very helpful.") {
		t.Fatalf("diff = %q", update.Changes[0].Diff)
	}
	if reloaded, err := Load(config.Path); err != nil || reloaded.Sources["shared"].Commit != old {
		t.Fatalf("preview must not rewrite commit: %v %#v", err, reloaded.Sources["shared"])
	}
	if _, err := Update(config, "shared", true); err != nil {
		t.Fatal(err)
	}
	accepted, err := Load(config.Path)
	if err != nil {
		t.Fatal(err)
	}
	if accepted.Sources["shared"].Commit != next {
		t.Fatalf("accepted commit = %q, want %s", accepted.Sources["shared"].Commit, next)
	}
	plan, err := PlanConfig(accepted)
	if err != nil {
		t.Fatal(err)
	}
	if plan.Changes[0].Status != FileChanged {
		t.Fatalf("plan after accept = %#v", plan.Changes[0])
	}
}

func TestUpdateRefusesUnpinnedSource(t *testing.T) {
	upstream := newUpstream(t)
	config := gitConsumer(t, upstream, "")
	if _, err := Update(config, "", false); err == nil || !strings.Contains(err.Error(), "unpinned") {
		t.Fatalf("update error = %v", err)
	}
}

func TestLocalSourceInsideGitWorkTreeReportsHead(t *testing.T) {
	library := newUpstream(t)
	head := gitRun(t, library, "rev-parse", "--short", "HEAD")
	root := t.TempDir()
	writeLibraryFile(t, root, ConfigFile, `
mogent { format = 2 }
sources {
  source "shared" { local = "`+library+`" }
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
	result, err := PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if codes := diagnosticCodes(result.Plan); len(codes) != 1 || codes[0] != "MOGENT207" {
		t.Fatalf("diagnostics = %v", codes)
	}
	info := result.Plan.Diagnostics[0]
	if info.Severity != SeverityInfo || !strings.Contains(info.Message, head) || !strings.Contains(info.Message, "(clean)") {
		t.Fatalf("info = %#v", info)
	}
	writeLibraryFile(t, library, "agents/intro/role.md", "# Role\n\nEdited.\n")
	result, err = PlanConfig(config)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Plan.Diagnostics[0].Message, "(dirty)") {
		t.Fatalf("dirty info = %q", result.Plan.Diagnostics[0].Message)
	}
}

func TestLoadRejectsBadGitPins(t *testing.T) {
	cases := []struct{ name, source, want string }{
		{"commit-on-local", "source \"x\" {\n    local = \"library\"\n    commit = \"0123456789012345678901234567890123456789\"\n  }", "commit is valid only for git sources"},
		{"short-commit", "source \"x\" {\n    git = \"https://example.test/lib.git\"\n    commit = \"abc123\"\n  }", "full 40-character"},
		{"credentials", `source "x" { git = "https://user:pw@example.test/lib.git" }`, "embedded credentials"},
		{"bare-path", `source "x" { git = "/srv/lib.git" }`, "http(s) or file URL"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			content := "mogent { format = 2 }\nsources {\n  " + test.source + "\n}\noutputs {\n  output \"o\" {\n    path = \"AGENTS.md\"\n    kind = \"markdown\"\n    source \"x\" {\n      from = \"x:agents\"\n      select = {}\n    }\n  }\n}\n"
			_, err := Load(writeConfig(t, content))
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Load error = %v, want %q", err, test.want)
			}
		})
	}
}
