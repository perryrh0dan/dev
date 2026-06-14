package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	containercmd "github.com/perryrh0dan/dev/cmd/container"
	imagecmd "github.com/perryrh0dan/dev/cmd/image"
	identitycmd "github.com/perryrh0dan/dev/cmd/identity"
	"github.com/perryrh0dan/dev/internal/config"
	"github.com/perryrh0dan/dev/internal/container"
	devssh "github.com/perryrh0dan/dev/internal/ssh"
	"github.com/spf13/cobra"
)

func NewRootCmd(cfg *config.Config, engine container.ContainerEngine) *cobra.Command {
	var ports []string
	var tag string

	rootCmd := &cobra.Command{
		Use:   "dev [volume|directory]",
		Short: "Docker-based dev environment manager",
		Long:  "Launch and manage Docker-based development environments.",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return devssh.ConnectAll(cfg.RemoteDevEnv)
			}

			target := args[0]

			// Expand ~ to home directory
			if strings.HasPrefix(target, "~/") {
				home, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("get home dir: %w", err)
				}
				target = filepath.Join(home, target[2:])
			}

			// Determine if target is a directory path or a volume name
			isDir := strings.HasPrefix(target, ".") ||
				strings.HasPrefix(target, "/") ||
				strings.HasPrefix(target, "~") ||
				strings.Contains(target, string(filepath.Separator))

			// Load or build environment config (only meaningful for named volumes, not directory mounts)
			envCfg := &config.Environment{Name: target}
			if !isDir {
				loaded, err := config.LoadEnvironmentByName(target)
				if err != nil {
					return err
				}
				if loaded != nil {
					envCfg = loaded
				}
			}

			// Apply flag overrides
			if len(ports) > 0 {
				envCfg.Port = ports
			}
			if tag != "" {
				envCfg.Tag = tag
			}

			// Persist updated environment config only for named volumes when flags were explicitly provided
			if !isDir && (len(ports) > 0 || tag != "") {
				if err := config.SaveEnvironment(*envCfg); err != nil {
					return err
				}
			}

			// Build image name
			image := cfg.Image
			if envCfg.Tag != "" {
				image = image + ":" + envCfg.Tag
			}

			// Memory limit: config value or default
			memory := cfg.Memory
			if memory == "" {
				memory = "24gb"
			}

			// Build mounts
			var mounts []container.Mount

			if isDir {
				absPath, err := filepath.Abs(target)
				if err != nil {
					return err
				}
				mounts = append(mounts, container.Mount{Type: "bind", Source: absPath, Target: "/root/workspace"})
			} else {
				// Named volume: ensure companion volumes exist
				for _, vol := range []string{target, target + "-history", target + "-resurrect", "dev-zoxide"} {
					if err := engine.CreateVolume(vol); err != nil {
						fmt.Fprintf(os.Stderr, "warning: could not create volume %q: %v\n", vol, err)
					}
				}

				mounts = append(mounts,
					container.Mount{Type: "volume", Source: target, Target: "/root/workspace"},
					container.Mount{Type: "volume", Source: target + "-history", Target: "/root/.history"},
					container.Mount{Type: "volume", Source: target + "-resurrect", Target: "/root/.local/share/tmux/resurrect"},
					container.Mount{Type: "volume", Source: "dev-zoxide", Target: "/root/.local/share/zoxide"},
				)
			}

			// Append configured bind mounts (skip empty)
			type dirMount struct{ src, dst string }
			dirMounts := []dirMount{
				{cfg.SSHDirectory, "/root/.ssh"},
				{cfg.SharedDirectory, "/root/shared"},
				{cfg.KubeDirectory, "/root/.kube"},
				{cfg.TimewarriorDirectory, "/root/.local/share/timewarrior"},
				{cfg.OpenCodeDirectory, "/root/.local/share/opencode"},
				{cfg.NgrokDirectory, "/root/.config/ngrok"},
				{cfg.GitHubCopilotDirectory, "/root/.config/github-copilot"},
			}
			for _, dm := range dirMounts {
				if dm.src != "" {
					mounts = append(mounts, container.Mount{Type: "bind", Source: dm.src, Target: dm.dst})
				}
			}

			// GPG: mount specific files rather than the whole directory
			if cfg.GPGDirectory != "" {
				for _, f := range []struct{ src, dst string }{
					{filepath.Join(cfg.GPGDirectory, "pubring.kbx"), "/root/.gnupg/pubring.kbx"},
					{filepath.Join(cfg.GPGDirectory, "trustdb.gpg"), "/root/.gnupg/trustdb.gpg"},
					{filepath.Join(cfg.GPGDirectory, "private-keys-v1.d"), "/root/.gnupg/private-keys-v1.d"},
				} {
					if _, err := os.Stat(f.src); err == nil {
						mounts = append(mounts, container.Mount{Type: "bind", Source: f.src, Target: f.dst})
					}
				}
			}

			// File mounts (skip empty)
			type fileMount struct{ src, dst string }
			fileMounts := []fileMount{
				{cfg.EnvFile, "/root/.local/share/.env"},
				{cfg.DictFile, "/root/dict.txt"},
				{cfg.NpmFile, "/root/.npmrc"},
			}
			for _, fm := range fileMounts {
				if fm.src != "" {
					mounts = append(mounts, container.Mount{Type: "bind", Source: fm.src, Target: fm.dst})
				}
			}

			// Docker socket (Docker-in-Docker)
			mounts = append(mounts, container.Mount{Type: "bind", Source: "/var/run/docker.sock", Target: "/var/run/docker.sock"})

			// Build env vars from active identity
			var env []string
			identity, err := config.LoadActiveIdentity()
			if err != nil {
				fmt.Fprintf(os.Stderr, "warning: could not load identity: %v\n", err)
			}
			if identity != nil {
				env = append(env,
					"GIT_USER="+identity.Name,
					"GIT_EMAIL="+identity.Email,
					"GIT_SIGNINGKEY="+identity.KeyID,
				)
			}

			// Container name: use volume name (not for directory mounts)
			name := ""
			if !isDir {
				name = target
			}

			return engine.RunInteractive(container.RunOptions{
				Image:      image,
				Name:       name,
				Mounts:     mounts,
				Ports:      envCfg.Port,
				Env:        env,
				Labels:     map[string]string{"dev": "yes"},
				Privileged: true,
				Memory:     memory,
				Remove:     true,
			})
		},
	}

	rootCmd.Flags().StringArrayVarP(&ports, "port", "p", nil, "Port mappings to expose (e.g. 8080:8080), can be repeated")
	rootCmd.Flags().StringVar(&tag, "tag", "", "Image tag to use (persisted per environment)")

	rootCmd.AddCommand(containercmd.NewContainerCmd(cfg, engine))
	rootCmd.AddCommand(imagecmd.NewImageCmd(cfg, engine))
	rootCmd.AddCommand(identitycmd.NewIdentityCmd())

	return rootCmd
}
