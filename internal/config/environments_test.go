package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadEnvironments_FileNotExist(t *testing.T) {
	envs, err := loadEnvironmentsFromPath(filepath.Join(t.TempDir(), "missing.json"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(envs) != 0 {
		t.Errorf("expected empty slice, got %d", len(envs))
	}
}

func TestSaveAndLoadEnvironment(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "environments.json")

	env := Environment{Name: "myproject", Port: []string{"8080"}, Tag: "latest"}
	if err := saveEnvironmentToPath(path, env); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	envs, err := loadEnvironmentsFromPath(path)
	if err != nil {
		t.Fatalf("load failed: %v", err)
	}
	if len(envs) != 1 {
		t.Fatalf("expected 1, got %d", len(envs))
	}
	if envs[0].Name != "myproject" || envs[0].Tag != "latest" {
		t.Errorf("unexpected environment: %+v", envs[0])
	}
}

func TestUpsertEnvironment_UpdatesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "environments.json")

	saveEnvironmentToPath(path, Environment{Name: "myproject", Port: []string{"8080"}, Tag: "v1"})
	saveEnvironmentToPath(path, Environment{Name: "myproject", Port: []string{"9090"}, Tag: "v2"})

	envs, _ := loadEnvironmentsFromPath(path)
	if len(envs) != 1 {
		t.Fatalf("expected 1 after upsert, got %d", len(envs))
	}
	if envs[0].Tag != "v2" {
		t.Errorf("expected tag v2, got %q", envs[0].Tag)
	}
}

func TestLoadEnvironmentByName(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "environments.json")
	os.WriteFile(path, []byte(`[{"name":"foo","port":["3000"],"tag":""}]`), 0644)

	env, err := loadEnvironmentByNameFromPath(path, "foo")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env == nil || env.Name != "foo" {
		t.Errorf("expected foo, got %v", env)
	}

	missing, err := loadEnvironmentByNameFromPath(path, "bar")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if missing != nil {
		t.Errorf("expected nil for missing env, got %v", missing)
	}
}
