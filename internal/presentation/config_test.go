package presentation_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/presentation"
)

func TestLoadDefaultsAndXDGConfig(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configRoot)
	config, err := presentation.Load()
	if err != nil {
		t.Fatal(err)
	}
	if config.Chars != "ascii" || !config.Align || config.Fit != "term" || config.Width != 0 {
		t.Fatalf("defaults = %#v", config)
	}
	directory := filepath.Join(configRoot, "mogent")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(directory, "config.yaml"), []byte("display:\n  chars: unicode\n  align: false\n  fit: none\n  width: 88\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	config, err = presentation.Load()
	if err != nil {
		t.Fatal(err)
	}
	if config.Chars != "unicode" || config.Align || config.Fit != "none" || config.Width != 88 {
		t.Fatalf("config = %#v", config)
	}
}

func TestLoadRejectsUnknownAndInvalidFields(t *testing.T) {
	configRoot := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", configRoot)
	directory := filepath.Join(configRoot, "mogent")
	if err := os.Mkdir(directory, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(directory, "config.yaml")
	for _, contents := range []string{
		"display:\n  pretty: true\n",
		"display:\n  chars: ansi\n",
		"display:\n  fit: screen\n",
		"display:\n  width: -1\n",
	} {
		if err := os.WriteFile(path, []byte(contents), 0o644); err != nil {
			t.Fatal(err)
		}
		if _, err := presentation.Load(); err == nil {
			t.Fatalf("expected error for %q", contents)
		}
	}
}
