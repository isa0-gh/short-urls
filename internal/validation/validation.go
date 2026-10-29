package validation

import (
	"net/url"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = newValidator()

type userInput struct {
	Username string `validate:"required,max=64"`
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8"`
}

type loginInput struct {
	Email    string `validate:"required,email"`
	Password string `validate:"required"`
}

type shortLinkInput struct {
	ShortCode string `validate:"required,max=64,shortcode"`
	URL       string `validate:"required,httpurl"`
}

type uuidInput struct {
	Value string `validate:"required,uuid"`
}

func newValidator() *validator.Validate {
	instance := validator.New()
	_ = instance.RegisterValidation("shortcode", func(fl validator.FieldLevel) bool {
		for _, character := range fl.Field().String() {
			if !(character >= 'a' && character <= 'z') &&
				!(character >= 'A' && character <= 'Z') &&
				!(character >= '0' && character <= '9') &&
				character != '-' && character != '_' {
				return false
			}
		}
		return true
	})
	_ = instance.RegisterValidation("httpurl", func(fl validator.FieldLevel) bool {
		parsed, err := url.ParseRequestURI(fl.Field().String())
		return err == nil && parsed.Host != "" &&
			(parsed.Scheme == "http" || parsed.Scheme == "https")
	})
	return instance
}

func User(username, email, password string) error {
	return validate.Struct(userInput{Username: username, Email: email, Password: password})
}

func Login(email, password string) error {
	return validate.Struct(loginInput{Email: email, Password: password})
}

func ShortLink(shortCode, destination string) error {
	return validate.Struct(shortLinkInput{ShortCode: shortCode, URL: destination})
}

func URL(value string) error {
	return validate.Struct(struct {
		URL string `validate:"required,httpurl"`
	}{URL: value})
}

func UUID(value string) error {
	return validate.Struct(uuidInput{Value: strings.TrimSpace(value)})
}
