package main

import (
	"ebank/models"
	"ebank/routes"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

//var store = sessions.NewCookieStore([]byte("EindiaBusiness"))

func main() {

	//
	if err := godotenv.Load(".env"); err != nil {
		fmt.Printf("ENV not Found")
		return
	}
	dataSourceName := os.Getenv("Pqdetails") // Get database details value from .env File with import os

	models.InitDB(dataSourceName)

	router := routes.InitRoutes()

	port := (":" + os.Getenv("Port")) // Get Port value from .env File with import os
	fmt.Printf("Starting server on %s\n", port)
	if err := http.ListenAndServe(port, router); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
		os.Exit(1)
	}
}
