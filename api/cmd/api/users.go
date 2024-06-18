package main

import (
	"ForumAPI/internal/data"
	"ForumAPI/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"time"
)

type userByIDForm struct {
	ID                  int      `form:"id"`
	Include             []string `form:"include[]"`
	validator.Validator `form:"-"`
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:   input.Name,
		Email:  input.Email,
		Status: data.StatusToConfirm,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if user.Validate(v); !v.Valid() {
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

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {

		data := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", data)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		TokenPlainText string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateTokenPlaintext(v, input.TokenPlainText); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetForToken(data.ScopeActivation, input.TokenPlainText)
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

	user.Status = data.StatusActivated

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.ScopeActivation, user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getUsersHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"users": "get_users"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getSingleUserHandler(w http.ResponseWriter, r *http.Request) {

	form := newUserByIDForm()
	err := app.decodeForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	form.Check(form.ID > 0, "id", "must contain a valid id")
	form.Check(validator.Unique(form.Include), "include[]", "duplicate values")
	for _, value := range form.Include {
		form.Check(validator.PermittedValue(value, validator.UserByIDValues...), "include[]", fmt.Sprintf("incorrect value %s", value))
	}

	if !form.Valid() {
		err = app.writeJSON(w, http.StatusBadRequest, envelope{"errors": form.Errors}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user, err := app.models.Users.GetByID(form.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if slices.Contains(form.Include, "following_tags") {
		user.FollowingTags, err = app.models.Tags.GetFollowingTagsByUserID(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Include, "favorite_threads") {
		user.FavoriteThreads, err = app.models.Threads.GetFavoriteThreadsByUserID(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Include, "categories_owned") {
		user.CategoriesOwned, err = app.models.Categories.GetOwnedCategoriesByUserID(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Include, "tags_owned") {
		user.TagsOwned, err = app.models.Tags.GetOwnedTagsByUserID(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Include, "threads_owned") {
		user.ThreadsOwned, err = app.models.Threads.GetOwnedThreadsByUserID(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Include, "posts") {
		user.Posts, err = app.models.Posts.GetPostsByAuthorID(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Include, "friends") {
		user.Friends, err = app.models.Users.GetFriendsByUserID(user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetByID(int(id))
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(user.Version) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Username             *string `json:"username"`
		Email                *string `json:"email"`
		Password             *string `json:"password"`
		NewPassword          *string `json:"new_password"`
		ConfirmationPassword *string `json:"confirmation_password"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Username != nil {
		user.Name = *input.Username
	}
	if input.Email != nil {
		user.Email = *input.Email
	}

	v := validator.New()

	if input.NewPassword != nil {
		err := user.Password.Set(*input.NewPassword)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		data.ValidateNewPassword(v, *input.NewPassword, *input.ConfirmationPassword)
	}

	if user.Validate(v); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
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

	err = app.writeJSON(w, http.StatusOK, envelope{"user": "updated_user"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"user": "user_removed"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
