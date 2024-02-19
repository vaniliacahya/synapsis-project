package domain

import (
	"github.com/gofiber/fiber/v2"
	"synapsis-project/database/databasesModel"
	"synapsis-project/structures/request"
	"synapsis-project/structures/response"
)

type ProductHandler interface {
	ListProduct() fiber.Handler
}

type ProductUseCase interface {
	ListProduct(param request.ListProductRequest) response.LogicReturn[response.ListProduct]
}

type ProductData interface {
	ListProduct(param request.ListProductRequest) ([]databasesModel.Product, int64, error)
}
