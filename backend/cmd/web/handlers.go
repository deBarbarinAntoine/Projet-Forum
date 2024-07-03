package main

import (
	"Projet-Forum/internal/data"
	"Projet-Forum/internal/validator"
	"encoding/json"
	"fmt"
	"github.com/alexedwards/flow"
	"net/http"
	"net/url"
	"time"
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

	// return to confirm page if there is an error
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

	// return to login page if there is an error
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

	// return to forgot-password page if there is an error
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

	// return to reset-password page if there is an error
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

	// render the template
	app.render(w, r, http.StatusOK, "user.tmpl", data)
}

func (app *application) updateUserPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newUserUpdateForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// creating the User struct to insert the new data into it
	user := &data.User{}

	// checking the data from the user
	var isEmpty = true
	if form.Username != nil {
		isEmpty = false
		form.StringCheck(*form.Username, 2, 70, false, "username")
		user.Name = *form.Username
	}
	if form.Password != nil || form.NewPassword != nil || form.ConfirmationPassword != nil {
		isEmpty = false
		form.ValidateNewPassword(*form.NewPassword, *form.ConfirmationPassword)
		user.Password = *form.NewPassword
	}
	if form.Email != nil {
		isEmpty = false
		form.ValidateEmail(*form.Email)
		user.Email = *form.Email
	}
	if form.Bio != nil {
		isEmpty = false
		form.StringCheck(*form.Bio, 2, 255, false, "bio")
		user.Bio = *form.Bio
	}
	if form.Birth != nil {
		isEmpty = false
		form.ValidateDate(*form.Birth, "birth")
		birthDate, err := time.Parse("2006-01-02", *form.Birth)
		if nil == err {
			user.BirthDate = birthDate
		}
	}
	if form.Signature != nil {
		isEmpty = false
		form.StringCheck(*form.Signature, 1, 255, false, "signature")
		user.Signature = *form.Signature
	}
	if isEmpty {
		form.AddNonFieldError("at least one field is required")
	}

	// return to update-user page if there is an error
	if !form.Valid() {

		// retrieving basic template data
		data := app.newTemplateData(r, false)

		data.NonFieldErrors = form.NonFieldErrors
		data.FieldErrors = form.FieldErrors

		// render the template
		app.render(w, r, http.StatusOK, "update-user.tmpl", data)
		return
	}

	// API request to send a reset password token
	v := validator.New()
	err = app.models.UserModel.Update(app.getToken(r), *form.Password, user, v)
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
		app.render(w, r, http.StatusOK, "update-user.tmpl", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Your data has been updated successfully!")
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *application) createCategory(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// render the template
	app.render(w, r, http.StatusOK, "category-create.tmpl", data)
}

func (app *application) createCategoryPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newCategoryForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// checking the data from the user
	if form.ParentCategoryID == nil {
		form.AddFieldError("parent_category_id", "must be provided")
	} else {
		form.CheckID(*form.ParentCategoryID, "parent_category_id")
	}
	if form.Name == nil {
		form.AddFieldError("name", "must be provided")
	} else {
		form.StringCheck(*form.Name, 2, 70, true, "name")
	}

	// return to category-create page if there is an error
	if !form.Valid() {
		data := app.newTemplateData(r, false)
		data.Form = form

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "category-create.tmpl", data)
		return
	}

	// creating the new category
	category := &data.Category{}
	category.Name = *form.Name
	category.Author.ID = app.getUserID(r)
	category.ParentCategory.ID = *form.ParentCategoryID

	// API request to create a category
	v := validator.New()
	err = app.models.CategoryModel.Create(app.getToken(r), category, v)
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
		app.render(w, r, http.StatusOK, "category-create.tmpl", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Category %s created successfully!", *form.Name))
	http.Redirect(w, r, fmt.Sprintf("/category/%d", category.ID), http.StatusSeeOther)
}

func (app *application) createThread(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// render the template
	app.render(w, r, http.StatusOK, "thread-create.tmpl", data)
}

func (app *application) createThreadPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newThreadForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// creating the new thread
	thread := &data.Thread{}

	// checking the data from the user
	if form.Title == nil {
		form.AddFieldError("title", "must be provided")
	} else {
		form.StringCheck(*form.Title, 2, 125, true, "title")
		thread.Title = *form.Title
	}
	if form.Description == nil {
		form.AddFieldError("description", "must be provided")
	} else {
		form.StringCheck(*form.Description, 1, 1_020, true, "name")
		thread.Description = *form.Description
	}
	if form.IsPublic == nil {
		thread.IsPublic = true
	} else {
		thread.IsPublic = *form.IsPublic
	}
	if form.CategoryID == nil {
		form.AddFieldError("category_id", "must be provided")
	} else {
		form.CheckID(*form.CategoryID, "category_id")
		thread.Category.ID = *form.CategoryID
	}

	// return to thread-create page if there is an error
	if !form.Valid() {
		data := app.newTemplateData(r, false)
		data.Form = form

		data.NonFieldErrors = form.NonFieldErrors
		data.FieldErrors = form.FieldErrors

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "thread-create.tmpl", data)
		return
	}

	// API request to create a category
	v := validator.New()
	err = app.models.ThreadModel.Create(app.getToken(r), thread, v)
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
		app.render(w, r, http.StatusOK, "thread-create.tmpl", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Thread created successfully!")
	http.Redirect(w, r, fmt.Sprintf("/thread/%d", thread.ID), http.StatusSeeOther)
}

func (app *application) createPost(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// render the template
	app.render(w, r, http.StatusOK, "post-create.tmpl", data)
}

