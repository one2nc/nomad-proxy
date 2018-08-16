package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type searchPayload struct {
	Prefix  string `json:"prefix"`
	Context string `json:"context"`
}

func parseSearch(r *http.Request) (*searchPayload, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	b := &searchPayload{}
	if err := json.Unmarshal(body, b); err != nil {
		return nil, err
	}

	return b, nil
}

// /v1/search overrides.
func search(r *http.Request) error {
	if r.Method != "POST" {
		return nil
	}

	b, err := parseSearch(r)
	if err != nil {
		return err
	}

	b.Prefix = *jobPrefix
	b.Context = "jobs"

	if err := newBody(r, &b); err != nil {
		return err
	}
	return nil
}
