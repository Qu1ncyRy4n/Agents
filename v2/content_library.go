package v2

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/content"
	"github.com/hashicorp/hcl/v2/gohcl"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/hashicorp/hcl/v2/hclsyntax"
)

type namedContent struct {
	Name      string
	Directory string `hcl:"directory"`
}

// Decode labeled content declarations separately, preserving the existing
// unlabeled authored-content schema and diagnostics in the same sidecar.
func decodeLibrary(file string) (libraryFile, []namedContent, error) {
	parser := hclparse.NewParser()
	parsed, diagnostics := parser.ParseHCLFile(file)
	if diagnostics.HasErrors() {
		return libraryFile{}, nil, diagnostics
	}
	body := parsed.Body.(*hclsyntax.Body)
	filtered := *body
	filtered.Blocks = nil
	var named []namedContent
	for _, block := range body.Blocks {
		if block.Type != "content" || len(block.Labels) == 0 {
			filtered.Blocks = append(filtered.Blocks, block)
			continue
		}
		if len(block.Labels) != 1 {
			return libraryFile{}, nil, fmt.Errorf("content requires one root label")
		}
		declaration := namedContent{Name: block.Labels[0]}
		if declaration.Name == "" || strings.ContainsAny(declaration.Name, "/:\\") || declaration.Name == "." || declaration.Name == ".." {
			return libraryFile{}, nil, fmt.Errorf("invalid content root name %q", declaration.Name)
		}
		// Name is not a decoded field; only the directory attribute is required.
		var attributes struct {
			Directory string `hcl:"directory"`
		}
		if diagnostics := gohcl.DecodeBody(block.Body, nil, &attributes); diagnostics.HasErrors() {
			return libraryFile{}, nil, diagnostics
		}
		declaration.Directory = attributes.Directory
		named = append(named, declaration)
	}
	var raw libraryFile
	if diagnostics := gohcl.DecodeBody(&filtered, nil, &raw); diagnostics.HasErrors() {
		return libraryFile{}, nil, diagnostics
	}
	return raw, named, nil
}

func loadContentDirectory(root, relative string) (*content.Inventory, error) {
	if err := content.ValidatePath(relative, true); err != nil {
		return nil, err
	}
	canonical, err := filepath.EvalSymlinks(root)
	if err != nil {
		return nil, err
	}
	current := canonical
	if relative != "." {
		for _, component := range strings.Split(relative, "/") {
			current = filepath.Join(current, component)
			info, err := os.Lstat(current)
			if err != nil {
				return nil, err
			}
			if info.Mode()&os.ModeSymlink != 0 {
				return nil, fmt.Errorf("content directory traverses symlink %q", current)
			}
		}
	}
	info, err := os.Stat(current)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("content directory %q is not a directory", relative)
	}
	return content.Discover(current)
}
