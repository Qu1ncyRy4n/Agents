package workspace

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDriftReportsAndRejectsDirectEditsExplicitly(t *testing.T) {
	_, manifestPath, _ := writeLocalizeFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := session.SaveAndBuild(); err != nil {
		t.Fatal(err)
	}
	outputPath := session.OutputPath()
	if err := os.WriteFile(outputPath, []byte("# Direct edit\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	report, err := session.Drift()
	if err != nil {
		t.Fatal(err)
	}
	if !report.DirectEdits || report.Status != StatusDirectEdits {
		t.Fatalf("report = %#v", report)
	}
	if err := session.RejectDrift(false); err == nil || !strings.Contains(err.Error(), "requires --force") {
		t.Fatalf("reject error = %v", err)
	}
	if err := session.RejectDrift(true); err != nil {
		t.Fatal(err)
	}
	contents, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	if string(contents) != session.Output {
		t.Fatalf("rejected output:\n%s", contents)
	}
}

func TestImportDriftLocalizesOnlyEditedManifestSection(t *testing.T) {
	temporary, manifestPath, _ := writeLocalizeFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := session.SaveAndBuild(); err != nil {
		t.Fatal(err)
	}
	edited := strings.Replace(session.Output, "Use fixtures.", "Use careful fixtures.", 1)
	if err := os.WriteFile(session.OutputPath(), []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}
	result, err := session.ImportDrift("Instructions/Tests", "")
	if err != nil {
		t.Fatal(err)
	}
	if result.LocalReference != "local:instructions/testing" {
		t.Fatalf("result = %#v", result)
	}
	local, err := os.ReadFile(filepath.Join(temporary, ".mogent", "library", "instructions", "testing.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(local), "Use careful fixtures.") || !strings.Contains(string(local), "## Errors") {
		t.Fatalf("local override:\n%s", local)
	}
	rebuilt, err := os.ReadFile(session.OutputPath())
	if err != nil {
		t.Fatal(err)
	}
	if string(rebuilt) != edited {
		t.Fatalf("rebuilt output lost imported edit:\n%s", rebuilt)
	}
}

func TestImportDriftRefusesEditsOutsideSelectedSection(t *testing.T) {
	_, manifestPath, _ := writeLocalizeFixture(t)
	session, err := New(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	if err := session.SaveAndBuild(); err != nil {
		t.Fatal(err)
	}
	edited := strings.Replace(session.Output, "# Instructions\n\n", "# Instructions\n\nOutside edit.\n\n", 1)
	edited = strings.Replace(edited, "Use fixtures.", "Use careful fixtures.", 1)
	if err := os.WriteFile(session.OutputPath(), []byte(edited), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := session.ImportDrift("Instructions/Tests", ""); err == nil || !strings.Contains(err.Error(), "outside manifest section") {
		t.Fatalf("import error = %v", err)
	}
	if _, err := os.Stat(filepath.Join(filepath.Dir(manifestPath), ".mogent", "library")); !os.IsNotExist(err) {
		t.Fatalf("refused import created local library: %v", err)
	}
}
