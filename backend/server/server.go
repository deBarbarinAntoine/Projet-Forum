package server

import (
	"Projet-Forum/internal/utils"
	"Projet-Forum/router"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

func Run() {
	// Loading environment variables
	err := godotenv.Load(utils.Path + ".envrc")
	if err != nil {
		log.Println("Error loading .envrc file", err)
	}

	// Initializing the routes
	router.Init()

	// Sending the assets to the clients
	fs := http.FileServer(http.Dir(utils.Path + "assets"))
	router.Mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Running the goroutine to automatically remove expired sessions every given time
	go utils.MonitorSessions()

	// Running the goroutine to change log file every given time
	go utils.LogInit()

	// Running the goroutine to automatically remove old TempUsers
	go utils.ManageTempUsers()

	// retrieving port from .envrc file
	port := os.Getenv("PORT")

	// Running the server
	log.Fatalln(http.ListenAndServe(port, router.Mux))
}
