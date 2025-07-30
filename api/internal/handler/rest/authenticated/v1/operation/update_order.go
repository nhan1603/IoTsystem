package operation

import (
	"net/http"

	"github.com/nhan1603/IoTsystem/api/internal/appconfig/httpserver"
)

// UpdateOrderResponse represents result of updating order
type UpdateOrderResponse struct {
	Success bool `json:"success"`
}

type UpdateOrderRequest struct {
	OrderID int    `json:"order_id"`
	Status  string `json:"status"`
}

func (h Handler) UpdateOrderStatus() http.HandlerFunc {
	return httpserver.HandlerErr(func(w http.ResponseWriter, r *http.Request) error {

		// var req UpdateOrderRequest
		// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
		// 	return err
		// }

		// if req.OrderID <= 0 || req.Status == "" {
		// 	return webErrInvalidOrder
		// }

		// err := h.orderCtrl.UpdateOrderStatus(r.Context(), req.OrderID, req.Status)
		// if err != nil {
		// 	return webInternalSerror
		// }

		// httpserver.RespondJSON(w, UpdateOrderResponse{
		// 	Success: true,
		// })

		return nil
	})
}
