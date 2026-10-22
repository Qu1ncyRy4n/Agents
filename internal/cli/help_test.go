package cli

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestWorkflowHelpLeadsToSelectionWithoutLegacyFlagDump(t *testing.T) {
	var expected string
	for _, args := range [][]string{nil, {"help"}, {"--help"}, {"-h"}} {
		var stdout, stderr bytes.Buffer
		if err := Run(args, &stdout, &stderr); err != nil {
			t.Fatal(err)
		}
		text := stdout.String()
		if expected == "" {
			expected = text
		} else if text != expected {
			t.Fatalf("help entry points disagree: %v", args)
		}
		previous := -1
		for _, section := range []string{"Project initialization", "Browse content", "Add modules", "Preview the selection", "Build outputs", "Further help"} {
			index := strings.Index(text, section)
			if index <= previous {
				t.Fatalf("missing or out-of-order section %q:\n%s", section, text)
			}
			previous = index
		}
		for _, want := range []string{"module add <name:path>", "module rm <name:path>", "mogent help module", "mogent help commands", "mogent help --all"} {
			if !strings.Contains(text, want) {
				t.Fatalf("workflow missing %q:\n%s", want, text)
			}
		}
		for _, advancedFlag := range []string{"--headings", "--under", "--preserve-html-comments", "--no-interactive", "--dry-run", "--known"} {
			if strings.Contains(text, advancedFlag) {
				t.Fatalf("top-level help dumps advanced flag %s", advancedFlag)
			}
		}
		if stderr.Len() != 0 {
			t.Fatalf("help wrote errors: %s", stderr.String())
		}
	}
}

func TestExpandedWorkflowHelpIncludesFlags(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"help", "--all"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	expanded := stdout.String()
	for _, want := range []string{"Project initialization", "Browse content", "--plan", "--source <name=path-or-url>", "--known", "--flat", "--check", "--dry-run", "--file <path>"} {
		if !strings.Contains(expanded, want) {
			t.Errorf("expanded help missing %q", want)
		}
	}
	stdout.Reset()
	if err := Run([]string{"--help"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	if len(expanded) <= stdout.Len() {
		t.Fatal("expanded help should contain more detail than the overview")
	}
}

func TestFocusedHelpDoesNotInitializeOrRequireProject(t *testing.T) {
	root := t.TempDir()
	t.Chdir(root)
	for _, test := range []struct {
		args []string
		want string
	}{
		{[]string{"init", "--help"}, "Advanced"},
		{[]string{"help", "source", "add"}, "source add"},
		{[]string{"source", "add", "--help"}, "remembered libraries"},
		{[]string{"help", "module"}, "mogent module add qmr:intro/author-ownership"},
		{[]string{"help", "edit"}, "This is the only edit operation currently available"},
		{[]string{"edit", "--help"}, "reference must already be selected"},
		{[]string{"build", "--help"}, "existing outputs"},
		{[]string{"show", "-h"}, "contribution preview"},
		{[]string{"source", "show", "--help"}, "mogent source show"},
	} {
		var stdout, stderr bytes.Buffer
		if err := Run(test.args, &stdout, &stderr); err != nil {
			t.Fatalf("%v: %v", test.args, err)
		}
		if !strings.Contains(stdout.String(), test.want) {
			t.Fatalf("%v missing %q:\n%s", test.args, test.want, stdout.String())
		}
	}
	entries, err := os.ReadDir(root)
	if err != nil || len(entries) != 0 {
		t.Fatalf("help wrote project files: %v %v", entries, err)
	}
}

func TestEveryCommandRetainsDiscoverableHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if err := Run([]string{"help", "commands"}, &stdout, &stderr); err != nil {
		t.Fatal(err)
	}
	reference := stdout.String()
	for _, command := range commands() {
		if !strings.Contains(reference, "mogent "+command.usage) {
			t.Errorf("reference missing %s", command.name)
		}
		stdout.Reset()
		if err := Run([]string{command.name, "--help"}, &stdout, &stderr); err != nil {
			t.Errorf("%s --help: %v", command.name, err)
		}
		if stdout.Len() == 0 {
			t.Errorf("%s has empty help", command.name)
		}
	}
}

func TestHelpRejectsUnknownTopics(t *testing.T) {
	for _, args := range [][]string{
		{"help", "typo"},
		{"help", "source", "typo"},
		{"help", "build", "extra"},
		{"help", "source", "list", "extra"},
	} {
		var stdout, stderr bytes.Buffer
		if err := Run(args, &stdout, &stderr); err == nil {
			t.Errorf("unknown help path accepted: %v", args)
		}
	}
}
