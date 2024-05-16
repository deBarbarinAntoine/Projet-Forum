package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")

	ErrDuplicateCredential = errors.New("models: duplicate credential")

	ErrInvalidCredentials = errors.New("models: invalid credentials")
)
