package assembly_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

func write(t *testing.T, root, name, body string) {
	t.Helper()
	path := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func git(t *testing.T, directory string, args ...string) string {
	t.Helper()
	output, err := exec.Command("git", append([]string{"-C", directory}, args...)...).CombinedOutput()
	if err != nil {
		t.Fatalf("git %v: %s: %v", args, output, err)
	}
	return strings.TrimSpace(string(output))
}

func librarySource() string {
	return `library { format = 1; name = "Team" }
md-section "guidance" {
  md-section "workflow" { source = "guidance/workflow.md"; vars = ["repo_name"] }
}
dir "skills" {
  source = "skills"
  curate = "choose"
  node "review" { role = "skill"; content = "SKILL.md"; default = true; tags = ["lang/go"] }
  node "testing" { default = true; tags = ["lang/go"] }
  node "deployment" { default = false }
}
preset "recommended" {
  include "self:guidance"
  include "self:skills" { accept_defaults = true }
}
`
}

func fixtureLibrary(t *testing.T, root string) string {
	t.Helper()
	lib := filepath.Join(root, "team")
	write(t, lib, "mogent.lib", librarySource())
	write(t, lib, "guidance/workflow.md", "# Workflow\n\nProject: ${repo_name}\n\nShell: ${PATH}\n")
	write(t, lib, "skills/review/SKILL.md", "---\nname: review\ndescription: Review\n---\n# Review\n")
	write(t, lib, "skills/review/scripts/check.sh", "#!/bin/sh\n")
	write(t, lib, "skills/testing/SKILL.md", "---\nname: testing\ndescription: Test\n---\n# Testing\n")
	write(t, lib, "skills/deployment/SKILL.md", "---\nname: deployment\ndescription: Deploy\n---\n# Deployment\n")
	return lib
}

func parse(t *testing.T, filename, source string) *syntax.File {
	t.Helper()
	f, err := syntax.Parse(filename, []byte(source))
	if err != nil {
		t.Fatal(err)
	}
	return f
}

func TestLoadLibraryBuildsTypedCatalogWithoutAssetRegistry(t *testing.T) {
	root := t.TempDir()
	lib := fixtureLibrary(t, root)
	loaded, err := assembly.LoadLibrary(lib)
	if err != nil {
		t.Fatal(err)
	}
	if loaded.Name != "Team" || loaded.ByPath["guidance/workflow"].Type != assembly.Section || loaded.ByPath["skills"].Type != assembly.Directory {
		t.Fatalf("catalog=%#v", loaded)
	}
	review := loaded.ByPath["skills/review"]
	if review.Role != "skill" || review.Content != "SKILL.md" || review.Default == nil || !*review.Default {
		t.Fatalf("review=%#v", review)
	}
	inventory, err := loaded.ResolvePhysical(review)
	if err != nil {
		t.Fatal(err)
	}
	if inventory.Nodes["review/scripts/check.sh"] == nil || inventory.Nodes["review"].Role != "skill" {
		t.Fatalf("inventory=%#v", inventory.Nodes)
	}
	if _, found := loaded.Presets["recommended"]; !found {
		t.Fatal("preset missing")
	}
}

func TestLoadLibraryRejectsInvalidTypedDeclarations(t *testing.T) {
	cases := []struct{ name, content, want string }{
		{"root-without-source", `library {format=1;name="x"}; dir "skills" {}`, "requires source"},
		{"section-vars-physical", `library {format=1;name="x"}; file "x" {source="x";vars=["a"]}`, "cannot declare vars"},
		{"markdown-extension", `library {format=1;name="x"}; md-file "x" {source="x.txt"}`, "must end in .md"},
		{"section-heading", `library {format=1;name="x"}; md-section "x" {heading=["x"]}`, "heading needs a source"},
		{"duplicate-preset", `library {format=1;name="x"}; node "a" {type="md-section"}; preset "a" {}`, "duplicate library:name"},
		{"invalid-curate", `library {format=1;name="x"}; md-section "a" {curate="maybe"}`, "unknown curate"},
		{"preset-mutates", `library {format=1;name="x"}; preset "a" {source "x" {local="."}}`, "unexpected"},
		{"missing-relation", `library {format=1;name="x"}; md-section "a" {requires=["missing"]}`, "missing relationship target"},
		{"requirement-cycle", `library {format=1;name="x"}; md-section "a" {requires=["b"]};md-section "b" {requires=["a"]}`, "requirement cycle"},
		{"unknown-group", `library {format=1;name="x"}; md-section "a" {exclusive_group="missing"}`, "unknown exclusive group"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			write(t, root, "mogent.lib", tc.content)
			_, err := assembly.LoadLibrary(root)
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error=%v want=%q", err, tc.want)
			}
		})
	}
}

