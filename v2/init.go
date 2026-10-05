package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
)

// Init writes a minimal user-side configuration for one library. location is
// a local directory path or an http(s) or file Git URL; a Git source starts
// unpinned and apply records its commit.
func Init(path, alias, location, root, output string, dryRun bool) (string, error) {
	if alias == "" || location == "" || root == "" || output == "" {
		return "", fmt.Errorf("config, source alias, source location, root, and output are required")
	}
	if _, err := os.Stat(path); err == nil {
		return "", fmt.Errorf("refusing to overwrite existing config %q", path)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("inspect config %q: %w", path, err)
	}
	sourceLine := "    local = " + strconv.Quote(location) + "\n"
	if isGitURL(location) {
		if err := validateGitURL(location); err != nil {
			return "", err
		}
		sourceLine = "    git = " + strconv.Quote(location) + "\n    ref = \"HEAD\"\n"
	}
	content := fmt.Sprintf("mogent {\n  format = 2\n}\n\nsources {\n  source %s {\n%s  }\n}\n\noutputs {\n  output \"agent-instructions\" {\n    path = %s\n    kind = \"markdown\"\n\n    source %s {\n      from = %s\n      select = {\n        all = true\n      }\n    }\n  }\n}\n", strconv.Quote(alias), sourceLine, strconv.Quote(output), strconv.Quote(alias), strconv.Quote(alias+":"+root))
	if !dryRun {
		if err := renderfs.WriteAtomically(filepath.Clean(path), []byte(content)); err != nil {
			return "", err
		}
	}
	return content, nil
}

func isGitURL(value string) bool {
	return strings.HasPrefix(value, "http://") || strings.HasPrefix(value, "https://") || strings.HasPrefix(value, "file://")
}
