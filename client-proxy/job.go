package main

import (
	"fmt"
	"net/http"
	"strings"
)

var EMPTY = ""

// /v1/job overrides.
func job(r *http.Request) error {
	// skip check if prefix is not passed.
	if *jobPrefix == EMPTY {
		return nil
	}
	if !strings.HasPrefix(r.URL.Path, "/v1/job/"+*jobPrefix) {
		return fmt.Errorf("jobs should begin with a prefix %v", *jobPrefix)
	}

	return nil
}
