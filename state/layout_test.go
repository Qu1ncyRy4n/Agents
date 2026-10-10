package state_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/state"
)

func TestDirectoryLayoutDetectsModeAndEmptyDirectoryEdits(t *testing.T) {
	root := t.TempDir()
	output := filepath.Join(root, "artifact")
	statePath := filepath.Join(root, ".mogent/state.json")
	if err := os.MkdirAll(filepath.Join(output, "empty"), 0755); err != nil {
		t.Fatal(err)
	}
	script := filepath.Join(output, "check.sh")
	if err := os.WriteFile(script, []byte("#!/bin/sh\n"), 0755); err != nil {
		t.Fatal(err)
	}
	hashes, err := state.DirectoryHashes(output)
	if err != nil {
		t.Fatal(err)
	}
	layout, err := state.InspectLayout(output)
	if err != nil {
		t.Fatal(err)
	}
	if err := state.WriteDirectoryLayout(statePath, output, hashes, layout); err != nil {
		t.Fatal(err)
	}
	if status, err := state.InspectDirectory(output, statePath); err != nil || status != state.OutputClean {
		t.Fatalf("clean=%s %v", status, err)
	}
	if err := os.Chmod(script, 0644); err != nil {
		t.Fatal(err)
	}
	if status, err := state.InspectDirectory(output, statePath); err != nil || status != state.OutputModified {
		t.Fatalf("mode drift=%s %v", status, err)
	}
	if err := os.Chmod(script, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(output, "extra"), 0755); err != nil {
		t.Fatal(err)
	}
	if status, err := state.InspectDirectory(output, statePath); err != nil || status != state.OutputModified {
		t.Fatalf("directory drift=%s %v", status, err)
	}
}
