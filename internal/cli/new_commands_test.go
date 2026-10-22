package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeNewCLIFile(t *testing.T, root, name, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func newCLIFixture(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	writeNewCLIFile(t, root, "team/mogent.lib", `library {format=1;name="Team"}
md-section "guidance" {md-section "workflow" {source="guidance/workflow.md";vars=["repo_name"]}}
md-section "obsolete" {source="obsolete.md"}
dir "skills" {source="skills";curate="choose"
  node "review" {role="skill";content="SKILL.md";default=true;tags=["quality/review"]}
  node "testing" {default=true;tags=["language/go"]}
}
preset "recommended" {include "self:guidance";include "self:skills" {accept_defaults=true}}
`)
	writeNewCLIFile(t, root, "team/guidance/workflow.md", "# Workflow\n\nProject: ${repo_name}\n")
	writeNewCLIFile(t, root, "team/obsolete.md", "# Obsolete\n")
	writeNewCLIFile(t, root, "team/skills/review/SKILL.md", "---\nname: review\ndescription: Review\n---\n# Review\n")
	writeNewCLIFile(t, root, "team/skills/testing/SKILL.md", "---\nname: testing\ndescription: Test\n---\n# Testing\n")
	writeNewCLIFile(t, root, "mogent.plan", `plan {format=1}
let repo_name="Demo"
source "team" {local="team"}
include "team:recommended"
`)
	return root
}

func TestSourceSearchListsTagMatchesAndAddCommands(t *testing.T) {
	root := newCLIFixture(t)
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"source", "search", "team:skills", "--tag-all", "language", "--add-commands", "--file", filepath.Join(root, "mogent.plan")}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"REF", "team:skills/testing", "language/go", "Add commands:", "mogent module add team:skills/testing"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("search missing %q:\n%s", want, stdout.String())
		}
	}
	stdout.Reset()
	if err := Run([]string{"source", "search", "team:skills", "--tag-none", "language", "--file", filepath.Join(root, "mogent.plan")}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "team:skills/review") || strings.Contains(stdout.String(), "team:skills/testing") {
		t.Fatalf("negative search=%s", stdout.String())
	}
}

func TestCanonicalCatalogAndResultCommands(t *testing.T) {
	root := newCLIFixture(t)
	plan := filepath.Join(root, "mogent.plan")
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"source", "tree", "team:skills", "--file", plan}, &stdout, &stderr); err != nil || !strings.Contains(stdout.String(), "review") {
		t.Fatalf("source tree stdout=%q err=%v", stdout.String(), err)
	}
	stdout.Reset()
	if err := Run([]string{"source", "content", "team:guidance/workflow", "--file", plan}, &stdout, &stderr); err != nil || !strings.Contains(stdout.String(), "Project: Demo") {
		t.Fatalf("source content stdout=%q err=%v", stdout.String(), err)
	}
	stdout.Reset()
	if err := Run([]string{"show", "tree", "--file", plan}, &stdout, &stderr); err != nil || !strings.Contains(stdout.String(), "RESULT") {
		t.Fatalf("show tree stdout=%q err=%v", stdout.String(), err)
	}
}

func TestNewSequenceBuildShowWorkflow(t *testing.T) {
	root := newCLIFixture(t)
	plan := filepath.Join(root, "mogent.plan")
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"sequence", "--file", plan}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Sequenced") {
		t.Fatalf("stdout=%q", stdout.String())
	}
	build := filepath.Join(root, "mogent.build")
	contents, err := os.ReadFile(build)
	if err != nil || !strings.Contains(string(contents), `origin = "generated"`) {
		t.Fatalf("build=%q err=%v", contents, err)
	}
	stdout.Reset()
	if err := Run([]string{"show", "--file", build, "--build"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "team:guidance/workflow -> AGENTS.md") || !strings.Contains(stdout.String(), ".agents/skills/review") {
		t.Fatalf("show=%q", stdout.String())
	}
	stdout.Reset()
	if err := Run([]string{"build", "--file", build}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if content, err := os.ReadFile(filepath.Join(root, "AGENTS.md")); err != nil || !strings.Contains(string(content), "Project: Demo") {
		t.Fatalf("AGENTS=%q err=%v", content, err)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents/skills/review/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if err := Run([]string{"build", "--file", build, "--check"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Build check passed") {
		t.Fatalf("check=%q", stdout.String())
	}
}

func TestNewBuildProtectsEditedGeneratedFileAndShowsIssues(t *testing.T) {
	root := newCLIFixture(t)
	plan := filepath.Join(root, "mogent.plan")
	build := filepath.Join(root, "mogent.build")
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"sequence", "--file", plan}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(build)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(build, append(contents, []byte("// local edit\n")...), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"sequence", "--file", plan}, &stdout, &stderr); err == nil || !strings.Contains(err.Error(), "edited") {
		t.Fatalf("sequence error=%v", err)
	}
	writeNewCLIFile(t, root, "bad.plan", `plan {format=1}
source "team" {local="team"}
include "team:skills"
`)
	stdout.Reset()
	if err := Run([]string{"show", "--file", filepath.Join(root, "bad.plan"), "--issues"}, &stdout, &stderr); err == nil || !strings.Contains(stdout.String(), "UNDECIDED") {
		t.Fatalf("show error=%v stdout=%q", err, stdout.String())
	}
}

func TestModuleLocalize(t *testing.T) {
	root := newCLIFixture(t)
	plan := filepath.Join(root, "mogent.plan")
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"module", "localize", "team:skills/review", "--file", plan, "--dry-run"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "local:skills/review") {
		t.Fatalf("dry run=%q", stdout.String())
	}
	if _, err := os.Stat(filepath.Join(root, ".mogent/local/skills/review")); !os.IsNotExist(err) {
		t.Fatalf("dry-run wrote local source: %v", err)
	}
	stdout.Reset()
	if err := Run([]string{"module", "localize", "team:skills/review", "--file", plan}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, ".mogent/local/skills/review/SKILL.md")); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(plan)
	if err != nil || !strings.Contains(string(contents), `include "local:skills/review"`) {
		t.Fatalf("plan=%q err=%v", contents, err)
	}
}

func TestBuildDiscoversPlanAndSequenceProtectsAuthoredBuild(t *testing.T) {
	root := newCLIFixture(t)
	plan := filepath.Join(root, "mogent.plan")
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(root); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"build"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(root, "AGENTS.md")); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "mogent.build"), []byte(`build {format=1;origin="authored"}`), 0644); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"sequence", "--file", plan}, &stdout, &stderr); err == nil || !strings.Contains(err.Error(), "authored") {
		t.Fatalf("sequence error=%v", err)
	}
}
