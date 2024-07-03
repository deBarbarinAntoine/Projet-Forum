package main

import (
	"Projet-Forum/internal/data"
	"Projet-Forum/internal/validator"
	"bytes"
	"errors"
	"fmt"
	"github.com/alexedwards/flow"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"log/slog"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"time"
)

func (app *application) logout(r *http.Request) error {

	err := app.sessionManager.Clear(r.Context())
	if err != nil {
		return err
	}
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		return err
	}

	return nil
}

func (app *application) decodePostForm(r *http.Request, dst any) error {

	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.PostForm)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		status = http.StatusInternalServerError
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	app.logger.Error(err.Error(), slog.String("method", method), slog.String("URI", uri), slog.String("trace", trace))
	http.Error(w, http.StatusText(status), status)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *application) getToken(r *http.Request) string {
	tokens, ok := app.sessionManager.Get(r.Context(), userTokenSessionManager).(*data.Tokens)
	if !ok {
		return ""
	}
	return tokens.Authentication.Token
}

func (app *application) getUserID(r *http.Request) int {
	id, ok := app.sessionManager.Get(r.Context(), authenticatedUserIDSessionManager).(int)
	if !ok {
		return 0
	}
	return id
}

func newUserRegisterForm() *userRegisterForm {
	return &userRegisterForm{
		Validator: *validator.New(),
	}
}

func newUserConfirmForm() *userConfirmForm {
	return &userConfirmForm{
		Validator: *validator.New(),
	}
}

func newUserLoginForm() *userLoginForm {
	return &userLoginForm{
		Validator: *validator.New(),
	}
}

func newUserUpdateForm() *userUpdateForm {
	return &userUpdateForm{
		Validator: *validator.New(),
	}
}

func newForgotPasswordForm() *forgotPasswordForm {
	return &forgotPasswordForm{
		Validator: *validator.New(),
	}
}

func newResetPasswordForm() *resetPasswordForm {
	return &resetPasswordForm{
		Validator: *validator.New(),
	}
}

func newCategoryForm() *categoryForm {
	return &categoryForm{
		Validator: *validator.New(),
	}
}

func newThreadForm() *threadForm {
	return &threadForm{
		Validator: *validator.New(),
	}
}

func newPostForm() *postForm {
	return &postForm{
		Validator: *validator.New(),
	}
}

func newTagForm() *tagForm {
	return &tagForm{
		Validator: *validator.New(),
	}
}

func (app *application) newTemplateData(r *http.Request, allUser bool) templateData {

	// checking is the user is authenticated
	isAuthenticated := app.isAuthenticated(r)
	token := app.getToken(r)
	// retrieving the user data
	var user *data.User
	v := validator.New()
	if isAuthenticated {
		var query url.Values
		if !allUser {
			query = url.Values{
				"includes[]": {"posts"},
			}
		} else {
			query = url.Values{
				"includes[]": {"following_tags", "favorite_threads", "categories_owned", "tags_owned", "threads_owned", "friends", "posts"},
			}
		}
		user, _ = app.models.UserModel.GetByID(token, "me", query, v)
	}
	categories, metadata, _ := app.models.CategoryModel.Get(token, nil, v)
	tags, threads, _ := app.models.TagModel.GetPopular(token, v)

	// returning the templateData with all information
	return templateData{
		CurrentYear:       time.Now().Year(),
		Flash:             app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated:   isAuthenticated,
		CSRFToken:         nosurf.Token(r),
		User:              user,
		PopularTags:       tags,
		PopularThreads:    threads,
		CategoriesNavLeft: categories,
		CategoryList: struct {
			Metadata data.Metadata
			List     []*data.Category
		}{Metadata: metadata, List: categories},
	}
}

func (app *application) render(w http.ResponseWriter, r *http.Request, status int, page string, data templateData) {

	// retrieving the appropriate set of templates
	ts, ok := app.templateCache[page]
	if !ok {
		app.serverError(w, r, fmt.Errorf("the template %s does not exist", page))
		return
	}

	// creating a bytes Buffer
	buf := new(bytes.Buffer)

	// executing the template in the buffer to catch any possible parsing error,
	// so that the user doesn't see a half-empty page
	err := ts.ExecuteTemplate(buf, "base", data)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// if it's all okay, write the status in the header and write the buffer in the ResponseWriter
	w.WriteHeader(status)

	buf.WriteTo(w)
}

func getPathID(r *http.Request) (int, error) {

	// fetching the id param from the URL
	param := flow.Param(r.Context(), "id")

	// looking for errors
	if param == "" {
		return 0, fmt.Errorf("id required")
	}
	if param != "me" {
		id, err := strconv.Atoi(param)
		if err != nil || id < 1 {
			return 0, fmt.Errorf("invalid id")
		}

		// return the integer id
		return id, nil
	}

	// return -1 when id == "me"
	return -1, nil
}
