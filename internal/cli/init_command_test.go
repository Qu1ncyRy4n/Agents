package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitDefaultCreatesOnlyAuthoredBuild(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"init"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile("mogent.build")
	if err != nil || !strings.Contains(string(data), `origin = "authored"`) {
		t.Fatalf("build=%s err=%v", data, err)
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 1 {
		t.Fatalf("init wrote unexpected files: %v %v", entries, err)
	}
	if !strings.Contains(stdout.String(), "Created mogent.build") || strings.Contains(stdout.String(), "HCL") {
		t.Fatalf("init output=%s", stdout.String())
	}
	if err := Run([]string{"init"}, &stdout, &stderr); err == nil {
		t.Fatal("repeated init should refuse overwrite")
	}
	if err := Run([]string{"init", "--plan"}, &stdout, &stderr); err == nil {
		t.Fatal("init should not introduce competing authority")
	}
}

func TestInitPlanDryRunAndSourceInspection(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	writeNewCLIFile(t, root, "team/mogent.lib", "library {format=1;name=\"Team\"}\nmd-section \"workflow\" {source=\"workflow.md\"}\n")
	writeNewCLIFile(t, root, "team/workflow.md", "# Workflow\nRead the project notes.\n")
	var stdout, stderr bytes.Buffer
	args := []string{"init", "--plan", "--source", "team=team"}
	if err := Run(append(args, "--dry-run"), &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "plan {") || strings.Contains(stdout.String(), "include ") {
		t.Fatalf("preview=%s", stdout.String())
	}
	if _, err := os.Stat("mogent.plan"); !os.IsNotExist(err) {
		t.Fatalf("dry run wrote entry: %v", err)
	}
	if err := Run(args, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	stdout.Reset()
	if err := Run([]string{"show", "--source", "--file", "mogent.plan"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(stdout.String(), "workflow") {
		t.Fatalf("source inspection=%s", stdout.String())
	}
	if _, err := os.Stat("AGENTS.md"); !os.IsNotExist(err) {
		t.Fatalf("init/inspection installed output: %v", err)
	}
}

func TestInitExplicitLegacyHCLAndInvalidFlagMixes(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "mogent.hcl")
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"init", "--config", path, "--source", "team=library"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path)
	if err != nil || !strings.Contains(string(data), "format = 2") {
		t.Fatalf("legacy config=%s err=%v", data, err)
	}
	for _, args := range [][]string{
		{"--plan", "--config", path},
		{"--file", filepath.Join(root, "mogent.build"), "--manifest", "agents.yaml"},
		{"--config", path, "--manifest", "agents.yaml"},
		{"--template", "qmr-core"},
		{"--build"},
		{"--output", "custom.md"},
	} {
		if err := Run(append([]string{"init"}, args...), &stdout, &stderr); err == nil {
			t.Errorf("invalid init flags accepted: %v", args)
		}
	}
}
