package button

type Operation string

const (
	OperationCancelOrder  Operation = "CancelOrder"
	OperationCancelCart   Operation = "CancelCart"
	OperationConfirmCart  Operation = "ConfirmCart"
	OperationViewCategory Operation = "ViewCategory"
	OperationProduct      Operation = "Product"
	OperationBackToOrder  Operation = "BackToOrder"
)
