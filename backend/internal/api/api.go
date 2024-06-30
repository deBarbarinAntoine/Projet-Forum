package api

import (
	"Projet-Forum/internal/validator"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"sync"
	"time"
)

const (
	userTokenSessionKey = "user_token"
	StatusFailedRequest = 0
)

var (
	client                = http.Client{Timeout: time.Second * 5}
	permittedMethods      = []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
	lock                  = &sync.Mutex{}
	ErrUnmarshallAPIError = errors.New("error unmarshalling API error response")
)

type API struct {
	url         string
	clientToken string
	pemKey      []byte
}

var apiInstance *API

func GetInstance(url, clientToken string, pemKey []byte) *API {
	lock.Lock()
	defer lock.Unlock()
	if apiInstance == nil {
		apiInstance = &API{
			url:         url,
			clientToken: clientToken,
			pemKey:      pemKey,
		}
	}
	return apiInstance
}

func (api *API) Get(userToken, endpoint string, query url.Values) ([]byte, int, error) {

	// building the url request
	urlRequest := fmt.Sprintf("%s/v1%s", api.url, endpoint)

	// creating the request
	req, err := http.NewRequest(http.MethodGet, urlRequest, nil)
	if err != nil {
		return nil, StatusFailedRequest, err
	}
	req.Header.Set("Accept", "application/json")

	// adding the query if necessary
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	// fetching the user token
	var userAuth string
	if userToken != "" {
		userAuth = fmt.Sprintf(",Bearer %s", userToken)
	}

	// setting the authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s%s", api.clientToken, userAuth))

	// sending the request
	res, err := client.Do(req)
	if err != nil {
		return nil, StatusFailedRequest, err
	}
	defer res.Body.Close()

	// reading the body of the response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, StatusFailedRequest, err
	}

	return body, res.StatusCode, nil
}

func (api *API) Request(userToken, method, endpoint string, body []byte, isEncrypted bool) ([]byte, int, error) {

	// checking the method
	if !slices.Contains(permittedMethods, method) {
		return nil, StatusFailedRequest, errors.New("invalid method")
	}

	// building the url request
	urlRequest := fmt.Sprintf("%s/v1%s", api.url, endpoint)

	// encrypting the body if necessary
	if isEncrypted {
		var err error
		body, err = api.encryptPEM(body)
		if err != nil {
			return nil, StatusFailedRequest, err
		}
	}

	// creating the request
	req, err := http.NewRequest(method, urlRequest, bytes.NewBuffer(body))
	if err != nil {
		return nil, StatusFailedRequest, err
	}

	// setting content related headers
	if body != nil {
		req.Header.Set("Content-Type", "application/json")

		if isEncrypted {
			req.Header.Set("X-Encryption", "RSA")
		}
	}
	req.Header.Set("Accept", "application/json")

	// fetching the user token
	var userAuth string
	if userToken != "" {
		userAuth = fmt.Sprintf(",Bearer %s", userToken)
	}

	// setting the authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s%s", api.clientToken, userAuth))

	// sending the request
	res, err := client.Do(req)
	if err != nil {
		return nil, StatusFailedRequest, err
	}
	defer res.Body.Close()

	// reading the body of the response
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, StatusFailedRequest, err
	}

	return body, res.StatusCode, nil
}

func GetErr(statusCode int, body []byte, v *validator.Validator) error {

	// No error (status 2XX)
	if statusCode < 400 {
		return nil
	}

	// Unprocessable Entity (field errors)
	if statusCode == 422 {
		var apiErr = make(map[string]map[string]string)
		err := json.Unmarshal(body, &apiErr)
		if err != nil {
			return ErrUnmarshallAPIError
		}
		v.FieldErrors = apiErr["errors"]
		return nil
	}

	// All other errors
	var apiErr = make(map[string]string)
	err := json.Unmarshal(body, &apiErr)
	if err != nil {
		return ErrUnmarshallAPIError
	}
	v.AddNonFieldError(apiErr["error"])
	return nil
}
