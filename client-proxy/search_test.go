package main

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearch(t *testing.T) {
	for _, method := range []string{"GET", "PUT", "DELETE"} {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/v1/search", nil)
			length := req.ContentLength

			if err := search(req); err != nil {
				t.Fatal(err)
			}

			if length != req.ContentLength {
				t.Fatal("Content Length should not change")
			}
		})
	}

	t.Run("POST", func(t *testing.T) {
		t.Run("Should prepend prefix in Search Payload", func(t *testing.T) {
			body, err := ioutil.ReadFile("./testdata/search.json")
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(
				"POST",
				"/v1/search",
				ioutil.NopCloser(bytes.NewBuffer(body)),
			)

			if err := search(req); err != nil {
				t.Fatal(err)
			}

			b, err := parseSearch(req)
			if err != nil {
				t.Fatal(err)
			}

			if !strings.HasPrefix(b.Prefix, testPrefix) {
				t.Fatal("Should have altered Job Prefix")
			}

			if !strings.HasPrefix(b.Context, "jobs") {
				t.Fatal("Should have altered Job Name Prefix")
			}
		})

		t.Run("Should override malformed Job payload", func(t *testing.T) {
			body, err := ioutil.ReadFile("./testdata/jobs.json")
			if err != nil {
				t.Fatal(err)
			}

			req := httptest.NewRequest(
				"POST",
				"/v1/search",
				ioutil.NopCloser(bytes.NewBuffer(body)),
			)

			if err := search(req); err != nil {
				t.Fatal(err)
			}

			b, err := parseSearch(req)
			if err != nil {
				t.Fatal(err)
			}

			if !strings.HasPrefix(b.Prefix, testPrefix) {
				t.Fatal("Should have altered Job Prefix")
			}

			if !strings.HasPrefix(b.Context, "jobs") {
				t.Fatal("Should have altered Job Name Prefix")
			}
		})
	})
}
