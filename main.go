package main

import (
	"backend/config"
	"backend/router"
	"log"
	"net/http"
	"time"
)

func main() {
	config.LoadConfig()

	r := router.SetupRouter(
		&router.SCPRouter{},
		&router.FileRouter{},
	)

	port := config.AppConfig.Server.Port
	server := &http.Server{
		Addr:         port,
		Handler:      r,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 10,
	}
	log.Fatal(server.ListenAndServe())
}
