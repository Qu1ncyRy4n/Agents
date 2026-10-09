package v2

import (
	"strings"
	"testing"
)

// offerLibrary has one branch per offer value so each consumer rule can be
// exercised in isolation.
func offerLibrary(t *testing.T) *Library {
	t.Helper()
	root := t.TempDir()
	for _, name := range []string{"identity", "security", "review", "release", "triage", "experiment"} {
		writeLibraryFile(t, root, "agents/"+name+".md", "# "+strings.ToUpper(name[:1])+name[1:]+"\n")
	}
	writeLibraryFile(t, root, LibraryFile, `
library {
  format = 2
  id = "offers"
  name = "Offers"
}
content {
  markdown_root = "agents"
}
sections {
  section "org" {
    title = "Organization"
    section "base" {
      title = "Base"
      curate = "foundation"
      section "identity" {
        title = "Identity"
        source = "identity.md"
      }
      section "security" {
        title = "Security"
        source = "security.md"
      }
    }
    section "core" {
      title = "Core"
      curate = "choose"
      defaults = {
        review = true
        release = false
      }
      section "review" {
        title = "Review"
        source = "review.md"
      }
      section "release" {
        title = "Release"
        source = "release.md"
      }
    }
    section "extras" {
      title = "Extras"
      curate = "optional"
      default = false
      section "triage" {
        title = "Triage"
        source = "triage.md"
      }
    }
    section "lab" {
      title = "Lab"
      curate = "opt_in"
      section "experiment" {
        title = "Experiment"
        source = "experiment.md"
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

func offerConfig(t *testing.T, selection string) *Config {
	t.Helper()
	return loadPlanConfig(t, `
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
      select = `+selection+`
    }
  }
}
`)
}

func selectedPaths(plan *Plan) []string {
	var paths []string
	for _, output := range plan.Outputs {
		for _, source := range output.Sources {
			for _, section := range source.Sections {
				paths = append(paths, section.Path)
			}
		}
	}
	return paths
}

func TestBroadSelectionHonorsEachOffer(t *testing.T) {
	plan := Compile(offerConfig(t, `{
  org = {
    base   = true
    core   = { accept_defaults = true }
    extras = true
    lab    = true
  }
}`), map[string]*Library{"shared": offerLibrary(t)})
	if got := planErrors(plan); len(got) != 0 {
		t.Fatalf("errors = %#v", got)
	}
	want := "org/base/identity org/base/security org/core/review"
	if got := strings.Join(selectedPaths(plan), " "); got != want {
		t.Fatalf("selected = %q, want %q", got, want)
	}
	if codes := diagnosticCodes(plan); strings.Join(codes, " ") != "MOGENT203 MOGENT203" {
		t.Fatalf("diagnostics = %v", codes)
	}
}

func TestFoundationChildNeedsForceExcludeWithReason(t *testing.T) {
	cases := []struct{ name, selection, want string }{
		{"plain-false", `{ org = { base = { security = false } } }`, "requires force_exclude = true and a reason"},
		{"no-reason", `{ org = { base = { security = { force_exclude = true } } } }`, "requires a non-empty reason"},
		{"else-exclude", `{ org = { base = { identity = true, else = "exclude" } } }`, "else = \"exclude\" cannot drop foundation children"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			plan := Compile(offerConfig(t, test.selection), map[string]*Library{"shared": offerLibrary(t)})
			got := planErrors(plan)
			if len(got) != 1 || got[0].Code != "MOGENT107" || !strings.Contains(got[0].Message, test.want) {
				t.Fatalf("errors = %#v", got)
			}
		})
	}
}

func TestFoundationForceExcludeWithReasonWarnsAndDrops(t *testing.T) {
	plan := Compile(offerConfig(t, `{
  org = {
    base = {
      security = {
        force_exclude = true
        reason        = "Isolated fixture."
      }
    }
  }
}`), map[string]*Library{"shared": offerLibrary(t)})
	if got := planErrors(plan); len(got) != 0 {
		t.Fatalf("errors = %#v", got)
	}
	if got := strings.Join(selectedPaths(plan), " "); got != "org/base/identity" {
		t.Fatalf("selected = %q", got)
	}
	if codes := diagnosticCodes(plan); len(codes) != 1 || codes[0] != "MOGENT205" {
		t.Fatalf("diagnostics = %v", codes)
	}
}

func TestForceExcludeOutsideFoundationIsError(t *testing.T) {
	plan := Compile(offerConfig(t, `{
  org = {
    core = {
      review  = { force_exclude = true, reason = "no" }
      release = false
    }
  }
}`), map[string]*Library{"shared": offerLibrary(t)})
	got := planErrors(plan)
	if len(got) != 1 || got[0].Code != "MOGENT107" || !strings.Contains(got[0].Message, "use review = false") {
		t.Fatalf("errors = %#v", got)
	}
}

func TestAcceptDefaultsOutsideChooseIsError(t *testing.T) {
	plan := Compile(offerConfig(t, `{ org = { extras = { accept_defaults = true } } }`), map[string]*Library{"shared": offerLibrary(t)})
	got := planErrors(plan)
	if len(got) != 1 || got[0].Code != "MOGENT108" || !strings.Contains(got[0].Message, "curates as optional") {
		t.Fatalf("errors = %#v", got)
	}
}

func TestOptInChildSelectedExplicitlyWarns(t *testing.T) {
	plan := Compile(offerConfig(t, `{ org = { lab = { experiment = true } } }`), map[string]*Library{"shared": offerLibrary(t)})
	if got := planErrors(plan); len(got) != 0 {
		t.Fatalf("errors = %#v", got)
	}
	if got := strings.Join(selectedPaths(plan), " "); got != "org/lab/experiment" {
		t.Fatalf("selected = %q", got)
	}
	if codes := diagnosticCodes(plan); len(codes) != 1 || codes[0] != "MOGENT204" {
		t.Fatalf("diagnostics = %v", codes)
	}
}

func TestTagQueriesSelectLeavesAndWarnWhenUnmatched(t *testing.T) {
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
      from     = "shared:agents"
      select   = { else = "exclude" }
      tags_any = ["lang/go", "ghost/none"]
    }
  }
}
`)
	library := offerLibrary(t)
	library.ByPath["org/extras/triage"].Tags = []string{"lang/go/testing"}
	plan := Compile(config, map[string]*Library{"shared": library})
	if got := planErrors(plan); len(got) != 0 {
		t.Fatalf("errors = %#v", got)
	}
	if got := strings.Join(selectedPaths(plan), " "); got != "org/extras/triage" {
		t.Fatalf("selected = %q", got)
	}
	codes := diagnosticCodes(plan)
	if len(codes) != 1 || codes[0] != "MOGENT201" || !strings.Contains(plan.Diagnostics[0].Message, "ghost/none") {
		t.Fatalf("diagnostics = %#v", plan.Diagnostics)
	}
}