func (app *application) createPostPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newPostForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// creating the new thread
	post := &data.Post{}

	// checking the data from the user
	if form.Content == nil {
		form.AddFieldError("content", "must be provided")
	} else {
		form.StringCheck(*form.Content, 2, 1_020, true, "content")
		post.Content = *form.Content
	}
	if form.ThreadID == nil {
		form.AddFieldError("thread_id", "must be provided")
	} else {
		form.CheckID(*form.ThreadID, "thread_id")
		post.Thread.ID = *form.ThreadID
	}
	if form.ParentPostID == nil {
		form.AddFieldError("parent_post_id", "must be provided")
	} else {
		form.CheckID(*form.ParentPostID, "parent_post_id")
		post.IDParentPost = *form.ParentPostID
	}

	// return to post-create page if there is an error
	if !form.Valid() {
		data := app.newTemplateData(r, false)
		data.Form = form

		data.NonFieldErrors = form.NonFieldErrors
		data.FieldErrors = form.FieldErrors

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "post-create.tmpl", data)
		return
	}

	// API request to create a category
	v := validator.New()
	err = app.models.PostModel.Create(app.getToken(r), post, v)
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
		app.render(w, r, http.StatusOK, "post-create.tmpl", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Post created successfully!")
	http.Redirect(w, r, fmt.Sprintf("/thread/%d", post.Thread.ID), http.StatusSeeOther)
}

func (app *application) createTag(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// render the template
	app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
}

func (app *application) createTagPost(w http.ResponseWriter, r *http.Request) {

	// retrieving the form data
	form := newTagForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// creating the new thread
	tag := &data.Tag{}

	// checking the data from the user
	if form.Name == nil {
		form.AddFieldError("name", "must be provided")
	} else {
		form.StringCheck(*form.Name, 2, 70, true, "name")
		tag.Name = *form.Name
	}

	// return to tag-create page if there is an error
	if !form.Valid() {
		data := app.newTemplateData(r, false)
		data.Form = form

		data.NonFieldErrors = form.NonFieldErrors
		data.FieldErrors = form.FieldErrors

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "tag-create.tmpl", data)
		return
	}

	// API request to create a category
	v := validator.New()
	err = app.models.TagModel.Create(app.getToken(r), tag, v)
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
		app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Tag %s created successfully!", tag.Name))
	http.Redirect(w, r, fmt.Sprintf("/tag/%d", tag.ID), http.StatusSeeOther)
}

func (app *application) updatePost(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// retrieving the post id from the path
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// retrieving the post from the API
	v := validator.New()
	post, err := app.models.PostModel.GetByID(app.getToken(r), id, v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// inserting the post values in the TemplateData's Form
	data.Form = post

	// render the template
	app.render(w, r, http.StatusOK, "post-update.tmpl", data)
}

func (app *application) updateTag(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// retrieving the tag id from the path
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// retrieving the tag from the API
	v := validator.New()
	query := url.Values{
		"includes[]": {"threads"},
	}
	tag, err := app.models.TagModel.GetByID(app.getToken(r), id, query, v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// inserting the tag values in the TemplateData's Form
	data.Form = tag

	// render the template
	app.render(w, r, http.StatusOK, "tag-update.tmpl", data)
}

func (app *application) updateCategory(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// retrieving the category id from the path
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// retrieving the category from the API
	v := validator.New()
	category, err := app.models.CategoryModel.GetByID(app.getToken(r), id, nil, v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// inserting the category values in the TemplateData's Form
	data.Form = category

	// render the template
	app.render(w, r, http.StatusOK, "category-update.tmpl", data)
}

func (app *application) updateThread(w http.ResponseWriter, r *http.Request) {

	// retrieving basic template data
	data := app.newTemplateData(r, false)

	// retrieving the thread id from the path
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// retrieving the thread from the API
	v := validator.New()
	thread, err := app.models.ThreadModel.GetByID(app.getToken(r), id, nil, v)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// inserting the thread values in the TemplateData's Form
	data.Form = thread

	// render the template
	app.render(w, r, http.StatusOK, "thread-update.tmpl", data)
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

	// retrieving the form data
	form := newTagForm()
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// checking the data from the user
	if form.Name != nil {
		form.StringCheck(*form.Name, 2, 70, false, "name")
	}
	if form.AddThreads != nil {
		form.Check(validator.Unique(form.AddThreads), "threads_ids", "must be unique")
	}
	if form.AddThreads != nil {
		form.Check(validator.Unique(form.AddThreads), "threads_ids", "must be unique")
	}

	// return to tag-create page if there is an error
	if !form.Valid() {
		data := app.newTemplateData(r, false)
		data.Form = form

		data.NonFieldErrors = form.NonFieldErrors
		data.FieldErrors = form.FieldErrors

		// render the template
		app.render(w, r, http.StatusUnprocessableEntity, "tag-create.tmpl", data)
		return
	}

	// creating the API request body
	body, err := json.Marshal(form)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// retrieving the thread id from the path
	id, err := getPathID(r)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// API request to create a category
	v := validator.New()
	tag, err := app.models.TagModel.Update(app.getToken(r), id, body, v)
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
		app.render(w, r, http.StatusOK, "tag-create.tmpl", data)
		return
	}

	app.sessionManager.Put(r.Context(), "flash", fmt.Sprintf("Tag %s updated successfully!", tag.Name))
	http.Redirect(w, r, fmt.Sprintf("/tag/%d", tag.ID), http.StatusSeeOther)
}
