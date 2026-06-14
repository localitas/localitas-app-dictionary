package dictionary

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type handler struct {
	sources map[string]bool
}

func (h *handler) handleLookup(w http.ResponseWriter, r *http.Request) {
	word := r.URL.Query().Get("q")
	if word == "" {
		word = r.PathValue("word")
	}
	if word == "" {
		writeErr(w, http.StatusBadRequest, "query parameter 'q' or path parameter 'word' is required")
		return
	}

	result := LookupAll(r.Context(), word, h.sources)
	if len(result.Entries) == 0 {
		writeErr(w, http.StatusNotFound, "no definition found for '%s'", word)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, status int, format string, args ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": fmt.Sprintf(format, args...)})
}
