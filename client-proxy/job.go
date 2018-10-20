package main

import (
	"fmt"
	"net/http"
	"strings"
)

// /v1/job overrides.
func job(r *http.Request) error {
	if !strings.HasPrefix(r.URL.Path, "/v1/job/"+*jobPrefix) {
		return fmt.Errorf("jobs should begin with a prefix %v", *jobPrefix)
	}

	if r.Method != http.MethodDelete {
		return nil
	}

	return validateToken(r.URL.Path, r.Method, r.Header.Get(NomadToken))
}
