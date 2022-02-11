package handlers

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/homepi/homepi/src/core"
	"github.com/mrjosh/respond.go"
)

type SystemInfo struct {
	OperatingSystem string `json:"operating_system"`
	Arch            string `json:"arch"`
	Version         string `json:"version"`
	GoVersion       string `json:"go_version"`
	CompiledBy      string `json:"compiled_by"`
	BuildTime       string `json:"build_time"`
	BuildType       string `json:"build_type"`
	AvatarsPattern  string `json:"avatars_pattern"`
	BaseURI         string `json:"base_uri"`
	APIBaseURI      string `json:"api_base_uri"`
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

			respond.NewWithWriter(w).Succeed(&SystemInfo{
				OperatingSystem: runtime.GOOS,
				Arch:            runtime.GOARCH,
				Version:         ctx.Config.VersionInfo.Version,
				BuildType:       ctx.Config.VersionInfo.BuildType,
				BaseURI:         baseURI,
				AvatarsPattern:  "/uploads/avatars/{avatar_name}.png",
				APIBaseURI:      "/api/v1",
				GoVersion:       ctx.Config.VersionInfo.GoVersion,
				CompiledBy:      ctx.Config.VersionInfo.CompiledBy,
				BuildTime:       ctx.Config.VersionInfo.BuildTime,
			})
			return

		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}
