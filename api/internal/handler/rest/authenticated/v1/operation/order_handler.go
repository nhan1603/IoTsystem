package operation

import (
	"net/http"

	"github.com/nhan1603/IoTsystem/api/internal/appconfig/httpserver"
	"github.com/nhan1603/IoTsystem/api/internal/model"
)

// CreateOrderResponse represents result of creating order
type CreateOrderResponse struct {
	Success bool        `json:"success"`
	Data    model.Order `json:"data"`
}

type CreateOrderRequest struct {
	TotalAmount float64           `json:"total_amount"`
	Items       []model.OrderItem `json:"items,omitempty"`
}

func (h Handler) CreateOrder() http.HandlerFunc {
	return httpserver.HandlerErr(func(w http.ResponseWriter, r *http.Request) error {
		// ctx := r.Context()
		// var userData iam.HostProfile
		// ctxUserValue := ctx.Value(iam.UserProfileKey)
		// if ctxUserValue != nil {
		// 	userData = ctxUserValue.(iam.HostProfile)
		// } else {
		// 	return webErrInternalServer
		// }
		// userID := userData.ID

		// var req CreateOrderRequest
		// if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		// 	http.Error(w, "Invalid request body", http.StatusBadRequest)
		// 	return err
		// }

		// if req.TotalAmount <= 0 {
		// 	return webErrInvalidRequest
		// }

		// for _, item := range req.Items {
		// 	if item.ID < 0 || item.Quantity <= 0 || item.UnitPrice <= 0 {
		// 		return webErrInvalidRequest
		// 	}
		// }

		// orderData, err := h.orderCtrl.CreateOrder(ctx, model.Order{
		// 	UserID:      userID,
		// 	TotalAmount: req.TotalAmount,
		// 	Items:       req.Items,
		// })
		// if err != nil {
		// 	return webInternalSerror
		// }

		// httpserver.RespondJSON(w, CreateOrderResponse{
		// 	Success: true,
		// 	Data:    orderData,
		// })

		return nil
	})
}

// GetOrdersResponse represents result of getting all orders response
type GetOrdersResponse struct {
	Success bool          `json:"success"`
	Data    []model.Order `json:"data"`
}

func (h Handler) GetAllOrders() http.HandlerFunc {
	return httpserver.HandlerErr(func(w http.ResponseWriter, r *http.Request) error {
		// ctx := r.Context()
		// var userData iam.HostProfile
		// ctxUserValue := ctx.Value(iam.UserProfileKey)
		// if ctxUserValue != nil {
		// 	userData = ctxUserValue.(iam.HostProfile)
		// }
		// userID := userData.ID

		// orderData, err := h.orderCtrl.GetUserOrders(ctx, int(userID))
		// if err != nil {
		// 	return webInternalSerror
		// }

		// httpserver.RespondJSON(w, GetOrdersResponse{
		// 	Success: true,
		// 	Data:    orderData,
		// })

		return nil
	})
}
