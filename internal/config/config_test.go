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
		TenantID:   "test-tenant",
		ClientID:   "test-client",
		Endpoint:   "https://test.openai.azure.com/",
		Deployment: "gpt-4o",
		APIVersion: "2024-02-15-preview",
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.TenantID != cfg.TenantID {
		t.Errorf("TenantID: got %q, want %q", loaded.TenantID, cfg.TenantID)
	}
	if loaded.ClientID != cfg.ClientID {
		t.Errorf("ClientID: got %q, want %q", loaded.ClientID, cfg.ClientID)
	}
	if loaded.Endpoint != cfg.Endpoint {
		t.Errorf("Endpoint: got %q, want %q", loaded.Endpoint, cfg.Endpoint)
	}
	if loaded.Deployment != cfg.Deployment {
		t.Errorf("Deployment: got %q, want %q", loaded.Deployment, cfg.Deployment)
	}
	if loaded.APIVersion != cfg.APIVersion {
		t.Errorf("APIVersion: got %q, want %q", loaded.APIVersion, cfg.APIVersion)
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
