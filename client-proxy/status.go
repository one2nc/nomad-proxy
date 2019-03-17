package main

import (
	"encoding/json"
	"net/http"
)

func status(r *http.Request) (*http.Response, error) {
	resp := map[string]interface{}{
		"prefix":     *jobPrefix,
		"version":    Version,
		"datacenter": *datacenter,
	}

	body, err := json.Marshal(&resp)
	if err != nil {
		return nil, err
	}

	return newResponse(r, http.StatusOK, body), nil
}
