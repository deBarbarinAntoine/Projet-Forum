package data

import (
	"Projet-Forum/internal/api"
	"Projet-Forum/internal/validator"
	"encoding/json"
	"net/http"
)

type CategoryModel struct {
	uri         string
	endpoint    string
	clientToken string
	pemKey      []byte
}

func (c *CategoryModel) api() *api.API {
	return api.GetInstance(c.uri, c.clientToken, c.pemKey)
}

func (c *CategoryModel) Create(token string, category *Category, v *validator.Validator) error {

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
	res, status, err := c.api().Request(token, http.MethodPost, c.endpoint, reqBody, false)
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
