package auth

import (
	"net/http"

	"github.com/nhan1603/IoTsystem/api/internal/appconfig/httpserver"
)

const (
	// ErrCodeValidationFailed represents the error code for a failed validation
	ErrCodeValidationFailed = "validation_failed"
)

// Web errors
var (
	webErrInvalidEmailOrPassword = &httpserver.Error{Status: http.StatusBadRequest, Code: ErrCodeValidationFailed, Desc: "invalid email / password"}
)
