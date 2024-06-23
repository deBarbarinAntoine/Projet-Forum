package main

import (
	"ForumAPI/internal/data"
	"ForumAPI/internal/validator"
	"net/http"
)

type getCategoriesForm struct {
	Search string `form:"q"`
	data.Filters
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

	v.Check(input.Name != "", "name", "must be provided")
	v.Check(len(input.Name) <= 70, "name", "must not be longer than 70 bytes")
	v.Check(len(input.Name) > 2, "name", "must be longer than 2 bytes")

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

	err = app.writeJSON(w, http.StatusCreated, envelope{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getSingleCategoryHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"category": "get_single_category"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"category": "update_category"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteCategoryHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"category": "delete_category"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
