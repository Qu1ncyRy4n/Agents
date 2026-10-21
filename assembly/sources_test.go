package assembly_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

func TestAddSourcePreservesAuthoredBytesAndDoesNotFetch(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "mogent.build")
	original := "// retain spacing and comment\nbuild {format=1; origin=\"authored\"}\n"
	if err := os.WriteFile(path, []byte(original), 0600); err != nil {
		t.Fatal(err)
	}
	preview, err := assembly.AddSource(path, "team", "https://example.invalid/team.git", true)
	if err != nil || !strings.HasPrefix(string(preview), original) {
		t.Fatalf("preview=%s err=%v", preview, err)
	}
	stored, err := os.ReadFile(path)
	if err != nil || string(stored) != original {
		t.Fatalf("dry run changed original: %s %v", stored, err)
	}
	if _, err := assembly.AddSource(path, "team", filepath.Join(root, "missing"), false); err != nil {
		t.Fatal(err)
	}
	file, err := assembly.LoadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	sources, err := assembly.SourceDeclarations(file)
	if err != nil || len(sources) != 1 || sources[0].Local != "missing" {
		t.Fatalf("sources=%v err=%v", sources, err)
	}
	if _, err := assembly.AddSource(path, "team", "elsewhere", false); err == nil {
		t.Fatal("duplicate alias accepted")
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("mode not preserved: %v %v", info, err)
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 1 {
		t.Fatalf("unexpected writes: %v %v", entries, err)
	}
}

func TestSourceDeclarationsResolveBindingsWithoutLoading(t *testing.T) {
	file, err := syntax.Parse("test.plan", []byte("plan {format=1}\nlet root=\"unavailable\"\nsource \"team\" {local=\"${root}\"}\n"))
	if err != nil {
		t.Fatal(err)
	}
	sources, err := assembly.SourceDeclarations(file)
	if err != nil || len(sources) != 1 || sources[0].Local != "unavailable" {
		t.Fatalf("sources=%v err=%v", sources, err)
	}
}

func TestAddSourceRefusesGeneratedBuild(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mogent.build")
	if err := os.WriteFile(path, []byte("build {format=1;origin=\"generated\"}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, err := assembly.AddSource(path, "team", "library", false); err == nil || !strings.Contains(err.Error(), "plan") {
		t.Fatalf("generated edit: %v", err)
	}
}
