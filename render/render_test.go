package render_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/manifest"
	"github.com/Qu1ncyRy4n/Agents/render"
	"github.com/Qu1ncyRy4n/Agents/sourcecache"
)

func TestBuildUsesManifestHeadingsTemplatesAndExclusions(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "rules.md"), "# Rules\nHello {{ .repo_name }}.\n\n## Keep\nKeep this.\n\n## Drop\nDrop this.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\nvars:\n  repo_name: mogent\noutput: AGENTS.md\ndoc:\n  - heading: Local Rules\n    from:\n      - shared:rules\n    exclude:\n      - shared:rules/drop\n")

	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	want := "# Local Rules\n\nHello mogent.\n\n## Keep\n\nKeep this.\n"
	if result.Content != want {
		t.Fatalf("rendered content:\n%s\nwant:\n%s", result.Content, want)
	}
}

func TestBuildStripsHTMLCommentsExceptFencedContentUnlessPreserved(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "rules.md"), "# Rules\nVisible. <!-- inline -->\n<!-- multi\nline -->\n```html\n<!-- retain -->\n```\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\ndoc:\n  - Rules: shared:rules\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(result.Content, "inline") || strings.Contains(result.Content, "multi") || !strings.Contains(result.Content, "<!-- retain -->") {
		t.Fatalf("default comment rendering = %q", result.Content)
	}
	preserved, err := render.BuildWithOptions(value, loadedPath, render.Options{PreserveHTMLComments: true})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(preserved.Content, "<!-- inline -->") || !strings.Contains(preserved.Content, "<!-- multi") {
		t.Fatalf("preserved comment rendering = %q", preserved.Content)
	}
}

func TestBuildIgnoresRawDirectoryOnlySourceMarkdown(t *testing.T) {
	temporary := t.TempDir()
	markdownLibrary := filepath.Join(temporary, "markdown")
	rawLibrary := filepath.Join(temporary, "raw")
	if err := os.MkdirAll(markdownLibrary, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(rawLibrary, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(markdownLibrary, "rules.md"), "# Rules\nWork carefully.\n")
	// This frontmatter is intentionally invalid Mogent metadata but is valid raw
	// source content for a directory output.
	writeFile(t, filepath.Join(rawLibrary, "SKILL.md"), "---\nname: raw-skill\ndescription: Raw skill.\n---\n# Raw Skill\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: markdown\n  raw: raw\noutputs:\n  - path: .agents/skills/\n    from: raw:skills\ndoc:\n  - Rules: shared:rules\n")

	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content != "# Rules\n\nWork carefully.\n" {
		t.Fatalf("content = %q", result.Content)
	}
}

func TestRenderOutputSelectsAllSourceTagsAndExclusions(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library", "agents")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "go.md"), "---\ntags: [lang/go]\n---\n# Go\nUse gofmt.\n")
	writeFile(t, filepath.Join(libraryPath, "windows.md"), "---\ntags: [os/windows]\n---\n# Windows\nUse PowerShell.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "roots:\n  qmr: library\nsources:\n  agents:\n    path:\n      root: qmr\n      subdir: agents\noutputs:\n  - path: AGENTS.md\n    include:\n      - all: agents\n    exclude:\n      - tags:\n          any: [os/windows]\n  - path: GO.md\n    include:\n      - source: agents:go\n      - tags:\n          any: [lang/go]\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	first, err := render.RenderOutput(value, loadedPath, value.Outputs[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(first.Content, "# Go") || strings.Contains(first.Content, "Windows") {
		t.Fatalf("all/exclude content = %q", first.Content)
	}
	second, err := render.RenderOutput(value, loadedPath, value.Outputs[1])
	if err != nil {
		t.Fatal(err)
	}
	if strings.Count(second.Content, "# Go") != 1 {
		t.Fatalf("source/tag deduplication content = %q", second.Content)
	}
}

func TestRenderOutputAllFormatsDirectoryHeadings(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library", "decision-first")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "rules.md"), "# Rules\nChoose clearly.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\noutputs:\n  - path: AGENTS.md\n    include:\n      - all: shared\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.RenderOutput(value, loadedPath, value.Outputs[0])
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(result.Content, "# Decision First") {
		t.Fatalf("canonical output = %q", result.Content)
	}
}

func TestRenderOutputAllAppliesSidecarOrderAndTitles(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.MkdirAll(filepath.Join(libraryPath, "group"), 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "alpha.md"), "# Alpha\nAlpha body.\n")
	writeFile(t, filepath.Join(libraryPath, "zeta.md"), "# Zeta\nZeta body.\n")
	writeFile(t, filepath.Join(libraryPath, "group", "rules.md"), "# Rules\n## First\nFirst body.\n## Second\nSecond body.\n")
	writeFile(t, filepath.Join(libraryPath, "library.mogent.yaml"), `schema: {id: mogent/1}
library: {name: Test}
tree:
  - source: zeta
    title: Start Here
  - source: group
    title: Curated Group
    children:
      - source: group/rules
        children:
          - source: group/rules/second
            title: Second Curated
          - source: group/rules/first
  - source: alpha
`)
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\noutputs:\n  - path: AGENTS.md\n    include:\n      - all: shared\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.RenderOutput(value, loadedPath, value.Outputs[0])
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"# Start Here", "# Curated Group", "### Second Curated"} {
		if !strings.Contains(result.Content, want) {
			t.Fatalf("canonical output missing %q:\n%s", want, result.Content)
		}
	}
	if strings.Index(result.Content, "# Start Here") > strings.Index(result.Content, "# Curated Group") || strings.Index(result.Content, "# Curated Group") > strings.Index(result.Content, "# Alpha") || strings.Index(result.Content, "### Second Curated") > strings.Index(result.Content, "### First") {
		t.Fatalf("sidecar order not applied:\n%s", result.Content)
	}
	sources, err := render.LoadSources(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, found := sources["shared"].ByPath["group/rules/second"]; !found {
		t.Fatal("sidecar title changed source path lookup")
	}
}

func TestBuildRejectsInvalidSidecar(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "rules.md"), "# Rules\nBody.\n")
	writeFile(t, filepath.Join(libraryPath, "library.mogent.yaml"), "schema: {id: mogent/1}\nlibrary: {name: Test}\ntree: []\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\ndoc:\n  - Rules: shared:rules\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := render.Build(value, loadedPath); err == nil || !strings.Contains(err.Error(), "sidecar validation failed") {
		t.Fatalf("Build error = %v, want sidecar validation failure", err)
	}
}

func TestBuildPreservesSiblingHeadingSourceOrder(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "rules.md"), "# Rules\n## Testing\nTest behavior.\n## Commit\nCommit carefully.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\ndoc:\n  - Rules: shared:rules\n")

	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	testing := strings.Index(result.Content, "## Testing")
	commit := strings.Index(result.Content, "## Commit")
	if testing < 0 || commit < 0 || testing > commit {
		t.Fatalf("sibling heading order = %q, want Testing before Commit", result.Content)
	}
}

