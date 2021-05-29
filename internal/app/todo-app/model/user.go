package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Name, validation.Length(3, 30)),
		validation.Field(&u.Email, is.Email),
		validation.Field(&u.Password, validation.Length(6, 30)),
	)
}
