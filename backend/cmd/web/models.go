package main

import (
	"Projet-Forum/internal/data"
	"Projet-Forum/internal/validator"
	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"html/template"
	"log/slog"
)

type config struct {
	isHTTPS bool
	apiURL  string
	port    int
}

type application struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	models         data.Models
	config         *config
}

type templateData struct {
	CurrentYear       int
	Message           string
	Form              any
	Flash             string
	IsAuthenticated   bool
	CSRFToken         string
	ActivationToken   string
	ResetToken        string
	FieldErrors       map[string]string
	NonFieldErrors    []string
	User              *data.User
	CategoriesNavLeft []*data.Category
	PopularTags       []*data.Tag
	PopularThreads    []*data.Thread
	CategoryList      struct {
		Metadata data.Metadata
		List     []*data.Category
	}
	ThreadList struct {
		Metadata data.Metadata
		List     []*data.Thread
	}
	PostList struct {
		Metadata data.Metadata
		List     []*data.Post
	}
	TagList struct {
		Metadata data.Metadata
		List     []*data.Tag
	}
	Category *data.Category
	Thread   *data.Thread
	Tag      *data.Tag
}

type userLoginForm struct {
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
}

type userRegisterForm struct {
	Username            string `form:"username"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	ConfirmPassword     string `form:"confirm_password"`
	validator.Validator `form:"-"`
}

type userConfirmForm struct {
	Token               string `form:"token"`
	validator.Validator `form:"-"`
}

type forgotPasswordForm struct {
	Email               string `form:"email"`
	validator.Validator `form:"-"`
}

type resetPasswordForm struct {
	Token               string `form:"token"`
	NewPassword         string `form:"new_password"`
	ConfirmPassword     string `form:"confirm_password"`
	validator.Validator `form:"-"`
}

type categoryForm struct {
	Name                string `form:"name"`
	Author              string `form:"author"`
	ParentCategory      *int   `form:"parent_category,omitempty"`
	validator.Validator `form:"-"`
}

type threadForm struct {
	Title               string `form:"title"`
	Description         string `form:"description"`
	IsPublic            bool   `form:"is_public"`
	Author              string `form:"author"`
	Category            string `form:"category"`
	validator.Validator `form:"-"`
}

type postForm struct {
	Content             string `form:"content"`
	Author              string `form:"author"`
	ThreadId            int    `form:"thread_id"`
	ParentPostId        *int   `form:"parent_post_id,omitempty"`
	validator.Validator `form:"-"`
}

type tagForm struct {
	Name                string `form:"name"`
	Author              string `form:"author"`
	validator.Validator `form:"-"`
}
