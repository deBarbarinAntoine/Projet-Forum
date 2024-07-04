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

type getCategoriesForm struct {
	Search string `form:"q"`
	data.Filters
	validator.Validator `form:"-"`
}

type categoryByIDForm struct {
	ID                  int      `form:"-"`
	Includes            []string `form:"includes[]"`
	PermittedFields     []string `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) getCategoriesHandler(w http.ResponseWriter, r *http.Request) {

	form := newGetCategoriesForm()

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

	categories, metadata, err := app.models.Categories.Get(form.Search, form.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"_metadata": metadata, "categories": categories}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name             string `json:"name"`
		ParentCategoryID int    `json:"parent_category_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.StringCheck(input.Name, 2, 70, true, "name")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)

	category := &data.Category{
		Name: input.Name,
		Author: struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		}{
			ID:   user.ID,
			Name: user.Name,
		},
		ParentCategory: struct {
			ID   int    `json:"id,omitempty"`
			Name string `json:"name,omitempty"`
		}{
			ID: input.ParentCategoryID,
		},
	}

	err = app.models.Categories.Insert(category)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateName):
			v.AddError("name", "a category with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getSingleCategoryHandler(w http.ResponseWriter, r *http.Request) {

	form := newCategoryByIDForm()

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

	category, err := app.models.Categories.GetByID(form.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {

			// DEBUG
			app.logger.Debug(fmt.Sprintf("error: %s", err.Error()))

			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	if slices.Contains(form.Includes, "categories") {
		category.Categories, err = app.models.Categories.GetByParentID(category.ID)
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
	if slices.Contains(form.Includes, "threads") {
		category.Threads, err = app.models.Threads.GetByCategory(category.ID)
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

	// DEBUG
	app.logger.Debug(fmt.Sprintf("category: %+v", category))

	err = app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category, err := app.models.Categories.GetByID(id)
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
	if !user.HasPermission(category.Author.ID) {
		app.notPermittedResponse(w, r)
		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(category.Version) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Name             *string `json:"name"`
		ParentCategoryID *int    `json:"parent_category_id"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Name != nil {

		v.StringCheck(*input.Name, 2, 70, true, "name")

		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}
		category.Name = *input.Name
	}
	if input.ParentCategoryID != nil {
		category.ParentCategory.ID = *input.ParentCategoryID
	}

	err = app.models.Categories.Update(&category)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		case errors.Is(err, data.ErrDuplicateName):
			v.AddError("name", "a category with this name already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	category, err := app.models.Categories.GetByID(id)
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
	if !user.HasPermission(category.Author.ID) {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.models.Categories.Delete(id)
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
		"message": fmt.Sprintf("deleted category with id %d", id),
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
