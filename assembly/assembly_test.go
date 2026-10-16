package assembly_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/assembly"
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
