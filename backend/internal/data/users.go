package data

import "net/http"

type UserModel struct {
	endpoint    string
	clientToken string
	pemKey      []byte
}

func (m *UserModel) Exists(id int) (bool, error) {
	var exists bool

	// TODO checking whether the user exists or not
	if id == 1 {
		exists = true
	}

	return exists, nil
}

func (m *UserModel) Register(email, password string) error {
	// TODO registering user
	return nil
}

func (m *UserModel) Activate(r *http.Request) error {
	// TODO activate user
	return nil
}

func (m *UserModel) Authenticate(email, password string) (int, error) {

	// TODO authentication with email and password
	if email == "thorgdar@gmail.com" {
		return 1, nil
	}

	return 0, ErrInvalidCredentials

}
