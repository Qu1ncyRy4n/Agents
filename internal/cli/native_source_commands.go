package cli

import (
	"fmt"
	"io"
)

func runSource(args []string, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return fmt.Errorf("source requires a subcommand; use source add, list (ls), register, show, or update local")
	}
	switch args[0] {
	case "register":
		return runSourceRegister(args[1:], stdout, stderr)
	case "add":
		return runNewSourceAdd(args[1:], stdout, stderr)
	case "list", "ls":
		return runNewSourceList(args[1:], stdout, stderr)
	case "show":
		return runNativeSourceShow(args[1:], stdout, stderr)
	case "update":
		if len(args) > 1 && args[1] == "local" {
			return runLocalSourceUpdate(args[2:], stdout, stderr)
		}
		return fmt.Errorf("source update supports only local inputs; use source update local")
	default:
		if suggestion := closestString(args[0], []string{"add", "list", "ls", "register", "show", "update"}, 2); suggestion != "" {
			return fmt.Errorf("unknown source subcommand %q; did you mean %q?", args[0], suggestion)
		}
		return fmt.Errorf("unknown source subcommand %q; use source add, list (ls), register, show, or update", args[0])
	}
}

func runNativeSourceShow(args []string, stdout, stderr io.Writer) error {
	return runShow(append([]string{"--source"}, args...), stdout, stderr)
}
