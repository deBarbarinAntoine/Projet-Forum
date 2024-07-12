package data

import (
	"Projet-Forum/internal/api"
	"Projet-Forum/internal/validator"
	"encoding/json"
	"fmt"
	"net/http"
)

type TokenModel struct {
	uri         string
	endpoint    string
	clientToken string
	pemKey      []byte
}

func (m *TokenModel) api() *api.API {
	return api.GetInstance(m.uri, m.clientToken, m.pemKey)
}

func (m *TokenModel) Authenticate(body map[string]string, v *validator.Validator) (*Tokens, error) {

	// converting the body to JSON format
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/authentication", m.endpoint)

	// making the request
	res, status, err := m.api().Request("", http.MethodPost, endpoint, reqBody, true)
	if err != nil {
		return nil, err
	}

	// checking for errors
	err = api.GetErr(status, res, v)
	if err != nil {
		return nil, err
	}
	var tokens *Tokens
	if v.Valid() {

		// retrieving the tokens
		err = json.Unmarshal(res, &tokens)
		if err != nil {
			return nil, err
		}
	}

	return tokens, nil
}

func (m *TokenModel) Refresh(token string, tokens *Tokens, v *validator.Validator) error {

	// creating the request body
	body := envelope{
		"token": tokens.Refresh.Token,
	}
	reqBody, err := json.Marshal(body)
	if err != nil {
		return err
	}

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/refresh", m.endpoint)

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
	if v.Valid() {

		// retrieving the tokens
		err = json.Unmarshal(res, &tokens)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *TokenModel) Logout(token string, v *validator.Validator) error {

	// building the endpoint's specific URL
	endpoint := fmt.Sprintf("%s/revoke/me", m.endpoint)

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
