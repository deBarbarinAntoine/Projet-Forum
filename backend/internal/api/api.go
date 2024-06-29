package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"time"
)

const (
	userTokenSessionKey = "user_token"
)

var (
	client           = http.Client{Timeout: time.Second * 5}
	permittedMethods = []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete}
)

func (app *main.application) requestAPI(ctx context.Context, method, endpoint string, query url.Values, body []byte, isEncrypted bool) ([]byte, *int, error) {

	// checking the method
	if !slices.Contains(permittedMethods, method) {
		return nil, nil, errors.New("invalid method")
	}

	// building the url request
	urlRequest := fmt.Sprintf("%s/v1%s", app.config.apiURL, endpoint)

	// encrypting the body if necessary
	if isEncrypted {
		var err error
		body, err = app.encryptPEM(body)
		if err != nil {
			return nil, nil, err
		}
	}

	// creating the request
	req, err := http.NewRequest(method, urlRequest, bytes.NewBuffer(body))
	if err != nil {
		return nil, nil, err
	}

	// setting content related headers
	if body != nil {
		req.Header.Set("Content-Type", "application/json")

		if isEncrypted {
			req.Header.Set("X-Encryption", "RSA")
		}
	}
	req.Header.Set("Accept", "application/json")

	// adding the query if necessary
	if query != nil {
		req.URL.RawQuery = query.Encode()
	}

	// fetching the user token
	userAuth := app.sessionManager.GetString(ctx, userTokenSessionKey)
	if userAuth != "" {
		userAuth = fmt.Sprintf(",Bearer %s", userAuth)
	}

	// setting the authorization header
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s%s", app.config.clientToken, userAuth))

	// sending the request
	res, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer res.Body.Close()

	// reading the body of the response
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return nil, nil, err
	}

	return body, &res.StatusCode, nil
}
