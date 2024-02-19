package registry

import (
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"os"
	"synapsis-project/helper"
	"synapsis-project/interface"
)

func InitRegistry(app *fiber.App, dbMysql gorm.DB, redis *redis.Client) {
	/**
	* DEFINE ENGINE
	*
	**/
	helper.LogrusDefiner()

	/**
	* DEFINE REGISTRY
	*
	**/

	_interface.Init(_interface.InitData{
		App:     app,
		DbMysql: dbMysql,
	})

	/**
	* RUN ENGINE
	*
	**/
	godotenv.Load()

	err := app.Listen(":" + os.Getenv("ENGINE_PORT"))
	if err != nil {
		panic(err)
		return
	}
}
