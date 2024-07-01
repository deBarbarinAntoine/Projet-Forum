package main

import (
	"Projet-Forum/internal/data"
	"Projet-Forum/internal/validator"
	"fmt"
	"github.com/alexedwards/flow"
	"net/http"
)

/* #############################################################################
/*	COMMON
/* #############################################################################*/

func (app *application) notFound(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// render the template
	app.render(w, r, http.StatusOK, "error.tmpl", data)
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, true)

	// render the template
	app.render(w, r, http.StatusOK, "index.tmpl", data)
}

func (app *application) about(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// render the template
	app.render(w, r, http.StatusOK, "about.tmpl", data)
}

func (app *application) categoryGet(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// fetching the category id in the path
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// fetching the category
	v := validator.New()
	data.Category, err = app.models.CategoryModel.GetByID(app.getToken(r), id, r.URL.Query(), v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking API request errors
	if !v.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// render the template
	app.render(w, r, http.StatusOK, "category.tmpl", data)
}

func (app *application) threadGet(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// fetching the thread id in the path
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// fetching the thread
	v := validator.New()
	data.Thread, err = app.models.ThreadModel.GetByID(app.getToken(r), id, r.URL.Query(), v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking API request errors
	if !v.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// render the template
	app.render(w, r, http.StatusOK, "thread.tmpl", data)
}

func (app *application) tagGet(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// fetching the tag id in the path
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// fetching the tag
	v := validator.New()
	data.Tag, err = app.models.TagModel.GetByID(app.getToken(r), id, r.URL.Query(), v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking API request errors
	if !v.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// render the template
	app.render(w, r, http.StatusOK, "tag.tmpl", data)
}

func (app *application) TagsGet(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// fetching the tags
	v := validator.New()
	var err error
	data.TagList.List, data.TagList.Metadata, err = app.models.TagModel.Get(app.getToken(r), r.URL.Query(), v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking API request errors
	if !v.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// render the template
	app.render(w, r, http.StatusOK, "tags.tmpl", data)
}

func (app *application) categoriesGet(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// fetching the categories
	v := validator.New()
	var err error
	data.CategoryList.List, data.CategoryList.Metadata, err = app.models.CategoryModel.Get(app.getToken(r), r.URL.Query(), v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking API request errors
	if !v.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// render the template
	app.render(w, r, http.StatusOK, "tag.tmpl", data)
}

func (app *application) search(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// fetching the categories
	v := validator.New()
	var err error
	data.CategoryList.List, data.CategoryList.Metadata, err = app.models.CategoryModel.Get(app.getToken(r), r.URL.Query(), v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// fetching the threads
	data.ThreadList.List, data.ThreadList.Metadata, err = app.models.ThreadModel.Get(app.getToken(r), r.URL.Query(), v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// fetching the tags
	data.TagList.List, data.TagList.Metadata, err = app.models.TagModel.Get(app.getToken(r), r.URL.Query(), v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// checking API request errors
	if !v.Valid() {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// render the template
	app.render(w, r, http.StatusOK, "search.tmpl", data)
}

/* #############################################################################
/*	USER ACCESS
/* #############################################################################*/

func (app *application) register(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// render the template
	app.render(w, r, http.StatusOK, "register.tmpl", data)
}

func (app *application) registerPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newUserRegisterForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// checking the data from the user
	form.StringCheck(form.Username, 2, 70, true, "username")
	form.ValidateEmail(form.Email)
	form.ValidateRegisterPassword(form.Password, form.ConfirmPassword)

	// return to register page if there is an error
	if !form.Valid() {

		// retrieving basic template data
		data := app.newTemplateData(r, false)
		data.Form = form

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "register.tmpl", data)
		return
	}

	// creating and registering the user
	v := validator.New()
	user := &data.User{
		Name:     form.Username,
		Email:    form.Email,
		Password: form.Password,
	}
	err = app.models.UserModel.Create(app.getToken(r), user, v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// looking for errors from the API
	if !v.Valid() {

		// retrieving basic template data
		data := app.newTemplateData(r, false)

		data.Form = form

		data.FieldErrors = v.FieldErrors
		data.NonFieldErrors = v.NonFieldErrors

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "register.tmpl", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "We've sent you a confirmation email!")

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) confirm(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// retrieving the activation token from the URL
	data.ActivationToken = flow.Param(r.Context(), "token")
	if data.ActivationToken == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// render the template
	app.render(w, r, http.StatusOK, "confirm.tmpl", data)
}

func (app *application) confirmPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newUserConfirmForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// checking the data from the user
	form.ValidateToken(form.Token)

	// looking for errors from the API
	if !form.Valid() {

		// retrieving basic template data
		data := app.newTemplateData(r, false)

		data.FieldErrors = form.FieldErrors
		data.ActivationToken = form.Token

		// render the template
		app.render(w, r, http.StatusOK, "confirm.tmpl", data)
		return
	}

	// API request to activate the user account
	v := validator.New()
	err = app.models.UserModel.Activate(form.Token, v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// looking for errors from the API
	if !v.Valid() {

		// retrieving basic template data
		data := app.newTemplateData(r, false)

		data.FieldErrors = form.FieldErrors
		data.ActivationToken = form.Token

		// render the template
		app.render(w, r, http.StatusOK, "confirm.tmpl", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your session has been activated successfully!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) login(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// render the template
	app.render(w, r, http.StatusOK, "login.tmpl", data)
}

func (app *application) loginPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newUserLoginForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// checking the data from the user
	form.Check(validator.NotBlank(form.Email), "email", "This field cannot be blank")
	form.Check(validator.Matches(form.Email, validator.EmailRX), "email", "This field must be a valid email address")
	form.ValidatePassword(form.Password)

	// looking for errors from the API
	if !form.Valid() {

		// retrieving basic template data
		data := app.newTemplateData(r, false)

		data.Form = form

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// building the API request body
	body := map[string]string{
		"email":    form.Email,
		"password": form.Password,
	}

	// API request to authenticate the user
	v := validator.New()
	tokens, err := app.models.TokenModel.Authenticate(body, v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// looking for errors from the API
	if !v.Valid() {

		// retrieving basic template data
		data := app.newTemplateData(r, false)

		data.FieldErrors = form.FieldErrors

		form.Password = ""
		data.Form = form

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "login.tmpl", data)
		return
	}

	// renewing the user session
	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// fetching the user id from the API
	user, err := app.models.UserModel.GetByID(tokens.Authentication.Token, "me", nil, v)
	if err != nil || !v.Valid() {
		app.serverError(w, r, err)
		return
	}

	// storing the user id and tokens in the user session
	app.sessionManager.Put(r.Context(), authenticatedUserIDSessionManager, user.ID)
	app.sessionManager.Put(r.Context(), userTokenSessionManager, tokens)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) forgotPassword(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// render the template
	app.render(w, r, http.StatusOK, "forgot-password.tmpl", data)
}

func (app *application) forgotPasswordPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newForgotPasswordForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// checking the data from the user
	form.ValidateEmail(form.Email)

	// looking for errors from the API
	if !form.Valid() {

		// retrieving basic template data
		data := app.newTemplateData(r, false)

		data.NonFieldErrors = form.NonFieldErrors
		data.FieldErrors = form.FieldErrors

		data.Form = form

		// render the template
		app.render(w, r, http.StatusOK, "forgot-password.tmpl", data)
		return
	}

	// API request to send a reset password token
	v := validator.New()
	err = app.models.UserModel.ForgotPassword(form.Email, v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// looking for errors from the API
	if !v.Valid() {

		// retrieving basic template data
		data := app.newTemplateData(r, false)

		data.NonFieldErrors = form.NonFieldErrors
		data.FieldErrors = form.FieldErrors

		data.Form = form

		// render the template
		app.render(w, r, http.StatusOK, "forgot-password.tmpl", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "We've sent you a mail to reset your password!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func (app *application) resetPassword(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// retrieving the reset token from the URL
	data.ResetToken = flow.Param(r.Context(), "token")
	if data.ResetToken == "" {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// render the template
	app.render(w, r, http.StatusOK, "reset-password.tmpl", data)
}

func (app *application) resetPasswordPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newResetPasswordForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// checking the data from the user
	form.ValidateNewPassword(form.NewPassword, form.ConfirmPassword)
	form.ValidateToken(form.Token)

	// looking for errors from the API
	if !form.Valid() {

		// retrieving basic template data
		data := app.newTemplateData(r, false)

		data.NonFieldErrors = form.NonFieldErrors
		data.FieldErrors = form.FieldErrors

		// render the template
		app.render(w, r, http.StatusOK, "reset-password.tmpl", data)
		return
	}

	// building the API request body
	body := map[string]string{
		"token":            form.Token,
		"new_password":     form.NewPassword,
		"confirm_password": form.ConfirmPassword,
	}

	// API request to send a reset password token
	v := validator.New()
	err = app.models.UserModel.ResetPassword(body, v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// looking for errors from the API
	if !v.Valid() {

		// retrieving basic template data
		data := app.newTemplateData(r, false)

		data.NonFieldErrors = form.NonFieldErrors
		data.FieldErrors = form.FieldErrors

		// render the template
		app.render(w, r, http.StatusOK, "reset-password.tmpl", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your password has been updated successfully!")
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

/* #############################################################################
/*	RESTRICTED
/* #############################################################################*/

func (app *application) dashboard(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, true)

	// render the template
	app.render(w, r, http.StatusOK, "dashboard.tmpl", data)
}

func (app *application) logoutPost(w http.ResponseWriter, r *http.Request) {

	// logging the user out
	err := app.logout(r)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "You've been logged out successfully!")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) updateUser(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	id, err := getId(r.URL.Query())
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}

	data.User.Fetch(id)

	// render the template
	app.render(w, r, http.StatusOK, "user.tmpl", data)
}

func (app *application) updateUserPost(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	id, err := getId(r.URL.Query())
	if err != nil {
		app.clientError(w, http.StatusNotFound)
		return
	}

	data.User.Fetch(id)

	// render the template
	app.render(w, r, http.StatusOK, "user.tmpl", data)
}

func (app *application) createCategory(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newCategoryForm()

	// render the template
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
		data := app.newTemplateData(r, false)
		data.Form = form

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "category-create.tmpl", data)
		return
	}

	// TODO create category
	id := app.models.CategoryModel.Create(form.Name, form.Author, form.ParentCategory)

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Category %s created successfully!", form.Name))

	http.Redirect(w, r, fmt.Sprintf("/category/%d", id), http.StatusSeeOther)
}

func (app *application) createThread(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newThreadForm()

	// render the template
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
		data := app.newTemplateData(r, false)
		data.Form = form

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "thread-create.tmpl", data)
		return
	}

	// TODO create thread
	id := app.models.ThreadModel.Create(form.Title, form.Description, form.IsPublic, form.Author, form.Category)

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Thread %s created successfully!", form.Title))

	http.Redirect(w, r, fmt.Sprintf("/thread/%d", id), http.StatusSeeOther)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newPostForm()

	// render the template
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
		data := app.newTemplateData(r, false)
		data.Form = form

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "post-create.tmpl", data)
		return
	}

	// TODO create post
	id := app.models.PostModel.Create(form.Content, form.Author, form.ThreadId, form.ParentPostId)

	http.Redirect(w, r, fmt.Sprintf("/thread/%d", id), http.StatusSeeOther)
}

func (app *application) createTag(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newTagForm()

	// render the template
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
		data := app.newTemplateData(r, false)
		data.Form = form

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "tag-create.tmpl", data)
		return
	}

	// TODO create tag
	id := app.models.TagModel.Create(form.Name, form.Author)

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Tag %s created successfully!", form.Name))

	http.Redirect(w, r, fmt.Sprintf("/tag/%d", id), http.StatusSeeOther)
}

func (app *application) updatePost(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newTagForm()

	// render the template
	app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
}

func (app *application) updateTag(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newTagForm()

	// render the template
	app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
}

func (app *application) updateCategory(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newTagForm()

	// render the template
	app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
}

func (app *application) updateThread(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newTagForm()

	// render the template
	app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
}

func (app *application) updateCategoryPut(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newTagForm()

	// render the template
	app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
}

func (app *application) updateThreadPut(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newTagForm()

	// render the template
	app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
}

func (app *application) updatePostPut(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newTagForm()

	// render the template
	app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
}

func (app *application) updateTagPut(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	data.Form = newTagForm()

	// render the template
	app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
}
