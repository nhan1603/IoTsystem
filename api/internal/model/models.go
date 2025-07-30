package model

import (
	"time"
)

// MenuItem represents a food item available for ordering
type MenuItem struct {
	ID          int64     `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Price       float64   `json:"price" db:"price"`
	Category    string    `json:"category" db:"category"`
	ImageUrl    string    `json:"imageUrl" db:"image_url"`
	IsAvailable bool      `json:"isAvailable" db:"is_available"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
}

// Order represents a food order placed by a user
type Order struct {
	ID          int64       `json:"id" db:"id"`
	UserID      int64       `json:"user_id" db:"user_id"`
	TotalAmount float64     `json:"total_amount" db:"total_amount"`
	Status      string      `json:"status" db:"status"`
	CreatedAt   time.Time   `json:"created_at" db:"created_at"`
	Items       []OrderItem `json:"items,omitempty" db:"-"`
}

// OrderItem represents an individual item within an order
type OrderItem struct {
	ID         int64    `json:"id" db:"id"`
	OrderID    int64    `json:"order_id" db:"order_id"`
	MenuItemID int64    `json:"menu_item_id" db:"menu_item_id"`
	Quantity   int      `json:"quantity" db:"quantity"`
	UnitPrice  float64  `json:"unit_price" db:"unit_price"`
	Subtotal   float64  `json:"subtotal" db:"subtotal"`
	Name       string   `json:"name" db:"name"`
	MenuItem   MenuItem `json:"menu_item,omitempty" db:"-"`
}

// PayPalTransaction represents a payment transaction through PayPal
type PayPalTransaction struct {
	ID                  int64     `json:"id" db:"id"`
	OrderID             int64     `json:"order_id" db:"order_id"`
	PayPalTransactionID string    `json:"paypal_transaction_id" db:"paypal_transaction_id"`
	PaymentStatus       string    `json:"payment_status" db:"payment_status"`
	PaymentAmount       string    `json:"payment_amount" db:"payment_amount"`
	Currency            string    `json:"currency" db:"currency"`
	PayerEmail          string    `json:"payer_email" db:"payer_email"`
	PaymentDate         time.Time `json:"payment_date" db:"payment_date"`
}

// OrderStatus constants
const (
	OrderStatusPending   = "pending"
	OrderStatusPaid      = "paid"
	OrderStatusPreparing = "preparing"
	OrderStatusReady     = "ready"
	OrderStatusCompleted = "completed"
	OrderStatusCancelled = "cancelled"
)

// PaymentStatus constants
const (
	PaymentStatusPending   = "pending"
	PaymentStatusCompleted = "completed"
	PaymentStatusFailed    = "failed"
	PaymentStatusRefunded  = "refunded"
)

// OrderRequest represents the input for creating a new order
type OrderRequest struct {
	UserID     int64              `json:"user_id" validate:"required"`
	Items      []OrderItemRequest `json:"items" validate:"required,dive"`
	PickupTime time.Time          `json:"pickup_time" validate:"required,future"`
}

// OrderItemRequest represents the input for creating order items
type OrderItemRequest struct {
	MenuItemID int64 `json:"menu_item_id" validate:"required"`
	Quantity   int   `json:"quantity" validate:"required,min=1"`
}

// PayPalPaymentRequest represents the input for initiating a PayPal payment
type PayPalPaymentRequest struct {
	OrderID int64   `json:"order_id" validate:"required"`
	Amount  float64 `json:"amount" validate:"required,gt=0"`
}

// PayPalPaymentResponse represents the response from PayPal payment initiation
type PayPalPaymentResponse struct {
	PaymentID    string `json:"payment_id"`
	ApprovalURL  string `json:"approval_url"`
	PaymentToken string `json:"payment_token"`
}
