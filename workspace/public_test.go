package workspace_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Qu1ncyRy4n/Agents/workspace"
)

// TestPublicWorkspaceWorkflow dogfoods the same exported operations used by
// the CLI from the perspective of a separate Go package.
func TestPublicWorkspaceWorkflow(t *testing.T) {
	temporary := t.TempDir()
	libraryPath := filepath.Join(temporary, "library")
	if err := os.Mkdir(libraryPath, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(libraryPath, "rules.md"), []byte("# Identity\nAct carefully.\n\n# Workflow\nKeep changes focused.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	manifestPath := filepath.Join(temporary, "agents.yaml")
	manifest := "sources:\n  shared: library\ndoc:\n  - Identity: shared:identity\n"
	if err := os.WriteFile(manifestPath, []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}

	session, err := workspace.New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	status, err := session.Status()
	if err != nil {
		t.Fatal(err)
	}
	if status.ManifestPath != manifestPath || status.Output != workspace.StatusMissing {
		t.Fatalf("status = %#v", status)
	}
	nodes, err := session.SourceNodes(workspace.SourceListOptions{Search: "workflow"})
	if err != nil {
		t.Fatal(err)
	}
	if len(nodes) != 1 || nodes[0].Reference != "shared:workflow" {
		t.Fatalf("source nodes = %#v", nodes)
	}
	preview, err := session.AddSource(workspace.AddOptions{
		Reference: "shared:workflow",
		Heading:   "Workflow",
		Append:    true,
	}, true)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(preview.Preview, "# Workflow") {
		t.Fatalf("preview = %#v", preview)
	}
	contents, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != manifest {
		t.Fatalf("dry run changed manifest:\n%s", contents)
	}
}
