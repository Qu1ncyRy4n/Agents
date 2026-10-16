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
dir "skills" {source="skills";curate="choose"
  node "review" {role="skill";content="SKILL.md";default=true}
  node "testing" {default=true}
}
preset "recommended" {include "self:guidance";include "self:skills" {accept_defaults=true}}
`)
	writeNewCLIFile(t, root, "team/guidance/workflow.md", "# Workflow\n\nProject: ${repo_name}\n")
	writeNewCLIFile(t, root, "team/skills/review/SKILL.md", "---\nname: review\ndescription: Review\n---\n# Review\n")
	writeNewCLIFile(t, root, "team/skills/testing/SKILL.md", "---\nname: testing\ndescription: Test\n---\n# Testing\n")
	writeNewCLIFile(t, root, "mogent.plan", `plan {format=1}
let repo_name="Demo"
source "team" {local="team"}
include "team:recommended"
`)
	return root
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
