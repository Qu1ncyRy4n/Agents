package renderfs_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
)

func TestWriteAtomicallyCreatesParentAndReplacesContent(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "output.md")
	if err := renderfs.WriteAtomically(path, []byte("first\n")); err != nil {
		t.Fatal(err)
	}
	if err := renderfs.WriteAtomically(path, []byte("second\n")); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != "second\n" {
		t.Fatalf("content = %q", contents)
	}
}

func TestWriteAtomicallyCleansTemporaryFileAfterReplaceFailure(t *testing.T) {
	directory := t.TempDir()
	target := filepath.Join(directory, "existing-directory")
	if err := os.Mkdir(target, 0o755); err != nil {
		t.Fatal(err)
	}
	err := renderfs.WriteAtomically(target, []byte("content"))
	if err == nil || !strings.Contains(err.Error(), "replace output") {
		t.Fatalf("WriteAtomically error = %v, want replace error", err)
	}
	entries, err := os.ReadDir(directory)
	if err != nil {
		t.Fatal(err)
	}
	for _, entry := range entries {
		if strings.HasPrefix(entry.Name(), ".mogent-write-") {
			t.Fatalf("temporary file was not removed: %s", entry.Name())
		}
	}
}

func TestWriteAtomicallyPreservesExistingPermissions(t *testing.T) {
	path := filepath.Join(t.TempDir(), "private.md")
	if err := os.WriteFile(path, []byte("first\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := renderfs.WriteAtomically(path, []byte("second\n")); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("mode = %o, want 600", got)
	}
}
