package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	tmpDir := t.TempDir()

	origConfigDir := configDir
	origConfigPath := configPath
	configDir = tmpDir
	configPath = filepath.Join(tmpDir, "config.json")
	defer func() {
		configDir = origConfigDir
		configPath = origConfigPath
	}()

	cfg := &Config{
		Provider: "zen",
		Endpoint: "https://opencode.ai/zen/v1",
		Model:    "minimax-m2.5-free",
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Provider != cfg.Provider {
		t.Errorf("Provider: got %q, want %q", loaded.Provider, cfg.Provider)
	}
	if loaded.Endpoint != cfg.Endpoint {
		t.Errorf("Endpoint: got %q, want %q", loaded.Endpoint, cfg.Endpoint)
	}
	if loaded.Model != cfg.Model {
		t.Errorf("Model: got %q, want %q", loaded.Model, cfg.Model)
	}
}

func TestLoadNotFound(t *testing.T) {
	tmpDir := t.TempDir()

	origConfigDir := configDir
	origConfigPath := configPath
	configDir = tmpDir
	configPath = filepath.Join(tmpDir, "nonexistent.json")
	defer func() {
		configDir = origConfigDir
		configPath = origConfigPath
	}()

	_, err := Load()
	if err == nil {
		t.Error("Load should fail for nonexistent file")
	}
	if !IsConfigNotFound(err) {
		t.Errorf("IsConfigNotFound should return true, got false for: %v", err)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()

	origConfigDir := configDir
	origConfigPath := configPath
	configDir = tmpDir
	configPath = filepath.Join(tmpDir, "invalid.json")
	defer func() {
		configDir = origConfigDir
		configPath = origConfigPath
	}()

	os.WriteFile(configPath, []byte("not json"), 0644)

	_, err := Load()
	if err == nil {
		t.Error("Load should fail for invalid JSON")
	}
}
