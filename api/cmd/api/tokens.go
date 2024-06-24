package main

import (
	"ForumAPI/internal/data"
	"ForumAPI/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"time"
)

func (app *application) cleanExpiredTokens(frequency time.Duration) {
	defer func() {
		if err := recover(); err != nil {
			app.logger.Error(fmt.Sprintf("%v", err))
		}
	}()
	for {
		err := app.models.Tokens.DeleteExpired()
		if err != nil {
			app.logger.Error(err.Error())
		}
		time.Sleep(frequency)
	}
}

func (app *application) newAuthenticationToken(user *data.User) (envelope, error) {

	authToken, err := app.models.Tokens.New(user.ID, 24*time.Hour, data.TokenScope.Authentication)
	if err != nil {
		return nil, err
	}

	refreshToken, err := app.models.Tokens.New(user.ID, 48*time.Hour, data.TokenScope.Refresh)
	if err != nil {
		return nil, err
	}

	return envelope{
		"authentication_token": authToken,
		"refresh_token":        refreshToken,
	}, nil
}

func (app *application) createClientTokenHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name  string `json:"username"`
		Email string `json:"email"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:   input.Name,
		Email:  input.Email,
		Status: data.UserStatus.Client,
		Role:   data.UserRole.Client,
	}

	user.NoLogin()

	v := validator.New()

	data.ValidateEmail(v, user.Email)
	v.Check(user.Name != "", "name", "must be provided")
	v.Check(len(user.Name) > 2, "name", "must be more than 2 bytes long")
	v.Check(len(user.Name) <= 70, "name", "must not be more than 70 bytes long")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this email address already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrDuplicateUsername):
			v.AddError("username", "a user with this username already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	token, err := app.models.Tokens.New(user.ID, data.MaxDuration, data.TokenScope.Client)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := envelope{
		"user": envelope{
			"id":         user.ID,
			"username":   user.Name,
			"email":      user.Email,
			"created_at": user.CreatedAt,
		},
		"client_token": token,
	}

	err = app.writeJSON(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	data.ValidateEmail(v, input.Email)
	data.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.newAuthenticationToken(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, token, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) refreshAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		RefreshToken string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	user, err := app.models.Users.GetForToken(data.TokenScope.Refresh, input.RefreshToken)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.TokenScope.Refresh, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.newAuthenticationToken(user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, token, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
