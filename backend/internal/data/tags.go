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

type TagModel struct {
	uri         string
	endpoint    string
	clientToken string
	pemKey      []byte
}

func (m *TagModel) api() *api.API {
	return api.GetInstance(m.uri, m.clientToken, m.pemKey)
}

func (m *TagModel) Create(token string, tag *Tag, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"name":    tag.Name,
		"threads": tag.Threads,
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

		// retrieving the tag
		var response = make(map[string]*Tag)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return err
		}
		tag = response["tag"]
		if tag.ID < 1 {
			return errors.New("invalid tag id")
		}
	}

	return nil
}

func (m *TagModel) Update(token string, id int, body []byte, v *validator.Validator) (*Tag, error) {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d", m.endpoint, id)

	// making the request
	res, status, err := m.api().Request(token, http.MethodPut, endpoint, body, false)
	if err != nil {
		return nil, err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return nil, err
	}
	var tag *Tag
	if v.Valid() {

		// retrieving the tag
		var response = make(map[string]*Tag)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, err
		}
		tag = response["tag"]
		if tag.ID < 1 {
			return nil, errors.New("invalid tag id")
		}
	}

	return tag, nil
}

func (m *TagModel) Delete(token string, id int, v *validator.Validator) error {

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

func (m *TagModel) Get(token string, query url.Values, v *validator.Validator) ([]*Tag, Metadata, error) {

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
	var tags []*Tag
	var metadata Metadata
	if v.Valid() {

		// retrieving the results
		var response = make(map[string]any)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, Metadata{}, err
		}
		err = api.UnmarshallSlice(response["tags"], &tags)
		if err != nil {
			return nil, Metadata{}, err
		}
		err = api.Unmarshall(response["_metadata"], &metadata)
		if err != nil {
			return nil, Metadata{}, err
		}
	}

	return tags, metadata, nil
}

func (m *TagModel) GetByID(token string, id int, query url.Values, v *validator.Validator) (*Tag, error) {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d", m.endpoint, id)

	// adding the threads in the query string
	if query == nil {
		query = make(url.Values)
	}
	query.Add("includes[]", "threads")

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
	var tag *Tag
	if v.Valid() {

		// retrieving the tag
		var response = make(map[string]*Tag)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, err
		}
		tag = response["tag"]
	}

	return tag, nil
}

func (m *TagModel) GetPopular(token string, v *validator.Validator) ([]*Tag, []*Thread, error) {

	// making the request
	res, status, err := m.api().Get(token, "/popular", nil)
	if err != nil {
		return nil, nil, err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return nil, nil, err
	}
	var tags []*Tag
	var threads []*Thread
	if v.Valid() {

		// retrieving the tags and threads
		var response = make(map[string]map[string]any)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, nil, err
		}
		err = api.UnmarshallSlice(response["popular"]["tags"], &tags)
		if err != nil {
			return nil, nil, err
		}
		err = api.UnmarshallSlice(response["popular"]["threads"], &threads)
		if err != nil {
			return nil, nil, err
		}
	}

	return tags, threads, nil
}

func (m *TagModel) Follow(token string, id int, v *validator.Validator) error {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d/follow", m.endpoint, id)

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

func (m *TagModel) Unfollow(token string, id int, v *validator.Validator) error {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d/follow", m.endpoint, id)

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
