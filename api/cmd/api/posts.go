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

type getPostsForm struct {
	Search string `form:"q"`
	data.Filters
	validator.Validator `form:"-"`
}

type postByIDForm struct {
	ID                  int      `form:"-"`
	Includes            []string `form:"includes[]"`
	PermittedFields     []string `form:"-"`
	validator.Validator `form:"-"`
}

func (app *application) getPostsHandler(w http.ResponseWriter, r *http.Request) {

	form := newGetPostsForm()

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

	posts, metadata, err := app.models.Posts.Get(form.Search, form.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Posts.GetReactions(posts)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"_metadata": metadata, "posts": posts}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createPostHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Content      *string `json:"content"`
		ThreadID     *int    `json:"thread_id"`
		ParentPostID *int    `json:"parent_post_id"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Content == nil {
		v.AddError("content", "must be provided")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	v.StringCheck(*input.Content, 1, 1_020, true, "content")

	if input.ThreadID == nil {
		v.AddError("thread", "must be provided")
		app.failedValidationResponse(w, r, v.Errors)
		return
	}
	v.Check(*input.ThreadID > 0, "thread", "must be greater than zero")

	if input.ParentPostID != nil {
		v.Check(*input.ParentPostID > 0, "parent_post_id", "must be greater than zero")
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user := app.contextGetUser(r)

	post := &data.Post{
		Content: *input.Content,
		Author: data.User{
			ID:   user.ID,
			Name: user.Name,
		},
		Thread: data.Thread{
			ID: *input.ThreadID,
		},
	}

	if input.ParentPostID != nil {
		post.IDParentPost = *input.ParentPostID
	}

	err = app.models.Posts.Insert(post)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// DEBUG
	app.logger.Debug(fmt.Sprintf("created post: %+v", post))

	err = app.writeJSON(w, http.StatusCreated, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getSinglePostHandler(w http.ResponseWriter, r *http.Request) {

	form := newPostByIDForm()

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

	post, err := app.models.Posts.GetByID(form.ID)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	if slices.Contains(form.Includes, "popularity") || slices.Contains(form.Includes, "reactions") {
		posts := []*data.Post{post}
		err = app.models.Posts.GetReactions(posts)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
		post = posts[0]
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updatePostHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post, err := app.models.Posts.GetByID(id)
	if err != nil {
		if errors.Is(err, data.ErrRecordNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)
	if !user.HasPermission(post.Author.ID) {
		app.notPermittedResponse(w, r)
		return
	}

	if r.Header.Get("X-Expected-Version") != "" {
		if strconv.Itoa(post.Version) != r.Header.Get("X-Expected-Version") {
			app.editConflictResponse(w, r)
			return
		}
	}

	var input struct {
		Content *string `json:"content"`
		Thread  *int    `json:"thread"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if input.Content != nil {
		v.StringCheck(*input.Content, 1, 1_020, true, "name")
		post.Content = *input.Content
	}
	if input.Thread != nil {
		post.Thread.ID = *input.Thread
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Posts.Update(*post)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"post": post}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deletePostHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	post, err := app.models.Posts.GetByID(id)
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
	if !user.HasPermission(post.Author.ID) {
		app.notPermittedResponse(w, r)
		return
	}

	err = app.models.Posts.Delete(id)
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
		"message": fmt.Sprintf("deleted post with id %d", id),
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) reactToPostHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		Reaction string `json:"reaction"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	err = app.models.Posts.React(user, id, input.Reaction)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrDuplicateEntry):
			v := validator.New()
			v.AddError("post", "you already reacted to this post")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := envelope{
		"message": fmt.Sprintf("added reaction %s to post with id %d", input.Reaction, id),
	}

	err = app.writeJSON(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) changeReactionPostHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		Reaction string `json:"reaction"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	err = app.models.Posts.UpdateReaction(user, id, input.Reaction)
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
		"message": fmt.Sprintf("updated reaction %s to post with id %d", input.Reaction, id),
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) removeReactionPostHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// DEBUG
	app.logger.Debug(fmt.Sprintf("asking to remove reaction to post with id %d", id))

	user := app.contextGetUser(r)

	err = app.models.Posts.RemoveReaction(user, id)
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
		"message": fmt.Sprintf("reaction removed from post with id %d", id),
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