func TestLoadLibraryRejectsBadOffers(t *testing.T) {
	leaf := "section \"leaf\" {\n        title = \"Leaf\"\n        source = \"leaf.md\"\n      }"
	cases := []struct{ name, section, want string }{
		{"legacy-inclusion", "section \"b\" {\n      title = \"B\"\n      inclusion { policy = \"baseline\" }\n      " + leaf + "\n    }", "removed inclusion block"},
		{"removed-offer", "section \"b\" {\n      title = \"B\"\n      offer = \"foundation\"\n      " + leaf + "\n    }", "uses removed offer; rename it to curate"},
		{"unknown-curate", "section \"b\" {\n      title = \"B\"\n      curate = \"baseline\"\n      " + leaf + "\n    }", "unknown curate value \"baseline\""},
		{"curate-on-leaf", "section \"b\" {\n      title = \"B\"\n      curate = \"foundation\"\n      source = \"leaf.md\"\n    }", "curate requires child sections"},
		{"choose-missing-default", "section \"b\" {\n      title = \"B\"\n      curate = \"choose\"\n      defaults = {}\n      " + leaf + "\n    }", "requires defaults for every direct child"},
		{"choose-unknown-default", "section \"b\" {\n      title = \"B\"\n      curate = \"choose\"\n      defaults = { ghost = true }\n      " + leaf + "\n    }", "no default for child \"leaf\""},
		{"optional-without-default", "section \"b\" {\n      title = \"B\"\n      curate = \"optional\"\n      " + leaf + "\n    }", "optional curate requires default"},
		{"foundation-with-default", "section \"b\" {\n      title = \"B\"\n      curate = \"foundation\"\n      default = true\n      " + leaf + "\n    }", "only an optional curate may declare default"},
		{"opt-in-with-defaults", "section \"b\" {\n      title = \"B\"\n      curate = \"opt_in\"\n      defaults = { leaf = true }\n      " + leaf + "\n    }", "only a choose curate may declare defaults"},
		{"defaults-without-curate", "section \"b\" {\n      title = \"B\"\n      defaults = { leaf = true }\n      " + leaf + "\n    }", "without curate"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			writeLibraryFile(t, root, "agents/leaf.md", "# Leaf\n")
			writeLibraryFile(t, root, LibraryFile, "library {\n  format = 2\n  id = \"t\"\n  name = \"T\"\n}\ncontent {\n  markdown_root = \"agents\"\n}\nsections {\n  section \"a\" {\n    title = \"A\"\n    "+test.section+"\n  }\n}\n")
			_, err := LoadLibrary(root)
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("LoadLibrary error = %v, want %q", err, test.want)
			}
		})
	}
}
