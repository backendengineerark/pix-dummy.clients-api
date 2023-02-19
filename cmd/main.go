package main

import (
	"github.com/backendengineerark/clients-api/configs"
	"github.com/backendengineerark/clients-api/internal/infra/webserver"
)

func main() {
	configs, err := configs.LoadConfig("./")
	if err != nil {
		panic(err)
	}

	webserver := webserver.NewWebServer(configs.AppPort)
	webserver.Start()
}
