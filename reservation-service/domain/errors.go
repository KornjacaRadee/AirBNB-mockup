package domain

import "errors"

var (
	errAccommodationNotFound error = errors.New("accommodation not found")
	errUserNotFound          error = errors.New("user not found")
	errInvalidCredentials    error = errors.New("incorrect username or password")
	errInvalidToken          error = errors.New("token invalid")
	errUnauthorized          error = errors.New("unauthorized")
)

func ErrAccommodationNotFound() error {
	return errAccommodationNotFound
}

func ErrUserNotFound() error {
	return errUserNotFound
}

func ErrInvalidCredentials() error {
	return errInvalidCredentials
}

func ErrInvalidToken() error {
	return errInvalidToken
}

func ErrUnauthorized() error {
	return errUnauthorized
}
