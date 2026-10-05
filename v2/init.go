package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/renderfs"
)

// InitTemplate is a named starting selection for a new mogent.hcl. Templates
// are ordinary configuration text, not hidden presets: the written file is the
// whole contract.
type InitTemplate struct {
	Name        string
	Description string
	outputs     func(alias, root, output string) string
}

// InitTemplates lists the available v2 starting points in a stable order.
func InitTemplates() []InitTemplate {
	return []InitTemplate{
		{
			Name:        "minimal",
			Description: "one Markdown output selecting every section the library offers",
			outputs: func(alias, root, output string) string {
				return fmt.Sprintf("outputs {\n  output \"agent-instructions\" {\n    path = %s\n    kind = \"markdown\"\n\n    source %s {\n      from = %s\n      select = {\n        all = true\n      }\n    }\n  }\n}\n", strconv.Quote(output), strconv.Quote(alias), strconv.Quote(alias+":"+root))
			},
		},
		{
			Name:        "qmr-core",
			Description: "QMR library core: intro and constraints as offered, workflow with accepted defaults, skills tree",
			outputs: func(alias, root, output string) string {
				return fmt.Sprintf("outputs {\n  output \"agent-instructions\" {\n    path = %s\n    kind = \"markdown\"\n\n    source %s {\n      from = %s\n      select = {\n        intro       = true\n        workflow    = { accept_defaults = true }\n        constraints = true\n      }\n    }\n  }\n\n  output \"agent-skills\" {\n    path = \".agents/skills/\"\n    kind = \"tree\"\n\n    source %s {\n      from   = %s\n      select = { all = true }\n    }\n  }\n}\n", strconv.Quote(output), strconv.Quote(alias), strconv.Quote(alias+":"+root), strconv.Quote(alias), strconv.Quote(alias+":skills"))
			},
		},
	}
}

// Init writes a user-side configuration for one library from the named
// template; an empty template name means minimal. location is a local
// directory path or an http(s) or file Git URL; a Git source starts unpinned
// and apply records its commit.
func Init(path, template, alias, location, root, output string, dryRun bool) (string, error) {
	if alias == "" || location == "" || root == "" || output == "" {
		return "", fmt.Errorf("config, source alias, source location, root, and output are required")
	}
	if template == "" {
		template = "minimal"
	}
	var chosen *InitTemplate
	for _, candidate := range InitTemplates() {
		if candidate.Name == template {
			copied := candidate
			chosen = &copied
		}
	}
	if chosen == nil {
		names := make([]string, 0)
		for _, candidate := range InitTemplates() {
			names = append(names, candidate.Name)
		}
		sort.Strings(names)
		return "", fmt.Errorf("unknown v2 template %q; use %s", template, strings.Join(names, ", "))
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
	content := fmt.Sprintf("mogent {\n  format = 2\n}\n\nsources {\n  source %s {\n%s  }\n}\n\n%s", strconv.Quote(alias), sourceLine, chosen.outputs(alias, root, output))
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
