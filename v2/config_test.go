package v2

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadParsesSourcesAndOutputs(t *testing.T) {
	path := writeConfig(t, `
mogent {
  format = 2
}

sources {
  source "personal" {
    local = "../library"
  }

  source "shared" {
    git = "https://example.com/library.git"
    ref = "main"
  }
}

outputs {
  output "agents" {
    paths = ["AGENTS.md", "CLAUDE.md"]
    kind = "markdown"

    source "personal" {
      from = "personal:agents"
      select = { all = true }
    }
  }
}
`)
	config, err := Load(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Sources) != 2 || config.Sources["shared"].Ref != "main" {
		t.Fatalf("sources = %#v", config.Sources)
	}
	if len(config.Outputs) != 1 || len(config.Outputs[0].Paths) != 2 || config.Outputs[0].Sources[0].Select == nil {
		t.Fatalf("outputs = %#v", config.Outputs)
	}
}

func TestLoadParsesTreeExcludes(t *testing.T) {
	config := mustLoadConfig(t, `
mogent { format = 2 }
sources {
  source "shared" { local = "library" }
}
outputs {
  output "skills" {
    path = ".agents/skills"
    kind = "tree"
    source "shared" {
      from = "shared:skills"
      select = { all = true }
      exclude = ["experimental", "legacy/old-skill"]
    }
  }
}
`)
	if got := strings.Join(config.Outputs[0].Sources[0].Exclude, " "); got != "experimental legacy/old-skill" {
		t.Fatalf("exclude = %q", got)
	}
}

func TestLoadRejectsInvalidStructuralConfig(t *testing.T) {
	cases := []struct {
		name    string
		content string
		want    string
	}{
		{
			name:    "source-without-location",
			content: "mogent {\n  format = 2\n}\n\nsources {\n  source \"bad\" {}\n}\n\noutputs {\n  output \"x\" {\n    path = \"AGENTS.md\"\n    kind = \"markdown\"\n    source \"bad\" {\n      from = \"bad:agents\"\n      select = {}\n    }\n  }\n}\n",
			want:    "declare exactly one of git or local",
		},
		{
			name:    "unsafe-output-path",
			content: "mogent {\n  format = 2\n}\n\nsources {\n  source \"local\" {\n    local = \"library\"\n  }\n}\n\noutputs {\n  output \"x\" {\n    path = \"../AGENTS.md\"\n    kind = \"markdown\"\n    source \"local\" {\n      from = \"local:agents\"\n      select = {}\n    }\n  }\n}\n",
			want:    "must be a safe relative path",
		},
		{
			name:    "undeclared-output-source",
			content: "mogent {\n  format = 2\n}\n\nsources {\n  source \"local\" {\n    local = \"library\"\n  }\n}\n\noutputs {\n  output \"x\" {\n    path = \"AGENTS.md\"\n    kind = \"markdown\"\n    source \"other\" {\n      from = \"other:agents\"\n      select = {}\n    }\n  }\n}\n",
			want:    "is not declared in sources",
		},
		{
			name:    "markdown-exclude",
			content: "mogent {\n  format = 2\n}\n\nsources {\n  source \"local\" { local = \"library\" }\n}\n\noutputs {\n  output \"x\" {\n    path = \"AGENTS.md\"\n    kind = \"markdown\"\n    source \"local\" {\n      from = \"local:agents\"\n      select = {}\n      exclude = [\"old\"]\n    }\n  }\n}\n",
			want:    "exclude is valid only for tree outputs",
		},
		{
			name:    "unsafe-tree-exclude",
			content: "mogent {\n  format = 2\n}\n\nsources {\n  source \"local\" { local = \"library\" }\n}\n\noutputs {\n  output \"x\" {\n    path = \".agents\"\n    kind = \"tree\"\n    source \"local\" {\n      from = \"local:skills\"\n      select = { all = true }\n      exclude = [\"../old\"]\n    }\n  }\n}\n",
			want:    "safe non-root relative path",
		},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := Load(writeConfig(t, test.content))
			if err == nil || !strings.Contains(err.Error(), test.want) {
				t.Fatalf("Load error = %v, want %q", err, test.want)
			}
		})
	}
}

func writeConfig(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ConfigFile)
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func mustLoadConfig(t *testing.T, content string) *Config {
	t.Helper()
	config, err := Load(writeConfig(t, content))
	if err != nil {
		t.Fatal(err)
	}
	return config
}
