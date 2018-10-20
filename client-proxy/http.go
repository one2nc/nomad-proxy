package main

import (
	"net/http"
	"io/ioutil"
	"bytes"
)

var statusMsgs = map[int]string{
	http.StatusOK: "200 OK",
	http.StatusInternalServerError: "500 Internal Server Error",
}

func newResponse(r *http.Request, code int, body []byte) *http.Response {
	return &http.Response{
		Status:        statusMsgs[code],
		StatusCode:    code,
		Proto:         r.Proto,
		ProtoMajor:    r.ProtoMajor,
		ProtoMinor:    r.ProtoMinor,
		Body:          ioutil.NopCloser(bytes.NewBuffer(body)),
		ContentLength: int64(len(body)),
		Request:       r,
		Header:        make(http.Header, 0), // Not sure how this will work
	}
}
