package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/Qu1ncyRy4n/Agents/v2"
)

func runPlan(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("plan", flag.ContinueOnError)
	flags.SetOutput(stderr)
	configPath := flags.String("config", v2.ConfigFile, "path to v2 HCL configuration")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("plan accepts no positional arguments")
	}
	config, err := v2.Load(*configPath)
	if err != nil {
		return err
	}
	libraries, err := v2.LoadLocalLibraries(config)
	if err != nil {
		return err
	}
	plan := v2.Compile(config, libraries)
	for _, output := range plan.Outputs {
		if _, err := fmt.Fprintf(stdout, "Output %s: %v\n", output.Name, output.Paths); err != nil {
			return fmt.Errorf("write plan output: %w", err)
		}
		for _, source := range output.Sources {
			if _, err := fmt.Fprintf(stdout, "  %s:\n", source.From); err != nil {
				return fmt.Errorf("write plan source: %w", err)
			}
			for _, section := range source.Sections {
				if _, err := fmt.Fprintf(stdout, "    %s\n", section.Path); err != nil {
					return fmt.Errorf("write plan section: %w", err)
				}
			}
		}
	}
	hasErrors := false
	for _, diagnostic := range plan.Diagnostics {
		if diagnostic.Severity == v2.SeverityError {
			hasErrors = true
		}
		if _, err := fmt.Fprintln(stderr, diagnostic.String()); err != nil {
			return fmt.Errorf("write plan diagnostic: %w", err)
		}
	}
	if hasErrors {
		return fmt.Errorf("plan contains errors")
	}
	return nil
}

func runV2Build(configPath string, force, dryRun bool, stdout, stderr io.Writer) error {
	config, err := v2.Load(configPath)
	if err != nil {
		return err
	}
	plan, err := v2.BuildLocal(config, force, dryRun)
	if plan != nil {
		for _, diagnostic := range plan.Diagnostics {
			if _, writeErr := fmt.Fprintln(stderr, diagnostic.String()); writeErr != nil {
				return fmt.Errorf("write build diagnostic: %w", writeErr)
			}
		}
	}
	if err != nil {
		return err
	}
	message := "Wrote generated outputs\n"
	if dryRun {
		message = "Dry run: no files written\n"
	}
	if _, err := fmt.Fprint(stdout, message); err != nil {
		return fmt.Errorf("write build result: %w", err)
	}
	return nil
}
