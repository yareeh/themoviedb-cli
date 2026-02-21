package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	// Use temp dir for config
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg := &Config{
		AccessToken:     "test-token",
		SessionID:       "test-session",
		AccountID:       12345,
		AccountObjectID: "abc123def",
	}

	if err := Save(cfg); err != nil {
		t.Fatalf("Save error: %v", err)
	}

	loaded, err := Load()
	if err != nil {
		t.Fatalf("Load error: %v", err)
	}

	if loaded.AccessToken != cfg.AccessToken {
		t.Errorf("AccessToken = %q, want %q", loaded.AccessToken, cfg.AccessToken)
	}
	if loaded.SessionID != cfg.SessionID {
		t.Errorf("SessionID = %q, want %q", loaded.SessionID, cfg.SessionID)
	}
	if loaded.AccountID != cfg.AccountID {
		t.Errorf("AccountID = %d, want %d", loaded.AccountID, cfg.AccountID)
	}
	if loaded.AccountObjectID != cfg.AccountObjectID {
		t.Errorf("AccountObjectID = %q, want %q", loaded.AccountObjectID, cfg.AccountObjectID)
	}
}

func TestLoadMissing(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load error for missing file: %v", err)
	}
	if cfg.AccessToken != "" {
		t.Errorf("expected empty config, got AccessToken=%q", cfg.AccessToken)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	tmpDir := t.TempDir()
	origHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", origHome)

	configDir := filepath.Join(tmpDir, ".config", "themoviedb-cli")
	os.MkdirAll(configDir, 0700)
	os.WriteFile(filepath.Join(configDir, "config.json"), []byte("not json"), 0600)

	_, err := Load()
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestDir(t *testing.T) {
	d := Dir()
	if d == "" {
		t.Fatal("Dir() returned empty string")
	}
}
