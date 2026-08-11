package starter

import "testing"

func TestTemplatesRequireExplicitSourcesAndReturnIndependentManifests(t *testing.T) {
	template, err := Get("go")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := template.Manifest(map[string]string{"shared": "library"}, "AGENTS.md"); err == nil {
		t.Fatal("expected missing go source to fail")
	}
	first, err := template.Manifest(map[string]string{"shared": "library", "go": "library"}, "AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	second, err := template.Manifest(map[string]string{"shared": "library", "go": "library"}, "AGENTS.md")
	if err != nil {
		t.Fatal(err)
	}
	first.Doc[0].Heading = "Changed"
	if second.Doc[0].Heading == "Changed" {
		t.Fatal("template manifests share mutable entries")
	}
}

func TestTemplateListIsStable(t *testing.T) {
	values := List()
	want := []string{"go", "minimal", "personal-go-nix", "research-python"}
	if len(values) != len(want) {
		t.Fatalf("templates = %#v", values)
	}
	for index := range want {
		if values[index].Name != want[index] {
			t.Fatalf("template[%d] = %q, want %q", index, values[index].Name, want[index])
		}
	}
}
