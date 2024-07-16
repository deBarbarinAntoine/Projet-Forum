package main

type contextKey string

const (
	isAuthenticatedContextKey = contextKey("isAuthenticated")
	nonceContextKey           = contextKey("nonce")
)
