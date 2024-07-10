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
	isHTTPS     bool
	apiURL      string
	port        int64
	dsn         string
	secret      string
	clientToken string
	pemPath     string
}

type application struct {
	logger         *slog.Logger
	templateCache  map[string]*template.Template
	formDecoder    *form.Decoder
	sessionManager *scs.SessionManager
	models         data.Models
	config         *config
}

type overlayEnum struct {
	Default        string
	Login          string
	Register       string
	ForgotPassword string
	ResetPassword  string
}

var Overlay = overlayEnum{
	Default:        "default",
	Login:          "login",
	Register:       "register",
	ForgotPassword: "forgot-password",
	ResetPassword:  "reset-password",
}

type templateData struct {
	Title             string
	Overlay           string
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

type userUpdateForm struct {
	Username             *string `form:"username,omitempty"`
	Email                *string `form:"email,omitempty"`
	Password             *string `form:"password,omitempty"`
	NewPassword          *string `form:"new_password,omitempty"`
	ConfirmationPassword *string `form:"confirmation_password,omitempty"`
	Avatar               *string `form:"avatar,omitempty"`
	Birth                *string `form:"birth,omitempty"`
	Bio                  *string `form:"bio,omitempty"`
	Signature            *string `form:"signature,omitempty"`
	validator.Validator  `form:"-"`
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
	Name                *string `form:"name,omitempty"`
	ParentCategoryID    *int    `form:"parent_category_id,omitempty"`
	validator.Validator `form:"-"`
}

type threadForm struct {
	Title               *string `form:"title,omitempty"`
	Description         *string `form:"description,omitempty"`
	IsPublic            *bool   `form:"is_public,omitempty"`
	CategoryID          *int    `form:"category_id,omitempty"`
	validator.Validator `form:"-"`
}

type postForm struct {
	Content             *string `form:"content,omitempty"`
	ThreadID            *int    `form:"thread_id,omitempty"`
	ParentPostID        *int    `form:"parent_post_id,omitempty"`
	validator.Validator `form:"-"`
}

type tagForm struct {
	Name                *string `form:"name,omitempty"`
	AddThreads          []int   `form:"add_threads,omitempty"`
	RemoveThreads       []int   `form:"remove_threads,omitempty"`
	validator.Validator `form:"-"`
}

type reactToPostForm struct {
	Reaction            string   `form:"reaction"`
	AllowedReactions    []string `form:"-"`
	validator.Validator `form:"-"`
}

type friendResponseForm struct {
	Status              string   `form:"status"`
	FriendStatuses      []string `form:"-"`
	validator.Validator `form:"-"`
}
