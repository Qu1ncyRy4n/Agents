package cli

import (
	"flag"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/starter"
	"github.com/Qu1ncyRy4n/Agents/workspace"
)

type sourceBindings map[string]string

func (s sourceBindings) String() string {
	var values []string
	for alias, path := range s {
		values = append(values, alias+"="+path)
	}
	sort.Strings(values)
	return strings.Join(values, ",")
}

type variableBindings map[string]string

func (v variableBindings) String() string {
	var values []string
	for name, value := range v {
		values = append(values, name+"="+value)
	}
	sort.Strings(values)
	return strings.Join(values, ",")
}

func (v variableBindings) Set(value string) error {
	name, content, found := strings.Cut(value, "=")
	name = strings.TrimSpace(name)
	if !found || name == "" {
		return fmt.Errorf("variable %q must use name=value", value)
	}
	if _, exists := v[name]; exists {
		return fmt.Errorf("variable %q was provided twice", name)
	}
	v[name] = content
	return nil
}

func (s sourceBindings) Set(value string) error {
	alias, path, found := strings.Cut(value, "=")
	alias, path = strings.TrimSpace(alias), strings.TrimSpace(path)
	if !found || alias == "" || path == "" {
		return fmt.Errorf("source binding %q must use alias=path-or-url", value)
	}
	if _, exists := s[alias]; exists {
		return fmt.Errorf("source alias %q was provided twice", alias)
	}
	s[alias] = path
	return nil
}

func runInit(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("init", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "new manifest path")
	templateName := flags.String("template", "", "starter template name")
	output := flags.String("output", "AGENTS.md", "generated output path")
	listTemplates := flags.Bool("list-templates", false, "list starter templates")
	dryRun := flags.Bool("dry-run", false, "preview the starter manifest without writing")
	build := flags.Bool("build", false, "also build generated output")
	force := flags.Bool("force", false, "allow --build to replace a reviewed untracked output")
	sources := make(sourceBindings)
	variables := make(variableBindings)
	flags.Var(sources, "source", "available local source binding alias=path; repeat for each required alias")
	flags.Var(variables, "var", "manifest template variable name=value; repeat as needed")
	args = reorderArgs(args, map[string]bool{
		"-manifest": true, "--manifest": true,
		"-template": true, "--template": true,
		"-output": true, "--output": true,
		"-source": true, "--source": true,
		"-var": true, "--var": true,
		"-list-templates": false, "--list-templates": false,
		"-dry-run": false, "--dry-run": false,
		"-build": false, "--build": false,
		"-force": false, "--force": false,
	})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("init accepts no positional arguments")
	}
	if *listTemplates || *templateName == "" {
		if _, err := fmt.Fprintln(stdout, "Starter templates:"); err != nil {
			return fmt.Errorf("write init guide: %w", err)
		}
		for _, value := range starter.List() {
			if _, err := fmt.Fprintf(stdout, "- %s  %s\n  requires: %s\n", value.Name, value.Description, strings.Join(value.Required, ", ")); err != nil {
				return fmt.Errorf("write init guide: %w", err)
			}
		}
		if *templateName == "" {
			if _, err := fmt.Fprintln(stdout, "\nNext: mogent init --template <name> --source alias=path [--source alias=path ...] --dry-run"); err != nil {
				return fmt.Errorf("write init guide: %w", err)
			}
			return nil
		}
	}
	template, err := starter.Get(*templateName)
	if err != nil {
		return err
	}
	value, err := template.Manifest(map[string]string(sources), *output)
	if err != nil {
		return err
	}
	value.Vars, err = initVariables(*manifestFile, variables)
	if err != nil {
		return err
	}
	result, err := workspace.Initialize(workspace.InitOptions{
		ManifestPath: *manifestFile,
		Manifest:     value,
		DryRun:       *dryRun,
		Build:        *build,
		ForceOutput:  *force,
	})
	if err != nil {
		return err
	}
	if *dryRun {
		if _, err := fmt.Fprintln(stdout, "Dry run: no files written\n\nStarter manifest:"); err != nil {
			return fmt.Errorf("write init result: %w", err)
		}
		if _, err := fmt.Fprint(stdout, result.ManifestYAML); err != nil {
			return fmt.Errorf("write init result: %w", err)
		}
		return nil
	}
	if _, err := fmt.Fprintf(stdout, "Wrote starter manifest %s\n", result.ManifestPath); err != nil {
		return fmt.Errorf("write init result: %w", err)
	}
	if result.Built {
		if _, err := fmt.Fprintf(stdout, "Built %s\n", result.OutputPath); err != nil {
			return fmt.Errorf("write init result: %w", err)
		}
	} else if _, err := fmt.Fprintf(stdout, "Next: mogent source list --manifest %s\n      mogent build --manifest %s\n", result.ManifestPath, result.ManifestPath); err != nil {
		return fmt.Errorf("write init result: %w", err)
	}
	return nil
}

func initVariables(manifestFile string, supplied variableBindings) (map[string]any, error) {
	values := make(map[string]any, len(supplied)+2)
	for name, value := range supplied {
		values[name] = value
	}
	manifestPath, err := filepath.Abs(manifestFile)
	if err != nil {
		return nil, fmt.Errorf("resolve manifest path for repository variables: %w", err)
	}
	repository := filepath.Dir(manifestPath)
	if _, found := values["repo_name"]; !found {
		values["repo_name"] = filepath.Base(repository)
	}
	if _, found := values["repo_url"]; !found {
		values["repo_url"] = discoverRepositoryURL(repository)
	}
	return values, nil
}

func discoverRepositoryURL(repository string) string {
	output, err := exec.Command("git", "-C", repository, "config", "--get", "remote.origin.url").Output()
	if err != nil {
		// A new manifest need not live in a Git checkout. Materialize an empty
		// value so later builds remain independent of ambient Git configuration.
		return ""
	}
	return strings.TrimSpace(string(output))
}
