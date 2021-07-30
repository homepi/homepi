package core

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfigMap(t *testing.T) {

	exampleConfigFile := `
env: development
ignoreArchWarning: true
database:
  driver: "sqlite"
  path: "./db/data/homepi.db"
  user: homepi
  pass: "super-secure-password"
  name: homepi
jwt:
  access_token:
    secret: "super-secure-jwt-access-token"
    expires_at: 240
  refresh_token:
    secret: "super-secure-jwt-refresh-token"
    expires_at: 1440
`

	t.Run("TestWithEnvironmentVariables", func(t *testing.T) {

		t.Parallel()

		cfg, err := LoadENV()
		assert.Error(t, err, fmt.Errorf("the environment variable \"HPI_JWT_ACCESS_TOKEN\" is missing"))

		err = os.Setenv("HPI_JWT_ACCESS_TOKEN", "super-secure-access-token")
		assert.NoError(t, err)

		err = os.Setenv("HPI_JWT_REFRESH_TOKEN", "super-secure-refresh-token")
		assert.NoError(t, err)

		err = os.Setenv("HPI_DB_DRIVER", "sqlite")
		assert.NoError(t, err)

		err = os.Setenv("HPI_DB_PATH", "./db/data/homepi.db")
		assert.NoError(t, err)

		cfg, err = LoadENV()

		assert.Equal(t, cfg.DB.Driver, "sqlite")
		assert.Equal(t, cfg.DB.Path, "./db/data/homepi.db")
		assert.Equal(t, cfg.JWT.AccessToken.Value, "super-secure-access-token")
		assert.Equal(t, cfg.JWT.RefreshToken.Value, "super-secure-refresh-token")

	})

	t.Run("TestWithYAMLConfig", func(t *testing.T) {

		t.Parallel()

		cfg, err := LoadYAMLFromReader(strings.NewReader(exampleConfigFile))
		assert.NoError(t, err)

		assert.Equal(t, cfg.Environment, "development")
		assert.Equal(t, cfg.DB.Driver, "sqlite")
		assert.Equal(t, cfg.DB.Path, "./db/data/homepi.db")
		assert.Equal(t, cfg.JWT.AccessToken.Value, "super-secure-jwt-access-token")
		assert.Equal(t, cfg.JWT.RefreshToken.Value, "super-secure-jwt-refresh-token")
		assert.Equal(t, cfg.JWT.AccessToken.ExpiresAt, 240)
		assert.Equal(t, cfg.JWT.RefreshToken.ExpiresAt, 1440)

	})

}
