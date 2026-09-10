package mint

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// diOwnerRE matches DI owner lines, not ordinary references to DI IDs. This
// distinction lets DRs and docs cite a DI without being treated as owning a
// second copy of that handle.
var diOwnerRE = regexp.MustCompile(`(?m)^ID:\s*DI-(` +
	`[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz]` +
	`(?:-[bdfghjklmnprstvz][aiou][bdfghjklmnprstvz][aiou][bdfghjklmnprstvz])?` +
	`)\s*$`)

// ScanCorpus returns the occupied handle set under repoRoot. It is deliberately
// path-agnostic: every file and directory name reserves any proquint-shaped
// substring it contains, and every exact DI owner line reserves its DI handle.
//
// Intent: Make mint a whole-working-tree occupied-set guard rather than a
// layout-specific artifact parser, so future paths cannot accidentally hide a
// handle-looking name from collision checks.
func ScanCorpus(repoRoot string) (map[string]string, error) {
	handles := make(map[string]string)
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return nil, err
	}
	if err := scanPathNames(root, handles); err != nil {
		return nil, err
	}
	if err := scanDIHandles(root, handles); err != nil {
		return nil, err
	}
	return handles, nil
}

// scanPathNames walks the whole working tree and records every proquint-shaped
// substring in each visited entry basename. The scan uses basenames so an
// ancestor's handle-looking name is recorded once rather than per descendant.
func scanPathNames(repoRoot string, handles map[string]string) error {
	return filepath.WalkDir(repoRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk %s: %w", path, walkErr)
		}
		if path == repoRoot {
			return nil
		}
		if entry.Name() == ".git" {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		relPath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return fmt.Errorf("relpath %s: %w", path, err)
		}
		rememberNameHandles(handles, entry.Name(), filepath.ToSlash(relPath))
		return nil
	})
}

// scanDIHandles records exact DI owner lines from readable regular files. It
// intentionally ignores arbitrary proquint-looking body text, which may cite a
// handle without owning it.
func scanDIHandles(repoRoot string, handles map[string]string) error {
	return filepath.WalkDir(repoRoot, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return fmt.Errorf("walk %s: %w", path, walkErr)
		}
		if path == repoRoot {
			return nil
		}
		if entry.Name() == ".git" {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		relPath, err := filepath.Rel(repoRoot, path)
		if err != nil {
			return fmt.Errorf("relpath %s: %w", path, err)
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read %s: %w", filepath.ToSlash(relPath), err)
		}
		for _, match := range diOwnerRE.FindAllStringSubmatch(string(data), -1) {
			rememberHandle(handles, match[1], filepath.ToSlash(relPath)+"#DI-"+match[1])
		}
		return nil
	})
}

// rememberNameHandles records all proquint-1 and proquint-2 substrings in one
// basename. The loops advance one byte at a time so overlaps are not missed.
func rememberNameHandles(handles map[string]string, name, relPath string) {
	for start := 0; start+5 <= len(name); start++ {
		candidate := name[start : start+5]
		if isProquint1String(candidate) {
			rememberHandle(handles, candidate, fmt.Sprintf("%s#name[%d:%d]", relPath, start, start+5))
		}
	}
	for start := 0; start+11 <= len(name); start++ {
		candidate := name[start : start+11]
		if isProquint2String(candidate) {
			rememberHandle(handles, candidate, fmt.Sprintf("%s#name[%d:%d]", relPath, start, start+11))
		}
	}
}

func isProquint1String(candidate string) bool {
	return len(candidate) == 5 &&
		strings.ContainsRune(proquintCons, rune(candidate[0])) &&
		strings.ContainsRune(proquintVows, rune(candidate[1])) &&
		strings.ContainsRune(proquintCons, rune(candidate[2])) &&
		strings.ContainsRune(proquintVows, rune(candidate[3])) &&
		strings.ContainsRune(proquintCons, rune(candidate[4]))
}

func isProquint2String(candidate string) bool {
	return len(candidate) == 11 && candidate[5] == '-' &&
		isProquint1String(candidate[:5]) && isProquint1String(candidate[6:])
}

// Repeated occurrences are expected under whole-name substring scanning, so
// the first owner is retained as useful diagnostic context.
func rememberHandle(handles map[string]string, handle, owner string) {
	if _, duplicate := handles[handle]; !duplicate {
		handles[handle] = owner
	}
}
