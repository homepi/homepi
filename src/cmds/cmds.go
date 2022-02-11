package cmds

import (
	"github.com/homepi/homepi/src/core"
	"github.com/spf13/cobra"
)

var versionInfo *core.VersionInfo

func RegisterAndRun(vi *core.VersionInfo, rootCmd *cobra.Command) error {
	vi.BuildType = "Release"
	if vi.BranchName == "develop" {
		vi.BuildType = "Nightly"
	}
	versionInfo = vi
	rootCmd.AddCommand(versionCommand())
	rootCmd.AddCommand(apiServerCommand())
	rootCmd.AddCommand(initCommand())
	rootCmd.AddCommand(getCommand())
	rootCmd.AddCommand(createCommand())
	rootCmd.AddCommand(deleteCommand())
	rootCmd.AddCommand(runAccessoryCommand())
	return rootCmd.Execute()
}
