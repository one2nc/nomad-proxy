package main

import (
	"fmt"
	"net/http/httptest"
	"testing"
)

func TestJob(t *testing.T) {
	t.Run("request path", func(t *testing.T) {
		t.Run("Valid prefix is allowed", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/job/"+*jobPrefix, nil)
			if err := job(req); err != nil {
				t.Fatal("Should not have failed")
			}
		})

		t.Run("InValid prefix is not allowed", func(t *testing.T) {
			req := httptest.NewRequest("GET", "/v1/job/hello", nil)
			// This should fail only if prefix is required and doesn't match hello.
			if err := job(req); !*skipPrefix && err == nil {
				t.Fatal("Should have failed")
			} else if err != nil && err.Error() != fmt.Sprintf("jobs should begin with a prefix %v", *jobPrefix) {
				t.Error("Should fail with right message", err.Error())
			}
		})
	})
}
