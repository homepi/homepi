package accessories

import (
	"net/http"

	"github.com/homepi/homepi/client"
	"github.com/homepi/homepi/src/db/models"
)

type (
	ClientService struct {
		c *client.Client
	}
	// Response of accessories endpoint
	ListAccessoriesResponse struct {
		Result  []models.User `json:"result"`
		Message string        `json:"message"`
		Status  string        `json:"status"`
	}
)

func NewAccessoriesClientService(client *client.Client) *ClientService {
	return &ClientService{c: client}
}

func (s *ClientService) ListAccessories() (response *ListAccessoriesResponse, err error) {
	err = s.c.MakeRequest(&response, client.ListAccessoriesEndpoint, http.MethodGet)
	return
}
