package validator

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"unicode/utf8"
)

const (
	MinPasswordLength = 8
	MaxPasswordLength = 72
)

type Validator struct {
	NonFieldErrors []string
	FieldErrors    map[string]string
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

func (v *Validator) Check(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
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

func (v *Validator) ValidatePassword(password string) {
	v.StringCheck(password, MinPasswordLength, MaxPasswordLength, true, "password")
}

func (v *Validator) ValidateRegisterPassword(password, confirmationPassword string) {
	v.StringCheck(password, MinPasswordLength, MaxPasswordLength, true, "password")
	v.Check(confirmationPassword != "", "confirmation_password", "must be provided")
	v.Check(password == confirmationPassword, "confirmation_password", "must be the same")
}

func (v *Validator) ValidateNewPassword(newPassword, confirmationPassword string) {
	v.StringCheck(newPassword, MinPasswordLength, MaxPasswordLength, true, "new_password")
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

func MinChars(value string, n int) bool {
	return utf8.RuneCountInString(value) >= n
}

func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

func Matches(value string, rx *regexp.Regexp) bool {
	return rx.MatchString(value)
}

func PermittedValue[T comparable](value T, permittedValues ...T) bool {
	return slices.Contains(permittedValues, value)
}
