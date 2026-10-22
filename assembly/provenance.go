package assembly

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"

	"github.com/Qu1ncyRy4n/Agents/state"
)

// SourceBaselines returns deterministic fingerprints for every connected
// library. The root fingerprint covers library-level catalog metadata and its
// root nodes; each node fingerprint covers its own semantic metadata, payload,
// and ordered child fingerprints. Paths are deliberately not derived from
// hashes: they remain the stable module identities used by consumers.
func SourceBaselines(sources map[string]Source) (map[string]state.SourceBaseline, error) {
	baselines := make(map[string]state.SourceBaseline, len(sources))
	aliases := make([]string, 0, len(sources))
	for alias := range sources {
		aliases = append(aliases, alias)
	}
	sort.Strings(aliases)
	for _, alias := range aliases {
		baseline, err := libraryBaseline(sources[alias].Lib)
		if err != nil {
			return nil, fmt.Errorf("fingerprint source %q: %w", alias, err)
		}
		baselines[alias] = baseline
	}
	return baselines, nil
}

func libraryBaseline(library *Library) (state.SourceBaseline, error) {
	if library == nil {
		return state.SourceBaseline{}, fmt.Errorf("nil library")
	}
	nodes := make(map[string]string, len(library.ByPath))
	roots := make([]string, 0, len(library.Nodes))
	for _, node := range library.Nodes {
		digest, err := fingerprintNode(node, nodes)
		if err != nil {
			return state.SourceBaseline{}, err
		}
		roots = append(roots, digest)
	}
	catalog := ""
	if contents, err := os.ReadFile(library.Path); err == nil {
		catalog = hashBytes(contents)
	} else if !os.IsNotExist(err) {
		return state.SourceBaseline{}, err
	}
	type rootRecord struct {
		Name       string   `json:"name"`
		Catalog    string   `json:"catalog_sha256,omitempty"`
		RootHashes []string `json:"root_hashes"`
	}
	encoded, err := json.Marshal(rootRecord{Name: library.Name, Catalog: catalog, RootHashes: roots})
	if err != nil {
		return state.SourceBaseline{}, err
	}
	return state.SourceBaseline{SHA256: hashBytes(encoded), Nodes: nodes}, nil
}

func fingerprintNode(node *Node, nodes map[string]string) (string, error) {
	if digest, found := nodes[node.Path]; found {
		return digest, nil
	}
	childHashes := make([]string, 0, len(node.Children))
	for _, child := range node.Children {
		digest, err := fingerprintNode(child, nodes)
		if err != nil {
			return "", err
		}
		childHashes = append(childHashes, digest)
	}
	payload, err := NodeHash(node)
	if err != nil {
		return "", err
	}
	type nodeRecord struct {
		Name           string   `json:"name"`
		Path           string   `json:"path"`
		Type           NodeType `json:"type"`
		Source         string   `json:"source,omitempty"`
		Title          string   `json:"title,omitempty"`
		Tags           []string `json:"tags,omitempty"`
		Role           string   `json:"role,omitempty"`
		Content        string   `json:"content,omitempty"`
		Curate         Curate   `json:"curate,omitempty"`
		Default        *bool    `json:"default,omitempty"`
		ReasonRequired bool     `json:"reason_required,omitempty"`
		Requires       []string `json:"requires,omitempty"`
		ConflictsWith  []string `json:"conflicts_with,omitempty"`
		ExclusiveGroup string   `json:"exclusive_group,omitempty"`
		Into           string   `json:"into,omitempty"`
		Variables      []string `json:"variables,omitempty"`
		Heading        []string `json:"heading,omitempty"`
		Payload        string   `json:"payload_sha256"`
		Children       []string `json:"children"`
	}
	encoded, err := json.Marshal(nodeRecord{
		Name: node.Name, Path: node.Path, Type: node.Type, Source: node.Source,
		Title: node.Title, Tags: node.Tags, Role: node.Role, Content: node.Content,
		Curate: node.Curate, Default: node.Default, ReasonRequired: node.ReasonRequired,
		Requires: node.Requires, ConflictsWith: node.ConflictsWith,
		ExclusiveGroup: node.ExclusiveGroup, Into: node.Into, Variables: node.Variables,
		Heading: node.Heading, Payload: payload, Children: childHashes,
	})
	if err != nil {
		return "", err
	}
	digest := hashBytes(encoded)
	nodes[node.Path] = digest
	return digest, nil
}
