package cmds

import (
	"github.com/spf13/cobra"
)

// TODO
func runAccessoryCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run an accessory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
}
