package auth

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/nhan1603/IoTsystem/api/internal/appconfig/httpserver"
	"github.com/nhan1603/IoTsystem/api/internal/model"
)

// CreateUserResponse represents result of creating user
type CreateUserResponse struct {
	Message string `json:"message"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Name     string `json:"username"`
	Password string `json:"password"`
}

func (h Handler) CreateUser() http.HandlerFunc {
	return httpserver.HandlerErr(func(w http.ResponseWriter, r *http.Request) error {
		var req CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return err
		}

		if strings.TrimSpace(req.Email) == "" || strings.TrimSpace(req.Password) == "" ||
			strings.TrimSpace(req.Name) == "" {
			return webErrInvalidEmailOrPassword
		}

		user := model.User{
			Email:    req.Email,
			Name:     req.Name,
			Password: req.Password,
		}

		err := h.authCtrl.CreateUser(r.Context(), &user)
		if err != nil {
			http.Error(w, "Failed to create user", http.StatusInternalServerError)
			return err
		}

		httpserver.RespondJSON(w, CreateUserResponse{
			Message: "User created successfully",
		})

		return nil
	})
}
