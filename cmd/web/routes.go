package main

import (
	"Projet-Forum/ui"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {

	mux := http.NewServeMux()

	mux.Handle("GET /static/", http.FileServerFS(ui.Files))

	dynamic := alice.New(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	mux.Handle("GET /{$}", dynamic.ThenFunc(app.index))
	mux.Handle("GET /about", dynamic.ThenFunc(app.about))
	mux.Handle("GET /login", dynamic.ThenFunc(app.login))
	mux.Handle("POST /login", dynamic.ThenFunc(app.loginPost))
	mux.Handle("GET /register", dynamic.ThenFunc(app.register))
	mux.Handle("POST /register", dynamic.ThenFunc(app.registerPost))
	mux.Handle("GET /confirm", dynamic.ThenFunc(app.confirmHandler))
	mux.Handle("GET /thread/{id}", dynamic.ThenFunc(app.getThread))
	mux.Handle("GET /tag/{id}", dynamic.ThenFunc(app.getTag))
	mux.Handle("GET /category/{id}", dynamic.ThenFunc(app.getCategory))
	mux.Handle("GET /profile", dynamic.ThenFunc(app.getProfile))

	protected := dynamic.Append(app.requireAuthentication)

	mux.Handle("GET /home", protected.ThenFunc(app.homeHandler))
	mux.Handle("POST /logout", protected.ThenFunc(app.logoutPost))
	mux.Handle("GET /post/{id}/create", dynamic.ThenFunc(app.createPost))
	mux.Handle("GET /tag/{id}/create", dynamic.ThenFunc(app.createTag))
	mux.Handle("GET /category/{id}/create", dynamic.ThenFunc(app.createCategory))
	mux.Handle("GET /thread/{id}/create", dynamic.ThenFunc(app.createThread))
	mux.Handle("POST /category", protected.ThenFunc(app.createCategoryPost))
	mux.Handle("POST /thread", protected.ThenFunc(app.createThreadPost))
	mux.Handle("POST /post", protected.ThenFunc(app.createPostPost))
	mux.Handle("POST /tag", protected.ThenFunc(app.createTagPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(mux)
}
