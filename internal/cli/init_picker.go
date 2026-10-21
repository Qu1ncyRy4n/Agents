package cli

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/charmbracelet/x/term"
)

type sourcePicker interface {
	Choose(header string, choices []string) (string, error)
	Input(prompt string) (string, error)
}

func initializeProject(options assembly.InitOptions, picker sourcePicker) (*assembly.InitResult, error) {
	if picker != nil && !options.DryRun && len(options.Sources) == 0 {
		preview := options
		preview.DryRun = true
		if _, err := assembly.Initialize(preview); err != nil {
			return nil, err
		}
		selected, err := chooseInitSources(picker)
		if err != nil {
			return nil, err
		}
		options.Sources = selected
	}
	return assembly.Initialize(options)
}

// Gum is an optional presentation executable, never a build dependency. Its UI
// goes to the terminal on stderr; only the selected value is captured on stdout.
type gumPicker struct {
	executable string
	input      io.Reader
	output     io.Writer
}

func (g gumPicker) run(args ...string) (string, error) {
	command := exec.Command(g.executable, args...)
	command.Stdin = g.input
	command.Stderr = g.output
	var selected bytes.Buffer
	command.Stdout = &selected
	if err := command.Run(); err != nil {
		return "", fmt.Errorf("source selection cancelled or failed; no project files written: %w", err)
	}
	return strings.TrimSuffix(strings.TrimSuffix(selected.String(), "\n"), "\r"), nil
}

func (g gumPicker) Choose(header string, choices []string) (string, error) {
	args := []string{"choose", "--header", header, "--"}
	return g.run(append(args, choices...)...)
}

func (g gumPicker) Input(prompt string) (string, error) {
	return g.run("input", "--prompt", prompt+": ")
}

func interactiveInitPicker(stdout, stderr io.Writer) sourcePicker {
	output, ok := stdout.(*os.File)
	if !ok || !term.IsTerminal(output.Fd()) || !term.IsTerminal(os.Stdin.Fd()) {
		return nil
	}
	errorOutput, ok := stderr.(*os.File)
	if !ok || !term.IsTerminal(errorOutput.Fd()) {
		return nil
	}
	executable, err := exec.LookPath("gum")
	if err != nil {
		_, _ = fmt.Fprintln(stderr, "Tip: install gum for interactive source selection, or use init --source name=path-or-url.")
		return nil
	}
	return gumPicker{executable: executable, input: os.Stdin, output: stderr}
}

func chooseInitSources(picker sourcePicker) (sourceBindings, error) {
	entries, err := knownLibraries()
	if err != nil {
		return nil, err
	}
	choices := make([]string, 0, len(entries)+3)
	locations := make(map[string]sourceBindings)
	for _, entry := range entries {
		label := entry.Alias
		if entry.Name != "" {
			label += " — " + entry.Name
		}
		label += " (" + entry.Location + ")"
		choices = append(choices, label)
		locations[label] = sourceBindings{entry.Alias: entry.Location}
	}
	const local, git, skip = "Local directory", "Git repository", "Skip for now"
	choices = append(choices, local, git, skip)
	selected, err := picker.Choose("Add a source — known libraries or a new location", choices)
	if err != nil {
		return nil, err
	}
	if binding, found := locations[selected]; found {
		return binding, nil
	}
	if selected == skip {
		return sourceBindings{}, nil
	}
	if selected != local && selected != git {
		return nil, fmt.Errorf("source picker returned an unknown choice")
	}
	alias, err := picker.Input("Source name")
	if err != nil {
		return nil, err
	}
	prompt := "Local directory path"
	if selected == git {
		prompt = "Git URL (https:// or file://)"
	}
	location, err := picker.Input(prompt)
	if err != nil {
		return nil, err
	}
	if selected == git && !strings.Contains(location, "://") || selected == local && strings.Contains(location, "://") {
		return nil, fmt.Errorf("location does not match selected source kind")
	}
	return sourceBindings{alias: location}, nil
}
