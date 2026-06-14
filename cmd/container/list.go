package container

import (
	"fmt"

	"github.com/perryrh0dan/dev/internal/container"
	"github.com/spf13/cobra"
)

func newListCmd(engine container.ContainerEngine) *cobra.Command {
	var onlyRunning bool

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List dev environments",
		RunE: func(cmd *cobra.Command, args []string) error {
			volumes, err := engine.ListVolumes()
			if err != nil {
				return fmt.Errorf("list volumes: %w", err)
			}

			runningNames, err := engine.ListRunningDevContainerNames()
			if err != nil {
				return fmt.Errorf("list running containers: %w", err)
			}

			running := make(map[string]bool, len(runningNames))
			for _, name := range runningNames {
				running[name] = true
			}

			for _, v := range volumes {
				status := "STOPPED"
				if running[v] {
					status = "RUNNING"
				}
				if onlyRunning && status == "STOPPED" {
					continue
				}
				fmt.Fprintf(cmd.OutOrStdout(), "%-20s  %s\n", v, status)
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&onlyRunning, "running", "r", false, "Only show currently running environments")
	return cmd
}
