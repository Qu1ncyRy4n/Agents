package cli

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/Qu1ncyRy4n/Agents/state"
	"github.com/Qu1ncyRy4n/Agents/workspace"
)

// runNativeDrift compares current output both with the last Mogent output hash
// and with the current build composition. It is intentionally read-only.
func runNativeDrift(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("drift", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or build")
	diff := flags.Bool("diff", false, "print current-build differences from disk")
	all := flags.Bool("all", false, "report every managed output root")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true, "--diff": false, "-diff": false, "--all": false, "-all": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("drift accepts no positional arguments")
	}
	entryPath, err := projectEntryFile(*path)
	if err != nil {
		return err
	}
	file, err := assembly.LoadFile(entryPath)
	if err != nil {
		return err
	}
	result, _, err := assembly.ResolvePlan(file, entryPath)
	if err != nil {
		return err
	}
	if result.HasErrors() {
		writeAssemblyDiagnostics(stderr, result)
		return fmt.Errorf("drift cannot compose current build")
	}
	manifest, err := result.Compose()
	if err != nil {
		return err
	}
	workspaceRoot := filepath.Dir(entryPath)
	statePath := filepath.Join(workspaceRoot, ".mogent", "state.json")
	if !*all {
		return writeNativeFileDrift(stdout, filepath.Join(workspaceRoot, "AGENTS.md"), statePath, manifest.Files["AGENTS.md"].Bytes, *diff)
	}
	if err := writeNativeFileDrift(stdout, filepath.Join(workspaceRoot, "AGENTS.md"), statePath, manifest.Files["AGENTS.md"].Bytes, *diff); err != nil {
		return err
	}
	return writeNativeSkillsDrift(stdout, workspaceRoot, statePath, manifest, *diff)
}

func nativeStatusRequested(args []string) bool {
	if hasFlag(args, "--manifest") || hasFlag(args, "-manifest") {
		return false
	}
	return hasFlag(args, "--file") || hasFlag(args, "-file") || newSourceRequested(nil)
}

// runNativeStatus is the concise all-output counterpart to drift. It uses the
// current composition but never writes or repairs managed outputs.
func runNativeStatus(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("status", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "project plan or build")
	args = reorderArgs(args, map[string]bool{"--file": true, "-file": true})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("status accepts no positional arguments")
	}
	entryPath, err := projectEntryFile(*path)
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(stdout, "Project: %s\n", entryPath); err != nil {
		return err
	}
	return runNativeDrift([]string{"--file", entryPath, "--all"}, stdout, stderr)
}

func writeNativeFileDrift(stdout io.Writer, outputPath, statePath string, expected []byte, diff bool) error {
	baseline, err := state.Inspect(outputPath, statePath)
	if err != nil {
		return err
	}
	actual, err := os.ReadFile(outputPath)
	if errors.Is(err, os.ErrNotExist) {
		actual = nil
	} else if err != nil {
		return err
	}
	current := string(expected) == string(actual)
	if _, err := fmt.Fprintf(stdout, "Output: %s\nOutput baseline: %s\nCurrent build: %s\n", outputPath, baseline, map[bool]string{true: "matches", false: "differs"}[current]); err != nil {
		return err
	}
	if !diff {
		return nil
	}
	if current {
		_, err = fmt.Fprintln(stdout, "No AGENTS.md diff: current build matches disk.")
	} else {
		_, err = fmt.Fprint(stdout, workspace.UnifiedDiff(string(expected), string(actual)))
	}
	return err
}

func writeNativeSkillsDrift(stdout io.Writer, root, statePath string, manifest *content.Manifest, diff bool) error {
	outputPath := filepath.Join(root, ".agents", "skills")
	baseline, err := state.InspectDirectory(outputPath, statePath)
	if err != nil {
		return err
	}
	expected := map[string]string{}
	for name, payload := range manifest.Files {
		if strings.HasPrefix(name, ".agents/skills/") {
			expected[strings.TrimPrefix(name, ".agents/skills/")] = state.Hash(payload.Bytes)
		}
	}
	actual, err := state.DirectoryHashes(outputPath)
	if errors.Is(err, os.ErrNotExist) {
		actual = map[string]string{}
	} else if err != nil {
		return err
	}
	current := equalHashes(expected, actual)
	if _, err := fmt.Fprintf(stdout, "Output: %s\nOutput baseline: %s\nCurrent build: %s\n", outputPath, baseline, map[bool]string{true: "matches", false: "differs"}[current]); err != nil {
		return err
	}
	if diff && !current {
		return writeDirectoryDiff(stdout, expected, actual)
	}
	if diff {
		_, err = fmt.Fprintln(stdout, "No .agents/skills diff: current build matches disk.")
	}
	return err
}

func equalHashes(first, second map[string]string) bool {
	if len(first) != len(second) {
		return false
	}
	for path, hash := range first {
		if second[path] != hash {
			return false
		}
	}
	return true
}

func writeDirectoryDiff(stdout io.Writer, expected, actual map[string]string) error {
	paths := map[string]bool{}
	for path := range expected {
		paths[path] = true
	}
	for path := range actual {
		paths[path] = true
	}
	ordered := make([]string, 0, len(paths))
	for path := range paths {
		ordered = append(ordered, path)
	}
	sort.Strings(ordered)
	for _, path := range ordered {
		if expected[path] == actual[path] {
			continue
		}
		status := "modified"
		if actual[path] == "" {
			status = "missing"
		}
		if expected[path] == "" {
			status = "unexpected"
		}
		if _, err := fmt.Fprintf(stdout, "  %s: %s\n", status, path); err != nil {
			return err
		}
	}
	return nil
}
