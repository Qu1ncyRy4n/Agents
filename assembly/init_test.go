package assembly_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

func TestInitializeCreatesAuthoredBuildWithoutInstallingOutputs(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "mogent.build")
	result, err := assembly.Initialize(assembly.InitOptions{Path: path})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := syntax.Parse(path, result.Content)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Kind != syntax.Build || !strings.Contains(string(result.Content), `origin = "authored"`) {
		t.Fatalf("starter=%s", result.Content)
	}
	stored, err := os.ReadFile(path)
	if err != nil || string(stored) != string(result.Content) {
		t.Fatalf("stored=%q %v", stored, err)
	}
	for _, name := range []string{"AGENTS.md", ".agents", ".mogent", "mogent.plan", "mogent.hcl"} {
		if _, err := os.Lstat(filepath.Join(root, name)); !os.IsNotExist(err) {
			t.Errorf("init unexpectedly created %s: %v", name, err)
		}
	}
}

func TestInitializePlanDryRunDeclaresSourcesWithoutFetching(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "mogent.plan")
	result, err := assembly.Initialize(assembly.InitOptions{Path: path, Kind: syntax.Plan, DryRun: true, Sources: map[string]string{
		"team":  "https://example.invalid/library.git",
		"local": filepath.Join(root, "not-created"),
	}})
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := syntax.Parse(path, result.Content)
	if err != nil {
		t.Fatal(err)
	}
	if parsed.Kind != syntax.Plan || strings.Contains(string(result.Content), "origin =") || !strings.Contains(string(result.Content), `local = "not-created"`) {
		t.Fatalf("plan=%s", result.Content)
	}
	if strings.Index(string(result.Content), `source "local"`) > strings.Index(string(result.Content), `source "team"`) {
		t.Fatal("source order is not deterministic")
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatalf("dry-run wrote files: %v %v", entries, err)
	}
}

func TestInitializeRefusesExistingTargetsAndCompetingAuthorities(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "mogent.build")
	if err := os.WriteFile(path, []byte("keep me"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := assembly.Initialize(assembly.InitOptions{Path: path}); err == nil {
		t.Fatal("existing target overwritten")
	}
	if _, err := assembly.Initialize(assembly.InitOptions{Path: filepath.Join(root, "mogent.plan"), Kind: syntax.Plan}); err == nil {
		t.Fatal("second authority silently created")
	}
	data, err := os.ReadFile(path)
	if err != nil || string(data) != "keep me" {
		t.Fatal("existing bytes changed")
	}
}

func TestInitializeRejectsInvalidKindsSourcesAndFilenameRoles(t *testing.T) {
	for _, options := range []assembly.InitOptions{
		{Kind: syntax.Library},
		{Kind: syntax.Plan, Path: filepath.Join(t.TempDir(), "mogent.build")},
		{Sources: map[string]string{"bad:alias": "library"}},
		{Sources: map[string]string{"team": "ssh://example.org/library"}},
		{Sources: map[string]string{"team": "https://user:secret@example.org/library"}},
	} {
		if options.Path == "" {
			options.Path = filepath.Join(t.TempDir(), "mogent.build")
		}
		if _, err := assembly.Initialize(options); err == nil {
			t.Errorf("invalid options accepted: %#v", options)
		}
	}
}
