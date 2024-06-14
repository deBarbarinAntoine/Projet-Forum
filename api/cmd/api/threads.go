package main

import "net/http"

func (app *application) getThreadsHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"threads": "get_threads"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) createThreadHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusCreated, envelope{"thread": "create_thread"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) getSingleThreadHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"thread": "get_single_thread"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateThreadHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"thread": "update_thread"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteThreadHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"thread": "delete_thread"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) followThreadHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusCreated, envelope{"thread": "follow_thread"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) unfollowThreadHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"thread": "unfollow_thread"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
