// apperrors package with app's errors
package apperrors

import (
	"errors"
)

var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrNoSuchUser       = errors.New("no such user")
	ErrUserAlreadyExist = errors.New("user already exists")
)
