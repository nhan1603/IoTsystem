package httpserver

import (
	"net/http"
)

// CORSConfig holds the CORS configuration
type CORSConfig struct {
	allowedOrigins   []string
	allowedMethods   []string
	allowedHeaders   []string
	exposedHeaders   []string
	allowCredentials bool
	maxAge           int
}

// NewCORSConfig initializes and returns a CORSConfig
func NewCORSConfig(origins []string, opts ...CORSOption) CORSConfig {
	cfg := CORSConfig{
		allowedOrigins: origins,
		allowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},
		allowedHeaders: []string{
			"Accept",
			"Authorization",
			"Content-Type",
			"x-datadog-origin",
			"x-datadog-parent-id",
			"x-datadog-sampling-priority",
			"x-datadog-trace-id",
			"x-datadog-sampled",
		},
		exposedHeaders:   []string{"Link"},
		allowCredentials: true,
		maxAge:           300,
	}
	for _, o := range opts {
		o(&cfg)
	}

	return cfg
}

// CORSOption enables tweaking the CORSConfig
type CORSOption func(*CORSConfig)

// CORSAddExtraAllowedMethods adds to default allowed methods
func CORSAddExtraAllowedMethods(methods ...string) CORSOption {
	return func(c *CORSConfig) {
		c.allowedMethods = append(c.allowedMethods, methods...)
	}
}

// CORSAddExtraAllowedHeaders adds to default allowed headers
func CORSAddExtraAllowedHeaders(headers ...string) CORSOption {
	return func(c *CORSConfig) {
		c.allowedHeaders = append(c.allowedHeaders, headers...)
	}
}

// CORSAddExtraExposedHeaders adds to default exposed headers
func CORSAddExtraExposedHeaders(headers ...string) CORSOption {
	return func(c *CORSConfig) {
		c.exposedHeaders = append(c.exposedHeaders, headers...)
	}
}

// CORSDisableCredentials disables the allowCredentials option
func CORSDisableCredentials() CORSOption {
	return func(c *CORSConfig) {
		c.allowCredentials = false
	}
}

// CORSMaxAge overrides the CORS max age
func CORSMaxAge(v int) CORSOption {
	return func(c *CORSConfig) {
		c.maxAge = v
	}
}
