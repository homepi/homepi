package cmds

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"

	"strings"

	"github.com/homepi/homepi/pkg/libstr"
	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db"
	"github.com/homepi/homepi/src/db/models"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

func NewInput(i io.Reader, name, description string, required bool) (string, error) {
	reader := bufio.NewReader(i)
	fmt.Print(description)
	value, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	value = strings.TrimSpace(value)
	if required {
		if value != "" {
			return strings.TrimSpace(value), nil
		}
		return "", errors.New(name + " is required")
	}
	return strings.TrimSpace(value), nil
}

func NewUserCommand() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "user",
		Short: "User commands [get|add|delete|update]",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(&cobra.Command{
		Use:   "get",
		Short: "Get a user",
		RunE: func(cmd *cobra.Command, args []string) error {

			cfg, err := core.LoadConfig()
			if err != nil {
				return err
			}

			database, err := db.NewConnection(cfg)
			if err != nil {
				log.Fatal(fmt.Errorf("error configuring http validator CLI: %v", err))
			}

			username, err := NewInput(cmd.InOrStdin(), "username", "Type username of a user : ", true)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStderr(), err)
				return nil
			}

			fmt.Fprintf(cmd.OutOrStdout(), "==> Finding user with [%s] username ...\n", username)

			user := &models.User{}
			// PRELOAD
			if err := database.Where("username =?", username).Find(&user).Error; err != nil {
				fmt.Fprintln(cmd.OutOrStderr(), fmt.Errorf("could not find user with username [%s]", username))
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), "ID \t Fullname \t Username \t Email \t IsActive \t Role \t JoinedAt \t LastLogin")
			fmt.Fprintf(
				cmd.OutOrStdout(),
				"%d \t %s \t %s \t %s \t %v \t %s \t %s \t %s\n",
				user.ID, user.Fullname, user.Username, user.Email, user.IsActive, user.Role.Title, user.JoinedAt.Format("Y/m/d"), user.LastLogin.Format("Y/m/d"),
			)
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		RunE: func(cmd *cobra.Command, args []string) error {

			var (
				database *gorm.DB
				err      error
				role     string
				user     = &models.User{}
			)

			cfg, err := core.LoadConfig()
			if err != nil {
				return err
			}

			database, err = db.NewConnection(cfg)
			if err != nil {
				log.Fatal(fmt.Errorf("error configuring http validator CLI: %v", err))
			}

			user.Username, err = NewInput(cmd.InOrStdin(), "username", "Username : ", true)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStderr(), err)
				return nil
			}
			user.Fullname = user.Username

			user.Email, err = NewInput(cmd.InOrStdin(), "email", "Email : ", true)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStderr(), err)
				return nil
			}

			password, err := NewInput(cmd.InOrStdin(), "password", "Password: Press enter to use random password ", false)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStderr(), err)
				return nil
			}
			if password == "" {
				password = libstr.RandomLetters(20)
			}

			if err := user.SetPassword(password); err != nil {
				return fmt.Errorf("could not set the password: %v", err)
			}

			role, err = NewInput(cmd.InOrStdin(), "role", "Role [ Can be root|user ] Press enter to use user role : ", false)
			if err != nil {
				fmt.Fprintln(cmd.OutOrStderr(), err)
				return nil
			}
			if role == "" {
				role = "user"
			}
			if role != "admin" && role != "user" {
				fmt.Fprintln(cmd.OutOrStderr(), fmt.Errorf("role can be admin or user"))
				return nil
			}

			fmt.Fprintln(cmd.OutOrStdout(), "==> Creating user ...")

			if err := database.Create(user).Error; err != nil {
				return fmt.Errorf("could not create root user : %v", err)
			}

			fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf(
				"    user created successfully! \n"+
					"    Credentials : [ User=%s | Pass=%s ]",
				user.Username,
				password,
			))
			return nil
		},
	})
	cmd.SuggestionsMinimumDistance = 1
	return cmd
}
