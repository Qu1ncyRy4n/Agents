package cli

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/syntax"
)

func hasFlag(args []string, wanted string) bool {
	for _, arg := range args {
		if arg == wanted || strings.HasPrefix(arg, wanted+"=") {
			return true
		}
	}
	return false
}

func runSequence(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("sequence", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "mogent.plan", "plan file")
	args = reorderArgs(args, map[string]bool{"-file": true, "--file": true})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("sequence accepts no positional arguments")
	}
	return sequenceFile(*path, stdout)
}

func sequenceFile(planPath string, stdout io.Writer) error {
	plan, err := assembly.LoadFile(planPath)
	if err != nil {
		return err
	}
	contents, result, err := assembly.Sequence(plan, planPath)
	if err != nil {
		return err
	}
	if result.HasErrors() {
		writeAssemblyDiagnostics(stderrDiscard{}, result)
		return fmt.Errorf("sequence contains errors")
	}
	buildPath := filepath.Join(filepath.Dir(planPath), "mogent.build")
	if existing, err := os.ReadFile(buildPath); err == nil {
		existingFile, parseErr := syntax.Parse(buildPath, existing)
		if parseErr != nil {
			return fmt.Errorf("refusing to replace unreadable %s: %w", buildPath, parseErr)
		}
		origin, found, originErr := buildOrigin(existingFile)
		if originErr != nil || !found || origin == "authored" {
			return fmt.Errorf("refusing to replace authored %s; select it explicitly with build --file", buildPath)
		}
		if err := assembly.VerifyBuildChecksum(existing); err != nil {
			return fmt.Errorf("refusing to replace %s: %w", buildPath, err)
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := os.WriteFile(buildPath, contents, 0644); err != nil {
		return err
	}
	_, err = fmt.Fprintf(stdout, "Sequenced %s -> %s\n", planPath, buildPath)
	return err
}

func buildOrigin(file *syntax.File) (string, bool, error) {
	if file.Kind != syntax.Build || len(file.Entries) == 0 {
		return "", false, nil
	}
	for _, child := range file.Entries[0].Body {
		if child.Kind == syntax.Field && child.Name == "origin" {
			value, err := child.Value.ResolveString(nil)
			return value, true, err
		}
	}
	return "", false, nil
}

// stderrDiscard keeps API diagnostics out of successful command stdout while
// callers that need rich display use show. CLI error messages still name blockers.
type stderrDiscard struct{}

func (stderrDiscard) Write([]byte) (int, error) { return 0, nil }

func writeAssemblyDiagnostics(output io.Writer, result *assembly.Result) {
	for _, diagnostic := range result.Diagnostics {
		_, _ = fmt.Fprintln(output, diagnostic.String())
	}
}

func runNewBuild(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("build", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "", "plan or build file")
	check := flags.Bool("check", false, "verify without writing workspace outputs")
	args = reorderArgs(args, map[string]bool{"-file": true, "--file": true, "-check": false, "--check": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 || *path == "" {
		return fmt.Errorf("new build requires --file mogent.plan or mogent.build")
	}
	file, err := assembly.LoadFile(*path)
	if err != nil {
		return err
	}
	buildPath := *path
	if file.Kind == syntax.Plan {
		if *check {
			result, _, err := assembly.ResolvePlan(file, *path)
			if err != nil {
				return err
			}
			if result.HasErrors() {
				writeAssemblyDiagnostics(stderr, result)
				return fmt.Errorf("build check contains errors")
			}
			_, err = result.Compose()
			return err
		}
		if err := sequenceFile(*path, stdout); err != nil {
			return err
		}
		buildPath = filepath.Join(filepath.Dir(*path), "mogent.build")
		file, err = assembly.LoadFile(buildPath)
		if err != nil {
			return err
		}
	}
	if file.Kind != syntax.Build {
		return fmt.Errorf("build file must contain a build header")
	}
	contents, err := os.ReadFile(buildPath)
	if err != nil {
		return err
	}
	if err := assembly.VerifyBuildChecksum(contents); err != nil {
		return err
	}
	result, _, err := assembly.ResolvePlan(file, buildPath)
	if err != nil {
		return err
	}
	if result.HasErrors() {
		writeAssemblyDiagnostics(stderr, result)
		return fmt.Errorf("build contains errors")
	}
	manifest, err := result.Compose()
	if err != nil {
		return err
	}
	if *check {
		if err := assembly.CheckWorkspace(filepath.Dir(buildPath), manifest); err != nil {
			return err
		}
		_, err := fmt.Fprintln(stdout, "Build check passed")
		return err
	}
	if err := assembly.Apply(filepath.Dir(buildPath), manifest); err != nil {
		return err
	}
	_, err = fmt.Fprintf(stdout, "Built %s\n", filepath.Dir(buildPath))
	return err
}

func runShow(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("show", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "mogent.plan", "plan or build file")
	sourceView := flags.Bool("source", false, "show source catalog")
	buildView := flags.Bool("build", false, "show ordered contributions")
	issues := flags.Bool("issues", false, "show only diagnostics")
	args = reorderArgs(args, map[string]bool{"-file": true, "--file": true, "-source": false, "--source": false, "-build": false, "--build": false, "-issues": false, "--issues": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() > 1 || *sourceView && *buildView {
		return fmt.Errorf("show accepts at most one node and only one of --source or --build")
	}
	file, err := assembly.LoadFile(*path)
	if err != nil {
		return err
	}
	result, sources, err := assembly.ResolvePlan(file, *path)
	if err != nil {
		return err
	}
	if *issues {
		writeAssemblyDiagnostics(stdout, result)
		if result.HasErrors() {
			return fmt.Errorf("show found errors")
		}
		return nil
	}
	if *sourceView {
		aliases := make([]string, 0, len(sources))
		for alias := range sources {
			aliases = append(aliases, alias)
		}
		sort.Strings(aliases)
		for _, alias := range aliases {
			source := sources[alias]
			if _, err := fmt.Fprintf(stdout, "%s: %s\n", alias, source.Lib.Name); err != nil {
				return err
			}
			paths := make([]string, 0, len(source.Lib.ByPath))
			for path := range source.Lib.ByPath {
				paths = append(paths, path)
			}
			sort.Strings(paths)
			for _, path := range paths {
				node := source.Lib.ByPath[path]
				if _, err := fmt.Fprintf(stdout, "  %s [%s]\n", path, node.Type); err != nil {
					return err
				}
			}
		}
		return nil
	}
	for _, contribution := range result.Contributions {
		name := contribution.Source
		if contribution.Node != nil {
			name += ":" + contribution.Node.Path
		} else {
			name = "inline"
		}
		if _, err := fmt.Fprintf(stdout, "%s -> %s\n", name, contribution.Into); err != nil {
			return err
		}
	}
	writeAssemblyDiagnostics(stdout, result)
	if result.HasErrors() {
		return fmt.Errorf("show found errors")
	}
	return nil
}

func runEdit(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 || args[0] != "include" {
		return fmt.Errorf("edit requires include; use edit include alias:path --localize")
	}
	flags := flag.NewFlagSet("edit include", flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("file", "mogent.plan", "plan file")
	localize := flags.Bool("localize", false, "create editable local source")
	dryRun := flags.Bool("dry-run", false, "preview without writing")
	args = reorderArgs(args[1:], map[string]bool{"-file": true, "--file": true, "-localize": false, "--localize": false, "-dry-run": false, "--dry-run": false})
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 1 || !*localize {
		return fmt.Errorf("edit include currently requires one reference and --localize")
	}
	plan, err := assembly.LoadFile(*path)
	if err != nil {
		return err
	}
	result, err := assembly.Localize(plan, *path, flags.Arg(0), *dryRun)
	if err != nil {
		return err
	}
	if *dryRun {
		_, err = fmt.Fprintf(stdout, "Dry run: local source %s\n\n%s", result.LocalPath, result.Plan)
		return err
	}
	_, err = fmt.Fprintf(stdout, "Localized %s as %s\n", result.SourceReference, result.LocalReference)
	return err
}
