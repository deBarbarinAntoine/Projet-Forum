package main

import (
	"github.com/alexedwards/flow"
	"net/http"
)

func (app *application) routes() http.Handler {

	router := flow.New()

	/* #############################################################################
	/* # COMMON MIDDLEWARES
	/* ############################################################################# */

	router.Use(app.recoverPanic, app.enableCORS, app.rateLimit, app.authenticate)

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	/* #############################################################################
	/* # HEALTHCHECK
	/* ############################################################################# */

	router.HandleFunc("/v1/healthcheck", app.healthcheckHandler, http.MethodGet)

	/* #############################################################################
	/* # USERS
	/* ############################################################################# */

	router.HandleFunc("/v1/users", app.getUsersHandler, http.MethodGet)
	router.HandleFunc("/v1/users", app.registerUserHandler, http.MethodPost)
	router.HandleFunc("/v1/users/activated", app.activateUserHandler, http.MethodPut)

	router.HandleFunc("/v1/users/:id", app.getSingleUserHandler, http.MethodGet)
	router.HandleFunc("/v1/users/:id", app.updateUserHandler, http.MethodPut)
	router.HandleFunc("/v1/users/:id", app.deleteUserHandler, http.MethodDelete)

	router.HandleFunc("/v1/users/:id/friend", app.friendRequestHandler, http.MethodPost)
	router.HandleFunc("/v1/users/:id/friend", app.friendResponseHandler, http.MethodPut)
	router.HandleFunc("/v1/users/:id/friend", app.friendDeleteHandler, http.MethodDelete)

	/* #############################################################################
	/* # TOKENS
	/* ############################################################################# */

	router.HandleFunc("/v1/tokens/authentication", app.createAuthenticationTokenHandler, http.MethodPost)
	router.HandleFunc("/v1/tokens/refresh", app.refreshAuthenticationTokenHandler, http.MethodPost)

	/* #############################################################################
	/* # CATEGORIES
	/* ############################################################################# */

	router.HandleFunc("/v1/categories", app.getCategoriesHandler, http.MethodGet)
	router.HandleFunc("/v1/categories", app.createCategoryHandler, http.MethodPost)
	router.HandleFunc("/v1/categories/:id", app.getSingleCategoryHandler, http.MethodGet)
	router.HandleFunc("/v1/categories/:id", app.updateCategoryHandler, http.MethodPut)
	router.HandleFunc("/v1/categories/:id", app.deleteCategoryHandler, http.MethodDelete)

	/* #############################################################################
	/* # THREADS
	/* ############################################################################# */

	router.HandleFunc("/v1/threads", app.getThreadsHandler, http.MethodGet)
	router.HandleFunc("/v1/threads", app.createThreadHandler, http.MethodPost)
	router.HandleFunc("/v1/threads/:id", app.getSingleThreadHandler, http.MethodGet)
	router.HandleFunc("/v1/threads/:id", app.updateThreadHandler, http.MethodPut)
	router.HandleFunc("/v1/threads/:id", app.deleteThreadHandler, http.MethodDelete)

	router.HandleFunc("/v1/threads/:id/follow", app.followThreadHandler, http.MethodPost)
	router.HandleFunc("/v1/threads/:id/follow", app.unfollowThreadHandler, http.MethodDelete)

	/* #############################################################################
	/* # TAGS
	/* ############################################################################# */

	router.HandleFunc("/v1/tags", app.getTagsHandler, http.MethodGet)
	router.HandleFunc("/v1/tags", app.createTagHandler, http.MethodPost)
	router.HandleFunc("/v1/tags/:id", app.getSingleTagHandler, http.MethodGet)
	router.HandleFunc("/v1/tags/:id", app.updateTagHandler, http.MethodPut)
	router.HandleFunc("/v1/tags/:id", app.deleteTagHandler, http.MethodDelete)

	router.HandleFunc("/v1/tags/:id/favorite", app.addFavoriteTagHandler, http.MethodPost)
	router.HandleFunc("/v1/tags/:id/favorite", app.removeFavoriteTagHandler, http.MethodDelete)

	/* #############################################################################
	/* # POSTS
	/* ############################################################################# */

	router.HandleFunc("/v1/posts", app.getPostsHandler, http.MethodGet)
	router.HandleFunc("/v1/posts", app.createPostHandler, http.MethodPost)
	router.HandleFunc("/v1/posts/:id", app.getSinglePostHandler, http.MethodGet)
	router.HandleFunc("/v1/posts/:id", app.updatePostHandler, http.MethodPut)
	router.HandleFunc("/v1/posts/:id", app.deletePostHandler, http.MethodDelete)

	router.HandleFunc("/v1/posts/:id/react", app.reactToPostHandler, http.MethodPost)
	router.HandleFunc("/v1/posts/:id/react", app.changeReactionPostHandler, http.MethodPatch)
	router.HandleFunc("/v1/posts/:id/react", app.removeReactionPostHandler, http.MethodDelete)

	/* #############################################################################
	/* # DATA MANIPULATION
	/* ############################################################################# */

	router.HandleFunc("/v1/popular", app.getPopularHandler, http.MethodGet)
	router.HandleFunc("/v1/recommendations/:id", app.getRecommendations, http.MethodGet)
	router.HandleFunc("/v1/search", app.searchHandler, http.MethodGet)

	/* #############################################################################
	/* # DEBUG
	/* ############################################################################# */

	//router.Handle("/debug/vars", expvar.Handler(), http.MethodGet)

	return router
}
