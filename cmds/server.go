package cmds

import (
	"fmt"
	"runtime"

	"github.com/gin-gonic/gin"
	"github.com/homepi/homepi/api/app/http/v1/validators"
	"github.com/homepi/homepi/api/services/accessory"
	"github.com/homepi/homepi/api/services/auth"
	"github.com/homepi/homepi/api/services/system"
	"github.com/homepi/homepi/api/services/user"
	"github.com/homepi/homepi/api/services/webhook"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

type ServerFlags struct {
	Host       string
	Port       int
	IgnoreArch bool
}

func NewApiServerCommand() *cobra.Command {
	cFlags := new(ServerFlags)
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run api server",
		Long: "Runs fortress api server\n\n" +
			"Required environment variables\n" +
			"HPI_SQLITE3_PATH is for sqlite3 database file path\n" +
			"HPI_ACCESS_TOKEN_SECRET is for jwt access token\n" +
			"HPI_REFRESH_TOKEN_SECRET is for jwt refresh token\n\n" +
			"Optional environment variables\n" +
			"HPI_ACCESS_TOKEN_EXPIRE_TIME is a duration that an access_token could be valid (default \"240 minutes\")\n" +
			"HPI_REFRESH_TOKEN_EXPIRE_TIME is a duration that an refresh_token could be valid (default \"1440 minutes\")\n",
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Fprintln(cmd.OutOrStdout(), cmd.Root().Long+"\n")

			if !cFlags.IgnoreArch {
				if runtime.GOOS != "linux" && runtime.GOARCH != "arm" {
					return fmt.Errorf(
						"This application can not run on [%s/%s], Required: All versions of Raspberry\n"+
							"To ignore this error, use `--ignore-arch` flag but be aware that someof functionalities may not work!",
						runtime.GOOS,
						runtime.GOARCH,
					)
				}
			}

			if database, ok := cmd.Context().Value("db").(*gorm.DB); ok {

				if err := validators.Configure(); err != nil {
					return fmt.Errorf("Error configuring http validator CLI: %v\n", err)
				}

				authSvc, err := auth.NewAuthService(database)
				if err != nil {
					return err
				}

				var (
					sysSvc       = system.NewSystemService(database)
					webhookSvc   = webhook.NewWebhookService(database)
					userSvc      = user.NewUserService(database, authSvc)
					accessorySvc = accessory.NewAccessoryService(database)
				)

				router := gin.Default()
				router.Use(sysSvc.CORSMiddleware)

				api := router.Group("api")
				api.GET("manifest.json", sysSvc.Manifest)

				v1 := api.Group("v1")
				v1.GET("webhooks/:id/:hash/run.json", webhookSvc.RunWebhook)

				authGroup := v1.Group("auth")
				authGroup.POST("create.json", authSvc.CreateAuthToken)
				authGroup.POST("refresh.json", authSvc.RefreshAuthToken)
				authReqGroup := v1.Group("")
				authReqGroup.Use(authSvc.Authentication)

				usersGroup := authReqGroup.Group("users")
				usersGroup.GET("all.json", userSvc.GetAllUsers)
				usersGroup.POST("search.json", userSvc.SearchUser)

				pinsGroup := authReqGroup.Group("pins.json")
				pinsGroup.GET("", accessorySvc.GetGpioPins)

				sysGroup := authReqGroup.Group("sys")
				sysGroup.GET("health.json", sysSvc.GetHealthCharts)

				userAuthGroup := authReqGroup.Group("user")
				userAuthGroup.POST("create.json", userSvc.CreateUser)
				userAuthGroup.GET("me.json", userSvc.GetMe)
				userAuthGroup.GET("logs.json", userSvc.GetLogs)
				userAuthGroup.POST("avatar.json", userSvc.UpdateAvatar)
				userAuthGroup.GET("accessories.json", userSvc.GetAccessories)
				userAuthGroup.GET("webhooks.json", webhookSvc.GetWebhooks)

				webhookAdminGroup := authReqGroup.Group("webhooks")
				webhookAdminGroup.POST("create.json", webhookSvc.CreateWebhook)

				accessoriesGroup := authReqGroup.Group("accessories")
				accessoriesGroup.POST("create.json", accessorySvc.CreateAccessory)
				accessoriesGroup.GET(":accessory_id/get.json", accessorySvc.GetAccessory)
				accessoriesGroup.GET(":accessory_id/run.json", accessorySvc.RunAccessory)
				accessoriesGroup.DELETE(":accessory_id/delete.json", accessorySvc.RemoveAccessory)

				if err := router.Run(fmt.Sprintf("%s:%d", cFlags.Host, cFlags.Port)); err != nil {
					return fmt.Errorf("could not serve router: %v", err)
				}

				return nil
			}
			return fmt.Errorf("Could not get database from context!")
		},
	}
	cmd.SuggestionsMinimumDistance = 1
	cmd.Flags().StringVarP(&cFlags.Host, "host", "H", "0.0.0.0", "HTTP listener hostname")
	cmd.Flags().IntVarP(&cFlags.Port, "port", "p", 55283, "HTTP listener port")
	cmd.PersistentFlags().BoolVar(&cFlags.IgnoreArch, "ignore-arch", false, "Ignore arch error")
	return cmd
}
