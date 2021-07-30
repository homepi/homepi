package handlers

import (
	"net/http"

	"github.com/homepi/homepi/pkg/gpio"
	"github.com/homepi/homepi/src/core"
	"github.com/mrjosh/respond.go"
)

func HandleListGPIOPins(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			ctx.WrapAuthAdminHandler(getGPIOPins(ctx)).ServeHTTP(w, r)
			return
		default:
			respond.NewWithWriter(w).MethodNotAllowed()
		}
	})
}

func getGPIOPins(ctx *core.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond.NewWithWriter(w).Succeed(gpio.GetPins(ctx.Database))
	})
}
