package main

import (
	"fmt"
	"synapsis-project/config"
	"synapsis-project/database"
	"synapsis-project/registry"
)

func main() {
	//DB config
	fmt.Println("Connect to database")

	dbCfgMysql := config.GetDBConfig()
	dbMysql := database.InitMySqlDB(dbCfgMysql)
	dbRedis := config.InitRedisConfig()

	app := config.GetServer()

	fmt.Println("Starting Server ...")
	registry.InitRegistry(app, *dbMysql, dbRedis)

}
