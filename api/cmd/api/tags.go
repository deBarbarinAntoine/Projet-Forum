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

type getTagsForm struct {
	Search string `form:"q"`
	data.Filters
	validator.Validator `form:"-"`
}

type tagByIDForm struct {
	ID                  int      `form:"-"`
	Includes            []string `form:"includes[]"`
	PermittedFields     []string `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) getTagsHandler(w http.ResponseWriter, r *http.Request) {

	form := newGetTagsForm()

	err := app.decodeForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if form.Page == 0 {
		form.Page = 1
	}
	if form.PageSize == 0 {
		form.PageSize = 100
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

	tags, metadata, err := app.models.Tags.Get(form.Search, form.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"_metadata": metadata, "tags": tags}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createTagHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name    string `json:"name"`
		Threads []int  `json:"threads"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.StringCheck(input.Name, 2, 125, true, "title")
	v.Check(validator.Unique(input.Threads), "threads", "duplicate values")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)
	var threads []data.Thread

	for _, id := range input.Threads {
		threads = append(threads, data.Thread{ID: id})
	}

	tag := &data.Tag{
		Name: input.Name,
		Author: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			ID:   user.ID,
			Name: user.Name,
		},
		Threads: threads,
	}

	err = app.models.Tags.Insert(tag)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEntry):
			v.AddError("threads", "thread already linked to this tag")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("threads", "thread not found")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"tag": tag}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getSingleTagHandler(w http.ResponseWriter, r *http.Request) {

	form := newTagByIDForm()

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

	tag, err := app.models.Tags.GetByID(form.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	if slices.Contains(form.Includes, "threads") {
		tag.Threads, err = app.models.Threads.GetByTag(tag.ID)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				break
			default:
				app.serverErrorResponse(w, r, err)
				return
			}
		}
	}
	if slices.Contains(form.Includes, "popularity") {
		tag.Popularity, err = app.models.Tags.GetPopularity(tag.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tag": tag}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateTagHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tag, err := app.models.Tags.GetByID(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)
	if !user.HasPermission(tag.Author.ID) {
		app.notPermittedResponse(w, r)
		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(tag.Version) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Name          *string `json:"name"`
		AddThreads    *[]int  `json:"add_threads"`
		RemoveThreads *[]int  `json:"remove_threads"`
	}
	var addThreads []data.Thread

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Name != nil {
		v.StringCheck(*input.Name, 2, 125, true, "name")
		tag.Name = *input.Name
	}
	if input.AddThreads != nil {
		for _, id := range *input.AddThreads {
			addThreads = append(addThreads, data.Thread{ID: id})
		}
		tag.Threads = addThreads
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Tags.Update(tag, *input.RemoveThreads)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateName):
			v.AddError("name", "a tag with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"tag": tag}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteTagHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	tag, err := app.models.Tags.GetByID(id)
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
	if !user.HasPermission(tag.Author.ID) {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.models.Tags.Delete(id)
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
		"message": fmt.Sprintf("deleted tag with id %d", id),
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) followTagHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	err = app.models.Tags.Follow(user, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrDuplicateEntry):
			v := validator.New()
			v.AddError("favorite", "you already follow this tag")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := envelope{
		"message": fmt.Sprintf("tag with id %d successfully added to your following list", id),
	}

	err = app.writeJSON(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) unfollowTagHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	err = app.models.Tags.Unfollow(user, id)
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
		"message": fmt.Sprintf("tag with id %d successfully removed from following list", id),
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
