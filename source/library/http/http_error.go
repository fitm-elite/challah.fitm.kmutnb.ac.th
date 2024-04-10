package http

import "errors"

var (
	ErrorDoRequest error = errors.New("Unable to create a request to reach that destination.")
)
