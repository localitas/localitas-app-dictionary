package dictionary

import (
	"net/http"

	"github.com/localitas/localitas-go/httputil"
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
		writeErr(w, r, http.StatusBadRequest, "query parameter 'q' or path parameter 'word' is required")
		return
	}

	result := LookupAll(r.Context(), word, h.sources)
	if len(result.Entries) == 0 {
		writeErr(w, r, http.StatusNotFound, "no definition found for '%s'", word)
		return
	}

	writeJSON(w, r, http.StatusOK, result)
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, v interface{}) {
	httputil.WriteResponse(w, r, status, v)
}

func writeErr(w http.ResponseWriter, r *http.Request, status int, format string, args ...interface{}) {
	httputil.WriteError(w, r, status, format, args...)
}
