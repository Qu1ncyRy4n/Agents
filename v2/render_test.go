package v2

import "testing"

func TestRenderMarkdownUsesSidecarHierarchyAndStripsMatchingSourceHeading(t *testing.T) {
	library := planLibrary(t)
	plan := &Plan{Outputs: []PlannedOutput{{
		Name: "agents",
		Kind: "markdown",
		Sources: []PlannedSource{{
			Name:     "shared",
			Sections: []*Section{library.ByPath["org/core/review"], library.ByPath["org/core/release"]},
		}},
	}}}
	outputs, err := RenderMarkdown(plan, map[string]*Library{"shared": library})
	if err != nil {
		t.Fatal(err)
	}
	want := "# Organization\n\n## Core\n\n### Review\n\n### Release\n"
	if got := outputs["agents"]; got != want {
		t.Fatalf("rendered output:\n%s\nwant:\n%s", got, want)
	}
}

func TestRenderMarkdownRejectsMismatchedSourceHeading(t *testing.T) {
	library := planLibrary(t)
	library.ByPath["org/core/review"].Title = "Code Review"
	plan := &Plan{Outputs: []PlannedOutput{{
		Name:    "agents",
		Kind:    "markdown",
		Sources: []PlannedSource{{Name: "shared", Sections: []*Section{library.ByPath["org/core/review"]}}},
	}}}
	_, err := RenderMarkdown(plan, map[string]*Library{"shared": library})
	if err == nil || !contains(err.Error(), "# Code Review") {
		t.Fatalf("RenderMarkdown error = %v", err)
	}
}

func contains(value, substring string) bool {
	for index := 0; index+len(substring) <= len(value); index++ {
		if value[index:index+len(substring)] == substring {
			return true
		}
	}
	return false
}
