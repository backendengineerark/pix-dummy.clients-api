package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/backendengineerark/clients-api/configs"
	"github.com/backendengineerark/clients-api/internal/infra/webserver"
	"github.com/streadway/amqp"

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
		log.Printf("Fail to connect to database %s:%d", configs.DBHost, configs.DBPort)
		panic(err)
	}
	defer db.Close()
	log.Printf("Success to connect to database %s:%d", configs.DBHost, configs.DBPort)

	rabbitMQChannel, err := OpenRabbitMQChannel(configs)
	if err != nil {
		log.Printf("Fail to connect to rabbitmq %s:%d", configs.RabbitmqHost, configs.RabbitmqPort)
		panic(err)
	}
	log.Printf("Success to connect to rabbitmq %s:%d", configs.RabbitmqHost, configs.RabbitmqPort)

	webserver := webserver.NewWebServer(configs.AppPort, db, rabbitMQChannel)
	webserver.Start()
}

func OpenDatabaseConnection(configs *configs.Conf) (*sql.DB, error) {
	db, err := sql.Open(configs.DBDriver, fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName))
	if err != nil {
		return nil, err
	}
	return db, nil
}

func OpenRabbitMQChannel(configs *configs.Conf) (*amqp.Channel, error) {
	conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", configs.RabbitmqUser, configs.RabbitmqPassword, configs.RabbitmqHost, configs.RabbitmqPort))
	if err != nil {
		return nil, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}
	return ch, nil
}
