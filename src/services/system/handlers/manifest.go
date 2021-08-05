package handlers

import (
	"fmt"
	"net/http"
	"os"

	"github.com/homepi/homepi/src/core"
	"github.com/mrjosh/respond.go"
)

type ManifestType struct {
	Version        string `json:"version"`
	AvatarsPattern string `json:"avatars_pattern"`
	BaseURI        string `json:"base_uri"`
	APIBaseURI     string `json:"api_base_uri"`
}

func HandleHostInfo(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			var (
				err      error
				hostname = ctx.Config.Hostname
			)
			if hostname == "" {
				hostname, err = os.Hostname()
				if err != nil {
					hostname = "127.0.0.1"
				}
			}
			port := ctx.Config.Port
			baseURI := fmt.Sprintf("http://%s:%d", hostname, port)
			respond.NewWithWriter(w).Succeed(&ManifestType{
				Version:        "v1",
				BaseURI:        baseURI,
				AvatarsPattern: "/uploads/avatars/{avatar_name}.png",
				APIBaseURI:     "/api/v1",
			})
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}
