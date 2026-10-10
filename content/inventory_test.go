package content_test

import (
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/content"
)

func write(t *testing.T, root, name, body string) {
	t.Helper()
	file := filepath.Join(root, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(body), 0755); err != nil {
		t.Fatal(err)
	}
}

func TestDiscoverKeepsSkillPayloadAndDocumentViewsSeparate(t *testing.T) {
	root := t.TempDir()
	skill := "---\nname: review\ndescription: Review changes.\n---\n# Review\n\n## Procedure\nInspect changes.\n### Verify\nRun checks.\n"
	write(t, root, "review/SKILL.md", skill)
	write(t, root, "review/scripts/check.sh", "#!/bin/sh\nexit 0\n")
	write(t, root, "settings.json", "{\"x\":1}\n")
	i, err := content.Discover(root)
	if err != nil {
		t.Fatal(err)
	}
	if i.Nodes["review"].Role != "skill" || i.Nodes["review/scripts/check.sh"].Bundle != "review" {
		t.Fatalf("missing bundle: %#v", i.Nodes)
	}
	before := i.Paths()
	doc, err := i.Document("review/SKILL.md")
	if err != nil {
		t.Fatal(err)
	}
	h, err := doc.Resolve([]string{"Review", "Procedure"})
	if err != nil {
		t.Fatal(err)
	}
	if h.Kind != content.HeadingKind || len(h.Children) != 1 {
		t.Fatalf("heading=%#v", h)
	}
	if string(doc.Source) != skill || !reflect.DeepEqual(i.Paths(), before) {
		t.Fatal("view changed physical payload/inventory")
	}
	rendered, err := doc.Render([]string{"Review", "Procedure"}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if string(rendered.Bytes) != "# Procedure\nInspect changes.\n## Verify\nRun checks.\n" {
		t.Fatalf("rendered=%q", rendered.Bytes)
	}
	if _, err := i.Document("settings.json"); err == nil {
		t.Fatal("config became a document view")
	}
}

func TestDiscoverRejectsSourceSymlinks(t *testing.T) {
	root := t.TempDir()
	write(t, root, "good.md", "# Good\n")
	if err := os.Symlink("good.md", filepath.Join(root, "link.md")); err != nil {
		t.Fatal(err)
	}
	if _, err := content.Discover(root); err == nil || !strings.Contains(err.Error(), "symlink") {
		t.Fatalf("err=%v", err)
	}
}

func TestDocumentParserUsesRealHeadingsAndDiagnosesDuplicates(t *testing.T) {
	doc, err := content.ParseDocument([]byte("Intro\n\nTitle\n=====\n\n## Examples\n```md\n# Not a heading\n```\n\n## Examples\nOther.\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(doc.Headings) != 1 || doc.Headings[0].Title != "Title" || len(doc.Headings[0].Children) != 2 {
		t.Fatalf("headings=%#v", doc.Headings)
	}
	if _, err := doc.Resolve([]string{"Title", "Examples"}); err == nil || !strings.Contains(err.Error(), "ambiguous") {
		t.Fatalf("err=%v", err)
	}
	rendered, err := doc.Render([]string{"Title"}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(rendered.Bytes), "# Title\n") || !strings.Contains(string(rendered.Bytes), "# Not a heading") || strings.Contains(string(rendered.Bytes), "=====") {
		t.Fatalf("rendered=%q", rendered.Bytes)
	}
}

func TestDocumentReferencesAndValidation(t *testing.T) {
	doc, err := content.ParseDocument([]byte("# Guide\n[Checklist](references/check.md)\n![Figure](assets/figure.png)\n"))
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := doc.Render(nil, 1)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(rendered.References, []string{"references/check.md", "assets/figure.png"}) {
		t.Fatalf("refs=%v", rendered.References)
	}
	if _, err := content.ParseDocument([]byte("---\nname: broken\n# Body\n")); err == nil {
		t.Fatal("unterminated frontmatter accepted")
	}
	for _, p := range []string{"../x", "a/../x", "/x", "a\\b", ""} {
		if err := content.ValidatePath(p, true); err == nil {
			t.Errorf("unsafe path %q", p)
		}
	}
	if err := content.ValidatePath(".", true); err != nil {
		t.Fatal(err)
	}
}

func TestSetextTitleBeginningWithHashIsNotAnATXHeading(t *testing.T) {
	doc, err := content.ParseDocument([]byte("#literal\n========\n\nBody.\n"))
	if err != nil {
		t.Fatal(err)
	}
	rendered, err := doc.Render([]string{"#literal"}, 1)
	if err != nil {
		t.Fatal(err)
	}
	if string(rendered.Bytes) != "# #literal\n\nBody.\n" {
		t.Fatalf("Setext derivative=%q", rendered.Bytes)
	}
}
