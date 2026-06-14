package identity

import (
	"github.com/spf13/cobra"
)

func NewIdentityCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "identity",
		Short: "Manage git identities",
	}
	cmd.AddCommand(newAddCmd())
	cmd.AddCommand(newActivateCmd())
	return cmd
}
