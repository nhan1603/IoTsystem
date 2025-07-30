package iam

import (
	"context"
	"errors"
	"os"
)

// Config represents config of iam aka user, includes config config
type Config struct {
	JWTKey string
}

// NewConfig returns iam config
func NewConfig() Config {
	return Config{
		JWTKey: os.Getenv("JWT_KEY"),
	}
}

// Validate validates iam config
func (c Config) Validate() error {
	if c.JWTKey == "" {
		return errors.New("required env variable 'JWT_KEY' not found")
	}

	return nil
}

// IAmConfigKey is key for getter/setter value from/to context
var IAmConfigKey = "iam_config"

// SetConfigToContext sets the iam config in the given context
func SetConfigToContext(ctx context.Context, config Config) context.Context {
	return context.WithValue(ctx, IAmConfigKey, config)
}

// ConfigFromContext gets the iam config from the given context
func ConfigFromContext(ctx context.Context) Config {
	if v, ok := ctx.Value(IAmConfigKey).(Config); ok {

		return v
	}
	return Config{}
}
