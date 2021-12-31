package cmds

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/homepi/homepi/src/core"
	"github.com/homepi/homepi/src/handlers"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

type ServerFlags struct {
	Host              string
	Port              int
	IgnoreArch        bool
	ConfigFile        string
	LogOutputFilename string
	NoTLS             bool
	TLSCertFile       string
	TLSKeyFile        string
}

func apiServerCommand() *cobra.Command {
	cFlags := new(ServerFlags)
	cmd := &cobra.Command{
		Use:   "server",
		Short: "Run api server",
		Long:  "Runs fortress api server",
		RunE: func(cmd *cobra.Command, args []string) error {

			logrus.SetOutput(cmd.OutOrStdout())
			if cFlags.LogOutputFilename != "" {
				logFile, err := os.OpenFile(cFlags.LogOutputFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
				if err != nil {
					logrus.WithFields(logrus.Fields{
						"path": cFlags.LogOutputFilename,
					}).Errorf("could not open/create log file: %v", err)
					return errors.New("")
				}
				logrus.SetOutput(logFile)
			}

			cfg, err := core.LoadConfig(cFlags.ConfigFile)
			if err != nil {
				return err
			}

			cfg.Port = cFlags.Port

			logrus.SetFormatter(&logrus.TextFormatter{})
			if cfg.Environment == "production" {
				logrus.SetFormatter(&logrus.JSONFormatter{})
			}

			logrus.WithFields(logrus.Fields{"path": cFlags.ConfigFile}).Infof("Config file loaded")

			if !cFlags.IgnoreArch && !cfg.IgnoreArch && cfg.Environment == "production" {
				if runtime.GOOS != "linux" && runtime.GOARCH != "arm" {
					logrus.WithFields(logrus.Fields{
						"error":  fmt.Sprintf("This application can not run on [%s/%s]", runtime.GOOS, runtime.GOARCH),
						"reason": "Required all versions of raspberry pi",
					}).Log(logrus.ErrorLevel, "Could not run server on this machine")
					return errors.New("could not run server on this machine")
				}
			}

			logrus.WithFields(logrus.Fields{
				"driver": cfg.DB.Driver,
				"path":   cfg.DB.Path,
			}).Infof("Creating new database connection")

			logrus.Info("Creating new http handler")
			handler, err := handlers.NewHandler(cfg)
			if err != nil {
				return fmt.Errorf("could not create new http handler: %v", err)
			}

			server := &http.Server{
				Addr:              fmt.Sprintf(":%d", cFlags.Port),
				Handler:           handler,
				ReadHeaderTimeout: 10 * time.Second,
				ReadTimeout:       30 * time.Second,
				IdleTimeout:       5 * time.Minute,
			}

			protocol := "http"
			if !cFlags.NoTLS {
				protocol = "https"
			}

			logrus.WithFields(logrus.Fields{
				"addr": fmt.Sprintf("%s://%s:%d", protocol, cFlags.Host, cFlags.Port),
			}).Infof("HomePi server is up and running")

			if !cFlags.NoTLS {

				if cFlags.TLSCertFile == "" || cFlags.TLSKeyFile == "" {
					log.Fatal("You must provide CertFile and KeyFile for TLS option, Usage: --tls-cert-file {path} --tls-key-file {path}")
				}

				if err := server.ListenAndServeTLS(cFlags.TLSCertFile, cFlags.TLSKeyFile); err != nil {
					return fmt.Errorf("could not serve TLS router: %v", err)
				}

			} else {
				if err := server.ListenAndServe(); err != nil {
					return fmt.Errorf("could not serve router: %v", err)
				}
			}

			return nil
		},
	}
	cmd.SuggestionsMinimumDistance = 1
	cmd.PersistentFlags().StringVarP(&cFlags.Host, "host", "H", "0.0.0.0", "HTTP listener hostname")
	cmd.PersistentFlags().IntVarP(&cFlags.Port, "port", "p", 55283, "HTTP listener port")
	cmd.PersistentFlags().BoolVar(&cFlags.IgnoreArch, "ignore-arch", false, "Ignore arch error")
	cmd.PersistentFlags().StringVarP(&cFlags.ConfigFile, "config-file", "c", "", "Server config file")
	cmd.PersistentFlags().StringVar(&cFlags.LogOutputFilename, "log-file", "", "Log output filename")
	cmd.PersistentFlags().BoolVar(&cFlags.NoTLS, "no-tls", true, "Start server without tls enabled")
	cmd.PersistentFlags().StringVar(&cFlags.TLSCertFile, "tls-cert-file", "", "TLS Certificate file")
	cmd.PersistentFlags().StringVar(&cFlags.TLSKeyFile, "tls-key-file", "", "TLS Key file")
	return cmd
}
