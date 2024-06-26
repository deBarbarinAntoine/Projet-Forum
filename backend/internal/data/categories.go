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

type CategoryModel struct {
	uri         string
	endpoint    string
	clientToken string
	pemKey      []byte
}

func (m *CategoryModel) api() *api.API {
	return api.GetInstance(m.uri, m.clientToken, m.pemKey)
}

func (m *CategoryModel) Create(token string, category *Category, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"name":               category.Name,
		"parent_category_id": category.ParentCategory.ID,
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

		// retrieving the category
		var response = make(map[string]*Category)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return err
		}
		category = response["category"]
	}

	return nil
}

func (m *CategoryModel) Update(token string, category *Category, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"name":               category.Name,
		"parent_category_id": category.ParentCategory.ID,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d", m.endpoint, category.ID)

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

		// retrieving the category
		var response = make(map[string]*Category)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return err
		}
		category = response["category"]
	}

	return nil
}

func (m *CategoryModel) Delete(token string, id int, v *validator.Validator) error {

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

func (m *CategoryModel) Get(token string, query url.Values, v *validator.Validator) ([]*Category, Metadata, error) {

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
	var categories []*Category
	var metadata Metadata
	if v.Valid() {

		// retrieving the results
		var response = make(map[string]any)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, Metadata{}, err
		}
		var ok bool
		if categories, ok = response["categories"].([]*Category); !ok {
			return nil, Metadata{}, errors.New("invalid response from Categories")
		}
		if metadata, ok = response["_metadata"].(Metadata); !ok {
			return nil, Metadata{}, errors.New("invalid response from Metadata")
		}
	}

	return categories, metadata, nil
}

func (m *CategoryModel) GetByID(token string, id int, query url.Values, v *validator.Validator) (*Category, error) {

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
	var category *Category
	if v.Valid() {

		// retrieving the category
		var response = make(map[string]*Category)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, err
		}
		category = response["category"]
	}

	return category, nil
}
