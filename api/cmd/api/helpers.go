package main

import (
	"ForumAPI/internal/data"
	"ForumAPI/internal/validator"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/alexedwards/flow"
	"github.com/go-playground/form/v4"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	ErrUserIDNotFound = errors.New("user id not found")
)

type envelope map[string]any

/* #######################################################################
/* # Query form data
/* ####################################################################### */

func newUserByIDForm() *userByIDForm {
	return &userByIDForm{
		Validator: *validator.New(),
	}
}

func newCategoryByIDForm() *categoryByIDForm {
	return &categoryByIDForm{
		PermittedFields: []string{"categories", "threads"},
		Validator:       *validator.New(),
	}
}

func newGetCategoriesForm() *getCategoriesForm {
	return &getCategoriesForm{
		Validator: *validator.New(),
		Filters: data.Filters{
			SortSafelist: []string{"Created_at", "Updated_at", "Name", "Id_author", "Id_parent_categories", "-Created_at", "-Updated_at", "-Name", "-Id_author", "-Id_parent_categories"},
		},
	}
}

func newThreadByIDForm() *threadByIDForm {
	return &threadByIDForm{
		PermittedFields: []string{"posts", "tags", "popularity"},
		Validator:       *validator.New(),
	}
}

func newGetThreadsForm() *getThreadsForm {
	return &getThreadsForm{
		Validator: *validator.New(),
		Filters: data.Filters{
			SortSafelist: []string{"Created_at", "Updated_at", "Title", "Is_public", "Status", "Id_author", "Id_categories", "-Created_at", "-Updated_at", "-Title", "-Is_public", "-Status", "-Id_author", "-Id_categories"},
		},
	}
}

func newTagByIDForm() *tagByIDForm {
	return &tagByIDForm{
		PermittedFields: []string{"threads", "popularity"},
		Validator:       *validator.New(),
	}
}

func newGetTagsForm() *getTagsForm {
	return &getTagsForm{
		Validator: *validator.New(),
		Filters: data.Filters{
			SortSafelist: []string{"Created_at", "Updated_at", "Name", "Id_author", "-Created_at", "-Updated_at", "-Name", "-Id_author"},
		},
	}
}

func newPostByIDForm() *postByIDForm {
	return &postByIDForm{
		PermittedFields: []string{"reactions", "popularity"},
		Validator:       *validator.New(),
	}
}

func newGetPostsForm() *getPostsForm {
	return &getPostsForm{
		Validator: *validator.New(),
		Filters: data.Filters{
			SortSafelist: []string{"Created_at", "Updated_at", "Id_author", "Id_threads", "-Created_at", "-Updated_at", "-Id_author", "-Id_threads"},
		},
	}
}

/* #######################################################################
/* # Other helper functions
/* ####################################################################### */

func (app *application) decodeForm(r *http.Request, dst any) error {

	err := r.ParseForm()
	if err != nil {
		return err
	}

	err = app.formDecoder.Decode(dst, r.Form)
	if err != nil {
		var invalidDecoderError *form.InvalidDecoderError

		if errors.As(err, &invalidDecoderError) {
			panic(err)
		}

		return err
	}

	return nil
}

func (app *application) readIDParam(r *http.Request) (int, error) {

	idParam := flow.Param(r.Context(), "id")
	if idParam == "me" {
		id := app.contextGetUser(r).ID
		if id == 0 {
			return 0, ErrUserIDNotFound
		}
		return id, nil
	}
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}

	return int(id), nil
}

func (app *application) writeJSON(w http.ResponseWriter, status int, data envelope, headers http.Header) error {

	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst any) error {

	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)

	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		var maxBytesError *http.MaxBytesError

		switch {
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly-formed JSON")

		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains incorrect JSON type (at character %d)", unmarshalTypeError.Offset)

		case errors.Is(err, io.EOF):
			return errors.New("body must not be empty")

		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", fieldName)

		case errors.As(err, &maxBytesError):
			return fmt.Errorf("body must not be longer than %d bytes", maxBytesError.Limit)

		case errors.As(err, &invalidUnmarshalError):
			panic(err)

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if !errors.Is(err, io.EOF) {
		return errors.New("body must contain a single JSON value")
	}

	return nil
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)

	if s == "" {
		return defaultValue
	}

	return s
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	csv := qs.Get(key)

	if csv == "" {
		return defaultValue
	}

	return strings.Split(csv, ",")
}

func (app *application) readInt(queryStr url.Values, key string, defaultValue int, v *validator.Validator) int {
	s := queryStr.Get(key)

	if s == "" {
		return defaultValue
	}

	i, err := strconv.Atoi(s)
	if err != nil {
		v.AddError(key, "must be an integer value")
		return defaultValue
	}

	return i
}

func (app *application) background(fn func()) {

	app.wg.Add(1)
	go func() {

		defer app.wg.Done()

		defer func() {
			if err := recover(); err != nil {
				app.logger.Error(fmt.Sprintf("%v", err))
			}
		}()

		fn()

	}()
}
