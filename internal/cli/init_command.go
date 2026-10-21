package cli

import (
	"flag"
	"fmt"
	"io"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/starter"
	"github.com/Qu1ncyRy4n/Agents/syntax"
	"github.com/Qu1ncyRy4n/Agents/v2"
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
	manifestFile := flags.String("manifest", "", "legacy YAML manifest path")
	configFile := flags.String("config", "", "legacy HCL v2 configuration path (explicit opt-in)")
	entryFile := flags.String("file", "", "plan/build path (default mogent.build, or mogent.plan with --plan)")
	plan := flags.Bool("plan", false, "create selection intent instead of an authored build")
	noInteractive := flags.Bool("no-interactive", false, "skip the optional Gum source picker")
	root := flags.String("root", "agents", "v2 server-side Markdown root")
	templateName := flags.String("template", "", "starter template name")
	output := flags.String("output", "AGENTS.md", "generated output path")
	listTemplates := flags.Bool("list-templates", false, "list starter templates")
	dryRun := flags.Bool("dry-run", false, "preview initialization without writing")
	build := flags.Bool("build", false, "also build generated output")
	force := flags.Bool("force", false, "allow --build to replace a reviewed untracked output")
	sources := make(sourceBindings)
	variables := make(variableBindings)
	flags.Var(sources, "source", "source declaration alias=local-path|git-url; repeat as needed (does not select content)")
	flags.Var(variables, "var", "manifest template variable name=value; repeat as needed")
	args = reorderArgs(args, map[string]bool{
		"-manifest": true, "--manifest": true,
		"-file": true, "--file": true,
		"-plan": false, "--plan": false,
		"-no-interactive": false, "--no-interactive": false,
		"-config": true, "--config": true,
		"-root": true, "--root": true,
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
	legacy := *configFile != "" || *manifestFile != "" || *listTemplates
	if *configFile != "" && *manifestFile != "" {
		return fmt.Errorf("init accepts only one of --config and --manifest")
	}
	if legacy && (*plan || *entryFile != "") {
		return fmt.Errorf("--plan and --file cannot be combined with legacy --config, --manifest, or --list-templates")
	}
	if !legacy {
		var legacyFlag string
		flags.Visit(func(f *flag.Flag) {
			switch f.Name {
			case "template", "root", "output", "var", "build", "force":
				legacyFlag = f.Name
			}
		})
		if legacyFlag != "" {
			return fmt.Errorf("--%s requires explicit legacy --config mogent.hcl or --manifest agents.yaml; plain init creates mogent.build (or use --plan)", legacyFlag)
		}
		kind := syntax.Build
		if *plan {
			kind = syntax.Plan
		}
		var picker sourcePicker
		if !*dryRun && !*noInteractive && len(sources) == 0 {
			picker = interactiveInitPicker(stdout, stderr)
		}
		result, err := initializeProject(assembly.InitOptions{Path: *entryFile, Kind: kind, Sources: sources, DryRun: *dryRun}, picker)
		if err != nil {
			return err
		}
		if *dryRun {
			_, err = fmt.Fprint(stdout, string(result.Content))
		} else {
			_, err = fmt.Fprintf(stdout, "Created %s\nSources are declarations only; choose includes before building.\n", result.Path)
			if err == nil {
				parsed, parseErr := syntax.Parse(result.Path, result.Content)
				if parseErr != nil {
					return parseErr
				}
				declarations, declarationErr := assembly.SourceDeclarations(parsed)
				if declarationErr != nil {
					return declarationErr
				}
				if len(declarations) == 0 {
					err = writeNoSources(stdout, result.Path)
				} else {
					_, err = fmt.Fprintf(stdout, "Next: mogent source show --file %s\n", result.Path)
				}
			}
		}
		return err
	}
	if *configFile != "" && !*listTemplates {
		if len(sources) != 1 || len(variables) != 0 {
			return fmt.Errorf("v2 init requires --config, exactly one --source alias=local-path|git-url, and optional --template, --output, or --root")
		}
		for alias, local := range sources {
			content, err := v2.Init(*configFile, *templateName, alias, local, *root, *output, *dryRun)
			if err != nil {
				return err
			}
			if *dryRun {
				_, err = fmt.Fprint(stdout, content)
			} else {
				_, err = fmt.Fprintf(stdout, "Wrote starter config %s\nNext: mogent plan --config %s\n", *configFile, *configFile)
			}
			return err
		}
	}
	if *manifestFile == "" && !*listTemplates {
		return fmt.Errorf("legacy YAML init requires --manifest agents.yaml")
	}
	if *listTemplates || *templateName == "" {
		if _, err := fmt.Fprintln(stdout, "Legacy HCL v2 templates (mogent init --config mogent.hcl --template <name> --source alias=path|url):"); err != nil {
			return fmt.Errorf("write template list: %w", err)
		}
		for _, template := range v2.InitTemplates() {
			if _, err := fmt.Fprintf(stdout, "  %s: %s\n", template.Name, template.Description); err != nil {
				return fmt.Errorf("write template list: %w", err)
			}
		}
		if _, err := fmt.Fprintln(stdout, "\nStarter templates:"); err != nil {
			return fmt.Errorf("write init guide: %w", err)
		}
		for _, value := range starter.List() {
			if _, err := fmt.Fprintf(stdout, "- %s  %s\n  requires: %s\n", value.Name, value.Description, strings.Join(value.Required, ", ")); err != nil {
				return fmt.Errorf("write init guide: %w", err)
			}
		}
		if *templateName == "" {
			if _, err := fmt.Fprintln(stdout, "\nFor legacy YAML: mogent init --manifest agents.yaml --template <name> --source alias=path --dry-run\nFor the current format: mogent init [--plan] [--source alias=path]"); err != nil {
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
