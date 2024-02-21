package _interface

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"synapsis-project/interface/data"
	dispatch2 "synapsis-project/interface/dispatch"
	"synapsis-project/interface/usecase"
)

type InitData struct {
	App     *fiber.App
	DbMysql gorm.DB
	Redis   *redis.Client
}

func Init(d InitData) {
	dataInterface := data.New(d.DbMysql)
	useCaseInterface := usecase.New(dataInterface, d.Redis)
	handlerInterface := dispatch2.New(useCaseInterface, dataInterface)
	dispatch2.Routes(d.App, handlerInterface)
}
