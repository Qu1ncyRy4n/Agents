// Package registry remembers library locations for discovery. Assembly and
// builds never read it: choosing an entry materializes a project declaration.
package registry

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Qu1ncyRy4n/Agents/assembly"
	"github.com/Qu1ncyRy4n/Agents/renderfs"
	"gopkg.in/yaml.v3"
)

type Entry struct {
	Alias    string `yaml:"alias"`
	Name     string `yaml:"name,omitempty"`
	Location string `yaml:"location"`
}

type file struct {
	Format  int     `yaml:"format"`
	Sources []Entry `yaml:"sources"`
}

func DefaultPath() (string, error) {
	root := os.Getenv("XDG_CONFIG_HOME")
	if root == "" {
		var err error
		root, err = os.UserConfigDir()
		if err != nil {
			return "", err
		}
	}
	return filepath.Join(root, "mogent", "registry.yaml"), nil
}

// Load does not access registered checkouts or the network. Missing registries
// are empty, so relocating/deleting a checkout cannot break project builds.
func Load(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(bytes.NewReader(data))
	decoder.KnownFields(true)
	var value file
	if err := decoder.Decode(&value); err != nil {
		return nil, fmt.Errorf("read library registry: %w", err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("registry must contain one YAML document")
	}
	if value.Format != 1 {
		return nil, fmt.Errorf("unsupported library registry format %d", value.Format)
	}
	seen := make(map[string]bool)
	for _, entry := range value.Sources {
		if err := validate(entry); err != nil {
			return nil, err
		}
		if seen[entry.Alias] {
			return nil, fmt.Errorf("duplicate registered alias %q", entry.Alias)
		}
		seen[entry.Alias] = true
	}
	sort.Slice(value.Sources, func(i, j int) bool { return value.Sources[i].Alias < value.Sources[j].Alias })
	return value.Sources, nil
}

func validate(entry Entry) error {
	if entry.Alias == "" || strings.TrimSpace(entry.Alias) != entry.Alias || strings.ContainsAny(entry.Alias, ":/\\\n\r\t") {
		return fmt.Errorf("invalid registered alias %q", entry.Alias)
	}
	if strings.Contains(entry.Location, "://") {
		parsed, err := url.Parse(entry.Location)
		if err != nil || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
			return fmt.Errorf("registered Git URL must not contain credentials, query, or fragment")
		}
		if parsed.Scheme != "http" && parsed.Scheme != "https" && parsed.Scheme != "file" {
			return fmt.Errorf("registered Git URL must use http(s) or file")
		}
		if parsed.Scheme == "file" && parsed.Path == "" || parsed.Scheme != "file" && parsed.Host == "" {
			return fmt.Errorf("incomplete registered Git URL")
		}
	} else if !filepath.IsAbs(entry.Location) {
		return fmt.Errorf("registered local path must be absolute")
	}
	return nil
}

// Register records an explicit location, without cloning Git sources. A local
// checkout is validated and its library name retained for the picker.
func Register(path, alias, location string) (Entry, error) {
	entry := Entry{Alias: alias, Location: location}
	if !strings.Contains(location, "://") {
		if strings.HasPrefix(location, "~/") {
			home, err := os.UserHomeDir()
			if err != nil {
				return Entry{}, err
			}
			location = filepath.Join(home, location[2:])
		}
		absolute, err := filepath.Abs(location)
		if err != nil {
			return Entry{}, err
		}
		entry.Location = absolute
	}
	if err := validate(entry); err != nil {
		return Entry{}, err
	}
	entries, err := Load(path)
	if err != nil {
		return Entry{}, err
	}
	for _, existing := range entries {
		if existing.Alias == alias {
			if existing.Location == entry.Location {
				return existing, nil
			}
			return Entry{}, fmt.Errorf("registered alias %q already points to %s", alias, existing.Location)
		}
	}
	if !strings.Contains(entry.Location, "://") {
		library, err := assembly.LoadLibrary(entry.Location)
		if err != nil {
			return Entry{}, err
		}
		entry.Name = library.Name
	}
	entries = append(entries, entry)
	sort.Slice(entries, func(i, j int) bool { return entries[i].Alias < entries[j].Alias })
	data, err := yaml.Marshal(file{Format: 1, Sources: entries})
	if err != nil {
		return Entry{}, err
	}
	if err := renderfs.WriteAtomicallyMode(path, data, 0600); err != nil {
		return Entry{}, err
	}
	return entry, nil
}
