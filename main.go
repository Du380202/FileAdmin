package main

import (
	"backend/config"
	_ "backend/router"
	"backend/setup"
	"log"
	"net/http"
	"time"
)

func main() {
	config.LoadConfig()
	config.ConnectMySQL()

	port := config.AppConfig.Server.Port
	server := &http.Server{
		Addr:         port,
		Handler:      setup.R,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	log.Fatal(server.ListenAndServe())
}
