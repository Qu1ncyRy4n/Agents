package cli

import (
	"fmt"
	"io"
	"strings"
)

func unknownCommandError(command string) error {
	if suggestion := closestString(command, commandCandidates(), 3); suggestion != "" {
		return fmt.Errorf("unknown command %q; did you mean %q? Use help for usage", command, suggestion)
	}
	return fmt.Errorf("unknown command %q; use init, source, module, show, build, or help", command)
}

func commandCandidates() []string {
	definitions := commands()
	candidates := make([]string, 0, len(definitions)+1)
	for _, command := range definitions {
		candidates = append(candidates, command.name)
	}
	return append(candidates, "help")
}

func closestString(target string, candidates []string, maxDistance int) string {
	best, bestDistance := "", maxDistance+1
	for _, candidate := range candidates {
		distance := editDistance(strings.ToLower(target), strings.ToLower(candidate))
		if distance < bestDistance {
			best, bestDistance = candidate, distance
		}
	}
	if bestDistance > maxDistance {
		return ""
	}
	return best
}

func editDistance(first, second string) int {
	a, b := []rune(first), []rune(second)
	previous := make([]int, len(b)+1)
	for index := range previous {
		previous[index] = index
	}
	for index, left := range a {
		current := make([]int, len(b)+1)
		current[0] = index + 1
		for otherIndex, right := range b {
			cost := 0
			if left != right {
				cost = 1
			}
			current[otherIndex+1] = minimum(current[otherIndex]+1, previous[otherIndex+1]+1, previous[otherIndex]+cost)
		}
		previous = current
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

func runCompletion(args []string, stdout, stderr io.Writer) error {
	if len(args) != 1 || (args[0] != "bash" && args[0] != "zsh") {
		return fmt.Errorf("completion requires one shell: bash or zsh")
	}
	commands := strings.Join(commandCandidates(), " ")
	if args[0] == "bash" {
		_, err := fmt.Fprintf(stdout, "# Mogent bash completion\ncomplete -W %q mogent\n", commands)
		return err
	}
	_, err := fmt.Fprintf(stdout, "#compdef mogent\n_arguments '*:command:(%s)'\n", commands)
	return err
}
