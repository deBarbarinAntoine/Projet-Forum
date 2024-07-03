package main

import (
	"Projet-Forum/internal/data"
	"Projet-Forum/internal/validator"
	"context"
	"fmt"
	"github.com/justinas/nosurf"
	"log/slog"
	"net/http"
	"time"
)

const (
	userTokenSessionManager           = "user_token"
	authenticatedUserIDSessionManager = "authenticated_user_id"
)

func commonHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Security-Policy", "default-src 'self' ui-avatars.com; script-src 'self' fonts.googleapis.com fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		w.Header().Set("Server", "Golang server")

		next.ServeHTTP(w, r)
	})
}

func (app *application) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			ip     = r.RemoteAddr
			proto  = r.Proto
			method = r.Method
			uri    = r.URL.RequestURI()
		)

		app.logger.Debug("received request", slog.String("ip", ip), slog.String("protocol", proto), slog.String("method", method), slog.String("URI", uri))

		next.ServeHTTP(w, r)
	})
}

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverError(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func noSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   true,
	})

	return csrfHandler
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// getting the userID from the session
		id := app.sessionManager.GetInt(r.Context(), authenticatedUserIDSessionManager)

		// if user not authenticated
		if id == 0 {
			next.ServeHTTP(w, r)
			return
		}

		// getting the user tokens from the session
		tokens, ok := app.sessionManager.Get(r.Context(), userTokenSessionManager).(*data.Tokens)
		if ok {

			// checking the authentication imminent expiry
			if time.Until(tokens.Authentication.Expiry) < 1*time.Hour {

				// checking the refresh token validity
				if time.Until(tokens.Refresh.Expiry) > 0 {

					// request new tokens from API with refresh token
					v := validator.New()
					err := app.models.TokenModel.Refresh(tokens.Refresh.Token, tokens, v)
					if err != nil {
						app.serverError(w, r, err)
						return
					}

					// if refresh token invalid -> logout
					if !v.Valid() {
						err = app.logout(r)
						if err != nil {
							app.serverError(w, r, err)
							return
						}
						next.ServeHTTP(w, r)
						return
					}

					// replacing the tokens in the user session
					app.sessionManager.Put(r.Context(), userTokenSessionManager, tokens)
				} else {

					// if refresh token expired -> logout
					err := app.logout(r)
					if err != nil {
						app.serverError(w, r, err)
						return
					}
					next.ServeHTTP(w, r)
					return
				}
			}

			// setting the user as authenticated in the context
			ctx := context.WithValue(r.Context(), isAuthenticatedContextKey, true)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) requireAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if !app.isAuthenticated(r) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		w.Header().Add("Cache-Control", "no-store")

		next.ServeHTTP(w, r)
	})
}
