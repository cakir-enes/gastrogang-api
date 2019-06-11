package storage

import "errors"

var ConnectionFailed = errors.New("Connection failed.")
var UserAlreadyExists = errors.New("Username already exists")
var UserDoesntExist = errors.New("User not found")
