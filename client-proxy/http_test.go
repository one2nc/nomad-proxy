package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
)
func TestNewResponse(t *testing.T) {
	req := httptest.NewRequest("GET", "/v2/jobs", nil)
	body := "test response"
	code := http.StatusOK
	resp := newResponse(req, code, []byte(body))

	assert.Equal(t, statusMsgs[code], resp.Status)
	assert.Equal(t, code, resp.StatusCode)
	assert.Equal(t, req.Proto, resp.Proto)
	assert.Equal(t, req.ProtoMinor, resp.ProtoMinor)
	assert.Equal(t, req.ProtoMajor, resp.ProtoMajor)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.Nil(t, err)

	assert.Equal(t, body, string(respBody))
}
