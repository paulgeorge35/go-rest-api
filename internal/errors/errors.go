package errors

import "errors"

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
	ErrEmailAlreadyExists = errors.New("email already exists")
	ErrInvalidToken       = errors.New("invalid token")
	ErrTokenExpired       = errors.New("token expired")
	ErrUnauthorized       = errors.New("unauthorized")
	ErrFailedToCreateUser = errors.New("failed to create user")
	ErrFailedToSendEmail  = errors.New("failed to send email")
	ErrInvalidEmailOrPass = errors.New("invalid email or password")
	ErrSessionNotFound    = errors.New("session not found")
	ErrSessionExpired     = errors.New("session expired")
	ErrInvalidSession     = errors.New("invalid session")
)
