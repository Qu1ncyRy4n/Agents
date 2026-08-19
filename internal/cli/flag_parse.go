package cli

import (
	"flag"
	"fmt"
	"strings"
)

func reorderArgs(args []string, valueFlags map[string]bool) []string {
	var flagArgs []string
	var positional []string
	for index := 0; index < len(args); index++ {
		arg := args[index]
		if !strings.HasPrefix(arg, "-") || arg == "-" {
			positional = append(positional, arg)
			continue
		}
		flagArgs = append(flagArgs, arg)
		name := arg
		if before, _, found := strings.Cut(arg, "="); found {
			name = before
		}
		if valueFlags[name] && !strings.Contains(arg, "=") && index+1 < len(args) {
			index++
			flagArgs = append(flagArgs, args[index])
		}
	}
	return append(flagArgs, positional...)
}

func parseFlags(flags *flag.FlagSet, args []string) error {
	if err := suggestUnknownFlags(flags, args); err != nil {
		return err
	}
	return flags.Parse(args)
}

func suggestUnknownFlags(flags *flag.FlagSet, args []string) error {
	var known []string
	flags.VisitAll(func(flag *flag.Flag) {
		known = append(known, flag.Name)
	})
	for _, arg := range args {
		if !strings.HasPrefix(arg, "-") || arg == "-" || arg == "--" {
			continue
		}
		name := strings.TrimLeft(arg, "-")
		if before, _, found := strings.Cut(name, "="); found {
			name = before
		}
		if name == "h" || name == "help" || flags.Lookup(name) != nil {
			continue
		}
		if suggestion := closestString(name, known, 4); suggestion != "" {
			return fmt.Errorf("unknown flag %q; did you mean --%s?", arg, suggestion)
		}
		return fmt.Errorf("unknown flag %q", arg)
	}
	return nil
}
