package cli

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/registry"
)

func TestCatalogTreeFlatAndNodeScope(t *testing.T) {
	root := newCLIFixture(t)
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	var stdout, stderr bytes.Buffer
	for _, test := range []struct {
		args []string
		want []string
		omit string
	}{
		{[]string{"--tree", "--chars", "unicode"}, []string{"PATH", "TYPE", "├── agents", "│   └── guidance", "│       └── workflow", "md-section", "└── skills", "dir"}, "guidance/workflow"},
		{[]string{"--flat"}, []string{"guidance/workflow", "skills/review", "md-section"}, "├"},
		{[]string{"team:guidance", "--chars", "ascii"}, []string{"`-- guidance", "    `-- workflow"}, "skills"},
	} {
		stdout.Reset()
		args := append([]string{"source", "show", "--file", filepath.Join(root, "mogent.plan")}, test.args...)
		if err := Run(args, &stdout, &stderr); err != nil {
			t.Fatal(err)
		}
		for _, want := range test.want {
			if !strings.Contains(stdout.String(), want) {
				t.Errorf("missing %q:\n%s", want, stdout.String())
			}
		}
		if strings.Contains(stdout.String(), test.omit) {
			t.Errorf("unexpected %q:\n%s", test.omit, stdout.String())
		}
	}
}

func TestProjectSourceCommandsAndRegistryIndependence(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
	var stdout, stderr bytes.Buffer
	run := func(args ...string) string {
		t.Helper()
		stdout.Reset()
		if err := Run(args, &stdout, &stderr); err != nil {
			t.Fatal(err)
		}
		return stdout.String()
	}
	run("init", "--no-interactive")
	for _, args := range [][]string{{"source", "list"}, {"source", "show"}} {
		text := run(args...)
		for _, want := range []string{"No sources configured", "local library or Git", "source add"} {
			if !strings.Contains(text, want) {
				t.Fatalf("empty guidance missing %q: %s", want, text)
			}
		}
	}
	writeNewCLIFile(t, root, "team/mogent.lib", "library {format=1;name=\"Team\"}\nmd-section \"workflow\" {source=\"workflow.md\"}\n")
	writeNewCLIFile(t, root, "team/workflow.md", "# Workflow\nRead notes.\n")
	run("source", "register", "team", "team")
	if text := run("source", "ls", "--known"); !strings.Contains(text, "Team") {
		t.Fatalf("known=%s", text)
	}
	run("source", "add", "team", "--dry-run")
	if text := run("source", "list"); !strings.Contains(text, "No sources") {
		t.Fatalf("dry run modified project: %s", text)
	}
	run("source", "add", "team")
	if text := run("source", "ls"); !strings.Contains(text, "team") || !strings.Contains(text, "local") {
		t.Fatalf("project list=%s", text)
	}
	// Builds and source browsing must not read the discovery registry.
	registryPath, err := registry.DefaultPath()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(registryPath, []byte("broken: ["), 0600); err != nil {
		t.Fatal(err)
	}
	if text := run("source", "show"); !strings.Contains(text, "workflow") {
		t.Fatalf("catalog=%s", text)
	}
	if _, err := os.Stat("AGENTS.md"); !os.IsNotExist(err) {
		t.Fatalf("unexpected output installed: %v", err)
	}
	data, err := os.ReadFile("mogent.build")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile("mogent.build", append(data, []byte("\ninclude \"team:workflow\"\n")...), 0644); err != nil {
		t.Fatal(err)
	}
	run("build")
	run("build", "--check")
	if data, err := os.ReadFile("AGENTS.md"); err != nil || !strings.Contains(string(data), "Read notes.") {
		t.Fatalf("built=%s err=%v", data, err)
	}
	// Listing declarations must not open even unavailable sources.
	run("source", "add", "offline", "https://example.invalid/library.git")
	if text := run("source", "list"); !strings.Contains(text, "https://example.invalid/library.git") {
		t.Fatalf("offline declaration=%s", text)
	}
}

type scriptedPicker struct {
	choice string
	inputs []string
	err    error
	calls  int
}

func (p *scriptedPicker) Choose(_ string, choices []string) (string, error) {
	p.calls++
	if p.choice == "first" {
		return choices[0], p.err
	}
	return p.choice, p.err
}

