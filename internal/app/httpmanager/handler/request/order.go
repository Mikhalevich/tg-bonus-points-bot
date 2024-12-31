package request

type OrderStatus struct {
	Status string `json:"status" example:"ready" doc:"Order status"`
}
