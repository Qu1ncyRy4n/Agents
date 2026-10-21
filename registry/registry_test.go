package registry_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/registry"
)

func TestRegistryRoundTripAndCollision(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "config", "registry.yaml")
	entries, err := registry.Load(path)
	if err != nil || len(entries) != 0 {
		t.Fatalf("missing registry=%v %v", entries, err)
	}
	if _, err := registry.Register(path, "team", "https://example.invalid/team.git"); err != nil {
		t.Fatal(err)
	}
	if _, err := registry.Register(path, "team", "https://example.invalid/other.git"); err == nil {
		t.Fatal("alias silently replaced")
	}
	entries, err = registry.Load(path)
	if err != nil || len(entries) != 1 || entries[0].Alias != "team" {
		t.Fatalf("entries=%v err=%v", entries, err)
	}
	if _, err := registry.Register(path, "team", entries[0].Location); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("registry mode=%v %v", info, err)
	}
}

func TestRegistryLocalNameAndXDGIsolation(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(root, "config"))
	if err := os.WriteFile(filepath.Join(root, "mogent.lib"), []byte("library {format=1;name=\"My Team\"}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	path, err := registry.DefaultPath()
	if err != nil || path != filepath.Join(root, "config", "mogent", "registry.yaml") {
		t.Fatalf("path=%s err=%v", path, err)
	}
	entry, err := registry.Register(path, "team", root)
	if err != nil || entry.Name != "My Team" || entry.Location != root {
		t.Fatalf("entry=%v err=%v", entry, err)
	}
}

func TestRegistryRejectsMalformedAndCredentialEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "registry.yaml")
	for _, body := range []string{
		"format: 2\nsources: []\n",
		"format: 1\nunexpected: true\n",
		"format: 1\nsources:\n- alias: team\n  location: relative\n",
		"format: 1\nsources: []\n---\nformat: 1\n",
	} {
		if err := os.WriteFile(path, []byte(body), 0600); err != nil {
			t.Fatal(err)
		}
		if _, err := registry.Load(path); err == nil {
			t.Errorf("malformed registry accepted: %s", body)
		}
	}
	if _, err := registry.Register(filepath.Join(t.TempDir(), "registry.yaml"), "team", "https://user:password@example.org/team.git"); err == nil {
		t.Fatal("credentials accepted")
	}
}
