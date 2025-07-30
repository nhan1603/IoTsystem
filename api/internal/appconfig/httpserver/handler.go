package httpserver

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func Handler(_ context.Context, routerFn func(r chi.Router)) http.Handler {
	r := chi.NewRouter()

	corsConf := NewCORSConfig(
		[]string{"*"},
		CORSAddExtraAllowedMethods(http.MethodPut, http.MethodDelete),
	)

	// Note: Cannot be used wihin a `Group`. See README.md https://github.com/go-chi/cors
	r.Use(cors.New(cors.Options{
		AllowedOrigins:   corsConf.allowedOrigins,
		AllowedMethods:   corsConf.allowedMethods,
		AllowedHeaders:   corsConf.allowedHeaders,
		ExposedHeaders:   corsConf.exposedHeaders,
		AllowCredentials: corsConf.allowCredentials,
		MaxAge:           corsConf.maxAge, // Maximum value not ignored by any of major browsers
	}).Handler)

	r.Group(routerFn)

	return r
}
