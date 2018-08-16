package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

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
		return nil, fmt.Errorf("Cannot parse body to Job")
	}

	return b, nil
}

// /v1/jobs overrides.
func jobs(r *http.Request) error {
	if r.Method == "GET" {
		val := r.URL.Query()
		val.Set("prefix", *jobPrefix)
		r.URL.RawQuery = val.Encode()
	} else if r.Method == "POST" || r.Method == "PUT" {
		b, err := parseJob(r)
		if err != nil {
			return err
		}

		b.Job["ID"] = fmt.Sprintf("%v-%v", *jobPrefix, b.Job["ID"])
		b.Job["Name"] = b.Job["ID"]
		if err := newBody(r, &b); err != nil {
			return err
		}
	}

	return nil
}
