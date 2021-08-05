package cmds

import (
	"fmt"
	"os"

	"github.com/homepi/homepi/client"
	"github.com/homepi/homepi/client/services/accessories"
	"github.com/homepi/homepi/client/services/users"
	"github.com/spf13/cobra"
)

func NewGetRoleCommand(cFlags *GetResourcesFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "role",
		Short: "Get roles",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := createNewHTTPApiClient(cFlags)
			if err != nil {
				return fmt.Errorf("could not create HttpApiClient: [%v]", err)
			}
			usersClient := users.NewUsersClientService(c)
			roles, err := usersClient.ListRoles()
			if err != nil {
				return fmt.Errorf("could not get roles list: [%v]", err)
			}
			if roles.Status == client.FailedResponse {
				return fmt.Errorf("could not get roles list: [message=%s]", roles.Message)
			}
			NewTableWriter(cmd.OutOrStdout(), roles.Result, cFlags.Output)
			return nil
		},
	}
	return cmd
}

func NewGetUserCommand(cFlags *GetResourcesFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "user",
		Short: "Get users",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := createNewHTTPApiClient(cFlags)
			if err != nil {
				return fmt.Errorf("could not create HttpApiClient: [%v]", err)
			}
			usersClient := users.NewUsersClientService(c)
			users, err := usersClient.ListUsers()
			if err != nil {
				return fmt.Errorf("could not get users list: [%v]", err)
			}
			if users.Status == client.FailedResponse {
				return fmt.Errorf("could not get users list: [message=%s]", users.Message)
			}
			NewTableWriter(cmd.OutOrStdout(), users.Result, cFlags.Output)
			return nil
		},
	}
	return cmd
}

func NewGetAccessoryCommand(cFlags *GetResourcesFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accessory",
		Short: "Get accessories",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := createNewHTTPApiClient(cFlags)
			if err != nil {
				return fmt.Errorf("could not create HttpApiClient: [%v]", err)
			}
			accessoriesClient := accessories.NewAccessoriesClientService(c)
			accessories, err := accessoriesClient.ListAccessories()
			if err != nil {
				return fmt.Errorf("could not get accessories list: [%v]", err)
			}
			if accessories.Status == client.FailedResponse {
				return fmt.Errorf("could not get accessories list: [message=%s]", accessories.Message)
			}
			NewTableWriter(cmd.OutOrStdout(), accessories.Result, cFlags.Output)
			return nil
		},
	}
	return cmd
}

type GetResourcesFlags struct {
	BaseURL string
	Output  string
}

func createNewHTTPApiClient(cFlags *GetResourcesFlags) (c *client.Client, err error) {
	if apiURL := os.Getenv("HPI_API_SERVER_URL"); apiURL != "" {
		cFlags.BaseURL = apiURL
	}
	c, err = client.NewClient(cFlags.BaseURL)
	return
}

func getCommand() (cmd *cobra.Command) {
	cFlags := new(GetResourcesFlags)
	cmd = &cobra.Command{
		Use:   "get",
		Short: "Get homepi resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(NewGetUserCommand(cFlags))
	cmd.AddCommand(NewGetRoleCommand(cFlags))
	cmd.AddCommand(NewGetAccessoryCommand(cFlags))
	cmd.PersistentFlags().StringVarP(&cFlags.Output, "output", "o", "", "Set output format")
	cmd.PersistentFlags().StringVarP(&cFlags.BaseURL, "url", "u", "http://127.0.0.1:55283", "Base url of api server")
	return
}
