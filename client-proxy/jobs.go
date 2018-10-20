package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const NomadToken = "X-Nomad-Token"

type jobPayload struct {
	Job map[string]interface{} `json:"Job"`
}

func parseJob(r *http.Request) (*jobPayload, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	b := &jobPayload{}
	if err := json.Unmarshal(body, b); err != nil {
		return nil, err
	}

	if b.Job == nil {
		return nil, fmt.Errorf("cannot parse body to Job")
	}

	return b, nil
}

// /v1/jobs overrides.
func jobs(r *http.Request) error {
	if r.Method == "GET" {
		val := r.URL.Query()
		val.Set("prefix", *jobPrefix)
		r.URL.RawQuery = val.Encode()
	}

	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		return nil
	}

	// add a check for 2fa
	if err := validateToken(r.URL.Path, r.Method, r.Header.Get(NomadToken)); err != nil {
		return err
	}

	b, err := parseJob(r)
	if err != nil {
		return err
	}

	jID := b.Job["ID"].(string)
	if !strings.HasPrefix(jID, *jobPrefix) {
		b.Job["ID"] = fmt.Sprintf("%v_%v", *jobPrefix, b.Job["ID"])
	}

	b.Job["Name"] = b.Job["ID"]
	if err := newBody(r, &b); err != nil {
		return err
	}

	return nil
}
