package config

import "errors"

var (
	errPostgresPasswordIsEmpty = errors.New("postgres password from env is empty")
	errSigningKeyIsEmpty       = errors.New("signing key from env is empty")
	errSaltIsEmpty             = errors.New("salt from env is empty")
)
