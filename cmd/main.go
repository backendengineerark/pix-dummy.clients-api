package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/backendengineerark/clients-api/configs"
	"github.com/backendengineerark/clients-api/internal/infra/webserver"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig("./")
	if err != nil {
		panic(err)
	}

	db, err := OpenDatabaseConnection(configs)
	if err != nil {
		panic(err)
	}
	defer db.Close()
	log.Printf("Success to connect to database %s:%d", configs.DBHost, configs.DBPort)

	webserver := webserver.NewWebServer(configs.AppPort, db)
	webserver.Start()
}

func OpenDatabaseConnection(configs *configs.Conf) (*sql.DB, error) {
	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		return nil, err
	}
	return db, nil
}
