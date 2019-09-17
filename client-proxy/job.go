package main

import (
	"fmt"
	"net/http"
	"strings"
)

// /v1/job overrides.
func job(r *http.Request) error {
	// check prefix only if prefix is required, else skip.
	if !*skipPrefix && !strings.HasPrefix(r.URL.Path, "/v1/job/"+*jobPrefix) {
		return fmt.Errorf("jobs should begin with a prefix %v", *jobPrefix)
	}

	return nil
}