func TestLoadLibraryInheritsAncestorTags(t *testing.T) {
	root := t.TempDir()
	write(t, root, "mogent.lib", `library {format=1;name="x"};md-section "root" {tags=["topic"];md-section "child" {source="child.md";tags=["lang/go"]}}`)
	write(t, root, "child.md", "# Child\n")
	library, err := assembly.LoadLibrary(root)
	if err != nil || strings.Join(library.ByPath["root/child"].Tags, ",") != "topic,lang/go" {
		t.Fatalf("library=%#v err=%v", library, err)
	}
}

func TestResolvePlanExpandsPresetDefaultsAndRespectsScopedSelection(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	plan := parse(t, filepath.Join(root, "mogent.plan"), `plan {format=1}
let repo_name="Example"
source "team" {local="team"}
include "team:recommended"
`)
	result, sources, err := assembly.ResolvePlan(plan, filepath.Join(root, "mogent.plan"))
	if err != nil {
		t.Fatal(err)
	}
	if result.HasErrors() || len(result.Contributions) != 3 || len(sources) != 1 {
		t.Fatalf("result=%#v diagnostics=%#v", result.Contributions, result.Diagnostics)
	}
	if result.Contributions[0].Node.Path != "guidance/workflow" || result.Contributions[1].Node.Path != "skills/review" || result.Contributions[2].Node.Path != "skills/testing" {
		t.Fatalf("contributions=%#v", result.Contributions)
	}
	if result.Contributions[1].Into != ".agents/skills/review" || result.Contributions[2].Into != ".agents/skills/testing" {
		t.Fatalf("destinations=%#v", result.Contributions)
	}
	if _, found := result.Contributions[1].Values["repo_name"]; found {
		t.Fatal("physical skill received guidance values")
	}
}

func TestResolvePlanReportsChoicesAndUnknownNodesWithoutDiscardingIndependentNodes(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	plan := parse(t, filepath.Join(root, "mogent.plan"), `plan {format=1}
source "team" {local="team"}
include "team:skills"
include "team:guidance/workflow" {values {repo_name="Example"}}
include "team:skills" ["missing"]
`)
	result, _, err := assembly.ResolvePlan(plan, filepath.Join(root, "mogent.plan"))
	if err != nil {
		t.Fatal(err)
	}
	if !result.HasErrors() || len(result.Contributions) != 1 || result.Contributions[0].Node.Path != "guidance/workflow" {
		t.Fatalf("result=%#v diagnostics=%#v", result.Contributions, result.Diagnostics)
	}
	var undecided, missing bool
	for _, diagnostic := range result.Diagnostics {
		undecided = undecided || strings.Contains(diagnostic.Message, "UNDECIDED")
		missing = missing || strings.Contains(diagnostic.Message, "has no child")
	}
	if !undecided || !missing {
		t.Fatalf("diagnostics=%#v", result.Diagnostics)
	}
	foundUndecided, foundSelected := false, false
	for _, node := range result.Nodes {
		foundUndecided = foundUndecided || node.Status == assembly.Undecided && node.Path == "skills"
		foundSelected = foundSelected || node.Status == assembly.Selected && node.Path == "guidance/workflow"
	}
	if !foundUndecided || !foundSelected {
		t.Fatalf("result nodes=%#v", result.Nodes)
	}
}

