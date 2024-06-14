package main

import (
	"expvar"
	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)

	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	/* #############################################################################
	/* # HEALTHCHECK
	/* ############################################################################# */

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	/* #############################################################################
	/* # USERS
	/* ############################################################################# */

	router.HandlerFunc(http.MethodGet, "/v1/users", app.getUsersHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.getSingleUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/:id", app.updateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.deleteUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/users/:id/friend", app.friendRequestHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/:id/friend", app.friendResponseHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id/friend", app.friendDeleteHandler)

	/* #############################################################################
	/* # TOKENS
	/* ############################################################################# */

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tokens/refresh", app.refreshAuthenticationTokenHandler)

	/* #############################################################################
	/* # CATEGORIES
	/* ############################################################################# */

	router.HandlerFunc(http.MethodGet, "/v1/categories", app.getCategoriesHandler)
	router.HandlerFunc(http.MethodPost, "/v1/categories", app.createCategoryHandler)
	router.HandlerFunc(http.MethodGet, "/v1/categories/:id", app.getSingleCategoryHandler)
	router.HandlerFunc(http.MethodPut, "/v1/categories/:id", app.updateCategoryHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/categories/:id", app.deleteCategoryHandler)

	/* #############################################################################
	/* # THREADS
	/* ############################################################################# */

	router.HandlerFunc(http.MethodGet, "/v1/threads", app.getThreadsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/threads", app.createThreadHandler)
	router.HandlerFunc(http.MethodGet, "/v1/threads/:id", app.getSingleThreadHandler)
	router.HandlerFunc(http.MethodPut, "/v1/threads/:id", app.updateThreadHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/threads/:id", app.deleteThreadHandler)

	router.HandlerFunc(http.MethodPost, "/v1/threads/:id/follow", app.followThreadHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/threads/:id/follow", app.unfollowThreadHandler)

	/* #############################################################################
	/* # TAGS
	/* ############################################################################# */

	router.HandlerFunc(http.MethodGet, "/v1/tags", app.getTagsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/tags", app.createTagHandler)
	router.HandlerFunc(http.MethodGet, "/v1/tags/:id", app.getSingleTagHandler)
	router.HandlerFunc(http.MethodPut, "/v1/tags/:id", app.updateTagHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/tags/:id", app.deleteTagHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tags/:id/favorite", app.addFavoriteTagHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/tags/:id/favorite", app.removeFavoriteTagHandler)

	/* #############################################################################
	/* # POSTS
	/* ############################################################################# */

	router.HandlerFunc(http.MethodGet, "/v1/posts", app.getPostsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/posts", app.createPostHandler)
	router.HandlerFunc(http.MethodGet, "/v1/posts/:id", app.getSinglePostHandler)
	router.HandlerFunc(http.MethodPut, "/v1/posts/:id", app.updatePostHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/posts/:id", app.deletePostHandler)

	router.HandlerFunc(http.MethodPost, "/v1/posts/:id/react", app.reactToPostHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/posts/:id/react", app.changeReactionPostHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/posts/:id/react", app.removeReactionPostHandler)

	/* #############################################################################
	/* # DATA MANIPULATION
	/* ############################################################################# */

	router.HandlerFunc(http.MethodGet, "/v1/popular", app.getPopularHandler)
	router.HandlerFunc(http.MethodGet, "/v1/recommendations/:id", app.getRecommendations)
	router.HandlerFunc(http.MethodGet, "/v1/search", app.searchHandler)

	/* #############################################################################
	/* # DEBUG
	/* ############################################################################# */

	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	/* #############################################################################
	/* # COMMON MIDDLEWARES
	/* ############################################################################# */

	basic := alice.New(app.metrics, app.recoverPanic, app.enableCORS, app.rateLimit, app.authenticate)

	return basic.Then(router)
}
