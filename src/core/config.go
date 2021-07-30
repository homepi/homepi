package core

import (
	"errors"
	"io"
	"os"

	"github.com/joeshaw/envdecode"
	"gopkg.in/yaml.v3"
)

type (
	ConfMap struct {
		Environment string          `env:"HPI_ENVIRONMENT,default=development" yaml:"env"`
		TimeZone    string          `env:"HPI_TIMEZONE,default=Asia/Tehran" yaml:"timezone"`
		IgnoreArch  bool            `env:"HPI_IGNORE_ARCH_WARN,default=false" yaml:"ignoreArchWarning"`
		DB          DatabaseConfMap `yaml:"database"`
		JWT         JWTConfigMap    `yaml:"jwt"`
	}

	DatabaseConfMap struct {
		Driver string `env:"HPI_DB_DRIVER,default=sqlite" yaml:"driver"`
		Path   string `env:"HPI_DB_PATH,default=data/homepi.db" yaml:"path"`
		User   string `env:"HPI_DB_USER" yaml:"user"`
		Pass   string `env:"HPI_DB_PASS" yaml:"pass"`
		Name   string `env:"HPI_DB_NAME" yaml:"name"`
	}

	JWTConfigMap struct {
		AccessToken struct {
			Value     string `env:"HPI_JWT_ACCESS_TOKEN,default=super-secure-access-token" yaml:"secret"`
			ExpiresAt int    `env:"HPI_JWT_ACCESS_TOKEN_EXPIRES_AT,default=240" yaml:"expires_at"`
		} `yaml:"access_token"`

		RefreshToken struct {
			Value     string `env:"HPI_JWT_REFRESH_TOKEN,default=super-secure-refresh-token" yaml:"secret"`
			ExpiresAt int    `env:"HPI_JWT_REFRESH_TOKEN_EXPIRES_AT,default=1440" yaml:"expires_at"`
		} `yaml:"refresh_token"`
	}
)

func LoadENV() (*ConfMap, error) {
	var cfg ConfMap
	if err := envdecode.StrictDecode(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func LoadYAMLFromReader(r io.Reader) (cfg *ConfMap, err error) {
	err = yaml.NewDecoder(r).Decode(&cfg)
	return
}

func LoadConfig(filename ...string) (*ConfMap, error) {
	if filename != nil && filename[0] != "" {
		if len(filename) > 1 {
			return nil, errors.New("only 1 config filename as an argument is valid")
		}
		f, err := os.Open(filename[0])
		if err != nil {
			return nil, err
		}
		return LoadYAMLFromReader(f)
	}
	return LoadENV()
}
