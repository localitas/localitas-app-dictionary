package dictionary

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

// serveDictionary stands in for dictionaryapi.dev, which a test cannot make
// fail on demand: its first failures requests fail with status, then it
// defines "axolotl" and knows no other word.
func serveDictionary(t *testing.T, failures int, status int) *atomic.Int32 {
	t.Helper()
	var calls atomic.Int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if int(calls.Add(1)) <= failures {
			http.Error(w, "down", status)
			return
		}
		if r.URL.Path != "/axolotl" {
			http.NotFound(w, r)
			return
		}
		w.Write([]byte(`[{"word":"axolotl","meanings":[{"partOfSpeech":"noun","definitions":[{"definition":"A salamander."}]}]}]`))
	}))
	t.Cleanup(srv.Close)

	oldURL, oldUpstream := apiBaseURL, upstream
	apiBaseURL, upstream.Delay = srv.URL+"/", time.Millisecond
	t.Cleanup(func() { apiBaseURL, upstream = oldURL, oldUpstream })
	return &calls
}

func lookupWord(word string) *httptest.ResponseRecorder {
	h := &handler{sources: map[string]bool{"dictionary": true}}
	w := httptest.NewRecorder()
	h.handleLookup(w, httptest.NewRequest("GET", "/api/lookup?q="+word, nil))
	return w
}

// The dictionary failing for a moment is asked again, and the word defined.
func TestLookup_AMomentaryUpstreamFailureIsTriedAgain(t *testing.T) {
	calls := serveDictionary(t, 1, http.StatusBadGateway)

	if w := lookupWord("axolotl"); w.Code != http.StatusOK {
		t.Fatalf("lookup = %d %s, want 200", w.Code, w.Body)
	}
	if calls.Load() != 2 {
		t.Errorf("dictionary asked %d times, want 2", calls.Load())
	}
}

// A dictionary that keeps failing is a bad gateway, not a word with no
// definition.
func TestLookup_AnUpstreamThatKeepsFailingIs502(t *testing.T) {
	serveDictionary(t, 100, http.StatusServiceUnavailable)

	if w := lookupWord("axolotl"); w.Code != http.StatusBadGateway {
		t.Errorf("lookup = %d %s, want 502", w.Code, w.Body)
	}
}

// A word the dictionary does not know is not found, asked about once.
func TestLookup_AnUnknownWordIs404(t *testing.T) {
	calls := serveDictionary(t, 0, 0)

	if w := lookupWord("qwzx"); w.Code != http.StatusNotFound {
		t.Errorf("lookup = %d %s, want 404", w.Code, w.Body)
	}
	if calls.Load() != 1 {
		t.Errorf("dictionary asked %d times, want 1", calls.Load())
	}
}
