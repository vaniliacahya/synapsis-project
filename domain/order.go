package domain

import (
	"github.com/gofiber/fiber/v2"
	"synapsis-project/database/databasesModel"
	"synapsis-project/structures/request"
	"synapsis-project/structures/response"
)

type OrderHandler interface {
	Order() fiber.Handler
}

type OrderUseCase interface {
	Order(param request.OrderRequest) response.LogicReturn[response.SummaryOrder]
}

type OrderData interface {
	InsertOrder(body databasesModel.Order) error
}
