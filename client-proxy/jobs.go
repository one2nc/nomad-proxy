package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	NomadToken  = "X-Nomad-Token"
	Datacenters = "Datacenters"
	Prefix      = "prefix"
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
		return nil, fmt.Errorf("cannot parse body to Job")
	}

	return b, nil
}

// /v1/jobs overrides.
func jobs(r *http.Request) error {
	if r.Method == "GET" {
		val := r.URL.Query()
		existingPrefix := val.Get(Prefix)
		prefix := *jobPrefix
		if strings.HasPrefix(existingPrefix, *jobPrefix) {
			prefix = existingPrefix
		}

		val.Set("prefix", prefix)
		r.URL.RawQuery = val.Encode()
		return nil
	}

	if r.Method == http.MethodDelete {
		return nil
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

	// Validate datacenter
	if b.Job[Datacenters] == nil {
		return errors.New("datacenter is missing in job")
	}

	datacentersRaw := b.Job[Datacenters].([]interface{})
	if len(datacentersRaw) > 1 {
		return errors.New("only 1 datacenter is supported")
	}

	dc := datacentersRaw[0].(string)
	if dc != *datacenter {
		return fmt.Errorf("invalid datacenter in job, should be %v", *datacenter)
	}

	if err := newBody(r, &b); err != nil {
		return err
	}

	return nil
}
