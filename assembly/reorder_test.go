package assembly_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/assembly"
)

func TestReorderGroupsSourcesWithoutChangingIncludeOrder(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mogent.build")
	input := "// retain\nbuild {format=1;origin=\"authored\"}\ninclude \"team:first\"\nsource \"team\" {local=\"team\"}\nlet name=\"Demo\"\ninclude \"team:second\"\nsource \"other\" {local=\"other\"}\n"
	if err := os.WriteFile(path, []byte(input), 0600); err != nil {
		t.Fatal(err)
	}
	preview, changed, err := assembly.Reorder(path, true)
	if err != nil || !changed {
		t.Fatalf("changed=%t err=%v", changed, err)
	}
	text := string(preview)
	for _, order := range [][2]string{{`let name`, `source "team"`}, {`source "team"`, `source "other"`}, {`source "other"`, `include "team:first"`}, {`include "team:first"`, `include "team:second"`}} {
		if strings.Index(text, order[0]) >= strings.Index(text, order[1]) {
			t.Fatalf("wrong order %q before %q:\n%s", order[0], order[1], text)
		}
	}
	stored, err := os.ReadFile(path)
	if err != nil || string(stored) != input {
		t.Fatalf("dry run changed input: %s %v", stored, err)
	}
	if _, changed, err := assembly.Reorder(path, false); err != nil || !changed {
		t.Fatalf("write changed=%t err=%v", changed, err)
	}
	if _, changed, err := assembly.Reorder(path, false); err != nil || changed {
		t.Fatalf("idempotent changed=%t err=%v", changed, err)
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("mode=%v err=%v", info, err)
	}
}

func TestReorderRejectsGeneratedBuild(t *testing.T) {
	path := filepath.Join(t.TempDir(), "mogent.build")
	if err := os.WriteFile(path, []byte("build {format=1;origin=\"generated\"}\n"), 0644); err != nil {
		t.Fatal(err)
	}
	if _, _, err := assembly.Reorder(path, true); err == nil || !strings.Contains(err.Error(), "generated") {
		t.Fatalf("reorder generated=%v", err)
	}
}
