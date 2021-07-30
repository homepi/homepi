package users

import (
	"net/http"

	"github.com/homepi/homepi/client"
	"github.com/homepi/homepi/src/db/models"
)

type (
	ClientService struct {
		c *client.Client
	}
	// Response of roles endpoint
	ListRolesResponse struct {
		Result  []models.Role `json:"result"`
		Message string        `json:"message"`
		Status  string        `json:"status"`
	}
	// Response of users endpoint
	ListUsersResponse struct {
		Result  []models.User `json:"result"`
		Message string        `json:"message"`
		Status  string        `json:"status"`
	}
)

func NewUsersClientService(client *client.Client) *ClientService {
	return &ClientService{c: client}
}

func (s *ClientService) ListRoles() (response *ListRolesResponse, err error) {
	err = s.c.MakeRequest(&response, client.ListRoleEndpoint, http.MethodGet)
	return
}

func (s *ClientService) ListUsers() (response *ListUsersResponse, err error) {
	err = s.c.MakeRequest(&response, client.ListUsersEndpoint, http.MethodGet)
	return
}
