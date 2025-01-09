package button

type Operation string

const (
	OperationCancelOrder  Operation = "CancelOrder"
	OperationConfirmOrder Operation = "ConfirmOrder"
	OperationViewCategory Operation = "ViewCategory"
	OperationProduct      Operation = "Product"
	OperationBackToOrder  Operation = "BackToOrder"
)
