package container

import (
	"fmt"
	"strings"

	"github.com/perryrh0dan/dev/internal/config"
	"github.com/spf13/cobra"
)

func newShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <name>",
		Short: "Show stored environment config for a named dev environment",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			env, err := config.LoadEnvironmentByName(name)
			if err != nil {
				return err
			}
			if env == nil {
				env = &config.Environment{Name: name}
			}

			tag := env.Tag
			if tag == "" {
				tag = "(none)"
			}
			ports := strings.Join(env.Port, ", ")
			if ports == "" {
				ports = "(none)"
			}

			fmt.Fprintf(cmd.OutOrStdout(), "name:   %s\n", env.Name)
			fmt.Fprintf(cmd.OutOrStdout(), "tag:    %s\n", tag)
			fmt.Fprintf(cmd.OutOrStdout(), "ports:  %s\n", ports)
			return nil
		},
	}
}
