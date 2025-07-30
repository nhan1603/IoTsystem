package auth

import (
	"log"
	"net/http"
	"strings"

	"github.com/nhan1603/IoTsystem/api/internal/appconfig/httpserver"
	"github.com/nhan1603/IoTsystem/api/internal/controller/auth"
)

// AuthRequest holds the input payload for authentication request
type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse represents result of authentication
type AuthResponse struct {
	UserID  int64  `json:"userID"`
	Token   string `json:"token"`
	Message string `json:"message"`
}

// AuthenticateOperationUser handles operator user authentication
func (h Handler) AuthenticateOperationUser() http.HandlerFunc {
	return httpserver.HandlerErr(func(w http.ResponseWriter, r *http.Request) error {
		log.Println("[AuthenticateOperationUser] START processing requests")
		var req AuthRequest

		if err := httpserver.ParseJSON(r.Body, &req); err != nil {
			log.Printf("[AuthenticateOperationUser] failed to parse json. Err: %+v\n", err)
			return err
		}

		if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" {
			return webErrInvalidEmailOrPassword
		}

		u, token, err := h.authCtrl.CheckAuth(r.Context(), auth.LoginInput{
			Email:    strings.ToLower(req.Email),
			Password: req.Password,
		})
		if err != nil {
			log.Printf("[AuthenticateOperationUser] failed to authenticate operation user. Err: %+v\n", err)
			return err
		}

		httpserver.RespondJSON(w, AuthResponse{
			UserID:  u.ID,
			Token:   token,
			Message: "authenticated successfully",
		})

		return nil
	})
}
