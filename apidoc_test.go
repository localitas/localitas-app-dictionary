package dictionary

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandleSwagger(t *testing.T) {
	w := httptest.NewRecorder()
	HandleSwagger(w, httptest.NewRequest("GET", "/swagger.json", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	var spec APIDoc
	if err := json.Unmarshal(w.Body.Bytes(), &spec); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if spec.AppName != "Dictionary" {
		t.Errorf("expected app_name Dictionary, got %q", spec.AppName)
	}
	if len(spec.Endpoints) == 0 {
		t.Error("expected at least one endpoint")
	}
	hasLookup := false
	for _, ep := range spec.Endpoints {
		if strings.Contains(ep.Path, "/api/lookup") {
			hasLookup = true
		}
	}
	if !hasLookup {
		t.Error("expected /api/lookup endpoint")
	}
}

func TestToEntry(t *testing.T) {
	resp := &APIResponse{
		Word:     "test",
		Phonetic: "/tɛst/",
		Meanings: []struct {
			PartOfSpeech string `json:"partOfSpeech"`
			Definitions  []struct {
				Definition string   `json:"definition"`
				Example    string   `json:"example"`
				Synonyms   []string `json:"synonyms"`
				Antonyms   []string `json:"antonyms"`
			} `json:"definitions"`
		}{
			{
				PartOfSpeech: "noun",
				Definitions: []struct {
					Definition string   `json:"definition"`
					Example    string   `json:"example"`
					Synonyms   []string `json:"synonyms"`
					Antonyms   []string `json:"antonyms"`
				}{
					{Definition: "A procedure", Example: "take a test"},
				},
			},
		},
	}

	entry := resp.ToEntry()
	if entry.Word != "test" {
		t.Errorf("expected word 'test', got %q", entry.Word)
	}
	if len(entry.Meanings) != 1 {
		t.Fatalf("expected 1 meaning, got %d", len(entry.Meanings))
	}
	if entry.Meanings[0].PartOfSpeech != "noun" {
		t.Errorf("expected 'noun', got %q", entry.Meanings[0].PartOfSpeech)
	}
	if entry.Meanings[0].Definitions[0].Example != "take a test" {
		t.Errorf("expected example")
	}
	if entry.Source != "dictionary" {
		t.Errorf("expected source 'dictionary', got %q", entry.Source)
	}
}

func TestCleanUrbanText(t *testing.T) {
	input := "This is a [word] with [brackets]"
	got := cleanUrbanText(input)
	want := "This is a word with brackets"
	if got != want {
		t.Errorf("cleanUrbanText = %q, want %q", got, want)
	}
}

func TestHandleLookup_MissingWord(t *testing.T) {
	h := &handler{}
	w := httptest.NewRecorder()
	h.handleLookup(w, httptest.NewRequest("GET", "/api/lookup", nil))
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", w.Code)
	}
}

func TestLookupResult_EmptyEntries(t *testing.T) {
	r := &LookupResult{Word: "test", Entries: make([]*DictionaryEntry, 0)}
	b, _ := json.Marshal(r)
	s := string(b)
	if !jsonContains(s, `"entries":[]`) {
		t.Error("expected empty array, not null")
	}
}

func jsonContains(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
