package data

import (
	"Projet-Forum/internal/api"
	"Projet-Forum/internal/validator"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type ThreadModel struct {
	uri         string
	endpoint    string
	clientToken string
	pemKey      []byte
}

func (m *ThreadModel) api() *api.API {
	return api.GetInstance(m.uri, m.clientToken, m.pemKey)
}

func (m *ThreadModel) Create(token string, thread *Thread, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"title":       thread.Title,
		"description": thread.Description,
		"is_public":   thread.IsPublic,
		"category_id": thread.Category.ID,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// making the request
	res, status, err := m.api().Request(token, http.MethodPost, m.endpoint, reqBody, false)
	if err != nil {
		return err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return err
	}
	if v.Valid() {

		// retrieving the thread
		var response = make(map[string]*Thread)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return err
		}
		thread = response["thread"]
	}

	return nil
}

func (m *ThreadModel) Update(token string, thread *Thread, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"title":       thread.Title,
		"description": thread.Description,
		"is_public":   thread.IsPublic,
		"category_id": thread.Category.ID,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d", m.endpoint, thread.ID)

	// making the request
	res, status, err := m.api().Request(token, http.MethodPut, endpoint, reqBody, false)
	if err != nil {
		return err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return err
	}
	if v.Valid() {

		// retrieving the thread
		var response = make(map[string]*Thread)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return err
		}
		thread = response["thread"]
	}

	return nil
}

func (m *ThreadModel) Delete(token string, id int, v *validator.Validator) error {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d", m.endpoint, id)

	// making the request
	res, status, err := m.api().Request(token, http.MethodDelete, endpoint, nil, false)
	if err != nil {
		return err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return err
	}

	return nil
}

func (m *ThreadModel) Get(token string, query url.Values, v *validator.Validator) ([]*Thread, Metadata, error) {

	// making the request
	res, status, err := m.api().Get(token, m.endpoint, query)
	if err != nil {
		return nil, Metadata{}, err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return nil, Metadata{}, err
	}
	var threads []*Thread
	var metadata Metadata
	if v.Valid() {

		// retrieving the results
		var response = make(map[string]any)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, Metadata{}, err
		}
		var ok bool
		if threads, ok = response["threads"].([]*Thread); !ok {
			return nil, Metadata{}, errors.New("invalid response from Threads")
		}
		if metadata, ok = response["_metadata"].(Metadata); !ok {
			return nil, Metadata{}, errors.New("invalid response from Metadata")
		}
	}

	return threads, metadata, nil
}

func (m *ThreadModel) GetByID(token string, id int, query url.Values, v *validator.Validator) (*Thread, error) {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d", m.endpoint, id)

	// making the request
	res, status, err := m.api().Get(token, endpoint, query)
	if err != nil {
		return nil, err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return nil, err
	}
	var thread *Thread
	if v.Valid() {

		// retrieving the thread
		var response = make(map[string]*Thread)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, err
		}
		thread = response["thread"]
	}

	return thread, nil
}

func (m *TagModel) AddToFavorite(token string, id int, v *validator.Validator) error {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d/favorite", m.endpoint, id)

	// making the request
	res, status, err := m.api().Request(token, http.MethodPost, endpoint, nil, false)
	if err != nil {
		return err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return err
	}

	return nil
}

func (m *TagModel) RemoveFromFavorite(token string, id int, v *validator.Validator) error {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d/favorite", m.endpoint, id)

	// making the request
	res, status, err := m.api().Request(token, http.MethodDelete, endpoint, nil, false)
	if err != nil {
		return err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return err
	}

	return nil
}
