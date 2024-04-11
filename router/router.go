package router

import (
	"Projet-Forum/controllers"
	"net/http"
)

var Mux = http.NewServeMux()

func Init() {
	Mux.HandleFunc("GET /{$}", controllers.IndexHandlerGetBundle)
	Mux.HandleFunc("POST /add", controllers.IndexHandlerPutBundle)
	Mux.HandleFunc("GET /login", controllers.LoginHandlerGetBundle)
	Mux.HandleFunc("POST /login", controllers.LoginHandlerPostBundle)
	Mux.HandleFunc("GET /register", controllers.RegisterHandlerGetBundle)
	Mux.HandleFunc("POST /register", controllers.RegisterHandlerPostBundle)
	Mux.HandleFunc("GET /home", controllers.HomeHandlerGetBundle)
	Mux.HandleFunc("GET /logs", controllers.LogHandlerGetBundle)
	Mux.HandleFunc("GET /confirm", controllers.ConfirmHandlerGetBundle)
	Mux.HandleFunc("GET /logout", controllers.LogoutHandlerGetBundle)

	// Handling StatusNotFound errors
	Mux.HandleFunc("/", controllers.ErrorHandlerBundle)
}
