package request

type ListProductRequest struct {
	IdCustomer        string   `query:"id_customer" json:"id_customer"`
	IdProductCategory []string `query:"id_product_category" json:"id_product_category"`
	Limit             int      `query:"limit" json:"limit"`
	Offset            int      `query:"offset" json:"offset"`
	IdProduct         []string `json:"id_product"`
}
