package main

import (
	"github.com/marcioaso/consult/api"
	"github.com/marcioaso/consult/config"
)

func main() {
	config.LoadConfig()
	port := config.AppConfig.AppPort
	api.SetupServer(port)
}
