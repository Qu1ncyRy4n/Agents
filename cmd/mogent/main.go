package main

import (
	"fmt"
	"os"

	"github.com/Qu1ncyRy4n/Agents/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:], os.Stdout, os.Stderr); err != nil {
		fmt.Fprintln(os.Stderr, "mogent:", err)
		os.Exit(1)
	}
}
