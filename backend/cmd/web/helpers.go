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

func (app *application) clientError(r *http.Request, w http.ResponseWriter, status int) {

	// setting the templateData
	tmplData := app.newTemplateData(r, false, Overlay.Default)

	// setting the error title and message
	tmplData.Error.Title = fmt.Sprintf("Error %d", status)

	if status == http.StatusNotFound {
		tmplData.Error.Message = "We didn't find what you were looking for :("
	} else {
		tmplData.Error.Message = "Something went wrong!"
	}

	// rendering the error page
	app.render(w, r, status, "error.tmpl", tmplData)
}

func (app *application) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		status = http.StatusInternalServerError
		method = r.Method
		uri    = r.URL.RequestURI()
		trace  = string(debug.Stack())
	)

	// logging the error
	app.logger.Error(err.Error(), slog.String("method", method), slog.String("URI", uri), slog.String("trace", trace))

	// setting the templateData
	tmplData := app.newTemplateData(r, false, Overlay.Default)

	// setting the error title and message
	tmplData.Error.Title = fmt.Sprintf("Error %d", status)
	tmplData.Error.Message = "Something went wrong!"

	// rendering the error page
	app.render(w, r, status, "error.tmpl", tmplData)
}

func (app *application) isAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value(isAuthenticatedContextKey).(bool)
	if !ok {
		return false
	}

	return isAuthenticated
}

func (app *application) putToken(r *http.Request, tokens data.Tokens) {
	app.sessionManager.Put(r.Context(), authTokenSessionManager, tokens.Authentication.Token)
	app.sessionManager.Put(r.Context(), authExpirySessionManager, tokens.Authentication.Expiry.String())
	app.sessionManager.Put(r.Context(), refreshTokenSessionManager, tokens.Refresh.Token)
	app.sessionManager.Put(r.Context(), refreshExpirySessionManager, tokens.Refresh.Expiry.String())
}

func (app *application) getTokens(r *http.Request) (*data.Tokens, error) {
	var tokens data.Tokens
	authToken := app.getToken(r, authTokenSessionManager)
	refreshToken := app.getToken(r, refreshTokenSessionManager)
	authExpiry, ok := app.sessionManager.Get(r.Context(), authExpirySessionManager).(string)
	if !ok {
		return nil, fmt.Errorf("couldn't retreive authentication token's expiry from session manager")
	}
	refreshExpiry, ok := app.sessionManager.Get(r.Context(), refreshExpirySessionManager).(string)
	if !ok {
		return nil, fmt.Errorf("couldn't retreive refresh token's expiry from session manager")
	}
	var err error
	tokens.Authentication.Token = authToken
	tokens.Authentication.Expiry, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", authExpiry)
	if err != nil {
		return nil, err
	}
	tokens.Refresh.Token = refreshToken
	tokens.Refresh.Expiry, err = time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", refreshExpiry)
	if err != nil {
		return nil, err
	}
	return &tokens, nil
}

func (app *application) getToken(r *http.Request, key string) string {
	token, ok := app.sessionManager.Get(r.Context(), key).(string)
	if !ok {
		app.logger.Debug("could not get token from session manager")
		return ""
	}
	return token
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

func newReactToPostForm() *reactToPostForm {
	return &reactToPostForm{
		Validator:        *validator.New(),
		AllowedReactions: []string{"neutral", "laughing", "applause", "heart"},
	}
}

func newFriendResponseForm() *friendResponseForm {
	return &friendResponseForm{
		Validator:      *validator.New(),
		FriendStatuses: []string{"pending", "accepted", "rejected"},
	}
}

func (app *application) newTemplateData(r *http.Request, allUser bool, overlay string) templateData {

	// checking is the user is authenticated
	isAuthenticated := app.isAuthenticated(r)
	token := app.getToken(r, authTokenSessionManager)
	// retrieving the user data
	user := &data.User{}
	v := validator.New()
	if isAuthenticated {
		var query url.Values
		if !allUser {
			query = url.Values{
				"includes[]": {"posts"},
			}
		} else {
			query = url.Values{
				"includes[]": {"following_tags", "favorite_threads", "categories_owned", "tags_owned", "threads_owned", "friends", "posts", "reactions"},
			}
		}
		user, _ = app.models.UserModel.GetByID(token, "me", query, v)

		// DEBUG
		app.logger.Debug(fmt.Sprintf("user: %+v", user))
	}
	categories, metadata, err := app.models.CategoryModel.Get(token, nil, v)
	if err != nil {
		app.logger.Error(err.Error())
	}
	tags, threads, err := app.models.TagModel.GetPopular(token, v)
	if err != nil {
		app.logger.Error(err.Error())
	}

	// returning the templateData with all information
	return templateData{
		Overlay:         overlay,
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: isAuthenticated,
		CSRFToken:       nosurf.Token(r),
		Error: struct {
			Title   string
			Message string
		}{
			Title:   "Error 404",
			Message: "We didn't find what you were looking for :(",
		},
		User:              *user,
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
