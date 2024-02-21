package response

import "synapsis-project/database/databasesModel"

type ListCart struct {
	Count    int64                 `json:"count"`
	Products []databasesModel.Cart `json:"products"`
	Total    float64               `json:"total"`
}
