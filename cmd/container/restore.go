package container

import (
	"path/filepath"

	"github.com/perryrh0dan/dev/internal/container"
	"github.com/spf13/cobra"
)

func newRestoreCmd(engine container.ContainerEngine) *cobra.Command {
	return &cobra.Command{
		Use:   "restore <name> <path>",
		Short: "Restore a dev environment from a backup directory",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			backupPath, err := filepath.Abs(args[1])
			if err != nil {
				return err
			}

			if err := engine.CreateVolume(name); err != nil {
				return err
			}
			if err := engine.CreateVolume(name + "-history"); err != nil {
				return err
			}
			if err := engine.CreateVolume(name + "-resurrect"); err != nil {
				return err
			}
			if err := engine.CreateVolume("dev-zoxide"); err != nil {
				return err
			}

			script := `tar -xzf /input/backup.tar -C /restore/workspace && ` +
				`tar -xzf /input/backup-history.tar -C /restore/history && ` +
				`tar -xzf /input/backup-zoxide.tar -C /restore/zoxide && ` +
				`tar -xzf /input/backup-resurrect.tar -C /restore/resurrect`

			return engine.RunInteractive(container.RunOptions{
				Image:  "ubuntu:26.04",
				Remove: true,
				Mounts: []container.Mount{
					{Type: "volume", Source: name, Target: "/restore/workspace"},
					{Type: "volume", Source: name + "-history", Target: "/restore/history"},
					{Type: "volume", Source: "dev-zoxide", Target: "/restore/zoxide"},
					{Type: "volume", Source: name + "-resurrect", Target: "/restore/resurrect"},
					{Type: "bind", Source: backupPath, Target: "/input"},
				},
				Command: []string{"sh", "-c", script},
			})
		},
	}
}