func (p *scriptedPicker) Input(_ string) (string, error) {
	p.calls++
	if p.err != nil {
		return "", p.err
	}
	value := p.inputs[0]
	p.inputs = p.inputs[1:]
	return value, nil
}

func TestInitPickerKnownNewSkipAndCancellation(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
	registryPath, err := registry.DefaultPath()
	if err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Register(registryPath, "team", "https://example.invalid/team.git"); err != nil {
		t.Fatal(err)
	}
	for _, test := range []struct {
		name   string
		picker scriptedPicker
		want   string
	}{
		{"known", scriptedPicker{choice: "first"}, `git = "https://example.invalid/team.git"`},
		{"local", scriptedPicker{choice: "Local directory", inputs: []string{"personal", "../library"}}, `local = "../library"`},
		{"git", scriptedPicker{choice: "Git repository", inputs: []string{"remote", "https://example.invalid/remote.git"}}, `source "remote"`},
		{"skip", scriptedPicker{choice: "Skip for now"}, "Add a library"},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "mogent.build")
			result, err := initializeProject(assembly.InitOptions{Path: path}, &test.picker)
			if err != nil || !strings.Contains(string(result.Content), test.want) {
				t.Fatalf("result=%v err=%v", result, err)
			}
		})
	}
	path := filepath.Join(t.TempDir(), "mogent.build")
	picker := &scriptedPicker{err: errors.New("cancelled")}
	if _, err := initializeProject(assembly.InitOptions{Path: path}, picker); err == nil {
		t.Fatal("cancellation ignored")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("cancelled picker created entry: %v", err)
	}
	picker.calls = 0
	if _, err := initializeProject(assembly.InitOptions{Path: path, DryRun: true}, picker); err != nil || picker.calls != 0 {
		t.Fatalf("dry run prompted: %d %v", picker.calls, err)
	}
	if _, err := initializeProject(assembly.InitOptions{Path: path, Sources: map[string]string{"explicit": "library"}}, picker); err != nil || picker.calls != 0 {
		t.Fatalf("explicit source prompted: %d %v", picker.calls, err)
	}
	if _, err := initializeProject(assembly.InitOptions{Path: path}, picker); err == nil || picker.calls != 0 {
		t.Fatalf("existing target prompted: %d %v", picker.calls, err)
	}
}

func TestNativeModuleSelectionCommandsAndSourceAliasView(t *testing.T) {
	root := newCLIFixture(t)
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"source", "show", "team"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "team: Team") || !strings.Contains(stdout.String(), "workflow") {
		t.Fatalf("alias catalog=%s", stdout.String())
	}
	for _, command := range [][]string{
		{"add", "team:guidance/workflow", "--dry-run"},
		{"module", "add", "team:guidance/workflow", "--dry-run"},
	} {
		stdout.Reset()
		if err := Run(command, &stdout, &stderr); err != nil {
			t.Fatalf("%v: %v", command, err)
		}
		if !strings.Contains(stdout.String(), `include "team:guidance/workflow"`) {
			t.Fatalf("%v preview=%s", command, stdout.String())
		}
	}
	stdout.Reset()
	if err := Run([]string{"module", "add", "team:guidance/workflow"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Added team:guidance/workflow") {
		t.Fatalf("add=%s", stdout.String())
	}
	if err := Run([]string{"add", "team:guidance/workflow"}, &stdout, &stderr); err == nil || !strings.Contains(err.Error(), "already selected") {
		t.Fatalf("duplicate add=%v", err)
	}
	if err := Run([]string{"module", "add", "team:missing"}, &stdout, &stderr); err == nil || !strings.Contains(err.Error(), "unknown module") {
		t.Fatalf("unknown module add=%v", err)
	}
}

func TestModuleAddWritesTagQuery(t *testing.T) {
	root := newCLIFixture(t)
	plan := filepath.Join(root, "mogent.plan")
	var stdout, stderr bytes.Buffer
	args := []string{"module", "add", "team:skills", "--tag-all", "language", "--tag-none", "experimental", "--dry-run", "--file", plan}
	if err := Run(args, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{`include "team:skills" {`, `tags_all = ["language"]`, `tags_none = ["experimental"]`} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("tag preview missing %q:\n%s", want, stdout.String())
		}
	}
	if err := Run([]string{"module", "add", "team:skills", "--all", "--tag-all", "language", "--file", plan}, &stdout, &stderr); err == nil || !strings.Contains(err.Error(), "cannot be combined") {
		t.Fatalf("all/tags error=%v", err)
	}
}

