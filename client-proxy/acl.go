package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var sampleResponse = `{
  "AccessorID": "aa534e09-6a07-0a45-2295-a7f77063d429",
  "SecretID": "%s",
  "Name": "management token",
  "Type": "management",
  "Global": true,
  "CreateTime": "2017-08-23T23:25:41.429154233Z",
  "CreateIndex": 52,
  "ModifyIndex": 64
}`

func acl(r *http.Request) (*http.Response, error) {
	if r.Method != http.MethodGet {
		return nil, nil
	}

	var bag interface{}
	respBody := fmt.Sprintf(sampleResponse, r.Header.Get(NomadToken))
	if err := json.Unmarshal([]byte(respBody), &bag); err != nil {
		return nil, err
	}

	body, err := json.Marshal(bag)
	if err != nil {
		return nil, err
	}

	return newResponse(r, http.StatusOK, body), nil
}
