package data

import "errors"

var (
	ErrNoRecord = errors.New("data: no matching record found")

	ErrDuplicateCredential = errors.New("data: duplicate credential")

	ErrInvalidCredentials = errors.New("data: invalid credentials")
)