func TestNativeShowTreeGroupsHeadingsAndSources(t *testing.T) {
	root := newCLIFixture(t)
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"show", "--tree"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"RESULT", "SOURCE", "AGENTS.md", "guidance", "Workflow Guidance", ".agents/skills", "team"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("tree missing %q:\n%s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "team:guidance") {
		t.Fatalf("tree should use the source column rather than namespace labels:\n%s", stdout.String())
	}
}

func TestNativeShowContentPreviewAndLineLimit(t *testing.T) {
	root := newCLIFixture(t)
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	for _, test := range []struct {
		args []string
		want []string
	}{
		{[]string{"show", "team:guidance/workflow", "--content"}, []string{"Type: md-section", "# Workflow", "Project: Demo"}},
		{[]string{"show", "--content", "--lines", "1", "team:guidance/workflow"}, []string{"# Workflow", "... (2 more lines)"}},
	} {
		stdout.Reset()
		if err := Run(test.args, &stdout, &stderr); err != nil {
			t.Fatal(err)
		}
		for _, want := range test.want {
			if !strings.Contains(stdout.String(), want) {
				t.Errorf("%v missing %q:\n%s", test.args, want, stdout.String())
			}
		}
	}
	for _, args := range [][]string{
		{"show", "--content"},
		{"show", "team:guidance/workflow", "--headings"},
		{"show", "team:guidance/workflow", "--content", "--lines", "-1"},
	} {
		if err := Run(args, &stdout, &stderr); err == nil {
			t.Errorf("invalid preview accepted: %v", args)
		}
	}
}

