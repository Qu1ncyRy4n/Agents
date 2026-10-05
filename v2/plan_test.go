package v2

import (
	"strings"
	"testing"
)

func TestCompileSelectsExactDescendantWithoutUnrelatedExplicitChoice(t *testing.T) {
	config := loadPlanConfig(t, `
mogent { format = 2 }
sources {
  source "shared" { local = "library" }
}
outputs {
  output "agents" {
    path = "AGENTS.md"
    kind = "markdown"
    source "shared" {
      from = "shared:agents"
      select = {
        org = {
          core = {
            review = true
          }
        }
        else = "exclude"
      }
    }
  }
}
`)
	library := planLibrary(t)
	plan := Compile(config, map[string]*Library{"shared": library})
	if diagnostics := planErrors(plan); len(diagnostics) != 0 {
		t.Fatalf("errors = %#v", diagnostics)
	}
	sections := plan.Outputs[0].Sources[0].Sections
	if len(sections) != 1 || sections[0].Path != "org/core/review" {
		t.Fatalf("selected sections = %#v", sections)
	}
}

func TestCompileReportsIncompleteBroadExplicitSelection(t *testing.T) {
	config := loadPlanConfig(t, `
mogent { format = 2 }
sources {
  source "shared" { local = "library" }
}
outputs {
  output "agents" {
    path = "AGENTS.md"
    kind = "markdown"
    source "shared" {
      from = "shared:agents"
      select = {
        org = {
          core = true
        }
      }
    }
  }
}
`)
	plan := Compile(config, map[string]*Library{"shared": planLibrary(t)})
	if got := planErrors(plan); len(got) != 1 || got[0].Code != "MOGENT102" || !strings.Contains(got[0].Suggestion, "release") {
		t.Fatalf("errors = %#v", got)
	}
}

func TestCompileRootAllSelectsDescendantContent(t *testing.T) {
	config := loadPlanConfig(t, `
mogent { format = 2 }
sources {
  source "shared" { local = "library" }
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
	plan := Compile(config, map[string]*Library{"shared": planLibrary(t)})
	if got := planErrors(plan); len(got) != 1 || got[0].Code != "MOGENT102" {
		t.Fatalf("errors = %#v", got)
	}
}

func loadPlanConfig(t *testing.T, content string) *Config {
	t.Helper()
	return mustLoadConfig(t, content)
}

func planLibrary(t *testing.T) *Library {
	t.Helper()
	root := t.TempDir()
	writeLibraryFile(t, root, "agents/org/core/review.md", "# Review\n")
	writeLibraryFile(t, root, "agents/org/core/release.md", "# Release\n")
	writeLibraryFile(t, root, LibraryFile, `
library {
  format = 2
  id = "test"
  name = "Test"
}
content {
  markdown_root = "agents"
}
sections {
  section "org" {
    title = "Organization"
    section "core" {
      title = "Core"
      inclusion {
        policy = "explicit"
        defaults = {
          review = true
          release = false
        }
      }
      section "review" {
        title = "Review"
        source = "org/core/review.md"
      }
      section "release" {
        title = "Release"
        source = "org/core/release.md"
      }
    }
  }
}
`)
	library, err := LoadLibrary(root)
	if err != nil {
		t.Fatal(err)
	}
	return library
}

func planErrors(plan *Plan) []Diagnostic {
	var diagnostics []Diagnostic
	for _, diagnostic := range plan.Diagnostics {
		if diagnostic.Severity == SeverityError {
			diagnostics = append(diagnostics, diagnostic)
		}
	}
	return diagnostics
}
