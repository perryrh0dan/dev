package container

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/perryrh0dan/dev/internal/container"
	"github.com/spf13/cobra"
)

func newBackupCmd(engine container.ContainerEngine) *cobra.Command {
	return &cobra.Command{
		Use:   "backup <name>",
		Short: "Backup a dev environment to ./backup-<name>/",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			backupDir, err := filepath.Abs("backup-" + name)
			if err != nil {
				return err
			}

			if err := os.MkdirAll(backupDir, 0755); err != nil {
				return fmt.Errorf("create backup dir: %w", err)
			}

			script := `tar -czf /output/backup.tar ` +
				`--exclude=./node_modules --exclude=./.angular --exclude=./.nx ` +
				`--exclude=./dist --exclude=./.pnpm-store --exclude=./.next ` +
				`-C /backup/workspace . && ` +
				`tar -czf /output/backup-history.tar -C /backup/history . && ` +
				`tar -czf /output/backup-zoxide.tar -C /backup/zoxide . && ` +
				`tar -czf /output/backup-resurrect.tar -C /backup/resurrect .`

			return engine.RunInteractive(container.RunOptions{
				Image:  "ubuntu:26.04",
				Remove: true,
				Mounts: []container.Mount{
					{Type: "volume", Source: name, Target: "/backup/workspace"},
					{Type: "volume", Source: name + "-history", Target: "/backup/history"},
					{Type: "volume", Source: "dev-zoxide", Target: "/backup/zoxide"},
					{Type: "volume", Source: name + "-resurrect", Target: "/backup/resurrect"},
					{Type: "bind", Source: backupDir, Target: "/output"},
				},
				Command: []string{"sh", "-c", script},
			})
		},
	}
}
