package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type DBConfig struct {
	Driver   string
	Cluster  string
	Address  string
	Port     int
	Username string
	Password string
	Database string
}

var lock = &sync.Mutex{}
var dbConfig *DBConfig

func GetDBConfig() *DBConfig {
	lock.Lock()
	defer lock.Unlock()

	if dbConfig == nil {
		dbConfig = initDbMysqlConfig()
	}
	return dbConfig
}

func initDbMysqlConfig() *DBConfig {
	var configuration DBConfig
	err := godotenv.Load()
	if err != nil {
		log.Printf("Cannot read configuration, using system envxxx : %v", err.Error())
	}

	configuration.Database = os.Getenv("MYSQL_DB_DB")
	configuration.Username = os.Getenv("MYSQL_DB_USER")
	configuration.Password = os.Getenv("MYSQL_DB_PASS")
	configuration.Address = os.Getenv("MYSQL_DB_HOST")
	portConv, err := strconv.Atoi(os.Getenv("MYSQL_DB_PORT"))
	if err != nil {
		log.Println("Cannot parse DB Port variable, using default 3306")
		portConv = 3306
	}
	configuration.Port = portConv

	return &configuration
}
