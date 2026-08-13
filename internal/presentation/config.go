// Package presentation loads user-level CLI display preferences.
package presentation

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Config contains presentation-only preferences. Project composition remains
// in agents.yaml and must never depend on this file.
type Config struct {
	Chars string
	Align bool
	Fit   string
	Width int
}

type configFile struct {
	Display displayConfig `yaml:"display"`
}

type displayConfig struct {
	Chars string `yaml:"chars"`
	Align *bool  `yaml:"align"`
	Fit   string `yaml:"fit"`
	Width *int   `yaml:"width"`
}

// Load reads $XDG_CONFIG_HOME/mogent/config.yaml, or the platform user config
// directory when XDG_CONFIG_HOME is unset. A missing file uses stable defaults.
func Load() (Config, error) {
	result := Config{Chars: "ascii", Align: true, Fit: "term"}
	configDir, err := os.UserConfigDir()
	if err != nil {
		return Config{}, fmt.Errorf("resolve user config directory: %w", err)
	}
	path := filepath.Join(configDir, "mogent", "config.yaml")
	contents, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return result, nil
	}
	if err != nil {
		return Config{}, fmt.Errorf("read presentation config %q: %w", path, err)
	}
	decoder := yaml.NewDecoder(bytes.NewReader(contents))
	decoder.KnownFields(true)
	var file configFile
	if err := decoder.Decode(&file); errors.Is(err, io.EOF) {
		return result, nil
	} else if err != nil {
		return Config{}, fmt.Errorf("parse presentation config %q: %w", path, err)
	}
	var extra any
	if err := decoder.Decode(&extra); err != nil && !errors.Is(err, io.EOF) {
		return Config{}, fmt.Errorf("parse presentation config %q: %w", path, err)
	} else if err == nil {
		return Config{}, fmt.Errorf("parse presentation config %q: expected one YAML document", path)
	}
	if file.Display.Chars != "" {
		result.Chars = file.Display.Chars
	}
	if result.Chars != "ascii" && result.Chars != "unicode" {
		return Config{}, fmt.Errorf("presentation config display.chars is %q; use ascii or unicode", result.Chars)
	}
	if file.Display.Align != nil {
		result.Align = *file.Display.Align
	}
	if file.Display.Fit != "" {
		result.Fit = file.Display.Fit
	}
	if result.Fit != "term" && result.Fit != "none" {
		return Config{}, fmt.Errorf("presentation config display.fit is %q; use term or none", result.Fit)
	}
	if file.Display.Width != nil {
		result.Width = *file.Display.Width
	}
	if result.Width < 0 {
		return Config{}, fmt.Errorf("presentation config display.width must be zero or greater")
	}
	return result, nil
}
