package controllers

import (
	"Projet-Forum/internal/middlewares"
	"Projet-Forum/internal/models"
	"Projet-Forum/internal/utils"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

func indexHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	tmpl, err := template.ParseFiles(utils.Path + "templates/index.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "index", "indexHandlerGet")
	if err != nil {
		log.Fatalln(err)
	}
}

func indexHandlerPut(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	tmpl, err := template.ParseFiles(utils.Path + "templates/index.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	sessionID, _ := r.Cookie("updatedCookie")
	err = tmpl.ExecuteTemplate(w, "index", "indexHandlerPut"+sessionID.Value+"\nUsername: "+utils.SessionsData[sessionID.Value].Username+"\nIP address: "+utils.SessionsData[sessionID.Value].IpAddress)
	if err != nil {
		log.Fatalln(err)
	}
}

func loginHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	var message template.HTML
	if r.URL.Query().Has("err") {
		switch r.URL.Query().Get("err") {
		case "login":
			message = "<div class=\"message\">Wrong username or password!</div>"
		case "restricted":
			message = "<div class=\"message\">You need to login to access that area!</div>"
		}
	}
	tmpl, err := template.ParseFiles(utils.Path + "templates/login.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "login", message)
	if err != nil {
		log.Fatalln(err)
	}
}

func loginHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	credentials := models.Credentials{
		Username: r.FormValue("username"),
		Password: r.FormValue("password"),
	}
	if utils.CheckPwd(credentials) {
		utils.OpenSession(&w, credentials.Username, r)
		http.Redirect(w, r, "/home", http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/login?err=login", http.StatusSeeOther)
	}
}

func registerHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	var message template.HTML
	if r.URL.Query().Has("err") {
		switch r.URL.Query().Get("err") {
		case "username":
			message = "<div class=\"message\">Username already used!</div>"
		case "password":
			message = "<div class=\"message\">Both passwords need to be equal!</div>"
		case "email":
			message = "<div class=\"message\">Wrong email value!</div>"
		}
	}
	tmpl, err := template.ParseFiles(utils.Path + "templates/register.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "register", message)
	if err != nil {
		log.Fatalln(err)
	}
}

func registerHandlerPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	formValues := struct {
		username  string
		email     string
		password1 string
		password2 string
	}{
		username:  r.FormValue("username"),
		email:     strings.TrimSpace(strings.ToLower(r.FormValue("email"))),
		password1: r.FormValue("password1"),
		password2: r.FormValue("password2"),
	}
	_, exists := utils.SelectUser(formValues.username)
	switch {
	case exists:
		http.Redirect(w, r, "register?err=username", http.StatusSeeOther)
		return
	case formValues.password1 != formValues.password2:
		http.Redirect(w, r, "register?err=password", http.StatusSeeOther)
		return
	case !utils.CheckEmail(formValues.email):
		http.Redirect(w, r, "register?err=email", http.StatusSeeOther)
		return
	}
	hash, salt := utils.NewPwd(formValues.password1)
	newTempUser := models.TempUser{
		ConfirmID:    "",
		CreationTime: time.Now(),
		User: models.User{
			Id:       0,
			Username: formValues.username,
			HashedPwd: sql.NullString{
				String: hash,
				Valid:  true,
			},
			Salt: sql.NullString{
				String: salt,
				Valid:  true,
			},
			Email: formValues.email,
		},
	}
	utils.SendMail(&newTempUser, "creation")
	utils.TempUsers = append(utils.TempUsers, newTempUser)
	log.Printf("newTempUser: %#v\n", newTempUser)
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func homeHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	tmpl, err := template.ParseFiles(utils.Path + "templates/index.gohtml")
	if err != nil {
		log.Fatalln(err)
	}
	err = tmpl.ExecuteTemplate(w, "index", "homeHandlerGet --- Restricted area! ---")
	if err != nil {
		log.Fatalln(err)
	}
}

func logHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	w.Header().Set("Content-Type", "application/json")
	if r.URL.Query().Has("level") {
		json.NewEncoder(w).Encode(utils.FetchAttrLogs("level", r.URL.Query().Get("level")))
		return
	} else if r.URL.Query().Has("user") {
		json.NewEncoder(w).Encode(utils.FetchAttrLogs("user", r.URL.Query().Get("user")))
		return
	}
	json.NewEncoder(w).Encode(utils.RetrieveLogs())
}

func confirmHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	if r.URL.Query().Has("id") {
		id := r.URL.Query().Get("id")
		utils.PushTempUser(id)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}

func logoutHandlerGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	utils.Logout(&w, r)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())
	log.Println("HTTP Error", http.StatusNotFound)
	w.WriteHeader(http.StatusNotFound)
	utils.Logger.Warn("errorHandler", slog.Int("req_id", middlewares.LogId), slog.String("req_url", r.URL.String()), slog.Int("http_status", http.StatusNotFound))
	w.Write([]byte("Error " + fmt.Sprint(http.StatusNotFound) + " !"))
}

func createCategoryPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	var category models.Category
	var name string
	var parent string
	var idAuthor int

	if !category.Exists(r.FormValue("name")) {
		name = r.FormValue("name")
		if r.FormValue("parent") != "" {
			parent = r.FormValue("parent")
		}
		idAuthor = category.GetId("author")
	}

	// to avoid "not used" errors
	fmt.Println(name, parent, idAuthor)

	// fill post with data
	category.Create()
}

func createThreadPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	var thread models.Thread
	var name string
	var description string
	var public bool
	var idAuthor int
	var idCategory int

	if !thread.Exists(r.FormValue("name")) {
		name = r.FormValue("name")
		description = r.FormValue("description")
		if r.FormValue("public") == "public" {
			public = true
		} else {
			public = false
		}
		idAuthor = thread.GetId("author")
		idCategory = thread.GetId("category")
	}

	// to avoid "not used" errors
	fmt.Println(name, description, public, idAuthor, idCategory)

	// fill post with data
	thread.Create()
}

func createPostPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	var post models.Post
	var parent string

	content := r.FormValue("content")
	idAuthor := post.GetId("author")
	if post.Exists(r.FormValue("parent")) {
		parent = r.FormValue("parent")
	}
	idPost := post.GetId("thread")

	// to avoid "not used" errors
	fmt.Println(content, parent, idAuthor, idPost)

	// fill post with data
	post.Create()
}

func createTagPost(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	var tag models.Tag
	var name string
	var idAuthor int

	if !tag.Exists(r.FormValue("name")) {
		name = r.FormValue("name")
		idAuthor = tag.GetId("author")
	}

	// to avoid "not used" errors
	fmt.Println(name, idAuthor)

	// fill post with data
	tag.Create()
}

func threadGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	var thread models.Thread

	idThread := r.URL.Query().Get("id")
	thread.Fetch(idThread)

	// todo : execute template
}

func tagGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	var tag models.Tag

	idTag := r.URL.Query().Get("id")
	tag.Fetch(idTag)

	// todo : execute template
}

func categoryGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	var category models.Category

	idCategory := r.URL.Query().Get("id")
	category.Fetch(idCategory)

	// todo : execute template
}

func profileGet(w http.ResponseWriter, r *http.Request) {
	log.Println(utils.GetCurrentFuncName())

	var user models.User

	idUser := r.URL.Query().Get("id")
	user.Fetch(idUser)

	// todo : execute template
}
