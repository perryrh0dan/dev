package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Config holds host-side settings loaded from ~/.config/dev/config.json.
// Environment variables override individual fields after loading.
type Config struct {
	Image                  string   `json:"image"`
	Memory                 string   `json:"memory"`
	SSHDirectory           string   `json:"ssh_directory"`
	GPGDirectory           string   `json:"gpg_directory"`
	SharedDirectory        string   `json:"shared_directory"`
	KubeDirectory          string   `json:"kube_directory"`
	TimewarriorDirectory   string   `json:"timewarrior_directory"`
	OpenCodeDirectory      string   `json:"opencode_directory"`
	EnvFile                string   `json:"env_file"`
	DictFile               string   `json:"dict_file"`
	NgrokDirectory         string   `json:"ngrok_directory"`
	NpmFile                string   `json:"npm_file"`
	GitHubCopilotDirectory string   `json:"github_copilot_directory"`
	RemoteDevEnv           []string `json:"remote_dev_env"`
}

// Load reads ~/.config/dev/config.json and applies environment variable overrides.
func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	path := filepath.Join(home, ".config", "dev", "config.json")
	cfg, err := loadConfigFromPath(path)
	if err != nil {
		return nil, err
	}
	applyEnvOverrides(cfg)
	return cfg, nil
}

// loadConfigFromPath reads config from a specific path. Returns empty config if file is absent.
func loadConfigFromPath(path string) (*Config, error) {
	cfg := &Config{}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return cfg, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

// applyEnvOverrides overrides config fields with environment variable values when set.
func applyEnvOverrides(cfg *Config) {
	if v := os.Getenv("DOCKER_DEV_ENV"); v != "" {
		cfg.Image = v
	}
	if v := os.Getenv("MEMORY"); v != "" {
		cfg.Memory = v
	}
	if v := os.Getenv("SSH_DIRECTORY"); v != "" {
		cfg.SSHDirectory = v
	}
	if v := os.Getenv("GPG_DIRECTORY"); v != "" {
		cfg.GPGDirectory = v
	}
	if v := os.Getenv("SHARED_DIRECTORY"); v != "" {
		cfg.SharedDirectory = v
	}
	if v := os.Getenv("KUBE_DIRECTORY"); v != "" {
		cfg.KubeDirectory = v
	}
	if v := os.Getenv("TIMEWARRIOR_DIRECTORY"); v != "" {
		cfg.TimewarriorDirectory = v
	}
	if v := os.Getenv("OPENCODE_DIRECTORY"); v != "" {
		cfg.OpenCodeDirectory = v
	}
	if v := os.Getenv("ENV_FILE"); v != "" {
		cfg.EnvFile = v
	}
	if v := os.Getenv("DICT_FILE"); v != "" {
		cfg.DictFile = v
	}
	if v := os.Getenv("NGROK_DIRECTORY"); v != "" {
		cfg.NgrokDirectory = v
	}
	if v := os.Getenv("NPM_FILE"); v != "" {
		cfg.NpmFile = v
	}
	if v := os.Getenv("GITHUB_COPILOT_DIRECTORY"); v != "" {
		cfg.GitHubCopilotDirectory = v
	}
	if v := os.Getenv("REMOTE_DEV_ENV"); v != "" {
		parts := strings.Split(v, ",")
		var hosts []string
		for _, p := range parts {
			if h := strings.TrimSpace(p); h != "" {
				hosts = append(hosts, h)
			}
		}
		cfg.RemoteDevEnv = hosts
	}
}
