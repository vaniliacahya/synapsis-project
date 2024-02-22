package response

import "synapsis-project/database/databasesModel"

type SummaryOrder struct {
	Order       databasesModel.Order  `json:"order"`
	Description string                `json:"description"`
	Carts       []databasesModel.Cart `json:"carts"`
}
