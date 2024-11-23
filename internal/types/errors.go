package types

import "errors"

var (
    ErrUsernameTaken = errors.New("username already taken")
    ErrUserNotFound  = errors.New("user not found")
)
