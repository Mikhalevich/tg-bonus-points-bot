package request

type OrderStatus struct {
	Status string `json:"status" enum:"ready,completed" example:"ready" doc:"Order status"`
}
