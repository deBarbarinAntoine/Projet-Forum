package main

import (
	"Projet-Forum/ui"
	"github.com/alexedwards/flow"
	"io/fs"
	"net/http"
)

func (app *application) routes() http.Handler {

	// setting the files to put in the static handler
	staticFs, err := fs.Sub(ui.StaticFiles, "assets")
	if err != nil {
		panic(err)
	}

	router := flow.New()

	router.NotFound = http.HandlerFunc(app.notFound)                 // error 404 page
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowed) // error 405 page

	router.Handle("/static/...", http.StripPrefix("/static/", http.FileServerFS(staticFs)), http.MethodGet) // static files

	router.Use(app.recoverPanic, app.logRequest, commonHeaders, app.sessionManager.LoadAndSave, noSurf, app.authenticate)

	/* #############################################################################
	/*	COMMON
	/* #############################################################################*/

	router.HandleFunc("/", app.index, http.MethodGet)      // landing page
	router.HandleFunc("/home", app.index, http.MethodGet)  // landing page
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

	router.HandleFunc("/dashboard", app.dashboard, http.MethodGet) // dashboard page
	router.HandleFunc("/logout", app.logoutPost, http.MethodPost)  // logout route
	router.HandleFunc("/user", app.updateUser, http.MethodGet)     // update user page
	router.HandleFunc("/user", app.updateUserPut, http.MethodPut)  // update user treatment route

	router.HandleFunc("/post/create", app.createPost, http.MethodGet)         // post creation page
	router.HandleFunc("/tag/create", app.createTag, http.MethodGet)           // tag creation page
	router.HandleFunc("/category/create", app.createCategory, http.MethodGet) // category creation page
	router.HandleFunc("/thread/create", app.createThread, http.MethodGet)     // thread creation page

	router.HandleFunc("/category", app.createCategoryPost, http.MethodPost) // post creation treatment route
	router.HandleFunc("/thread", app.createThreadPost, http.MethodPost)     // tag creation treatment route
	router.HandleFunc("/post", app.createPostPost, http.MethodPost)         // category creation treatment route
	router.HandleFunc("/tag", app.createTagPost, http.MethodPost)           // thread creation treatment route

	router.HandleFunc("/post/:id/update", app.updatePost, http.MethodGet)         // post update page
	router.HandleFunc("/tag/:id/update", app.updateTag, http.MethodGet)           // tag update page
	router.HandleFunc("/category/:id/update", app.updateCategory, http.MethodGet) // category update page
	router.HandleFunc("/thread/:id/update", app.updateThread, http.MethodGet)     // thread update page

	router.HandleFunc("/category/:id/update", app.updateCategoryPut, http.MethodPut) // post update treatment route
	router.HandleFunc("/thread/:id/update", app.updateThreadPut, http.MethodPut)     // tag update treatment route
	router.HandleFunc("/post/:id/update", app.updatePostPut, http.MethodPut)         // category update treatment route
	router.HandleFunc("/tag/:id/update", app.updateTagPut, http.MethodPut)           // thread update treatment route

	/* #############################################################################
	/*	AJAX endpoints
	/* #############################################################################*/

	// Post reactions
	router.HandleFunc("/posts/:id/react", app.reactToPost, http.MethodPost)
	router.HandleFunc("/posts/:id/react", app.changeReactionPost, http.MethodPatch)
	router.HandleFunc("/posts/:id/react", app.removeReactionPost, http.MethodDelete)

	// Tag follow
	router.HandleFunc("/tags/:id/follow", app.followTag, http.MethodPost)
	router.HandleFunc("/tags/:id/follow", app.unfollowTag, http.MethodDelete)

	// Thread favorite
	router.HandleFunc("/threads/:id/favorite", app.addToFavoritesThread, http.MethodPost)
	router.HandleFunc("/threads/:id/favorite", app.removeFromFavoritesThread, http.MethodDelete)

	// Friends
	router.HandleFunc("/users/:id/friend", app.friendRequest, http.MethodPost)
	router.HandleFunc("/users/:id/friend", app.friendResponse, http.MethodPut)
	router.HandleFunc("/users/:id/friend", app.friendDelete, http.MethodDelete)

	return router
}
