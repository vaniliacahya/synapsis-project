package request

type AddCart struct {
	IdProduct string  `json:"id_product"`
	Qty       float64 `json:"qty"`
}

type AddCartRequest struct {
	AddCarts   AddCart `json:"add_carts"`
	IdCustomer string  `json:"id_customer" query:"id_customer"`
}

type DeleteCartRequest struct {
	Id         string `query:"id"`
	IdCustomer string `json:"id_customer"`
}