func TestResolvePlanFiltersTagsAndScopedExclusions(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	plan := parse(t, filepath.Join(root, "mogent.plan"), `plan {format=1}
source "team" {local="team"}
include "team:skills" {tags_all=["lang"];exclude=["review"]}
`)
	result, _, err := assembly.ResolvePlan(plan, filepath.Join(root, "mogent.plan"))
	if err != nil || result.HasErrors() || len(result.Contributions) != 1 || result.Contributions[0].Node.Path != "skills/testing" {
		t.Fatalf("result=%#v diagnostics=%#v err=%v", result.Contributions, result.Diagnostics, err)
	}
}

func TestMatchesTagsMatchesUnqualifiedTagComponents(t *testing.T) {
	node := &assembly.Node{Tags: []string{"scope/core", "topic/workflow"}}
	if !assembly.MatchesTags(node, nil, []string{"workflow"}, nil) {
		t.Fatal("unqualified component should match topic/workflow")
	}
	if assembly.MatchesTags(node, nil, []string{"work"}, nil) {
		t.Fatal("partial component must not become a fuzzy match")
	}
}

func TestResolvePlanFiltersNegativeTags(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	plan := parse(t, filepath.Join(root, "mogent.plan"), `plan {format=1}
source "team" {local="team"}
include "team:skills" {tags_none=["lang/go"]}
`)
	result, _, err := assembly.ResolvePlan(plan, filepath.Join(root, "mogent.plan"))
	if err != nil || result.HasErrors() || len(result.Contributions) != 1 || result.Contributions[0].Node.Path != "skills/deployment" {
		t.Fatalf("result=%#v diagnostics=%#v err=%v", result.Contributions, result.Diagnostics, err)
	}
}

func TestResolvePlanUnionsTagsAndExplicitChildren(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	plan := parse(t, filepath.Join(root, "mogent.plan"), `plan {format=1}
source "team" {local="team"}
include "team:skills" ["deployment"] {tags_all=["lang"]}
`)
	result, _, err := assembly.ResolvePlan(plan, filepath.Join(root, "mogent.plan"))
	if err != nil || result.HasErrors() || len(result.Contributions) != 3 {
		t.Fatalf("result=%#v diagnostics=%#v err=%v", result.Contributions, result.Diagnostics, err)
	}
	got := []string{result.Contributions[0].Node.Path, result.Contributions[1].Node.Path, result.Contributions[2].Node.Path}
	if strings.Join(got, ",") != "skills/review,skills/testing,skills/deployment" {
		t.Fatalf("selected=%v", got)
	}
}

func TestResolvePlanRequiresReasonForExperimentalAndFoundationExclusion(t *testing.T) {
	root := t.TempDir()
	write(t, root, "team/mogent.lib", `library {format=1;name="Team"}
md-section "guide" {source="guide.md";reason_required=true}
md-section "base" {curate="foundation"}
`)
	write(t, root, "team/guide.md", "# guide\n")
	for name, source := range map[string]string{
		"missing-include-reason":   `plan {format=1};source "team" {local="team"};include "team:guide"`,
		"missing-exclusion-reason": `plan {format=1};source "team" {local="team"};include "team:base" {exclude=["base"]}`,
	} {
		t.Run(name, func(t *testing.T) {
			result, _, err := assembly.ResolvePlan(parse(t, filepath.Join(root, "mogent.plan"), source), filepath.Join(root, "mogent.plan"))
			if err != nil || !result.HasErrors() {
				t.Fatalf("result=%#v err=%v", result, err)
			}
		})
	}
	result, _, err := assembly.ResolvePlan(parse(t, filepath.Join(root, "mogent.plan"), `plan {format=1};source "team" {local="team"};include "team:guide" {reason="needed"}`), filepath.Join(root, "mogent.plan"))
	if err != nil || result.HasErrors() || len(result.Contributions) != 1 || result.Contributions[0].Decision[0].Reason != "needed" {
		t.Fatalf("result=%#v err=%v", result, err)
	}
}

