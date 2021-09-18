package customerrors

import "errors"

var (
	ErrorNotFound = errors.New("key not found")
	ErrorInternal = errors.New("internal error")
)