func TestBuildSelectsDirectorySubtreeAndNarrowsItWithExclude(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library", "engineering")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "errors.md"), "# Error Handling\nWrap errors.\n")
	writeFile(t, filepath.Join(libraryPath, "tests.md"), "# Testing\nTest behavior.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\ndoc:\n  - heading: Engineering\n    from: [shared:engineering]\n    exclude: [shared:engineering/testing]\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	want := "# Engineering\n\n## Error Handling\n\nWrap errors.\n"
	if result.Content != want {
		t.Fatalf("rendered content:\n%s\nwant:\n%s", result.Content, want)
	}
}

func TestBuildRejectsDirectoryAndDescendantInSameEntry(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library", "engineering")
	if err := os.MkdirAll(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "tests.md"), "# Testing\nTest behavior.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\ndoc:\n  - heading: Engineering\n    from: [shared:engineering, shared:engineering/testing]\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := render.Build(value, loadedPath); err == nil || !strings.Contains(err.Error(), "overlapping source roots") {
		t.Fatalf("Build error = %v, want overlapping source roots", err)
	}
}

func TestBuildUsesVerifiedPinnedURLCacheOffline(t *testing.T) {
	temporary := t.TempDir()
	commit := strings.Repeat("d", 40)
	cache := filepath.Join(temporary, ".mogent", "sources", "remote", commit)
	if err := os.MkdirAll(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(cache, "rules.md"), "# Rules\nWork carefully.\n")
	hash, err := sourcecache.HashMarkdown(cache)
	if err != nil {
		t.Fatal(err)
	}
	lock := "version: 1\nsources:\n  remote:\n    url: https://example.com/library.git\n    commit: " + commit + "\n    content_sha256: " + hash + "\n"
	writeFile(t, filepath.Join(temporary, "mogent.lock.yaml"), lock)
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  remote: https://example.com/library.git\noutput: AGENTS.md\ndoc:\n  - Rules: remote:rules\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content != "# Rules\n\nWork carefully.\n" {
		t.Fatalf("content = %q", result.Content)
	}
}

