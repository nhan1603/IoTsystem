package auth

import "github.com/nhan1603/IoTsystem/api/internal/controller/auth"

// Handler is the web handler for this pkg
type Handler struct {
	authCtrl auth.Controller
}

// New instantiates a new Handler and returns it
func New(authCtrl auth.Controller) Handler {
	return Handler{authCtrl: authCtrl}
}
