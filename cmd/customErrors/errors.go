package customErrors

import "errors"

var (
	ErrorNotFound = errors.New("Key Not Found")
	ErrorInternal = errors.New("Internal Error")
)
