package main

import (
	"database/sql"
	"errors"
	"net/http"
)

func (app *application) getPopularHandler(w http.ResponseWriter, r *http.Request) {

	// get popular tags
	tags, err := app.models.Tags.GetPopular()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// get popular threads
	threads, err := app.models.Threads.GetPopular()
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	response := envelope{
		"popular": envelope{
			"tags":    tags,
			"threads": threads,
		},
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

//func (app *application) getRecommendations(w http.ResponseWriter, r *http.Request) {
//	err := app.writeJSON(w, http.StatusOK, envelope{"recommendations": "get_recommendations"}, nil)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//	}
//}
//
//func (app *application) searchHandler(w http.ResponseWriter, r *http.Request) {
//	err := app.writeJSON(w, http.StatusOK, envelope{"search": "results"}, nil)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//	}
//}
