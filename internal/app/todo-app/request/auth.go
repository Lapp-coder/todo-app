package request

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type SignUp struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SignIn struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (su *SignUp) Validate() error {
	return validation.ValidateStruct(
		su,
		validation.Field(&su.Name, validation.Length(2, 30)),
		validation.Field(&su.Email, is.Email),
		validation.Field(&su.Password, validation.Length(6, 30)),
	)
}

func (si *SignIn) Validate() error {
	return validation.ValidateStruct(
		si,
		validation.Field(&si.Email, is.Email),
		validation.Field(&si.Password, validation.Length(6, 30)),
	)
}
