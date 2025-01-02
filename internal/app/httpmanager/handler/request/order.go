package request

type OrderStatus struct {
	Status string `json:"status" enum:"ready,completed,rejected" example:"ready" doc:"Order status"`
}
