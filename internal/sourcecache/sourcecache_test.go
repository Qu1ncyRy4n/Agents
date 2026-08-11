package sourcecache

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveUsesVerifiedCacheWithoutNetwork(t *testing.T) {
	temporary := t.TempDir()
	manifestPath := filepath.Join(temporary, "agents.yaml")
	commit := strings.Repeat("a", 40)
	cache := cachePath(manifestPath, "shared", commit)
	if err := os.MkdirAll(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cache, "rules.md"), []byte("# Rules\nWork.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	hash, err := HashMarkdown(cache)
	if err != nil {
		t.Fatal(err)
	}
	lock := &Lock{Version: lockVersion, Sources: map[string]LockEntry{
		"shared": {URL: "https://example.com/library.git", Commit: commit, ContentSHA256: hash},
	}}
	if err := writeLock(lockPath(manifestPath), lock); err != nil {
		t.Fatal(err)
	}
	resolved, err := Resolve(manifestPath, "shared", "https://example.com/library.git")
	if err != nil {
		t.Fatal(err)
	}
	if resolved != cache {
		t.Fatalf("resolved = %q, want %q", resolved, cache)
	}
}

func TestResolveRejectsMissingAndChangedCaches(t *testing.T) {
	temporary := t.TempDir()
	manifestPath := filepath.Join(temporary, "agents.yaml")
	commit := strings.Repeat("b", 40)
	lock := &Lock{Version: lockVersion, Sources: map[string]LockEntry{
		"shared": {URL: "https://example.com/library.git", Commit: commit, ContentSHA256: strings.Repeat("c", 64)},
	}}
	if err := writeLock(lockPath(manifestPath), lock); err != nil {
		t.Fatal(err)
	}
	if _, err := Resolve(manifestPath, "shared", "https://example.com/library.git"); err == nil || !strings.Contains(err.Error(), "cache") {
		t.Fatalf("missing cache error = %v", err)
	}
	cache := cachePath(manifestPath, "shared", commit)
	if err := os.MkdirAll(cache, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(cache, "rules.md"), []byte("# Changed\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Resolve(manifestPath, "shared", "https://example.com/library.git"); err == nil || !strings.Contains(err.Error(), "failed content verification") {
		t.Fatalf("changed cache error = %v", err)
	}
}

func TestValidateRemoteRejectsCredentialsAndUnsafeAliases(t *testing.T) {
	for _, test := range []struct {
		alias string
		url   string
	}{
		{alias: "../shared", url: "https://example.com/library.git"},
		{alias: "shared", url: "https://user:secret@example.com/library.git"},
		{alias: "shared", url: "file:///tmp/library"},
	} {
		if err := validateRemote(test.alias, test.url); err == nil {
			t.Fatalf("validateRemote(%q, %q) succeeded", test.alias, test.url)
		}
	}
}

func TestMarkdownChangesReportsAddedModifiedAndRemoved(t *testing.T) {
	oldRoot := t.TempDir()
	newRoot := t.TempDir()
	for path, content := range map[string]string{"same.md": "same", "changed.md": "old", "removed.md": "gone"} {
		if err := os.WriteFile(filepath.Join(oldRoot, path), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	for path, content := range map[string]string{"same.md": "same", "changed.md": "new", "added.md": "here"} {
		if err := os.WriteFile(filepath.Join(newRoot, path), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	changes, err := markdownChanges(oldRoot, newRoot)
	if err != nil {
		t.Fatal(err)
	}
	got := make(map[string]string)
	for _, change := range changes {
		got[change.Path] = change.Status
	}
	for path, status := range map[string]string{"added.md": "added", "changed.md": "modified", "removed.md": "removed"} {
		if got[path] != status {
			t.Fatalf("changes = %#v", changes)
		}
	}
}

func TestUpdatePreviewWritesNothingAndAcceptInstallsCandidate(t *testing.T) {
	temporary := t.TempDir()
	manifestPath := filepath.Join(temporary, "agents.yaml")
	url := "https://example.com/library.git"
	oldCommit := strings.Repeat("1", 40)
	oldCache := cachePath(manifestPath, "shared", oldCommit)
	if err := os.MkdirAll(oldCache, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(oldCache, "rules.md"), []byte("# Rules\nOld.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	oldHash, err := HashMarkdown(oldCache)
	if err != nil {
		t.Fatal(err)
	}
	if err := writeLock(lockPath(manifestPath), &Lock{Version: lockVersion, Sources: map[string]LockEntry{
		"shared": {URL: url, Commit: oldCommit, ContentSHA256: oldHash},
	}}); err != nil {
		t.Fatal(err)
	}
	candidateRoot := t.TempDir()
	if err := os.WriteFile(filepath.Join(candidateRoot, "rules.md"), []byte("# Rules\nNew.\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	newHash, err := HashMarkdown(candidateRoot)
	if err != nil {
		t.Fatal(err)
	}
	newCommit := strings.Repeat("2", 40)
	originalFetcher := fetchCandidate
	fetchCandidate = func(sourceURL, ref string) (candidate, func(), error) {
		return candidate{LockEntry: LockEntry{URL: sourceURL, Commit: newCommit, ContentSHA256: newHash}, Path: candidateRoot}, func() {}, nil
	}
	t.Cleanup(func() { fetchCandidate = originalFetcher })

	preview, err := Update(manifestPath, "shared", url, "main", false)
	if err != nil {
		t.Fatal(err)
	}
	if preview.Wrote || len(preview.Changes) != 1 || preview.Changes[0].Status != "modified" {
		t.Fatalf("preview = %#v", preview)
	}
	if !strings.Contains(preview.Changes[0].Old, "Old.") || !strings.Contains(preview.Changes[0].New, "New.") {
		t.Fatalf("preview content = %#v", preview.Changes[0])
	}
	currentLock, err := load(lockPath(manifestPath))
	if err != nil {
		t.Fatal(err)
	}
	if currentLock.Sources["shared"].Commit != oldCommit {
		t.Fatal("preview changed lock")
	}
	if _, err := os.Stat(cachePath(manifestPath, "shared", newCommit)); !os.IsNotExist(err) {
		t.Fatalf("preview installed cache: %v", err)
	}

	if _, err := Update(manifestPath, "shared", url, "main", true); err == nil || !strings.Contains(err.Error(), "full commit") {
		t.Fatalf("moving-ref acceptance error = %v", err)
	}
	accepted, err := Update(manifestPath, "shared", url, newCommit, true)
	if err != nil {
		t.Fatal(err)
	}
	if !accepted.Wrote {
		t.Fatal("accepted update did not report write")
	}
	resolved, err := Resolve(manifestPath, "shared", url)
	if err != nil {
		t.Fatal(err)
	}
	if resolved != cachePath(manifestPath, "shared", newCommit) {
		t.Fatalf("resolved = %q", resolved)
	}
}
