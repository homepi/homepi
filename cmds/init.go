package cmds

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/homepi/homepi/api/app/components/strings"
	"github.com/homepi/homepi/api/db/models"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func NewInitCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize homepi",
		Long: "Initialize homepi\n\n" +
			"Required environment variables\n" +
			"SQLITE3_PATH is for sqlite3 database file path\n",
		RunE: func(cmd *cobra.Command, args []string) error {

			if db, ok := cmd.Context().Value("db").(*gorm.DB); ok {

				fmt.Fprintln(cmd.OutOrStdout(), "==> Migrating database tables...")

				dbModels := []interface{}{
					models.Role{},
					models.User{},
					models.RefreshedToken{},
					models.Accessory{},
					models.Log{},
					models.Webhook{},
				}

				for _, model := range dbModels {
					fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("    Created [%s]", reflect.ValueOf(model).Type().String()))
					db.AutoMigrate(model)
				}

				fmt.Fprintln(cmd.OutOrStdout(), "\n==> Creating roles ...")

				var adminRole *models.Role
				if err := db.FirstOrCreate(&adminRole, map[string]interface{}{
					"title":         "root",
					"administrator": true,
				}).Error; err != nil {
					return fmt.Errorf("Could not create root user : %v", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), "    root role created successfully!")

				if err := db.FirstOrCreate(&models.Role{}, map[string]interface{}{
					"title":         "user",
					"administrator": false,
				}).Error; err != nil {
					return fmt.Errorf("Could not create user : %v", err)
				}
				fmt.Fprintln(cmd.OutOrStdout(), "    user role created successfully!")

				var adminCount int64
				if err := db.Model(&models.User{}).Where("username =?", "root").Count(&adminCount).Error; err != nil {
					if !errors.Is(err, gorm.ErrRecordNotFound) {
						return fmt.Errorf("Could not get root user : %v", err)
					}
				}

				fmt.Fprintln(cmd.OutOrStdout(), "\n==> Creating root user ...")
				if adminCount == 0 {
					var (
						adminPassword = strings.RandomLetters(20)
						adminUser     = &models.User{
							Fullname: "Root",
							Username: "root",
							Email:    "root@homepi.local",
							RoleId:   adminRole.ID,
						}
					)
					adminUser.SetPassword(adminPassword)
					if err := db.FirstOrCreate(adminUser, map[string]interface{}{"username": "root"}).Error; err != nil {
						return fmt.Errorf("Could not create root user : %v", err)
					}
					fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf(
						"    root user created successfully! \n"+
							"    Credentials : [ User=root | Pass=%s ]",
						adminPassword,
					))
				} else {
					fmt.Fprintln(cmd.OutOrStdout(), "    Root user already exists! Skipped.")
				}
				return nil
			}
			return fmt.Errorf("Could not get database from context!")
		},
	}
	cmd.SuggestionsMinimumDistance = 1
	return cmd
}
