package validator

import (
	"fmt"
	"regexp"
	"slices"
	"unicode"
)

var (
	EmailRX        = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	UserByIDValues = []string{"following_tags", "favorite_threads", "categories_owned", "tags_owned", "threads_owned", "posts", "reactions", "friends"}
)

type Validator struct {
	Errors map[string]string
}

func New() *Validator {
	return &Validator{Errors: make(map[string]string)}
}

func (v *Validator) Valid() bool {
	return len(v.Errors) == 0
}

func (v *Validator) AddError(key, message string) {
	if _, exists := v.Errors[key]; !exists {
		v.Errors[key] = message
	}
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddError(key, message)
	}
}

func (v *Validator) StringCheck(str string, min, max int, isMandatory bool, key string) {
	if isMandatory {
		v.Check(str != "", key, "must be provided")
	}
	v.Check(len(str) >= min, key, fmt.Sprintf("must be minimum %d bytes long", min))
	v.Check(len(str) <= max, key, fmt.Sprintf("must not be more than %d bytes long", max))
}

func (v *Validator) CheckPassword(password string) {

	// setting booleans to check the criteria
	var (
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)

	// checking every character
	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	// adding errors if needed
	v.Check(hasUpper, "password", "must contain an uppercase character")
	v.Check(hasLower, "password", "must contain a lowercase character")
	v.Check(hasNumber, "password", "must contain a numeric character")
	v.Check(hasSpecial, "password", "must contain a special character")
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
