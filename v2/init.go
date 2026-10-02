package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
)

// Init writes a minimal user-side configuration for one local library.
func Init(path, alias, local, root, output string, dryRun bool) (string, error) {
	if alias == "" || local == "" || root == "" || output == "" {
		return "", fmt.Errorf("config, source alias, local source, root, and output are required")
	}
	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("refusing to overwrite existing config %q", path)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("inspect config %q: %w", path, err)
	}
	content := fmt.Sprintf("mogent {\n  format = 2\n}\n\nsources {\n  source %s {\n    local = %s\n  }\n}\n\noutputs {\n  output \"agent-instructions\" {\n    path = %s\n    kind = \"markdown\"\n\n    source %s {\n      from = %s\n      select = {\n        all = true\n      }\n    }\n  }\n}\n", strconv.Quote(alias), strconv.Quote(local), strconv.Quote(output), strconv.Quote(alias), strconv.Quote(alias+":"+root))
	if !dryRun {
		if err := renderfs.WriteAtomically(filepath.Clean(path), []byte(content)); err != nil {
			return "", err
		}
	}
	return content, nil
}
