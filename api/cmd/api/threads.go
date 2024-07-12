package main

import (
	"ForumAPI/internal/data"
	"ForumAPI/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strconv"
)

type getThreadsForm struct {
	Search string `form:"q"`
	data.Filters
	validator.Validator `form:"-"`
}

type threadByIDForm struct {
	ID                  int      `form:"-"`
	Includes            []string `form:"includes[]"`
	PermittedFields     []string `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) getThreadsHandler(w http.ResponseWriter, r *http.Request) {

	form := newGetThreadsForm()

	err := app.decodeForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if form.Page == 0 {
		form.Page = 1
	}
	if form.PageSize == 0 {
		form.PageSize = 10
	}
	if form.Sort == "" {
		form.Sort = form.SortSafelist[0]
	}

	data.ValidateFilters(&form.Validator, form.Filters)

	if !form.Valid() {
		err = app.writeJSON(w, http.StatusBadRequest, envelope{"errors": form.Errors}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	threads, metadata, err := app.models.Threads.Get(form.Search, form.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"_metadata": metadata, "threads": threads}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createThreadHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		IsPublic    bool   `json:"is_public"`
		CategoryID  int    `json:"category_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.StringCheck(input.Title, 2, 125, true, "title")
	v.StringCheck(input.Description, 0, 1_020, false, "description")
	v.Check(input.CategoryID > 0, "category_id", "must be greater than zero")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)

	thread := &data.Thread{
		Title:       input.Title,
		Description: input.Description,
		IsPublic:    input.IsPublic,
		Author: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			ID:   user.ID,
			Name: user.Name,
		},
		Category: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			ID: input.CategoryID,
		},
	}

	err = app.models.Threads.Insert(thread)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateName):
			v.AddError("name", "a thread with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"thread": thread}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getSingleThreadHandler(w http.ResponseWriter, r *http.Request) {

	form := newThreadByIDForm()

	err := app.decodeForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	form.ID, err = app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	form.Check(validator.Unique(form.Includes), "includes[]", "duplicate values")
	for _, field := range form.Includes {
		form.Check(validator.PermittedValue(field, form.PermittedFields...), "includes[]", fmt.Sprintf("incorrect value %s", field))
	}

	if !form.Valid() {
		err = app.writeJSON(w, http.StatusBadRequest, envelope{"errors": form.Errors}, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	thread, err := app.models.Threads.GetByID(form.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}

		// DEBUG
		app.logger.Debug("app.models.Threads.GetByID(form.ID)")

		app.serverErrorResponse(w, r, err)
		return
	}

	if slices.Contains(form.Includes, "posts") {
		thread.Posts, err = app.models.Posts.GetByThread(thread.ID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				break
			default:

				// DEBUG
				app.logger.Debug("app.models.Posts.GetByThread(thread.ID)")

				app.serverErrorResponse(w, r, err)
				return
			}
		}
		err = app.models.Posts.GetReactions(thread.Posts)
		if err != nil {

			// DEBUG
			app.logger.Debug("app.models.Posts.GetReactions(thread.Posts)")

			app.serverErrorResponse(w, r, err)
			return
		}

		// DEBUG
		app.logger.Debug(fmt.Sprintf("Posts: %+v", thread.Posts))
	}
	if slices.Contains(form.Includes, "tags") {
		thread.Tags, err = app.models.Tags.GetByThread(thread.ID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				break
			default:

				// DEBUG
				app.logger.Debug("app.models.Tags.GetByThread(thread.ID)")

				app.serverErrorResponse(w, r, err)
				return
			}
		}
	}
	if slices.Contains(form.Includes, "popularity") {
		thread.Popularity, err = app.models.Threads.GetPopularity(thread.ID)
		if err != nil {

			// DEBUG
			app.logger.Debug("app.models.Threads.GetPopularity(thread.ID)")

			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"thread": thread}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateThreadHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	thread, err := app.models.Threads.GetByID(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)
	if !user.HasPermission(thread.Author.ID) {
		app.notPermittedResponse(w, r)
		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(thread.Version) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		IsPublic    *bool   `json:"is_public"`
		CategoryID  *int    `json:"category_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Title != nil {
		v.StringCheck(*input.Title, 2, 125, true, "title")
		thread.Title = *input.Title
	}
	if input.Description != nil {
		v.StringCheck(*input.Description, 0, 1_020, true, "description")
		thread.Description = *input.Description
	}
	if input.CategoryID != nil {
		v.Check(*input.CategoryID > 0, "category_id", "must be greater than zero")
		thread.Category.ID = *input.CategoryID
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Threads.Update(thread)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateTitle):
			v.AddError("title", "a thread with this title already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"thread": thread}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteThreadHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	thread, err := app.models.Threads.GetByID(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user := app.contextGetUser(r)
	if !user.HasPermission(thread.Author.ID) {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.models.Threads.Delete(id)
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
		"message": fmt.Sprintf("deleted thread with id %d", id),
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) addToFavoritesThreadHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	err = app.models.Threads.AddToFavorites(user, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrDuplicateEntry):
			v := validator.New()
			v.AddError("favorite", "this thread is already in your favorite list")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := envelope{
		"message": fmt.Sprintf("thread with id %d successfully added to favorites", id),
	}

	err = app.writeJSON(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) removeFromFavoritesThreadHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	err = app.models.Threads.RemoveFromFavorites(user, id)
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
		"message": fmt.Sprintf("thread with id %d successfully removed from favorites", id),
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
