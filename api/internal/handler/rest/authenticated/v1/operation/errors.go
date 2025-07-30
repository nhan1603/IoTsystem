package operation

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
	webErrInvalidRequest = &httpserver.Error{Status: http.StatusBadRequest, Code: ErrCodeValidationFailed, Desc: "invalid request, empty cart"}
	webErrInvalidItem    = &httpserver.Error{Status: http.StatusBadRequest, Code: ErrCodeValidationFailed, Desc: "invalid data provided"}
	webErrInvalidOrder   = &httpserver.Error{Status: http.StatusBadRequest, Code: ErrCodeValidationFailed, Desc: "invalid request, invalid order"}
	webErrInternalServer = &httpserver.Error{Status: http.StatusInternalServerError, Code: "internal_error", Desc: "Something went wrong, please check again."}
)
