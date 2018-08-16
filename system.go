package main

import (
	"errors"
	"net/http"
)

// /v1/job overrides.
func system(r *http.Request) error {
	return errors.New("System calls are not allowed.")
}
