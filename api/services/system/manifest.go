package system

import (
	"fmt"
	"os"

	"github.com/MrJoshLab/go-respond"
	"github.com/gin-gonic/gin"
)

type ManifestType struct {
	Version        string `json:"version"`
	AvatarsPattern string `json:"avatars_pattern"`
	BaseUri        string `json:"base_uri"`
	ApiBaseUri     string `json:"api_base_uri"`
}

func (s *Service) Manifest(ctx *gin.Context) {
	hostname, err := os.Hostname()
	if err != nil {
		hostname = "localhost"
	}
	baseUri := fmt.Sprintf("http://%s:55283", hostname)
	ctx.JSON(respond.Default.Succeed(&ManifestType{
		Version:        "v1",
		BaseUri:        baseUri,
		AvatarsPattern: fmt.Sprintf("/uploads/avatars/{avatar_name}.png"),
		ApiBaseUri:     "/api/v1",
	}))
}