func TestNativeExceptAddsExclusionToSelectedAncestor(t *testing.T) {
	root := newCLIFixture(t)
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"add", "team:guidance"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if err := Run([]string{"except", "team:guidance/workflow", "--dry-run"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), `exclude = ["workflow"]`) {
		t.Fatalf("preview=%s", stdout.String())
	}
	stdout.Reset()
	if err := Run([]string{"exclude", "team:guidance/workflow"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Excepted team:guidance/workflow") {
		t.Fatalf("except=%s", stdout.String())
	}
	if err := Run([]string{"except", "team:guidance/workflow"}, &stdout, &stderr); err == nil || !strings.Contains(err.Error(), "already excepted") {
		t.Fatalf("duplicate except=%v", err)
	}
}

func TestNativeReorderPreservesSelectionOrder(t *testing.T) {
	root := newCLIFixture(t)
	plan := filepath.Join(root, "mogent.plan")
	contents, err := os.ReadFile(plan)
	if err != nil {
		t.Fatal(err)
	}
	contents = append(contents, []byte("source \"other\" {local=\"team\"}\ninclude \"team:guidance\"\n")...)
	if err := os.WriteFile(plan, contents, 0644); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"reorder", "--dry-run"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if strings.Index(stdout.String(), `source "other"`) > strings.Index(stdout.String(), `include "team:recommended"`) {
		t.Fatalf("reorder preview=%s", stdout.String())
	}
	if err := Run([]string{"reorder"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Reordered") {
		t.Fatalf("reorder=%s", stdout.String())
	}
}

func TestNativeDriftUsesStateAndCurrentComposition(t *testing.T) {
	root := newCLIFixture(t)
	plan := filepath.Join(root, "mogent.plan")
	file, err := assembly.LoadFile(plan)
	if err != nil {
		t.Fatal(err)
	}
	result, _, err := assembly.ResolvePlan(file, plan)
	if err != nil || result.HasErrors() {
		t.Fatalf("resolve=%v diagnostics=%v", err, result.Diagnostics)
	}
	manifest, err := result.Compose()
	if err != nil {
		t.Fatal(err)
	}
	if err := assembly.Apply(root, manifest); err != nil {
		t.Fatal(err)
	}
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"build", "--check"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Build check passed") {
		t.Fatalf("build check=%s", stdout.String())
	}
	if err := Run([]string{"drift", "--all", "--diff"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Output baseline: clean", "Current build: matches", "No AGENTS.md diff", ".agents/skills", "No .agents/skills diff"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("drift missing %q:\n%s", want, stdout.String())
		}
	}
	if err := os.WriteFile(filepath.Join(root, "AGENTS.md"), []byte("# direct edit\n"), 0644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if err := Run([]string{"drift", "--diff"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Output baseline: modified") || !strings.Contains(stdout.String(), "Current build: differs") || !strings.Contains(stdout.String(), "--- expected") {
		t.Fatalf("edited drift=%s", stdout.String())
	}
}

func TestNativeDriftReportsUpstreamModuleChanges(t *testing.T) {
	root := newCLIFixture(t)
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"build"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if err := Run([]string{"drift"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Source: team\nUpstream baseline: matches") {
		t.Fatalf("clean upstream status=%s", stdout.String())
	}
	writeNewCLIFile(t, root, "team/guidance/workflow.md", "# Workflow\n\nChanged upstream.\n")
	writeNewCLIFile(t, root, "team/added.md", "# Added\n")
	writeNewCLIFile(t, root, "team/mogent.lib", `library {format=1;name="Team"}
md-section "guidance" {md-section "workflow" {source="guidance/workflow.md";vars=["repo_name"]}}
md-section "added" {source="added.md"}
dir "skills" {source="skills";curate="choose"
  node "review" {role="skill";content="SKILL.md";default=true}
  node "testing" {default=true}
}
preset "recommended" {include "self:guidance";include "self:skills" {accept_defaults=true}}
`)
	stdout.Reset()
	if err := Run([]string{"drift"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"Source: team", "Upstream baseline: differs",
		"changed: team:guidance/workflow", "added: team:added", "removed: team:obsolete",
	} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("upstream drift missing %q:\n%s", want, stdout.String())
		}
	}
	if strings.Contains(stdout.String(), "changed: team:guidance\n") {
		t.Fatalf("upstream drift should descend to changed leaf:\n%s", stdout.String())
	}
}

func TestNativeStatusRemoveAndOptimise(t *testing.T) {
	root := newCLIFixture(t)
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"build"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if err := Run([]string{"status"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"Project: mogent.plan", "AGENTS.md", ".agents/skills", "Output baseline: clean", "Current build: matches"} {
		if !strings.Contains(stdout.String(), want) {
			t.Fatalf("status missing %q:\n%s", want, stdout.String())
		}
	}
	stdout.Reset()
	if err := Run([]string{"rm", "team:recommended", "--dry-run", "--diff"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Dry run: no files written") || !strings.Contains(stdout.String(), "Managed-output changes:") || !strings.Contains(stdout.String(), "Output: AGENTS.md") {
		t.Fatalf("remove diff=%s", stdout.String())
	}
	if err := Run([]string{"rm", "team:recommended"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"module", "rm", "team:recommended"}, &stdout, &stderr); err == nil || !strings.Contains(err.Error(), "not directly selected") {
		t.Fatalf("second remove=%v", err)
	}
	if err := Run([]string{"add", "team:guidance"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "team:guidance/workflow"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if err := Run([]string{"optimize"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Overlap: team:guidance already includes descendant team:guidance/workflow") || !strings.Contains(stdout.String(), "Dry run: no files written") || !strings.Contains(stdout.String(), "Proposed project change:") {
		t.Fatalf("optimize=%s", stdout.String())
	}
	contents, err := os.ReadFile(filepath.Join(root, "mogent.plan"))
	if err != nil || !strings.Contains(string(contents), `include "team:guidance/workflow"`) {
		t.Fatalf("optimize preview wrote project: %s err=%v", contents, err)
	}
	stdout.Reset()
	if err := Run([]string{"show", "--build"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "warning: team:guidance also directly includes team:guidance/workflow") {
		t.Fatalf("show=%s", stdout.String())
	}
	stdout.Reset()
	if err := Run([]string{"optimize", "--apply"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "Applied optimized selection:") {
		t.Fatalf("optimize apply=%s", stdout.String())
	}
	contents, err = os.ReadFile(filepath.Join(root, "mogent.plan"))
	if err != nil || strings.Contains(string(contents), `include "team:guidance/workflow"`) {
		t.Fatalf("optimize apply did not remove descendant: %s err=%v", contents, err)
	}
}
