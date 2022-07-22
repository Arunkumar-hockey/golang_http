package model

import(
	"time"
)


type CashFree struct {
	Token           string      `json:"token"`
	Stage           string      `json:"stage"`
	AppID           string      `json:"appId"`
	OrderID         string      `json:"order_id"`
	OrderAmount     float32     `json:"order_amount"`
	OrderNote       string      `json:"order_note"`
	OrderCurrency   string      `json:"order_currency"`
	OrderCreatedAt  time.Time   `json:"order_createdat"`
	CustomerName    string      `json:"customer_name"`
	CustomerPhone   string      `json:"customer_phone"`
	CustomerEmail   string      `json:"customer_email"`
	PaymentMode     string      `json:"payment_mode"`
}