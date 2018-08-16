package main

import (
	"fmt"
	"net/http"
	"strings"
)

// /v1/job overrides.
func job(r *http.Request) error {
	if !strings.HasPrefix(r.URL.Path, "/v1/job/"+*jobPrefix) {
		return fmt.Errorf("Jobs should begin with a prefix %v", *jobPrefix)
	}
	return nil
}
