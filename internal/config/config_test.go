package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_FileNotExist_ReturnsEmpty(t *testing.T) {
	dir := t.TempDir()
	cfg, err := loadConfigFromPath(filepath.Join(dir, "config.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Image != "" {
		t.Errorf("expected empty image, got %q", cfg.Image)
	}
}

func TestLoadConfig_ReadsFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	data := Config{Image: "registry.example.com/dev", SSHDirectory: "/home/user/.ssh"}
	b, _ := json.Marshal(data)
	os.WriteFile(path, b, 0644)

	cfg, err := loadConfigFromPath(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Image != "registry.example.com/dev" {
		t.Errorf("got %q, want %q", cfg.Image, "registry.example.com/dev")
	}
	if cfg.SSHDirectory != "/home/user/.ssh" {
		t.Errorf("got %q, want %q", cfg.SSHDirectory, "/home/user/.ssh")
	}
}

func TestApplyEnvOverrides_OverridesImage(t *testing.T) {
	t.Setenv("DOCKER_DEV_ENV", "registry.override.com/dev")
	cfg := &Config{Image: "registry.original.com/dev"}
	applyEnvOverrides(cfg)
	if cfg.Image != "registry.override.com/dev" {
		t.Errorf("got %q, want %q", cfg.Image, "registry.override.com/dev")
	}
}

func TestApplyEnvOverrides_RemoteDevEnv_Split(t *testing.T) {
	t.Setenv("REMOTE_DEV_ENV", "192.168.0.1,192.168.0.2")
	cfg := &Config{}
	applyEnvOverrides(cfg)
	if len(cfg.RemoteDevEnv) != 2 {
		t.Fatalf("expected 2 remotes, got %d", len(cfg.RemoteDevEnv))
	}
	if cfg.RemoteDevEnv[0] != "192.168.0.1" {
		t.Errorf("got %q", cfg.RemoteDevEnv[0])
	}
}

func TestApplyEnvOverrides_EmptyEnvDoesNotOverride(t *testing.T) {
	os.Unsetenv("DOCKER_DEV_ENV")
	cfg := &Config{Image: "registry.original.com/dev"}
	applyEnvOverrides(cfg)
	if cfg.Image != "registry.original.com/dev" {
		t.Errorf("env override should not happen when env var is empty, got %q", cfg.Image)
	}
}
