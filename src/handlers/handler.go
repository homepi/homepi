package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/db"
	accessoryHandler "github.com/homepi/homepi/src/services/accessory/handlers"
	authHandler "github.com/homepi/homepi/src/services/auth/handlers"
	systemHandler "github.com/homepi/homepi/src/services/system/handlers"
	userHandler "github.com/homepi/homepi/src/services/user/handlers"
	webhookHandler "github.com/homepi/homepi/src/services/webhook/handlers"
	"github.com/homepi/homepi/src/validators"
)

// Creates a new http handler with a new database connection
// and http custom validators with core.Context
func NewHandler(cfg *core.ConfMap) (http.Handler, error) {
	if err := validators.Configure(); err != nil {
		return nil, fmt.Errorf("error configuring http validator CLI: %v", err)
	}
	database, err := db.NewConnection(cfg)
	if err != nil {
		return nil, fmt.Errorf("error configuring database: [%v]", err)
	}
	return handler(&core.Context{
		Database: database,
		Config:   cfg,
	}), nil
}

// a new http handler
func handler(core *core.Context) http.Handler {

	mux := chi.NewMux()

	// register global middlewares
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(wrapHTTPLogHandler)
	mux.Use(middleware.Recoverer)
	mux.Use(wrapJSONHandler)
	mux.Use(wrapCORSHandler(core))

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	mux.Use(middleware.Timeout(20 * time.Second))

	mux.Handle("/uploads/avatars/{avatar}.png", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		avatarID := chi.URLParam(r, "avatar")
		avatarFile, err := os.Open(fmt.Sprintf("./uploads/avatars/%s.png", avatarID))
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if _, err := io.Copy(w, avatarFile); err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

	}))

	mux.Route("/api", func(r chi.Router) {

		// Prints api details
		r.Handle("/", systemHandler.HandleHostInfo(core))

		// version 1 routes
		r.Route("/v1", func(r chi.Router) {

			r.Handle("/auth/create.json", authHandler.HandleAuthTokens(core))
			r.Handle("/users.json", userHandler.HandleListUsers(core))
			r.Handle("/roles.json", userHandler.HandleListRoles(core))
			r.Handle("/pins.json", accessoryHandler.HandleListGPIOPins(core))
			r.Handle("/health.json", systemHandler.HandleGetHealth(core))

			r.Route("/users", func(r chi.Router) {
				r.Handle("/me.json", userHandler.HandleUsersMe(core))
				r.Handle("/logs.json", userHandler.HandleListLogs(core))
				r.Handle("/permissions.json", userHandler.HandleUserPermissions(core))
			})

			r.Handle("/accessories.json", accessoryHandler.HandleAccessories(core))
			r.Route("/accessories", func(r chi.Router) {
				r.Handle("/{id}/run.json", accessoryHandler.HandleRunAccessory(core))
				r.Handle("/{id}/get.json", accessoryHandler.HandleGetAccessory(core))
			})

			r.Handle("/webhooks.json", webhookHandler.HandleWebhooks(core))
			r.Route("/webhooks", func(r chi.Router) {
				r.Handle("/{id}/{hash}.json", webhookHandler.HandleRunWebhook(core))
			})

		})

	})

	return mux
}
