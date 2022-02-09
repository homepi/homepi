package core

import (
	"sync"

	"github.com/homepi/homepi/pkg/libstr"
)

const (
	CORSDisabled = false
	CORSEnabled  = true
)

func (ctx *Context) CORSConfig() *CORSConfig {
	return &CORSConfig{
		Enabled:        CORSEnabled,
		AllowedOrigins: ctx.Config.AllowedHosts,
		AllowedHeaders: ctx.Config.AllowedHeaders,
	}
}

type CORSConfig struct {
	Enabled        bool     `json:"enabled"`
	AllowedOrigins []string `json:"allowed_origins,omitempty"`
	AllowedHeaders []string `json:"allowed_headers,omitempty"`
	sync.RWMutex   `json:"-"`
}

// IsEnabled returns the value of CORSConfig.isEnabled
func (c *CORSConfig) IsEnabled() bool {
	return c.Enabled == CORSEnabled
}

// IsValidOrigin determines if the origin of the request is allowed to make
// cross-origin requests based on the CORSConfig.
func (c *CORSConfig) IsValidOrigin(origin string) bool {
	// If we aren't enabling CORS then all origins are valid
	if !c.IsEnabled() {
		return true
	}

	c.RLock()
	defer c.RUnlock()

	if len(c.AllowedOrigins) == 0 {
		return false
	}

	if len(c.AllowedOrigins) == 1 && (c.AllowedOrigins)[0] == "*" {
		return true
	}

	return libstr.StrListContains(c.AllowedOrigins, origin)
}
