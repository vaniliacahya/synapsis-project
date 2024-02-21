package domain

import (
	"github.com/gofiber/fiber/v2"
	"synapsis-project/database/databasesModel"
	"synapsis-project/structures/response"
)

type CustomerHandler interface {
	Register() fiber.Handler
	Login() fiber.Handler
}

type CustomerCase interface {
	AddCustomer(body databasesModel.Customer) response.LogicReturn[databasesModel.Customer]
	Login(body databasesModel.Customer) response.LogicReturn[response.LoginResponse]
}

type CustomerData interface {
	AddCustomer(body databasesModel.Customer) error
	GetCustomer(body databasesModel.Customer) (databasesModel.Customer, error)
}
