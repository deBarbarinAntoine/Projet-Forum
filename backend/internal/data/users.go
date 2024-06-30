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

type UserModel struct {
	uri         string
	endpoint    string
	clientToken string
	pemKey      []byte
}

func (m *UserModel) api() *api.API {
	return api.GetInstance(m.uri, m.clientToken, m.pemKey)
}

func (m *UserModel) Create(token string, user *User, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"name":     user.Name,
		"email":    user.Email,
		"password": user.Password,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// making the request
	res, status, err := m.api().Request(token, http.MethodPost, m.endpoint, reqBody, true)
	if err != nil {
		return err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return err
	}
	if v.Valid() {

		// retrieving the user
		var response = make(map[string]*User)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return err
		}
		user = response["user"]
	}

	return nil
}

func (m *UserModel) Update(token, previousPassword string, user *User, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"username":     user.Name,
		"email":        user.Email,
		"password":     previousPassword,
		"new_password": user.Password,
		"avatar":       user.Avatar,
		"birth":        user.BirthDate,
		"bio":          user.Bio,
		"signature":    user.Signature,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%d", m.endpoint, user.ID)

	// making the request
	res, status, err := m.api().Request(token, http.MethodPut, endpoint, reqBody, true)
	if err != nil {
		return err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return err
	}
	if v.Valid() {

		// retrieving the user
		var response = make(map[string]*User)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return err
		}
		user = response["user"]
	}

	return nil
}

func (m *UserModel) Delete(token string, id string, v *validator.Validator) error {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%s", m.endpoint, id)

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

func (m *UserModel) Get(token string, query url.Values, v *validator.Validator) ([]*User, Metadata, error) {

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
	var users []*User
	var metadata Metadata
	if v.Valid() {

		// retrieving the results
		var response = make(map[string]any)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, Metadata{}, err
		}
		var ok bool
		if users, ok = response["users"].([]*User); !ok {
			return nil, Metadata{}, errors.New("invalid response from Users")
		}
		if metadata, ok = response["_metadata"].(Metadata); !ok {
			return nil, Metadata{}, errors.New("invalid response from Metadata")
		}
	}

	return users, metadata, nil
}

func (m *UserModel) GetByID(token string, id string, query url.Values, v *validator.Validator) (*User, error) {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/%s", m.endpoint, id)

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
	var user *User
	if v.Valid() {

		// retrieving the user
		var response = make(map[string]*User)
		err = json.Unmarshal(res, &response)
		if err != nil {
			return nil, err
		}
		user = response["user"]
	}

	return user, nil
}

// FIXME -> ADD FRIENDS RELATED FUNCTIONS
