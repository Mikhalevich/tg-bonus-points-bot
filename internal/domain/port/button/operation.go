package button

type Operation string

const (
	OperationOrderCancel              Operation = "OrderCancel"
	OperationOrderHistoryPrevious     Operation = "OperationOrderHistoryPrevious"
	OperationOrderHistoryNext         Operation = "OperationOrderHistoryNext"
	OperationOrderHistoryFirst        Operation = "OperationOrderHistoryFirst"
	OperationOrderHistoryLast         Operation = "OperationOrderHistoryLast"
	OperationCartCancel               Operation = "CartCancel"
	OperationCartConfirm              Operation = "CartConfirm"
	OperationCartViewCategoryProducts Operation = "CartViewCategoryProducts"
	OperationCartViewCategories       Operation = "CartViewCategories"
	OperationCartAddProduct           Operation = "CartAddProduct"
)
