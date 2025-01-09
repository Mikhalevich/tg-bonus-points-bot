package button

type Operation string

const (
	OperationCancelOrderSendMessage Operation = "CancelOrderSendMessage"
	OperationCancelOrderEditMessage Operation = "CancelOrderEditMessage"
	OperationConfirmOrder           Operation = "ConfirmOrder"
	OperationViewCategory           Operation = "ViewCategory"
	OperationProduct                Operation = "Product"
	OperationBackToOrder            Operation = "BackToOrder"
)