func TestBuildUsesOnlyPinnedURLSubdirOffline(t *testing.T) {
	temporary := t.TempDir()
	commit := strings.Repeat("e", 40)
	cache := filepath.Join(temporary, ".mogent", "sources", "remote", commit)
	selected := filepath.Join(cache, "libraries", "cdint")
	if err := os.MkdirAll(selected, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(selected, "rules.md"), "# Rules\nSelected.\n")
	writeFile(t, filepath.Join(cache, "unrelated.md"), "# Rules\nDuplicate outside selected root.\n")
	hash, err := sourcecache.HashMarkdown(selected)
	if err != nil {
		t.Fatal(err)
	}
	lock := "version: 1\nsources:\n  remote:\n    url: https://example.com/repository.git\n    subdir: libraries/cdint\n    commit: " + commit + "\n    content_sha256: " + hash + "\n"
	writeFile(t, filepath.Join(temporary, "mogent.lock.yaml"), lock)
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  remote:\n    location: https://example.com/repository.git\n    subdir: libraries/cdint\ndoc:\n  - Rules: remote:rules\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	result, err := render.Build(value, loadedPath)
	if err != nil {
		t.Fatal(err)
	}
	if result.Content != "# Rules\n\nSelected.\n" {
		t.Fatalf("content = %q", result.Content)
	}
}

func TestBuildRequiresExplicitRepositoryTemplateVariables(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "identity.md"), "# Identity\nRepository: {{ .repo_name }}\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\ndoc:\n  - Identity: shared:identity\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = render.Build(value, loadedPath)
	if err == nil || !strings.Contains(err.Error(), "repo_name") || !strings.Contains(err.Error(), "manifest vars") {
		t.Fatalf("Build error = %v, want missing repo_name", err)
	}
}

func TestBuildRejectsUnresolvedAndNonDescendantExclusions(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	writeFile(t, filepath.Join(libraryPath, "rules.md"), "# Rules\n## Child\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  shared: library\noutput: AGENTS.md\ndoc:\n  - heading: Rules\n    from: [shared:rules]\n    exclude: [shared:rules]\n")

	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	_, err = render.Build(value, loadedPath)
	if err == nil || !strings.Contains(err.Error(), "cannot exclude an entry source root") {
		t.Fatalf("Build error = %v, want root exclusion error", err)
	}
}

func TestBuildRejectsCollisionsEmptyNodesAndURLs(t *testing.T) {
	temporary := t.TempDir()
	first := filepath.Join(temporary, "first")
	second := filepath.Join(temporary, "second")
	for _, path := range []string{first, second} {
		if err := os.Mkdir(path, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	writeFile(t, filepath.Join(first, "rules.md"), "# Rules\n## Child\nFirst.\n")
	writeFile(t, filepath.Join(second, "rules.md"), "# Rules\n## Child\nSecond.\n")
	manifestPath := filepath.Join(temporary, "agents.yaml")
	writeFile(t, manifestPath, "sources:\n  one: first\n  two: second\noutput: AGENTS.md\ndoc:\n  - heading: Rules\n    from: [one:rules, two:rules]\n")
	value, loadedPath, err := manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := render.Build(value, loadedPath); err == nil || !strings.Contains(err.Error(), "overlapping") {
		t.Fatalf("collision error = %v", err)
	}

	writeFile(t, filepath.Join(first, "empty.md"), "# Empty\n")
	writeFile(t, manifestPath, "sources:\n  one: first\noutput: AGENTS.md\ndoc:\n  - Empty: one:empty\n")
	value, loadedPath, err = manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := render.Build(value, loadedPath); err == nil || !strings.Contains(err.Error(), "empty source node") {
		t.Fatalf("empty node error = %v", err)
	}

	writeFile(t, manifestPath, "sources:\n  remote: https://example.com/library.git\noutput: AGENTS.md\ndoc:\n  - Remote: remote:rules\n")
	value, loadedPath, err = manifest.Load(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := render.Build(value, loadedPath); err == nil || !strings.Contains(err.Error(), "not pinned") {
		t.Fatalf("URL error = %v", err)
	}
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
