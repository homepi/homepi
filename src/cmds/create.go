package cmds

import (
	"github.com/spf13/cobra"
)

// TODO
func createCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create homepi resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "user",
		Short: "Create a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "accessory",
		Short: "Create an accessory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "role",
		Short: "Create a role",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
	return cmd
}
