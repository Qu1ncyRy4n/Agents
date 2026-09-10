package cli

import (
	"flag"
	"fmt"
	"io"

	"github.com/Qu1ncyRy4n/Agents/internal/mint"
)

func runMint(args []string, stdout, stderr io.Writer) error {
	flags := flag.NewFlagSet("mint", flag.ContinueOnError)
	flags.SetOutput(stderr)
	width := flags.Int("w", 1, "handle width: 1 (proquint-1) or 2 (proquint-2)")
	repoRoot := flags.String("r", ".", "repository root")
	dryRun := flags.Bool("n", false, "dry-run")
	seed := flags.Int64("s", 0, "entropy seed override")
	if err := parseFlags(flags, args); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return fmt.Errorf("mint accepts no positional arguments")
	}
	if *width != 1 && *width != 2 {
		return fmt.Errorf("-w must be 1 or 2, got %d", *width)
	}
	corpus, err := mint.ScanCorpus(*repoRoot)
	if err != nil {
		return err
	}
	handle, err := mint.Mint(*width, corpus, *seed, *dryRun)
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(stdout, handle)
	return err
}
