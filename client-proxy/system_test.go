package main

import (
	"net/http/httptest"
	"testing"
)

func TestSystem(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		req := httptest.NewRequest("PUT", "/", nil)

		if err := system(req); err == nil {
			t.Fatal("Should have failed")
		} else if err.Error() != "System calls are not allowed" {
			t.Error("Should not allow system calls")
		}
	})
}
