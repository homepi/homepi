package cmds

import (
	"bufio"
	"errors"
	"fmt"
	"io"

	libStr "strings"

	"github.com/homepi/homepi/api/app/components/strings"
	"github.com/homepi/homepi/api/db/models"
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
	value = libStr.TrimSpace(value)
	if required {
		if value != "" {
			return libStr.TrimSpace(value), nil
		}
		return "", errors.New(name + " is required")
	}
	return libStr.TrimSpace(value), nil
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
			if db, ok := cmd.Context().Value("db").(*gorm.DB); ok {
				username, err := NewInput(cmd.InOrStdin(), "username", "Type username of a user : ", true)
				if err != nil {
					fmt.Fprintln(cmd.OutOrStderr(), err)
					return nil
				}

				fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf("==> Finding user with [%s] username ...", username))

				user := new(models.User)
				if err := db.Where("username =?", username).Preload("Role").First(&user).Error; err != nil {
					fmt.Fprintln(cmd.OutOrStderr(), fmt.Errorf("Could not find user with username [%s]", username))
					return nil
				}

				fmt.Fprintln(cmd.OutOrStdout(), "ID \t Fullname \t Username \t Email \t IsActive \t Role \t JoinedAt \t LastLogin")
				fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf(
					"%d \t %s \t %s \t %s \t %v \t %s \t %s \t %s",
					user.ID, user.Fullname, user.Username, user.Email, user.IsActive, user.Role.Title, user.JoinedAt.Format("Y/m/d"), user.LastLogin.Format("Y/m/d"),
				))
				return nil
			}
			return fmt.Errorf("Could not get database from context")
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "create",
		Short: "Create a new user",
		RunE: func(cmd *cobra.Command, args []string) error {
			if db, ok := cmd.Context().Value("db").(*gorm.DB); ok {

				var (
					err  error
					role string
					user = new(models.User)
				)

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
					password = strings.RandomLetters(20)
				}
				user.SetPassword(password)

				role, err = NewInput(cmd.InOrStdin(), "role", "Role [ Can be root|user ] Press enter to use user role : ", false)
				if err != nil {
					fmt.Fprintln(cmd.OutOrStderr(), err)
					return nil
				}
				if role == "" {
					role = "user"
				}
				if role != "admin" && role != "user" {
					fmt.Fprintln(cmd.OutOrStderr(), fmt.Errorf("Role can be admin or user!"))
					return nil
				}

				fmt.Fprintln(cmd.OutOrStdout(), "==> Creating user ...")

				if err := db.Create(user).Error; err != nil {
					return fmt.Errorf("Could not create root user : %v", err)
				}

				fmt.Fprintln(cmd.OutOrStdout(), fmt.Sprintf(
					"    user created successfully! \n"+
						"    Credentials : [ User=%s | Pass=%s ]",
					user.Username,
					password,
				))
				return nil
			}
			return fmt.Errorf("Could not get database from context")
		},
	})
	cmd.SuggestionsMinimumDistance = 1
	return cmd
}
