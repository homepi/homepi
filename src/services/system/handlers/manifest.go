package handlers

import (
	"fmt"
	"net/http"

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

			var baseURI string
			if ctx.Config.Hostname != "" {
				baseURI = ctx.Config.Hostname
			} else {
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}
				baseURI = fmt.Sprintf("%s://%s", scheme, r.Host)
			}

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
