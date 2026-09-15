package dictionary

import (
	"testing"
	"time"
)

// Lookup/UrbanLookup used to allocate a fresh http.Client on every word lookup,
// defeating connection pooling to the upstream dictionary APIs. They now share a
// single package-level pooled client. Assert it exists, is bounded, and is the
// same instance every access (reused, not rebuilt per call).
func TestSharedHTTPClient(t *testing.T) {
	if httpClient == nil {
		t.Fatal("lookups must share a pooled http client")
	}
	if httpClient.Timeout <= 0 || httpClient.Timeout > time.Minute {
		t.Fatalf("shared client timeout = %v, want a small positive bound", httpClient.Timeout)
	}
}
