package mint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProquintEncoding(t *testing.T) {
	for _, test := range []struct {
		value uint16
		want  string
	}{{0x7F00, "lusab"}, {0x0001, "babad"}, {0x3F54, "gutih"}, {0xDCC1, "tugad"}, {0xFFFF, "zuzuz"}} {
		if got := uint16ToProquint(test.value); got != test.want {
			t.Errorf("uint16ToProquint(%#04x) = %q, want %q", test.value, got, test.want)
		}
	}
	for _, width := range []int{1, 2} {
		got, err := Mint(width, map[string]string{}, 42, false)
		if err != nil {
			t.Fatalf("Mint(%d): %v", width, err)
		}
		if (width == 1 && !isProquint1String(got)) || (width == 2 && !isProquint2String(got)) {
			t.Errorf("Mint(%d) = %q, not a valid proquint", width, got)
		}
	}
}

func TestMintRetriesAndDryRunReportsCollision(t *testing.T) {
	first, err := Mint(1, map[string]string{}, 42, false)
	if err != nil {
		t.Fatal(err)
	}
	if got, err := Mint(1, map[string]string{first: "taken"}, 42, false); err != nil || got == first {
		t.Fatalf("collision retry = %q, %v; want a different handle", got, err)
	}
	if _, err := Mint(1, map[string]string{first: "taken"}, 42, true); err == nil || !strings.Contains(err.Error(), "collides") {
		t.Fatalf("dry-run collision error = %v, want collision", err)
	}
}

func TestScanCorpus(t *testing.T) {
	root := t.TempDir()
	empty, err := ScanCorpus(root)
	if err != nil || len(empty) != 0 {
		t.Fatalf("empty scan = %v, %v", empty, err)
	}
	writeFile(t, root, "docs/thought-experiments/TE-vapoj-substrate.md", "# placeholder\n")
	writeFile(t, root, "arbitrary/not-a-coordination-place/prefixjodonpostfix.md", "# placeholder\n")
	writeFile(t, root, "nested/bahor-nisam.txt", "# placeholder\n")
	writeFile(t, root, "notes/plain.md", "ID: DI-sojum\nID: DI-fonuz\n")
	writeFile(t, root, "alpha/manifest.json", "{}\n")
	writeFile(t, root, "beta/manifest.json", "{}\n")
	writeFile(t, root, ".git/objects/hidden/TODO-piloh-hidden.md", "ID: DI-ravuz\n")
	if err := os.MkdirAll(filepath.Join(root, "scratch", "prefixgupatpostfix"), 0o755); err != nil {
		t.Fatal(err)
	}

	corpus, err := ScanCorpus(root)
	if err != nil {
		t.Fatal(err)
	}
	for _, handle := range []string{"vapoj", "jodon", "bahor", "nisam", "bahor-nisam", "sojum", "fonuz", "manif", "gupat"} {
		if _, ok := corpus[handle]; !ok {
			t.Errorf("missing handle %q in %v", handle, corpus)
		}
	}
	for _, handle := range []string{"piloh", "ravuz"} {
		if _, ok := corpus[handle]; ok {
			t.Errorf(".git handle %q leaked into corpus", handle)
		}
	}
}

func writeFile(t *testing.T, root, path, contents string) {
	t.Helper()
	full := filepath.Join(root, path)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(full, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}
