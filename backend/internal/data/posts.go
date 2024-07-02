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

type PostModel struct {
	uri         string
	endpoint    string
	clientToken string
	pemKey      []byte
}

func (m *PostModel) api() *api.API {
	return api.GetInstance(m.uri, m.clientToken, m.pemKey)
}

func (m *PostModel) Create(token string, post *Post, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"content":        post.Content,
		"thread":         post.Thread,
		"parent_post_id": post.IDParentPost,
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

		// retrieving the post
		var response = make(map[string]*Post)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return err
		}
		post = response["post"]
	}

	return nil
}

func (m *PostModel) Update(token string, post *Post, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"content":        post.Content,
		"thread":         post.Thread,
		"parent_post_id": post.IDParentPost,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d", m.endpoint, post.ID)

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

		// retrieving the post
		var response = make(map[string]*Post)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return err
		}
		post = response["post"]
	}

	return nil
}

func (m *PostModel) Delete(token string, id int, v *validator.Validator) error {

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

func (m *PostModel) Get(token string, query url.Values, v *validator.Validator) ([]*Post, Metadata, error) {

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
	var posts []*Post
	var metadata Metadata
	if v.Valid() {

		// retrieving the results
		var response = make(map[string]any)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, Metadata{}, err
		}
		var ok bool
		if posts, ok = response["posts"].([]*Post); !ok {
			return nil, Metadata{}, errors.New("invalid response from Posts")
		}
		if metadata, ok = response["_metadata"].(Metadata); !ok {
			return nil, Metadata{}, errors.New("invalid response from Metadata")
		}
	}

	return posts, metadata, nil
}

func (m *PostModel) GetByID(token string, id int, query url.Values, v *validator.Validator) (*Post, error) {

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
	var post *Post
	if v.Valid() {

		// retrieving the post
		var response = make(map[string]*Post)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, err
		}
		post = response["post"]
	}

	return post, nil
}

func (m *PostModel) React(token, reaction string, id int, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"reaction": reaction,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d/react", m.endpoint, id)

	// making the request
	res, status, err := m.api().Request(token, http.MethodPost, endpoint, reqBody, false)
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

func (m *PostModel) UpdateReaction(token, reaction string, id int, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"reaction": reaction,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d/react", m.endpoint, id)

	// making the request
	res, status, err := m.api().Request(token, http.MethodPatch, endpoint, reqBody, false)
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

func (m *PostModel) DeleteReaction(token string, id int, v *validator.Validator) error {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d/react", m.endpoint, id)

	// making the request
	res, status, err := m.api().Request(token, http.MethodPut, endpoint, nil, false)
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
