package main

import (
	"Projet-Forum/internal/models"
	"Projet-Forum/internal/validator"
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/form/v4"
	"github.com/justinas/nosurf"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
)

func (app *application) sendAuthPEM(credentials *models.Credentials) error {

	// mock data only for testing purposes
	authJson, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	cipher, err := app.encryptPEM(authJson)
	if err != nil {
		//app.serverErrorResponse(w, r, err)
		return err
	}
	client := http.Client{Timeout: 5 * time.Second}

	req, err := http.NewRequest(http.MethodPost, "http://localhost:3000/v1/tokens/authentication", strings.NewReader(hex.EncodeToString(cipher)))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", app.config.clientToken))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Encryption", "RSA")

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	app.logger.Debug(fmt.Sprintf("response: %s", string(responseBody)))
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

func newUserRegisterForm() *userRegisterForm {
	return &userRegisterForm{
		Validator: *validator.NewValidator(),
	}
}

func newUserLoginForm() *userLoginForm {
	return &userLoginForm{
		Validator: *validator.NewValidator(),
	}
}

func newCategoryForm() *categoryForm {
	return &categoryForm{
		Validator: *validator.NewValidator(),
	}
}

func newThreadForm() *threadForm {
	return &threadForm{
		Validator: *validator.NewValidator(),
	}
}

func newPostForm() *postForm {
	return &postForm{
		Validator: *validator.NewValidator(),
	}
}

func newTagForm() *tagForm {
	return &tagForm{
		Validator: *validator.NewValidator(),
	}
}

func (app *application) newTemplateData(r *http.Request) templateData {
	return templateData{
		CurrentYear:     time.Now().Year(),
		Flash:           app.sessionManager.PopString(r.Context(), "flash"),
		IsAuthenticated: app.isAuthenticated(r),
		CSRFToken:       nosurf.Token(r),
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

func getId(queryValues url.Values) (int, error) {
	if queryValues.Has("id") {
		return strconv.Atoi(queryValues.Get("id"))
	}
	return -1, fmt.Errorf("id not found")
}
