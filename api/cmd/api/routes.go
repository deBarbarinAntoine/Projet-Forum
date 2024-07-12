package main

import (
	"expvar"
	"github.com/alexedwards/flow"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := flow.New()

	/* #############################################################################
	/* # COMMON MIDDLEWARES
	/* ############################################################################# */

	router.Use(app.recoverPanic, app.enableCORS, app.rateLimit)

	/* #############################################################################
	/* # CLIENT TOKEN
	/* ############################################################################# */

	router.Group(func(group *flow.Mux) {
		group.Use(app.authenticateAPISecret)
		group.HandleFunc("/v1/tokens/client", app.createClientTokenHandler, http.MethodPost)
		group.HandleFunc("/v1/tokens/public-key", app.getPublicKeyPEM, http.MethodGet)
	})

	/* #############################################################################
	/* # BASIC ROUTES (WITH TOKEN HANDLING)
	/* ############################################################################# */

	router.Use(app.authenticateClient, app.authenticateUser)

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	/* #############################################################################
	/* # DEBUG (ONLY FOR DEVELOPMENT PHASE OR RESTRICT ACCESS THROUGH REVERSE PROXY)
	/* ############################################################################# */

	router.Handle("/debug/vars", expvar.Handler(), http.MethodGet)

	/* #############################################################################
	/* # HEALTHCHECK (OPTIONAL)
	/* ############################################################################# */

	router.HandleFunc("/v1/healthcheck", app.healthcheckHandler, http.MethodGet)

	/* #############################################################################
	/* # TOKENS
	/* ############################################################################# */

	router.HandleFunc("/v1/tokens/refresh", app.refreshAuthenticationTokenHandler, http.MethodPost)

	// ##################################
	// ENCRYPTED ROUTES
	// ##################################
	router.Group(func(group *flow.Mux) {
		group.Use(app.decryptRSA)

		group.HandleFunc("/v1/tokens/authentication", app.createAuthenticationTokenHandler, http.MethodPost)
	})

	// ##################################
	// PROTECTED ROUTES
	// ##################################
	router.Group(func(group *flow.Mux) {
		group.Use(app.requireActivatedUser, app.guardUserHandlers)

		group.HandleFunc("/v1/tokens/revoke/:id", app.revokeTokensHandler, http.MethodPost)
	})

	/* #############################################################################
	/* # USERS
	/* ############################################################################# */

	router.HandleFunc("/v1/users/activated", app.activateUserHandler, http.MethodPut)
	router.HandleFunc("/v1/users/forgot-password", app.forgotPasswordHandler, http.MethodPost)

	// ##################################
	// ENCRYPTED ROUTES
	// ##################################
	router.Group(func(group *flow.Mux) {
		group.Use(app.decryptRSA)

		group.HandleFunc("/v1/users", app.registerUserHandler, http.MethodPost)
		group.HandleFunc("/v1/users/password", app.resetPasswordHandler, http.MethodPut)
	})

	// ##################################
	// PROTECTED ROUTES
	// ##################################
	router.Group(func(group *flow.Mux) {
		group.Use(app.requireActivatedUser)

		group.HandleFunc("/v1/users", app.getUsersHandler, http.MethodGet)

		group.HandleFunc("/v1/users/:id", app.getSingleUserHandler, http.MethodGet)

		group.HandleFunc("/v1/users/:id/friend", app.friendRequestHandler, http.MethodPost)
		group.HandleFunc("/v1/users/:id/friend", app.friendResponseHandler, http.MethodPut)
		group.HandleFunc("/v1/users/:id/friend", app.friendDeleteHandler, http.MethodDelete)

		// CHECK PERMISSIONS FOR USER MANIPULATION
		group.Use(app.guardUserHandlers)
		group.HandleFunc("/v1/users/:id", app.deleteUserHandler, http.MethodDelete)

		// ENCRYPTED ROUTE
		group.Use(app.decryptRSA)
		group.HandleFunc("/v1/users/:id", app.updateUserHandler, http.MethodPut)

	})

	/* #############################################################################
	/* # CATEGORIES
	/* ############################################################################# */

	router.HandleFunc("/v1/categories", app.getCategoriesHandler, http.MethodGet)

	router.HandleFunc("/v1/categories/:id", app.getSingleCategoryHandler, http.MethodGet)

	// ##################################
	// PROTECTED ROUTES
	// ##################################
	router.Group(func(group *flow.Mux) {
		group.Use(app.requireActivatedUser)

		group.HandleFunc("/v1/categories", app.createCategoryHandler, http.MethodPost)

		group.HandleFunc("/v1/categories/:id", app.updateCategoryHandler, http.MethodPut)
		group.HandleFunc("/v1/categories/:id", app.deleteCategoryHandler, http.MethodDelete)
	})

	/* #############################################################################
	/* # THREADS
	/* ############################################################################# */

	router.HandleFunc("/v1/threads", app.getThreadsHandler, http.MethodGet)

	router.HandleFunc("/v1/threads/:id", app.getSingleThreadHandler, http.MethodGet)

	// ##################################
	// PROTECTED ROUTES
	// ##################################
	router.Group(func(group *flow.Mux) {
		group.Use(app.requireActivatedUser)

		group.HandleFunc("/v1/threads", app.createThreadHandler, http.MethodPost)

		group.HandleFunc("/v1/threads/:id", app.updateThreadHandler, http.MethodPut)
		group.HandleFunc("/v1/threads/:id", app.deleteThreadHandler, http.MethodDelete)

		group.HandleFunc("/v1/threads/:id/favorite", app.addToFavoritesThreadHandler, http.MethodPost)
		group.HandleFunc("/v1/threads/:id/favorite", app.removeFromFavoritesThreadHandler, http.MethodDelete)
	})

	/* #############################################################################
	/* # TAGS
	/* ############################################################################# */

	router.HandleFunc("/v1/tags", app.getTagsHandler, http.MethodGet)

	router.HandleFunc("/v1/tags/:id", app.getSingleTagHandler, http.MethodGet)

	// ##################################
	// PROTECTED ROUTES
	// ##################################
	router.Group(func(group *flow.Mux) {
		group.Use(app.requireActivatedUser)

		group.HandleFunc("/v1/tags", app.createTagHandler, http.MethodPost)

		group.HandleFunc("/v1/tags/:id", app.updateTagHandler, http.MethodPut)
		group.HandleFunc("/v1/tags/:id", app.deleteTagHandler, http.MethodDelete)

		group.HandleFunc("/v1/tags/:id/follow", app.followTagHandler, http.MethodPost)
		group.HandleFunc("/v1/tags/:id/follow", app.unfollowTagHandler, http.MethodDelete)
	})

	/* #############################################################################
	/* # POSTS
	/* ############################################################################# */

	router.HandleFunc("/v1/posts", app.getPostsHandler, http.MethodGet)

	router.HandleFunc("/v1/posts/:id", app.getSinglePostHandler, http.MethodGet)

	// ##################################
	// PROTECTED ROUTES
	// ##################################
	router.Group(func(group *flow.Mux) {
		group.Use(app.requireActivatedUser)

		group.HandleFunc("/v1/posts", app.createPostHandler, http.MethodPost)

		group.HandleFunc("/v1/posts/:id", app.updatePostHandler, http.MethodPut)
		group.HandleFunc("/v1/posts/:id", app.deletePostHandler, http.MethodDelete)

		group.HandleFunc("/v1/posts/:id/react", app.reactToPostHandler, http.MethodPost)
		group.HandleFunc("/v1/posts/:id/react", app.changeReactionPostHandler, http.MethodPatch)
		group.HandleFunc("/v1/posts/:id/react", app.removeReactionPostHandler, http.MethodDelete)
	})

	/* #############################################################################
	/* # DATA MANIPULATION
	/* ############################################################################# */

	router.HandleFunc("/v1/popular", app.getPopularHandler, http.MethodGet)
	//router.HandleFunc("/v1/recommendations/:id", app.getRecommendations, http.MethodGet)
	//router.HandleFunc("/v1/search", app.searchHandler, http.MethodGet)

	return router
}
