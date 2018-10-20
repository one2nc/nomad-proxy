package main

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"testing"
	)

func TestNewBody(t *testing.T) {
		req := httptest.NewRequest("GET", "/v2/jobs", nil)
		if err := newBody(req, 1); err != nil {
			t.Fatal(err)
		}

		if req.ContentLength != 1 {
			t.Fatal("Content Length should have been 1")
		}
}

func TestModifyRequest(t *testing.T) {
	t.Run("Should not care for an unknown URL", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/v2/jobs", nil)

		length := req.ContentLength
		qlen := len(req.URL.RawQuery)

		if err := modifyRequest(req); err != nil {
			t.Fatal(err)
		}

		if length != req.ContentLength {
			t.Fatal("Content Length should not change")
		}

		if qlen != len(req.URL.RawQuery) {
			t.Fatal("Query Length should not change")
		}
	})

	t.Run("Known Jobs", func(t *testing.T) {
		t.Run("/v1/jobs", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/jobs", nil)

			qlen := len(req.URL.RawQuery)

			if err := modifyRequest(req); err != nil {
				t.Fatal(err)
			}

			if qlen == len(req.URL.RawQuery) {
				t.Fatal("Query Length should have changed")
			}
		})

		t.Run("/v1/search", func(t *testing.T) {
			body, err := ioutil.ReadFile("./testdata/search.json")
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(
				"POST",
				"/v1/search",
				ioutil.NopCloser(bytes.NewBuffer(body)),
			)

			length := req.ContentLength

			if err := modifyRequest(req); err != nil {
				t.Fatal(err)
			}

			if length == req.ContentLength {
				t.Fatal("Content Length should have changed")
			}
		})

		t.Run("/v1/job", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/job/hello", nil)

			if err := modifyRequest(req); err == nil {
				t.Fatal("Non prefix Get should raise an error")
			}
		})

		t.Run("/v1/system", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/system", nil)

			if err := modifyRequest(req); err == nil {
				t.Fatal("System calls are not allowed")
			}
		})
	})
}