package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

			if req.URL.Query().Get("prefix") != *jobPrefix {
				t.Fatal("Should have attached prefix to the query params")
			}
		})

		t.Run("When request has others params", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/jobs?foo=bar", nil)
			if err := jobs(req); err != nil {
				t.Fatal(err)
			}

			if req.URL.Query().Get("prefix") != *jobPrefix {
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

			if req.URL.Query().Get("prefix") != *jobPrefix {
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

				if !strings.HasPrefix(b.Job["ID"].(string), *jobPrefix) {
					t.Fatal("Should have altered Job Prefix")
				}

				if !strings.HasPrefix(b.Job["Name"].(string), *jobPrefix) {
					t.Fatal("Should have altered Job Name Prefix")
				}
			})

			t.Run("Should not prepend prefix if already exists", func(t *testing.T) {
				body, err := ioutil.ReadFile("./testdata/jobs.json")
				if err != nil {
					t.Fatal(err)
				}

				var b jobPayload
				if err := json.Unmarshal(body, &b); err != nil {
					t.Fatal(err)
				}
				if *jobPrefix != EMPTY {
					b.Job["ID"] = fmt.Sprintf("%v_%v", *jobPrefix, b.Job["ID"])
				}

				newBody, err := json.Marshal(b)

				req := httptest.NewRequest(
					m,
					"/v1/jobs",
					ioutil.NopCloser(bytes.NewBuffer(newBody)),
				)

				if err := jobs(req); err != nil {
					t.Fatal(err)
				}

				nb, err := parseJob(req)
				if err != nil {
					t.Fatal(err)
				}

				if *jobPrefix != EMPTY {
					aORb := regexp.MustCompile(*jobPrefix)
					matches := aORb.FindAllStringIndex(nb.Job["ID"].(string), -1)

					if len(matches) != 1 {
						t.Fatal("Should have only one prefix")
					}
				}
				if !strings.HasPrefix(nb.Job["ID"].(string), *jobPrefix) {
					t.Fatal("Should have Job Prefix")
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
				} else if err.Error() != "cannot parse body to Job" {
					t.Fatal("Incorrect error string ", err.Error())
				}
			})

			t.Run("Should fail for invalid datacenter", func(t *testing.T) {
				body, err := ioutil.ReadFile("./testdata/jobs.json")
				if err != nil {
					t.Fatal(err)
				}

				var b jobPayload
				if err := json.Unmarshal(body, &b); err != nil {
					t.Fatal(err)
				}
				b.Job["ID"] = fmt.Sprintf("%v_%v", *jobPrefix, b.Job["ID"])
				b.Job[Datacenters] = []string{"invalid"}

				newBody, err := json.Marshal(b)

				req := httptest.NewRequest(
					m, "/v1/jobs", ioutil.NopCloser(bytes.NewBuffer(newBody)),
				)

				err = jobs(req)
				assert.NotNil(t, err)
				assert.True(t, strings.Contains(err.Error(), "invalid datacenter"))
			})
		})
	}
}

func TestMain(t *testing.M) {
	*jobPrefix = "testjob"
	*datacenter = "dc1"
	t.Run()
	*jobPrefix = EMPTY
	os.Exit(t.Run())
}
