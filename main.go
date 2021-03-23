package main

import (
	"context"
	"fmt"
	"os"

	"github.com/homepi/homepi/cmds"
	"github.com/spf13/cobra"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	rootCmd := &cobra.Command{
		Use: "homepi",
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
	rootCmd.SetArgs(os.Args[1:])
	rootCmd.AddCommand(cmds.NewApiServerCommand())
	rootCmd.AddCommand(cmds.NewInitCommand())
	rootCmd.AddCommand(cmds.NewUserCommand())
	dbPath := os.Getenv("HPI_SQLITE3_PATH")
	if dbPath == "" {
		fmt.Fprintln(os.Stderr, "HPI_SQLITE3_PATH is required!")
		os.Exit(1)
	}
	database, err := gorm.Open(sqlite.Open(dbPath), nil)
	if err != nil {
		fmt.Fprintln(os.Stderr, fmt.Sprintf("Error opening sqlite3 database CLI: %s", err))
		os.Exit(1)
	}
	mCtx := context.WithValue(context.Background(), "db", database)
	if err := rootCmd.ExecuteContext(mCtx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
