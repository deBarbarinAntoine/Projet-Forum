package main

import (
	"Projet-Forum/internal/data"
	"Projet-Forum/internal/utils"
	"Projet-Forum/internal/validator"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {

	log.Println(utils.GetCurrentFuncName())

	data := app.newTemplateData(r)
	data.Message = "Welcome to my website!"

	app.render(w, r, http.StatusOK, "index.tmpl", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the About page!"))
}

func (app *application) register(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = newUserRegisterForm()

	app.render(w, r, http.StatusOK, "register.tmpl", data)
}

func (app *application) registerPost(w http.ResponseWriter, r *http.Request) {

	form := newUserRegisterForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")
	form.CheckField(form.Password == form.ConfirmPassword, "password", "Different passwords")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	err = app.models.UserModel.Register(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, data.ErrDuplicateCredential) {
			form.AddNonFieldError("User already exists")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "We've sent you a comfirmation email!")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = newUserLoginForm()

	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {

	form := newUserLoginForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	var id int
	id, err = app.models.UserModel.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, data.ErrInvalidCredentials) {
			form.AddNonFieldError("Invalid credentials")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticatedUserID", id)

	http.Redirect(w, r, "/protected", http.StatusSeeOther)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the home page! (restricted entry)"))
}

func (app *application) confirmHandler(w http.ResponseWriter, r *http.Request) {

	err := app.models.UserModel.Activate(r)
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your session has been activated successfully!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) getCategory(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	id, err := getId(r.URL.Query())
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}

	data.Category.Fetch(id)

	app.render(w, r, http.StatusOK, "category.tmpl", data)
}

func (app *application) createCategory(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = newCategoryForm()

	app.render(w, r, http.StatusOK, "category-create.tmpl", data)
}

func (app *application) createCategoryPost(w http.ResponseWriter, r *http.Request) {

	form := newCategoryForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// TODO check form values

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "category-create.tmpl", data)
		return
	}

	// TODO create category
	id := app.models.CategoryModel.Create(form.Name, form.Author, form.ParentCategory)

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Category %s created successfully!", form.Name))

	http.Redirect(w, r, fmt.Sprintf("/category/%d", id), http.StatusSeeOther)
}

func (app *application) getThread(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	id, err := getId(r.URL.Query())
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}

	data.Thread.Fetch(id)

	app.render(w, r, http.StatusOK, "thread.tmpl", data)
}

func (app *application) createThread(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = newThreadForm()

	app.render(w, r, http.StatusOK, "thread-create.tmpl", data)
}

func (app *application) createThreadPost(w http.ResponseWriter, r *http.Request) {

	form := newThreadForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// TODO check form values

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "thread-create.tmpl", data)
		return
	}

	// TODO create thread
	id := app.models.ThreadModel.Create(form.Title, form.Description, form.IsPublic, form.Author, form.Category)

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Thread %s created successfully!", form.Title))

	http.Redirect(w, r, fmt.Sprintf("/thread/%d", id), http.StatusSeeOther)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = newPostForm()

	app.render(w, r, http.StatusOK, "post-create.tmpl", data)
}

func (app *application) createPostPost(w http.ResponseWriter, r *http.Request) {

	form := newPostForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// TODO check form values

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "post-create.tmpl", data)
		return
	}

	// TODO create post
	id := app.models.PostModel.Create(form.Content, form.Author, form.ThreadId, form.ParentPostId)

	http.Redirect(w, r, fmt.Sprintf("/thread/%d", id), http.StatusSeeOther)
}

func (app *application) getTag(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	id, err := getId(r.URL.Query())
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}

	data.Tag.Fetch(id)

	app.render(w, r, http.StatusOK, "tag.tmpl", data)
}

func (app *application) createTag(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = newTagForm()

	app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
}

func (app *application) createTagPost(w http.ResponseWriter, r *http.Request) {

	form := newTagForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// TODO check form values

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "tag-create.tmpl", data)
		return
	}

	// TODO create tag
	id := app.models.TagModel.Create(form.Name, form.Author)

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Tag %s created successfully!", form.Name))

	http.Redirect(w, r, fmt.Sprintf("/tag/%d", id), http.StatusSeeOther)
}

func (app *application) getProfile(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	id, err := getId(r.URL.Query())
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}

	data.User.Fetch(id)

	app.render(w, r, http.StatusOK, "user.tmpl", data)
}
