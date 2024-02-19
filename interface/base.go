package _interface

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"synapsis-project/interface/data"
	dispatch2 "synapsis-project/interface/dispatch"
	"synapsis-project/interface/usecase"
)

type InitData struct {
	App     *fiber.App
	DbMysql gorm.DB
}

func Init(d InitData) {
	dataInterface := data.New(d.DbMysql)
	useCaseInterface := usecase.New(dataInterface)
	handlerInterface := dispatch2.New(useCaseInterface, dataInterface)
	dispatch2.Routes(d.App, handlerInterface)
}
