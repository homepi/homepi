package cmds

import (
	"github.com/spf13/cobra"
)

func Run(args []string) error {

	rootCmd := &cobra.Command{
		Use:     "homepi",
		Version: "v0.0.1",
		Long: `
    __  __                     ____  _ 
   / / / /___  ________  ___  / __ \(_)
  / /_/ / __ \/ __  __ \/ _ \/ /_/ / / 
 / __  / /_/ / / / / / /  __/ ____/ /  
/_/ /_/\____/_/ /_/ /_/\___/_/   /_/`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.SetArgs(args)
	rootCmd.AddCommand(versionCommand())
	rootCmd.AddCommand(apiServerCommand())
	rootCmd.AddCommand(initCommand())
	rootCmd.AddCommand(getCommand())
	rootCmd.AddCommand(createCommand())
	rootCmd.AddCommand(deleteCommand())
	rootCmd.AddCommand(runAccessoryCommand())

	return rootCmd.Execute()
}
