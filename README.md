# Projet-Forum

<!-- TOC -->
* [Projet-Forum](#projet-forum)
  * [API](#api)
    * [Introduction](#introduction)
    * [Authentication](#authentication)
      * [Tokens](#tokens)
      * [Procedure](#procedure)
      * [Sending credentials](#sending-credentials)
    * [Base URL](#base-url)
    * [Rate Limiting](#rate-limiting)
    * [Error Handling](#error-handling)
    * [Versioning](#versioning)
  * [Setting up the project for contributors](#setting-up-the-project-for-contributors)
<!-- TOC -->

---

## API

### Introduction

This API is used to power up an online tech oriented forum named **Threadive**.
It handles all data, logics and manipulations of the actual forum: the website is only like an HTML/CSS interface.

### Authentication

The authentication system of the API is token-based, using the `Bearer {token}` model with `Authorization` Header.

#### Tokens

> There are various levels of tokens in the API used for authentication:
> - `api_secret`: this token is constant and is used to authenticate a client
> - `client`: this token is given when a client is successfully registered using the `api_secret` token. It is also constant, but every client has its own `client` token
> - `authentication`: this token is used by the users that are registered and logged in the forum.

#### Procedure

Here is the procedure for a new client:
1. Send a request:

   > **Register the client:**
   > ```http
   > POST /v1/tokens/client HTTP/1.1
   > Authorization: Bearer API_SECRET
   > ```
   > - Body:
   > ```json
   > {
   > "username": "client-test",
   > "email": "client@test.com"
   > }
   > ```

2. The response will be of that kind:
    ```json
    {
      "client_token": {
          "token": "vY9bSu6tuNgyrIu+D8akW87+e74M4DnHadwph+gGAbY5rDk2ErT/iNd8Dos3lCR3PnHk68vWxA/vLOivBpJTjQ",
          "expiry": "2316-10-04T16:32:02.54895195+02:00"
      },
      "user": {
          "created_at": "2024-06-24T16:44:45Z",
          "email": "client@test.com",
          "id": 46,
          "username": "client-test"
      }
    }
    ```

3. Consequently, all requests to the API *must* contain the `client` token in the `Authorization: Bearer {client_token}` form.
4. For a user to manipulate/modify data in the API, an `authentication` token is needed. In that case, the **Header** will be built that way:
   ```
   METHOD /v1/path HTTP/1.1
   Authorization: Bearer {client_token},Bearer {authentication_token}
   ```

#### Sending credentials

> [!WARNING]
> To send credentials *securely*, the `JSON body` of the request needs to be encrypted using an `RSA public key`.
>

1. Fetch the `RSA public key`:
    ```http request
    GET /v1/tokens/public-key HTTP/1.1
   Authorization: Bearer API_SECRET
    ```
   The API will answer with the `public.pem` key with `Content-Type: application/x-pem-file` Header.
2. Use a similar function to encode the `JSON body` of the requests involving credentials (authentication, update user, reset password):
    ```Go
   // encryptPEM encrypts data (credentials marshalled in JSON)
   func encryptPEM(data []byte) ([]byte, error) {

	publicKeyBlock, _ := pem.Decode(publicKeyRSA)

	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), data)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
    }
   ```


### Base URL

The base URL for all requests to the API from the client is (example with `port 3000`):
```
http://localhost:3000/v1/
```

### Rate Limiting

The default rate limiting policy is generic (not IP nor user based):
- 100 requests in burst
- 50 requests/s

That means 50 requests/s in a pool of 100 max.

### Error Handling

Here are the possible errors sent back by the API:


> `400 Bad Request`
> ```json
> {
>   "error": "error_message"
> }
> ```

> `401 Unauthorized`
> ```json
> {
>   "error": "invalid authentication credentials"
> }
> ```
> **_OR_**
> ```json
> {
>   "error": "invalid or missing client token"
> }
> ```
> **_OR_**
> ```json
> {
>   "error": "invalid or missing authentication token"
> }
> ```
> **_OR_**
> ```json
> {
>   "error": "you must be authenticated in to access this resource"
> }
> ```

> `403 Forbidden`
> ```json
> {
>   "error": "your user account must be activated to access this resource"
> }
> ```
> **_OR_**
> ```json
> {
>   "error": "your user account doesn't have the necessary permissions to access this resource"
> }
> ```

> `404 Not Found`
> ```json
> {
>   "error": "the requested resource could not be found"
> }
> ```

> `405 Method Not Allowed`
> ```json
> {
>   "error": "the {method_name} method is not supported for this resource"
> }
> ```

> `409 Conflict`
> ```json
> {
>   "error": "unable to update the record due to an edit conflict, please try again"
> }
> ```

> `422 Unprocessable Entity`
> ```json
> {
>   "errors": {
>     "field_name": "error_message"
>     }
> }
> ```

> `429 Too Many Requests`
> ```json
> {
>   "error": "rate limit exceeded"
> }
> ```

> `500 Internal Server Error`
> ```json
> {
>   "error": "the server encountered a problem and could not process your request"
> }
> ```

### Versioning

For now only version 1 exists, but later versions will be available changing the URLs from `/v1/` to the wanted version (`/v2/` for example).

---

## Setting up the project for contributors


Run the `data/cmd/main.go` program to set the whole environment.

> [!WARNING]
> **You need MySQL and an administrator account** to properly follow the script!
> 

After that, you're ready to go.

Good programming everyone!