func TestResolvePlanAppliesBroadCuration(t *testing.T) {
	root := t.TempDir()
	write(t, root, "team/mogent.lib", `library {format=1;name="Team"}
md-section "foundation" {curate="foundation";md-section "a" {source="a.md"};md-section "b" {source="b.md"}}
md-section "optional" {curate="optional";md-section "yes" {source="yes.md";default=true};md-section "no" {source="no.md";default=false}}
md-section "optin" {curate="opt_in";md-section "only" {source="only.md"}}
`)
	for _, name := range []string{"a", "b", "yes", "no", "only"} {
		write(t, root, "team/"+name+".md", "# "+name+"\n")
	}
	plan := parse(t, filepath.Join(root, "mogent.plan"), `plan {format=1};source "team" {local="team"};include "team:foundation";include "team:optional";include "team:optin"`)
	result, _, err := assembly.ResolvePlan(plan, filepath.Join(root, "mogent.plan"))
	if err != nil || result.HasErrors() || len(result.Contributions) != 3 {
		t.Fatalf("result=%#v diagnostics=%#v err=%v", result.Contributions, result.Diagnostics, err)
	}
	got := []string{result.Contributions[0].Node.Path, result.Contributions[1].Node.Path, result.Contributions[2].Node.Path}
	if strings.Join(got, ",") != "foundation/a,foundation/b,optional/yes" {
		t.Fatalf("selected=%v", got)
	}
	foundOptIn := false
	for _, node := range result.Nodes {
		foundOptIn = foundOptIn || node.Path == "optin" && node.Status == assembly.Undecided
	}
	if len(result.Diagnostics) != 1 || result.Diagnostics[0].Severity != assembly.Warning || !strings.Contains(result.Diagnostics[0].Message, "OPT_IN") || !foundOptIn {
		t.Fatalf("missing opt-in warning: diagnostics=%#v nodes=%#v", result.Diagnostics, result.Nodes)
	}
}

func TestResolvePlanValidatesFinalRelationships(t *testing.T) {
	root := t.TempDir()
	write(t, root, "team/mogent.lib", `library {format=1;name="Team"}
group "pick" {mode="at-most-one"}
md-section "a" {source="a.md";requires=["required"];exclusive_group="pick"}
md-section "b" {source="b.md";conflicts_with=["a"];exclusive_group="pick"}
md-section "required" {source="required.md"}
`)
	write(t, root, "team/a.md", "# A\n")
	write(t, root, "team/b.md", "# B\n")
	write(t, root, "team/required.md", "# R\n")
	plan := parse(t, filepath.Join(root, "mogent.plan"), `plan {format=1};source "team" {local="team"};include "team:a";include "team:b"`)
	result, _, err := assembly.ResolvePlan(plan, filepath.Join(root, "mogent.plan"))
	if err != nil || !result.HasErrors() {
		t.Fatalf("result=%#v err=%v", result, err)
	}
	var require, conflict, exclusive bool
	for _, diagnostic := range result.Diagnostics {
		require = require || strings.Contains(diagnostic.Message, "requires")
		conflict = conflict || strings.Contains(diagnostic.Message, "conflicts")
		exclusive = exclusive || strings.Contains(diagnostic.Message, "exclusive")
	}
	if !require || !conflict || !exclusive {
		t.Fatalf("diagnostics=%#v", result.Diagnostics)
	}
}

func TestResolvePlanExpandsLibraryRecipeAndRejectsCycles(t *testing.T) {
	root := t.TempDir()
	write(t, root, "team/mogent.lib", `library {format=1;name="Team"}
md-section "guide" {source="guide.md"}
recipe "recommended" {include "self:guide"}
recipe "cycle" {include "self:cycle"}
`)
	write(t, root, "team/guide.md", "# Guide\n")
	for name, planText := range map[string]string{
		"expand": `plan {format=1};source "team" {local="team"};include "team:recommended"`,
		"cycle":  `plan {format=1};source "team" {local="team"};include "team:cycle"`,
	} {
		t.Run(name, func(t *testing.T) {
			result, _, err := assembly.ResolvePlan(parse(t, filepath.Join(root, "mogent.plan"), planText), filepath.Join(root, "mogent.plan"))
			if err != nil {
				t.Fatal(err)
			}
			if name == "expand" && (result.HasErrors() || len(result.Contributions) != 1) {
				t.Fatalf("result=%#v", result)
			}
			if name == "cycle" && (!result.HasErrors() || !strings.Contains(result.Diagnostics[0].Message, "cycle")) {
				t.Fatalf("diagnostics=%#v", result.Diagnostics)
			}
		})
	}
}

