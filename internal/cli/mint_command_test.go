package cli_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/internal/cli"
)

func TestRunMintDefaultsAndRoot(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)

	var stdout, stderr bytes.Buffer
	if err := cli.Run([]string{"mint", "-s", "42"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	output := stdout.String()
	if len(strings.TrimSpace(output)) != 5 || !strings.HasSuffix(output, "\n") {
		t.Fatalf("default mint output = %q, want one proquint-1 line", output)
	}
	if err := os.WriteFile(filepath.Join(root, "occupied-"+strings.TrimSpace(output)), []byte(""), 0o644); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if err := cli.Run([]string{"mint", "-n", "-s", "42"}, &stdout, &stderr); err == nil {
		t.Fatal("default root did not reserve its occupied handle")
	}

	other := t.TempDir()
	stdout.Reset()
	if err := cli.Run([]string{"mint", "-r", other, "-w", "2", "-n", "-s", "42"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if got := strings.TrimSpace(stdout.String()); len(got) != 11 || strings.Count(got, "-") != 1 {
		t.Fatalf("root width-two mint output = %q", got)
	}
}

func TestRunMintRejectsInvalidArguments(t *testing.T) {
	for _, args := range [][]string{{"mint", "-w", "3"}, {"mint", "extra"}} {
		var stdout, stderr bytes.Buffer
		if err := cli.Run(args, &stdout, &stderr); err == nil {
			t.Errorf("Run(%v) succeeded, want error", args)
		}
	}
}
