package main

import (
	"Projet-Forum/ui"
	"github.com/alexedwards/flow"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := flow.New()

	router.Use(app.recoverPanic, app.logRequest, commonHeaders, app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	/* #############################################################################
	/*	COMMON
	/* #############################################################################*/

	router.NotFound = http.HandlerFunc(app.notFound) // error 404 page

	router.Handle("/static/...", http.StripPrefix("/static/", http.FileServerFS(ui.Files)), http.MethodGet) // static files

	router.HandleFunc("/", app.index, http.MethodGet)      // landing page
	router.HandleFunc("/about", app.about, http.MethodGet) // about page

	router.HandleFunc("/thread/:id", app.threadGet, http.MethodGet)     // thread page
	router.HandleFunc("/tag/:id", app.tagGet, http.MethodGet)           // tag page
	router.HandleFunc("/category/:id", app.categoryGet, http.MethodGet) // category page

	router.HandleFunc("/tags", app.TagsGet, http.MethodGet)             // all tags page
	router.HandleFunc("/categories", app.categoriesGet, http.MethodGet) // all categories page

	router.HandleFunc("/search", app.search, http.MethodGet) // search page

	/* #############################################################################
	/*	USER ACCESS
	/* #############################################################################*/

	router.HandleFunc("/login", app.login, http.MethodGet)      // login page
	router.HandleFunc("/login", app.loginPost, http.MethodPost) // login treatment route

	router.HandleFunc("/register", app.register, http.MethodGet)      // register page
	router.HandleFunc("/register", app.registerPost, http.MethodPost) // register treatment route

	router.HandleFunc("/confirm/:token", app.confirm, http.MethodGet) // confirmation page
	router.HandleFunc("/confirm", app.confirmPost, http.MethodPost)   // confirmation treatment route

	router.HandleFunc("/forgot-password", app.forgotPassword, http.MethodGet)      // forgot password page
	router.HandleFunc("/forgot-password", app.forgotPasswordPost, http.MethodPost) // forgot password treatment route

	router.HandleFunc("/reset-password/:token", app.resetPassword, http.MethodGet) // reset password page
	router.HandleFunc("/reset-password", app.resetPasswordPost, http.MethodPost)   // reset password treatment route

	/* #############################################################################
	/*	RESTRICTED
	/* #############################################################################*/

	router.Use(app.requireAuthentication)

	router.HandleFunc("/dashboard", app.dashboard, http.MethodGet)  // dashboard page
	router.HandleFunc("/logout", app.logoutPost, http.MethodPost)   // logout route
	router.HandleFunc("/user", app.updateUser, http.MethodGet)      // update user page
	router.HandleFunc("/user", app.updateUserPost, http.MethodPost) // update user treatment route

	router.HandleFunc("/post/:id/create", app.createPost, http.MethodGet)         // post creation page
	router.HandleFunc("/tag/:id/create", app.createTag, http.MethodGet)           // tag creation page
	router.HandleFunc("/category/:id/create", app.createCategory, http.MethodGet) // category creation page
	router.HandleFunc("/thread/:id/create", app.createThread, http.MethodGet)     // thread creation page

	router.HandleFunc("/category/:id/create", app.createCategoryPost, http.MethodPost) // post creation treatment route
	router.HandleFunc("/thread/:id/create", app.createThreadPost, http.MethodPost)     // tag creation treatment route
	router.HandleFunc("/post/:id/create", app.createPostPost, http.MethodPost)         // category creation treatment route
	router.HandleFunc("/tag/:id/create", app.createTagPost, http.MethodPost)           // thread creation treatment route

	router.HandleFunc("/post/:id/update", app.updatePost, http.MethodGet)         // post update page
	router.HandleFunc("/tag/:id/update", app.updateTag, http.MethodGet)           // tag update page
	router.HandleFunc("/category/:id/update", app.updateCategory, http.MethodGet) // category update page
	router.HandleFunc("/thread/:id/update", app.updateThread, http.MethodGet)     // thread update page

	router.HandleFunc("/category/:id/update", app.updateCategoryPut, http.MethodPut) // post update treatment route
	router.HandleFunc("/thread/:id/update", app.updateThreadPut, http.MethodPut)     // tag update treatment route
	router.HandleFunc("/post/:id/update", app.updatePostPut, http.MethodPut)         // category update treatment route
	router.HandleFunc("/tag/:id/update", app.updateTagPut, http.MethodPut)           // thread update treatment route

	return router
}
