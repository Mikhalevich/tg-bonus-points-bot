package button

type Operation string

const (
	OperationOrderCancel              Operation = "OrderCancel"
	OperationCartCancel               Operation = "CartCancel"
	OperationCartConfirm              Operation = "CartConfirm"
	OperationCartViewCategoryProducts Operation = "CartViewCategoryProducts"
	OperationCartViewCategories       Operation = "CartViewCategories"
	OperationCartAddProduct           Operation = "CartAddProduct"
)
