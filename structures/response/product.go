package response

import (
	"synapsis-project/database/databasesModel"
)

type ListProduct struct {
	Count    int64                    `json:"count"`
	Products []databasesModel.Product `json:"products"`
}
