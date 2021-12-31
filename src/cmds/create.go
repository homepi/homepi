package cmds

import (
	"fmt"

	"github.com/homepi/homepi/pkg/libstr"
	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db"
	"github.com/homepi/homepi/src/db/models"
	"github.com/spf13/cobra"
)

type CreateCommandFlags struct {
	ConfigFile string
}

func createCommand() *cobra.Command {

	cFlags := new(CreateCommandFlags)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create homepi resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(createUserCommand(cFlags))
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

	cmd.PersistentFlags().StringVarP(&cFlags.ConfigFile, "config-file", "c", "", "Server config file")
	return cmd
}

type CreateUserCommandFlags struct {
	Username string
	Name     string
	Email    string
	Password string
	Role     string
}

func createUserCommand(cFlags *CreateCommandFlags) *cobra.Command {

	ccFlags := new(CreateUserCommandFlags)

	cmd := &cobra.Command{
		Use:   "user",
		Short: "Create a user",
		RunE: func(cmd *cobra.Command, args []string) error {

			if cFlags.ConfigFile == "" {
				return fmt.Errorf("config file is required")
			}

			if ccFlags.Email == "" {
				return fmt.Errorf("email is required")
			}

			if ccFlags.Username == "" {
				return fmt.Errorf("username is required")
			}

			cfg, err := core.LoadConfig(cFlags.ConfigFile)
			if err != nil {
				return err
			}

			database, err := db.NewConnection(cfg)
			if err != nil {
				return fmt.Errorf("could not create database connection : %v", err)
			}

			var role *models.Role

			if ccFlags.Role == "admin" {

				if err := database.Where("title =?", "root").Find(&role).Error; err != nil {
					return fmt.Errorf("could not create root user : %v", err)
				}
			} else {

				if err := database.Where("title =?", "user").Find(&role).Error; err != nil {
					return fmt.Errorf("could not create root user : %v", err)
				}
			}

			name := ccFlags.Name
			if ccFlags.Name == "" {
				name = ccFlags.Username
			}

			var (
				password = libstr.RandomLetters(20)
				user     = &models.User{
					Fullname: name,
					Username: ccFlags.Username,
					Email:    ccFlags.Email,
					RoleID:   role.ID,
				}
			)

			if ccFlags.Password != "" {
				password = ccFlags.Password
			}

			if err := user.SetPassword(password); err != nil {
				return fmt.Errorf("could not create password for user root: [%v]", err)
			}

			if err := database.Create(user).Error; err != nil {
				return fmt.Errorf("could not create root user : %v", err)
			}

			showPass := ""
			if ccFlags.Password == "" {
				showPass = fmt.Sprintf(" | Pass=%s", password)
			}

			fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf(
				"user created successfully! \n"+
					"credentials : [ User=%s%s ]",
				user.Username,
				showPass,
			))

			return nil
		},
	}

	cmd.MarkPersistentFlagRequired("email")
	cmd.MarkPersistentFlagRequired("username")

	cmd.PersistentFlags().StringVarP(&ccFlags.Name, "name", "n", "", "user name")
	cmd.PersistentFlags().StringVarP(&ccFlags.Email, "email", "e", "", "user email")
	cmd.PersistentFlags().StringVarP(&ccFlags.Username, "username", "u", "", "user username")
	cmd.PersistentFlags().StringVarP(&ccFlags.Password, "password", "p", "", "user password")
	cmd.PersistentFlags().StringVarP(&ccFlags.Role, "role", "r", "user", "user role")
	return cmd
}
