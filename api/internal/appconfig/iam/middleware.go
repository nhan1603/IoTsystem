package iam

import (
	"context"
	"log"
	"net/http"

	"github.com/nhan1603/IoTsystem/api/internal/appconfig/httpserver"
)

// ErrAuthenticationFailed represents authentication failure
var ErrAuthenticationFailed = &httpserver.Error{Status: http.StatusUnauthorized, Code: "authentication_failed", Desc: "invalid access token"}

// AuthenticateUserMiddleware authenticates token
func AuthenticateUserMiddleware(rootCtx context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			strToken := r.Header.Get("Authorization")
			if strToken == "" {
				httpserver.RespondJSON(w, ErrAuthenticationFailed)
				return
			}

			iamCfg := ConfigFromContext(rootCtx)
			token, err := iamCfg.ValidateToken(strToken)
			if err != nil {
				log.Printf("Err when validating token: %+v\n", err)
				httpserver.RespondJSON(w, ErrAuthenticationFailed)
				return
			}

			// get user profile from token
			userProfile, err := GetUserProfileFomToken(token)
			if err != nil {
				log.Printf("Err when getting user profile: %+v\n", err)
				httpserver.RespondJSON(w, err)
				return
			}

			// set user profile to context
			ctx = SetUserProfileInContext(ctx, userProfile)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
