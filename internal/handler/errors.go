package handler

import "errors"

var (
	errInvalidInputBody   = errors.New("invalid input body")
	errFailedToGetUserID  = errors.New("failed to get user id")
	errInvalidParamID     = errors.New("invalid id param")
	errEmptyAuthHeader    = errors.New("empty auth header")
	errInvalidAuthHeader  = errors.New("invalid auth header")
	errEmptyToken         = errors.New("token is empty")
	errFailedToParseToken = errors.New("failed to parse token")
)
