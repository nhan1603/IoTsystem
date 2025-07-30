package operation

import (
	"github.com/nhan1603/IoTsystem/api/internal/controller/iot"
)

// Handler is the web handler for this pkg
type Handler struct {
	iotCtrl iot.Controller
}

// New instantiates a new Handler and returns it
func New(iotCtrl iot.Controller) Handler {
	return Handler{iotCtrl: iotCtrl}
}
