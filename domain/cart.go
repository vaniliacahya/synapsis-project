package domain

import (
	"github.com/gofiber/fiber/v2"
	"synapsis-project/database/databasesModel"
	"synapsis-project/structures/request"
	"synapsis-project/structures/response"
)

type CartHandler interface {
	AddCart() fiber.Handler
	ListCart() fiber.Handler
	DeleteCart() fiber.Handler
}

type CartUseCase interface {
	AddCart(request request.AddCartRequest) response.LogicReturn[response.ListCart]
	ListCart(param request.AddCartRequest) response.LogicReturn[response.ListCart]
	DeleteCart(param request.DeleteCartRequest) response.LogicReturn[response.ListCart]
}

type CartData interface {
	ListCart(param request.AddCartRequest) ([]databasesModel.Cart, int64, float64, error)
	AddCart(requestInsert databasesModel.Cart) error
	UpdateCart(requestInsert databasesModel.Cart) error
	DeleteCart(param request.DeleteCartRequest) error
}
