package user

import "errors"

var ErrorUserNotFound = errors.New("User not found")
var ErrorInvalidCredentials = errors.New("Invalid credentials")
