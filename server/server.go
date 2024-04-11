package server

import (
	"Projet-Forum/internal/utils"
	"Projet-Forum/router"
	"log"
	"net/http"
)

func Run() {
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

	// Running the server
	log.Fatalln(http.ListenAndServe(":8080", router.Mux))
}
