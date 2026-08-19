package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/internal/navigator"
)

func runTUI(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("tui", flag.ContinueOnError)
	flags.SetOutput(stderr)
	manifestFile := flags.String("manifest", "agents.yaml", "path to manifest")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("tui accepts no positional arguments")
	}
	return navigator.Run(*manifestFile, stdout)
}

func unknownCommandError(command string) error {
	commands := commandCandidates()
	if suggestion := closestString(command, commands, 3); suggestion != "" {
		return fmt.Errorf("unknown command %q; did you mean %q? Use help for usage", command, suggestion)
	}
	return fmt.Errorf("unknown command %q; use init, build, status, drift, coverage, source, add, localize, complete, completion, or tui", command)
}

func decorateSourceReferenceError(reference string, err error) error {
	message := err.Error()
	if strings.Contains(message, "invalid source reference") && strings.HasSuffix(strings.TrimSpace(reference), "/") {
		trimmed := strings.TrimSuffix(strings.TrimSpace(reference), "/")
		return fmt.Errorf("%w; remove the trailing slash, or use `mogent source list --search %s` to find descendant heading paths", err, shellToken(trimmed))
	}
	return err
}

func closestString(target string, candidates []string, maxDistance int) string {
	best := ""
	bestDistance := maxDistance + 1
	for _, candidate := range candidates {
		distance := editDistance(strings.ToLower(target), strings.ToLower(candidate))
		if distance < bestDistance {
			best = candidate
			bestDistance = distance
		}
	}
	if bestDistance > maxDistance {
		return ""
	}
	return best
}

func editDistance(first string, second string) int {
	if first == second {
		return 0
	}
	a := []rune(first)
	b := []rune(second)
	previous := make([]int, len(b)+1)
	current := make([]int, len(b)+1)
	for index := range previous {
		previous[index] = index
	}
	for i, ar := range a {
		current[0] = i + 1
		for j, br := range b {
			cost := 1
			if ar == br {
				cost = 0
			}
			current[j+1] = minimum(current[j]+1, previous[j+1]+1, previous[j]+cost)
		}
		previous, current = current, previous
	}
	return previous[len(b)]
}

func minimum(values ...int) int {
	best := values[0]
	for _, value := range values[1:] {
		if value < best {
			best = value
		}
	}
	return best
}

func shellToken(value string) string {
	if value == "" {
		return "''"
	}
	if strings.ContainsAny(value, " \t\n'\"`$\\") {
		return "'" + strings.ReplaceAll(value, "'", "'\"'\"'") + "'"
	}
	return value
}
