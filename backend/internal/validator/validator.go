package validator

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"
	"unicode"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
)

type Validator struct {
	NonFieldErrors []string          `json:"non_field_errors"`
	FieldErrors    map[string]string `json:"field_errors"`
}

var EmailRX = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func New() *Validator {
	return &Validator{
		FieldErrors: make(map[string]string),
	}
}

func (v *Validator) Valid() bool {
	return len(v.NonFieldErrors) == 0 && len(v.FieldErrors) == 0
}

func (v *Validator) AddNonFieldError(message string) {
	v.NonFieldErrors = append(v.NonFieldErrors, message)
}

func (v *Validator) AddFieldError(key, message string) {
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

func (v *Validator) Errors() []byte {
	errors, err := json.Marshal(v)
	if err != nil {
		panic(err) // FIXME -> maybe change that to handle errors better if necessary
	}
	return errors
}

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

func (v *Validator) ValidateDate(date, fieldName string) {
	_, err := time.Parse("01/02/2006", date)
	if err != nil {
		v.AddFieldError(fieldName, "invalid date")
	}
}

func (v *Validator) CheckID(id int, fieldName string) {
	v.Check(id > 0, fieldName, "ID must be greater than 0")
}

func (v *Validator) ValidateEmail(email string) {
	v.StringCheck(email, 5, 150, true, "email")
	v.Check(Matches(email, EmailRX), "email", "must be a valid email address")
}

func (v *Validator) StringCheck(str string, min, max int, isMandatory bool, key string) {
	if isMandatory {
		v.Check(str != "", key, "must be provided")
	}
	v.Check(len(str) >= min, key, fmt.Sprintf("must be minimum %d bytes long", min))
	v.Check(len(str) <= max, key, fmt.Sprintf("must not be more than %d bytes long", max))
}

func (v *Validator) CheckPassword(password, key string) {

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
	v.Check(hasUpper, key, "must contain an uppercase character")
	v.Check(hasLower, key, "must contain a lowercase character")
	v.Check(hasNumber, key, "must contain a numeric character")
	v.Check(hasSpecial, key, "must contain a special character")
}

func (v *Validator) ValidatePassword(password string) {
	v.StringCheck(password, MinPasswordLength, MaxPasswordLength, true, "password")
	v.CheckPassword(password, "password")
}

func (v *Validator) ValidateRegisterPassword(password, confirmationPassword string) {
	v.StringCheck(password, MinPasswordLength, MaxPasswordLength, true, "password")
	v.CheckPassword(password, "password")
	v.Check(confirmationPassword != "", "confirmation_password", "must be provided")
	v.Check(password == confirmationPassword, "confirmation_password", "must be the same")
}

func (v *Validator) ValidateNewPassword(newPassword, confirmationPassword string) {
	v.StringCheck(newPassword, MinPasswordLength, MaxPasswordLength, true, "new_password")
	v.CheckPassword(newPassword, "new_password")
	v.Check(confirmationPassword != "", "confirmation_password", "must be provided")
	v.Check(newPassword == confirmationPassword, "confirmation_password", "must be the same")
}

func (v *Validator) ValidateToken(token string) {
	v.Check(token != "", "token", "must be provided")
	v.Check(len(token) == 86, "token", "must be 86 bytes long")
}

func NotBlank(fieldName string) bool {
	return strings.TrimSpace(fieldName) != ""
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}

func Unique[T comparable](values []T) bool {
	uniqueValues := make(map[T]bool)

	for _, value := range values {
		uniqueValues[value] = true
	}

	return len(values) == len(uniqueValues)
}
