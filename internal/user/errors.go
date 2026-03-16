package user

import "errors"

var ErrorUserNotFound = errors.New("user not found")
var ErrorInvalidCredentials = errors.New("invalid credentials")
