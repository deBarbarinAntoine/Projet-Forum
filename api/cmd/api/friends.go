package main

import "net/http"

func (app *application) friendRequestHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"friend": "friend_added"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) friendResponseHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"friend": "friend_accepted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) friendDeleteHandler(w http.ResponseWriter, r *http.Request) {
	err := app.writeJSON(w, http.StatusOK, envelope{"friend": "friend_removed"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