func TestResolvePlanGuidanceValuesAreWhitelisted(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	for name, source := range map[string]string{
		"binding": `plan {format=1}
let repo_name="Example"
source "team" {local="team"}
include "team:guidance/workflow"`,
		"explicit": `plan {format=1}
let repo_name="Wrong"
source "team" {local="team"}
include "team:guidance/workflow" {values {repo_name="Example"}}`,
	} {
		t.Run(name, func(t *testing.T) {
			result, _, err := assembly.ResolvePlan(parse(t, filepath.Join(root, "mogent.plan"), source), filepath.Join(root, "mogent.plan"))
			if err != nil || result.HasErrors() || result.Contributions[0].Values["repo_name"] != "Example" {
				t.Fatalf("result=%#v err=%v", result, err)
			}
		})
	}
	for name, source := range map[string]string{
		"missing": `plan {format=1}
source "team" {local="team"}
include "team:guidance/workflow"`,
		"undeclared": `plan {format=1}
let nope="x"
source "team" {local="team"}
include "team:guidance/workflow" {values {nope="x"}}`,
		"physical": `plan {format=1}
source "team" {local="team"}
include "team:skills/review" {values {repo_name="x"}}`,
	} {
		t.Run(name, func(t *testing.T) {
			result, _, err := assembly.ResolvePlan(parse(t, filepath.Join(root, "mogent.plan"), source), filepath.Join(root, "mogent.plan"))
			if err != nil || !result.HasErrors() {
				t.Fatalf("result=%#v err=%v", result, err)
			}
		})
	}
}

func TestComposeProducesRootGuideAndCompleteBundles(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	plan := parse(t, filepath.Join(root, "mogent.plan"), `plan {format=1}
let repo_name="Example"
source "team" {local="team"}
include {
  into="AGENTS.md"
  text="""
    # Example
    """
}
include "team:recommended"
`)
	result, _, err := assembly.ResolvePlan(plan, filepath.Join(root, "mogent.plan"))
	if err != nil || result.HasErrors() {
		t.Fatalf("result=%#v err=%v", result, err)
	}
	manifest, err := result.Compose()
	if err != nil {
		t.Fatal(err)
	}
	agents := string(manifest.Files["AGENTS.md"].Bytes)
	if !strings.Contains(agents, "# Example") || !strings.Contains(agents, "Project: Example") || !strings.Contains(agents, "Shell: ${PATH}") {
		t.Fatalf("AGENTS.md=%q", agents)
	}
	if manifest.Files[".agents/skills/review/scripts/check.sh"].Bytes == nil || manifest.Files[".agents/skills/testing/SKILL.md"].Bytes == nil {
		t.Fatalf("files=%#v", manifest.Files)
	}
	if _, found := manifest.Files[".agents/skills/deployment/SKILL.md"]; found {
		t.Fatal("opted-out deployment included")
	}
}

func TestSequenceProducesVerifiableBuildThatRecomposes(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	planPath := filepath.Join(root, "mogent.plan")
	plan := parse(t, planPath, `plan {format=1}
let repo_name="Example"
source "team" {local="team"}
include "team:recommended"
`)
	first, result, err := assembly.Sequence(plan, planPath)
	if err != nil || result.HasErrors() {
		t.Fatalf("sequence err=%v result=%#v", err, result)
	}
	second, _, err := assembly.Sequence(plan, planPath)
	if err != nil || string(first) != string(second) {
		t.Fatalf("sequence not deterministic: %v\n%s\n%s", err, first, second)
	}
	if err := assembly.VerifyBuildChecksum(first); err != nil {
		t.Fatal(err)
	}
	if err := assembly.VerifyBuildChecksum([]byte(strings.Replace(string(first), `into = "AGENTS.md"`, `into = "AGENTS2.md"`, 1))); err == nil {
		t.Fatal("edited generated build accepted")
	}
	build := parse(t, filepath.Join(root, "mogent.build"), string(first))
	resolved, _, err := assembly.ResolvePlan(build, filepath.Join(root, "mogent.build"))
	if err != nil || resolved.HasErrors() {
		t.Fatalf("build resolve err=%v diagnostics=%#v", err, resolved.Diagnostics)
	}
	manifest, err := resolved.Compose()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(manifest.Files["AGENTS.md"].Bytes), "Project: Example") || manifest.Files[".agents/skills/review/SKILL.md"].Bytes == nil {
		t.Fatalf("manifest=%#v", manifest.Files)
	}
}

