package main

import (
	"bytes"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

const testPrefix = "testjob"

func TestJobs(t *testing.T) {
	t.Run("DELETE", func(t *testing.T) {
		req := httptest.NewRequest("DELETE", "/", nil)
		length := req.ContentLength

		if err := jobs(req); err != nil {
			t.Fatal(err)
		}

		if length != req.ContentLength {
			t.Fatal("Content Length should not change")
		}
	})

	t.Run("GET", func(t *testing.T) {
		t.Run("When request has no query params", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/jobs", nil)
			if err := jobs(req); err != nil {
				t.Fatal(err)
			}

			if req.URL.Query().Get("prefix") != "testjob" {
				t.Fatal("Should have attached prefix to the query params")
			}
		})

		t.Run("When request has others params", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/jobs?foo=bar", nil)
			if err := jobs(req); err != nil {
				t.Fatal(err)
			}

			if req.URL.Query().Get("prefix") != "testjob" {
				t.Fatal("Should have attached prefix to the query params")
			}

			if req.URL.Query().Get("foo") != "bar" {
				t.Fatal("Should not lose existing params")
			}
		})

		t.Run("When request has prefix params", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/jobs?prefix=balloo", nil)
			if err := jobs(req); err != nil {
				t.Fatal(err)
			}

			if req.URL.Query().Get("prefix") != testPrefix {
				t.Fatal("Should override prefix to the query params")
			}
		})
	})

	for _, m := range []string{"PUT", "POST"} {
		t.Run(m, func(t *testing.T) {
			t.Run("Should prepend Id in Job Payload", func(t *testing.T) {
				body, err := ioutil.ReadFile("./testdata/jobs.json")
				if err != nil {
					t.Fatal(err)
				}

				req := httptest.NewRequest(
					m,
					"/v1/jobs",
					ioutil.NopCloser(bytes.NewBuffer(body)),
				)

				if err := jobs(req); err != nil {
					t.Fatal(err)
				}

				b, err := parseJob(req)
				if err != nil {
					t.Fatal(err)
				}

				if !strings.HasPrefix(b.Job["ID"].(string), testPrefix) {
					t.Fatal("Should have altered Job Prefix")
				}

				if !strings.HasPrefix(b.Job["Name"].(string), testPrefix) {
					t.Fatal("Should have altered Job Name Prefix")
				}
			})

			t.Run("Should not touch malformed Job payload", func(t *testing.T) {
				body, err := ioutil.ReadFile("./testdata/search.json")
				if err != nil {
					t.Fatal(err)
				}

				req := httptest.NewRequest(
					m, "/v1/jobs", ioutil.NopCloser(bytes.NewBuffer(body)),
				)

				if err := jobs(req); err == nil {
					t.Fatal("Should have returned an error")
				} else if err.Error() != "Cannot parse body to Job" {
					t.Fatal("Incorrect error string")
				}
			})
		})
	}
}

func TestMain(t *testing.M) {
	*jobPrefix = testPrefix
	os.Exit(t.Run())
}
