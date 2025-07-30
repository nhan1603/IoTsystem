package operation

import (
	"net/http"

	"github.com/nhan1603/IoTsystem/api/internal/appconfig/httpserver"
	"github.com/nhan1603/IoTsystem/api/internal/model"
)

type MenuResponse struct {
	Items []model.MenuItem `json:"items"`
}

// Web errors
var (
	webInternalSerror = &httpserver.Error{Status: http.StatusInternalServerError, Code: "internal_error", Desc: "Something went wrong"}
)

func (h Handler) GetMenuItems() http.HandlerFunc {
	return httpserver.HandlerErr(func(w http.ResponseWriter, r *http.Request) error {

		// menuData, err := h.menuCtrl.GetAllItems(r.Context())
		// if err != nil {
		// 	return webInternalSerror
		// }

		// httpserver.RespondJSON(w, MenuResponse{
		// 	Items: menuData,
		// })

		return nil
	})
}