func TestApplyOwnsGuideAndSkillsButPreservesOtherAgentsFiles(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	planPath := filepath.Join(root, "mogent.plan")
	plan := parse(t, planPath, `plan {format=1}
let repo_name="Example"
source "team" {local="team"}
include "team:recommended"
`)
	result, _, err := assembly.ResolvePlan(plan, planPath)
	if err != nil || result.HasErrors() {
		t.Fatalf("resolve=%v diagnostics=%#v", err, result.Diagnostics)
	}
	manifest, err := result.Compose()
	if err != nil {
		t.Fatal(err)
	}
	write(t, root, ".agents/settings.json", "{\"keep\":true}\n")
	write(t, root, ".agents/skills/obsolete/SKILL.md", "obsolete\n")
	if err := assembly.Apply(root, manifest); err != nil {
		t.Fatal(err)
	}
	if stateBytes, err := os.ReadFile(filepath.Join(root, ".mogent", "state.json")); err != nil || !strings.Contains(string(stateBytes), "AGENTS.md") {
		t.Fatalf("authored state=%s err=%v", stateBytes, err)
	}
	if bytes, err := os.ReadFile(filepath.Join(root, "AGENTS.md")); err != nil || !strings.Contains(string(bytes), "Project: Example") {
		t.Fatalf("agents=%q err=%v", bytes, err)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents/settings.json")); err != nil {
		t.Fatalf("unrelated .agents file removed: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents/skills/obsolete")); !os.IsNotExist(err) {
		t.Fatalf("stale managed skill remains: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, ".agents/skills/review/scripts/check.sh")); err != nil {
		t.Fatal(err)
	}
	if err := assembly.CheckWorkspace(root, manifest); err != nil {
		t.Fatal(err)
	}
	write(t, root, ".agents/skills/review/extra.txt", "drift\n")
	if err := assembly.CheckWorkspace(root, manifest); err == nil || !strings.Contains(err.Error(), "unexpected") {
		t.Fatalf("check error=%v", err)
	}
	if err := os.Remove(filepath.Join(root, ".agents/skills/review/extra.txt")); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(filepath.Join(root, "AGENTS.md"), 0600); err != nil {
		t.Fatal(err)
	}
	if err := assembly.CheckWorkspace(root, manifest); err == nil || !strings.Contains(err.Error(), "mode") {
		t.Fatalf("mode check error=%v", err)
	}
}

func TestApplyGeneratedSwapsBuildWithManagedOutputs(t *testing.T) {
	root := t.TempDir()
	manifest := &content.Manifest{Files: map[string]content.Payload{"AGENTS.md": {Bytes: []byte("new\n"), Mode: 0644}}, Directories: map[string]os.FileMode{}, Bundles: map[string]string{}}
	write(t, root, "AGENTS.md", "old\n")
	write(t, root, "mogent.build", "old build\n")
	if err := assembly.ApplyGenerated(root, manifest, []byte("new build\n")); err != nil {
		t.Fatal(err)
	}
	for name, want := range map[string]string{"AGENTS.md": "new\n", "mogent.build": "new build\n"} {
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil || string(got) != want {
			t.Fatalf("%s=%q err=%v", name, got, err)
		}
	}
}

func TestRecordedBuildRejectsChangedLocalNode(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	planPath := filepath.Join(root, "mogent.plan")
	plan := parse(t, planPath, `plan {format=1}
let repo_name="Example"
source "team" {local="team"}
include "team:recommended"
`)
	buildBytes, _, err := assembly.Sequence(plan, planPath)
	if err != nil {
		t.Fatal(err)
	}
	write(t, root, "team/skills/review/scripts/check.sh", "changed\n")
	build := parse(t, filepath.Join(root, "mogent.build"), string(buildBytes))
	result, _, err := assembly.ResolvePlan(build, filepath.Join(root, "mogent.build"))
	if err != nil || !result.HasErrors() {
		t.Fatalf("result=%#v err=%v", result, err)
	}
	found := false
	for _, diagnostic := range result.Diagnostics {
		found = found || strings.Contains(diagnostic.Message, "changed")
	}
	if !found {
		t.Fatalf("diagnostics=%#v", result.Diagnostics)
	}
}

func TestRefreshLocalInputHashesEditsOnlyAuthoredHashes(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	planPath := filepath.Join(root, "mogent.build")
	build := `build {format=1;origin="authored"}
source "team" {local="team"}
include "team:guidance/workflow" {values {repo_name="Example"}}
resolved {input "team:guidance/workflow" {sha256="old"}}
// retain
`
	file := parse(t, planPath, build)
	updated, changes, err := assembly.RefreshLocalInputHashes(file, planPath)
	if err != nil || len(changes) != 1 || changes[0].Old != "old" || strings.Contains(string(updated), `sha256="old"`) || !strings.Contains(string(updated), "// retain") {
		t.Fatalf("updated=%q changes=%#v err=%v", updated, changes, err)
	}
	generated := parse(t, planPath, `build {format=1;origin="generated"}`)
	if _, _, err := assembly.RefreshLocalInputHashes(generated, planPath); err == nil || !strings.Contains(err.Error(), "sequence") {
		t.Fatalf("generated error=%v", err)
	}
}

func TestSequencePinsGitSource(t *testing.T) {
	root := t.TempDir()
	lib := fixtureLibrary(t, root)
	git(t, lib, "init", "-q")
	git(t, lib, "config", "user.email", "test@example.com")
	git(t, lib, "config", "user.name", "Test")
	git(t, lib, "add", ".")
	git(t, lib, "commit", "-qm", "library")
	commit := git(t, lib, "rev-parse", "HEAD")
	planPath := filepath.Join(root, "mogent.plan")
	plan := parse(t, planPath, `plan {format=1}
let repo_name="Example"
source "team" {git="file://`+lib+`";ref="HEAD"}
include "team:recommended"
`)
	build, _, err := assembly.Sequence(plan, planPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(build), `commit = "`+commit+`"`) {
		t.Fatalf("build missing pin:\n%s", build)
	}
	parsed := parse(t, filepath.Join(root, "mogent.build"), string(build))
	result, _, err := assembly.ResolvePlan(parsed, filepath.Join(root, "mogent.build"))
	if err != nil || result.HasErrors() {
		t.Fatalf("result=%#v err=%v", result, err)
	}
}

func TestLocalizeCopiesBundleAndAddsLocalReplacement(t *testing.T) {
	root := t.TempDir()
	fixtureLibrary(t, root)
	planPath := filepath.Join(root, "mogent.plan")
	plan := parse(t, planPath, `plan {format=1}
source "team" {local="team"}
include "team:skills" {accept_defaults=true}
`)
	result, err := assembly.Localize(plan, planPath, "team:skills/review", true)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(result.LocalPath); !os.IsNotExist(err) {
		t.Fatalf("dry-run created source: %v", err)
	}
	if !strings.Contains(string(result.Plan), `source "local"`) || !strings.Contains(string(result.Plan), `include "local:skills/review"`) {
		t.Fatalf("plan=%s", result.Plan)
	}
	result, err = assembly.Localize(plan, planPath, "team:skills/review", false)
	if err != nil {
		t.Fatal(err)
	}
	if bytes, err := os.ReadFile(filepath.Join(result.LocalPath, "scripts/check.sh")); err != nil || string(bytes) != "#!/bin/sh\n" {
		t.Fatalf("local bundle=%q err=%v", bytes, err)
	}
	updated, err := assembly.LoadFile(planPath)
	if err != nil {
		t.Fatal(err)
	}
	resolved, _, err := assembly.ResolvePlan(updated, planPath)
	if err != nil || resolved.HasErrors() {
		t.Fatalf("resolved=%#v err=%v", resolved, err)
	}
	if len(resolved.Contributions) != 3 || !resolved.Contributions[2].Replace || resolved.Contributions[2].Source != "local" {
		t.Fatalf("contributions=%#v", resolved.Contributions)
	}
	if _, err := assembly.Localize(updated, planPath, "team:skills/review", false); err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("second localize error=%v", err)
	}
}
