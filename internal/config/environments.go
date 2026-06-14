package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Environment holds per-environment port and image tag configuration.
type Environment struct {
	Name string   `json:"name"`
	Port []string `json:"port"`
	Tag  string   `json:"tag"`
}

// LoadEnvironments reads all environments from ~/.config/dev/environments.json.
func LoadEnvironments() ([]Environment, error) {
	path, err := environmentsPath()
	if err != nil {
		return nil, err
	}
	return loadEnvironmentsFromPath(path)
}

// SaveEnvironment upserts an environment record into ~/.config/dev/environments.json.
func SaveEnvironment(env Environment) error {
	path, err := environmentsPath()
	if err != nil {
		return err
	}
	return saveEnvironmentToPath(path, env)
}

// LoadEnvironmentByName returns the environment config for a given name, or nil if not found.
func LoadEnvironmentByName(name string) (*Environment, error) {
	path, err := environmentsPath()
	if err != nil {
		return nil, err
	}
	return loadEnvironmentByNameFromPath(path, name)
}

func environmentsPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("get home dir: %w", err)
	}
	return filepath.Join(home, ".config", "dev", "environments.json"), nil
}

func loadEnvironmentsFromPath(path string) ([]Environment, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return []Environment{}, nil
	}
	if err != nil {
		return nil, err
	}
	var envs []Environment
	if err := json.Unmarshal(data, &envs); err != nil {
		return nil, err
	}
	return envs, nil
}

func saveEnvironmentToPath(path string, env Environment) error {
	envs, err := loadEnvironmentsFromPath(path)
	if err != nil {
		return err
	}
	updated := false
	for i, e := range envs {
		if e.Name == env.Name {
			envs[i] = env
			updated = true
			break
		}
	}
	if !updated {
		envs = append(envs, env)
	}
	return writeJSON(path, envs)
}

func loadEnvironmentByNameFromPath(path, name string) (*Environment, error) {
	envs, err := loadEnvironmentsFromPath(path)
	if err != nil {
		return nil, err
	}
	for _, e := range envs {
		if e.Name == name {
			e := e
			return &e, nil
		}
	}
	return nil, nil
}
