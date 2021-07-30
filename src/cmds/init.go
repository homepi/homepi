package cmds

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
	"github.com/homepi/homepi/pkg/libstr"
	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db"
	"github.com/homepi/homepi/src/db/models"
	"github.com/spf13/cobra"
)

type InitCommandFlags struct {
	ConfigFile string
}

func initCommand() *cobra.Command {
	cFlags := new(InitCommandFlags)
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize homepi",
		Long:  "Initialize homepi database",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := core.LoadConfig(cFlags.ConfigFile)
			if err != nil {
				return err
			}

			database, err := db.NewConnection(cfg)
			if err != nil {
				return fmt.Errorf("could not create database connection : %v", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "==> Migrating database tables...")

			dbModels := []interface{}{
				models.Role{},
				models.User{},
				models.APIToken{},
				models.RefreshToken{},
				models.Accessory{},
				models.Log{},
				models.Webhook{},
			}

			for _, model := range dbModels {
				if err := database.AutoMigrate(model); err != nil {
					return fmt.Errorf("    err on creating [%s] REASON: [%v]", reflect.ValueOf(model).Type().String(), err)
				}
				fmt.Fprintf(cmd.OutOrStdout(), "    Created [%s]\n", reflect.ValueOf(model).Type().String())
			}

			fmt.Fprintln(cmd.OutOrStdout(), "\n==> Creating roles ...")

			adminRole := &models.Role{Title: "root", Administrator: true}
			if err := database.Create(adminRole).Error; err != nil {
				return fmt.Errorf("could not create root user : %v", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "    root role created successfully!")

			userRole := &models.Role{Title: "user"}
			if err := database.Create(userRole).Error; err != nil {
				return fmt.Errorf("could not create user: %v", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "    user role created successfully!")

			var c int64
			if err := database.Model(&models.User{}).Where("username =?", "root").Count(&c).Error; err != nil {
				return fmt.Errorf("could not get root user : %v", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), "\n==> Creating root user ...")
			if c == 0 {

				var (
					adminPassword = libstr.RandomLetters(20)
					adminUser     = &models.User{
						Fullname: "Root",
						Username: "root",
						Email:    "root@homepi.local",
						RoleID:   adminRole.ID,
					}
				)

				if err := adminUser.SetPassword(adminPassword); err != nil {
					return fmt.Errorf("could not create password for user root: [%v]", err)
				}

				if err := database.Create(adminUser).Error; err != nil {
					return fmt.Errorf("could not create root user : %v", err)
				}

				fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf(
					"    root user created successfully! \n"+
						"    credentials : [ User=%s | Pass=%s ]",
					adminUser.Username,
					adminPassword,
				))

				fmt.Fprintln(cmd.OutOrStdout(), "\n==> Creating api token ...")
				apiToken := &models.APIToken{
					ID:     uuid.New().ID(),
					UserID: adminUser.ID,
					RoleID: adminRole.ID,
				}

				if err := database.Create(apiToken).Error; err != nil {
					return fmt.Errorf("could not create api token : %v", err)
				}

				fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf(
					"    api token generated successfully for user %s! \n"+
						"    token : [%s]",
					adminUser.Username,
					apiToken.Token,
				))

			} else {

				fmt.Fprintln(cmd.OutOrStdout(), "    Root user already exists! Skipped.")
			}

			return nil
		},
	}
	cmd.SuggestionsMinimumDistance = 1
	cmd.PersistentFlags().StringVarP(&cFlags.ConfigFile, "config-file", "c", "", "Config file")
	return cmd
}
