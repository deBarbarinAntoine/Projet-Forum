package main

import (
	"ForumAPI/internal/data"
	"ForumAPI/internal/validator"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

func (app *application) getFriendsByUser(user *data.User) error {

	sentFriends, receivedFriends, err := app.models.Users.GetFriendsByUserID(user.ID)
	if err != nil {
		return err
	}

	for _, friend := range sentFriends {
		switch friend.Status {
		case data.FriendStatus.Pending:
			user.Invitations.Sent = append(user.Invitations.Sent, friend)
		case data.FriendStatus.Accepted:
			user.Friends = append(user.Friends, friend)
		case data.FriendStatus.Rejected:
			continue
		default:
			return errors.New("invalid friend status")
		}
	}

	for _, friend := range receivedFriends {
		switch friend.Status {
		case data.FriendStatus.Pending:
			user.Invitations.Received = append(user.Invitations.Sent, friend)
		case data.FriendStatus.Accepted:
			user.Friends = append(user.Friends, friend)
		case data.FriendStatus.Rejected:
			continue
		default:
			return errors.New("invalid friend status")
		}
	}

	return nil
}

func (app *application) friendRequestHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	err = app.models.Users.RequestFriend(user, id)
	if err != nil {
		v := validator.New()
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		case errors.Is(err, data.ErrDuplicateFriend):
			v.AddError("friend", "existent relationship")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := envelope{
		"message": fmt.Sprintf("requested friend with id %d", id),
	}

	err = app.writeJSON(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) friendResponseHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var input struct {
		Status string `json:"status"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	v.StringCheck(input.Status, 2, 50, true, "status")

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	input.Status = strings.TrimSpace(strings.ToLower(input.Status))

	user := app.contextGetUser(r)

	switch input.Status {
	case data.FriendStatus.Rejected:
		err = app.models.Users.RejectFriend(id, user)
	case data.FriendStatus.Accepted:
		err = app.models.Users.AcceptFriend(id, user)
	default:
		app.badRequestResponse(w, r, errors.New("invalid status"))
		return
	}
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := envelope{
		"message": fmt.Sprintf("%sed friend with id %d", input.Status, id),
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) friendDeleteHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := app.contextGetUser(r)

	err = app.models.Users.RemoveFriend(user, id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	response := envelope{
		"message": fmt.Sprintf("removed friend with id %d", id),
	}

	err = app.writeJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
