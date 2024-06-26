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
	ID                  int      `form:"-"`
	Includes            []string `form:"includes[]"`
	validator.Validator `form:"-"`
}

func (app *application) cleanExpiredUnactivatedUsers(frequency, timeout time.Duration) {
	defer func() {
		if err := recover(); err != nil {
			app.logger.Error(fmt.Sprintf("%v", err))
		}
	}()
	time.Sleep(timeout)
	for {
		err := app.models.Users.DeleteExpired()
		if err != nil {
			app.logger.Error(err.Error())
		}
		time.Sleep(frequency)
	}
}

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"username"`
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
		Status: data.UserStatus.ToConfirm,
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

	token, err := app.models.Tokens.New(user.ID, 3*24*time.Hour, data.TokenScope.Activation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {

		mailData := map[string]any{
			"activationToken": token.Plaintext,
			"userID":          user.ID,
		}

		err = app.mailer.Send(user.Email, "user_welcome.tmpl", mailData)
		if err != nil {
			app.logger.Error(err.Error())
		}
	})

	response := envelope{
		"user": envelope{
			"id":        user.ID,
			"name":      user.Name,
			"email":     user.Email,
			"createdAt": user.CreatedAt,
			"message":   fmt.Sprintf("a mail has been sent to %s", user.Email),
		},
	}

	err = app.writeJSON(w, http.StatusAccepted, response, nil)
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

	user, err := app.models.Users.GetForToken(data.TokenScope.Activation, input.TokenPlainText)
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

	err = app.models.Users.Activate(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(data.TokenScope.Activation, user.ID)
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

	// TODO -> for administration purposes

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
	form.ID, err = app.readIDParam(r)
	if err != nil {
		form.AddError("id", "must be an integer")
	}

	form.Check(form.ID > 0, "id", "must contain a valid id")
	form.Check(validator.Unique(form.Includes), "includes[]", "duplicate values")
	for _, value := range form.Includes {
		form.Check(validator.PermittedValue(value, validator.UserByIDValues...), "includes[]", fmt.Sprintf("incorrect value %s", value))
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
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	if slices.Contains(form.Includes, "following_tags") {
		user.FollowingTags, err = app.models.Tags.GetByFollowingUserID(user.ID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.notFoundResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Includes, "favorite_threads") {
		user.FavoriteThreads, err = app.models.Threads.GetFavoriteThreadsByUserID(user.ID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.notFoundResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Includes, "categories_owned") {
		user.CategoriesOwned, err = app.models.Categories.GetOwnedCategoriesByUserID(user.ID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.notFoundResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Includes, "tags_owned") {
		user.TagsOwned, err = app.models.Tags.GetByAuthorID(user.ID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.notFoundResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Includes, "threads_owned") {
		user.ThreadsOwned, err = app.models.Threads.GetOwnedThreadsByUserID(user.ID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.notFoundResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Includes, "posts") {
		user.Posts, err = app.models.Posts.GetByAuthorID(user.ID)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.notFoundResponse(w, r)
				return
			}
			app.serverErrorResponse(w, r, err)
			return
		}
	}
	if slices.Contains(form.Includes, "friends") {
		err = app.getFriendsByUser(user)
		if err != nil {
			if errors.Is(err, data.ErrRecordNotFound) {
				app.notFoundResponse(w, r)
				return
			}
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
		switch {
		case errors.Is(err, ErrUserIDNotFound):
			app.notFoundResponse(w, r)
		default:
			app.badRequestResponse(w, r, err)
		}
		return
	}

	user, err := app.models.Users.GetByID(id)
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
		Avatar               *string `json:"avatar"`
		Birth                *string `json:"birth"`
		Bio                  *string `json:"bio"`
		Signature            *string `json:"signature"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Username != nil {
		user.Name = *input.Username
	}
	if input.Email != nil {
		user.Email = *input.Email
	}
	if input.Avatar != nil {
		user.Avatar = *input.Avatar
	}
	if input.Birth != nil {
		birth, err := time.Parse("2006-01-02", *input.Birth)
		if err != nil {
			v.AddError("birth", "must be a valid date in the format YYYY-MM-DD")
		} else {
			user.BirthDate = birth
		}
	}
	if input.Bio != nil {
		user.Bio = *input.Bio
	}
	if input.Signature != nil {
		user.Signature = *input.Signature
	}

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

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	err = app.models.Tokens.DeleteAllForUser("*", id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Users.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := envelope{
		"message": fmt.Sprintf("user with id %d deleted", id),
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
