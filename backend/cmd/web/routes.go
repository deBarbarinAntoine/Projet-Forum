package main

import (
	"Projet-Forum/ui"
	"github.com/alexedwards/flow"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := flow.New()

	router.Handle("GET /static/", http.FileServerFS(ui.Files))

	router.Use(app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	router.Handle("GET /{$}", dynamic.ThenFunc(app.index))
	router.Handle("GET /about", dynamic.ThenFunc(app.about))
	router.Handle("GET /login", dynamic.ThenFunc(app.login))
	router.Handle("POST /login", dynamic.ThenFunc(app.loginPost))
	router.Handle("GET /register", dynamic.ThenFunc(app.register))
	router.Handle("POST /register", dynamic.ThenFunc(app.registerPost))
	router.Handle("GET /confirm", dynamic.ThenFunc(app.confirmHandler))
	router.Handle("GET /thread/{id}", dynamic.ThenFunc(app.getThread))
	router.Handle("GET /tag/{id}", dynamic.ThenFunc(app.getTag))
	router.Handle("GET /category/{id}", dynamic.ThenFunc(app.getCategory))
	router.Handle("GET /profile", dynamic.ThenFunc(app.getProfile))

	protected := dynamic.Append(app.requireAuthentication)

	router.Handle("GET /home", protected.ThenFunc(app.homeHandler))
	router.Handle("POST /logout", protected.ThenFunc(app.logoutPost))
	router.Handle("GET /post/{id}/create", dynamic.ThenFunc(app.createPost))
	router.Handle("GET /tag/{id}/create", dynamic.ThenFunc(app.createTag))
	router.Handle("GET /category/{id}/create", dynamic.ThenFunc(app.createCategory))
	router.Handle("GET /thread/{id}/create", dynamic.ThenFunc(app.createThread))
	router.Handle("POST /category", protected.ThenFunc(app.createCategoryPost))
	router.Handle("POST /thread", protected.ThenFunc(app.createThreadPost))
	router.Handle("POST /post", protected.ThenFunc(app.createPostPost))
	router.Handle("POST /tag", protected.ThenFunc(app.createTagPost))

	standard := alice.New(app.recoverPanic, app.logRequest, commonHeaders)

	return standard.Then(router)
}
