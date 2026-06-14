package dictionary

import (
	"html/template"
	"log"
	"net/http"
	"strings"
)

type App struct {
	BasePath string
	Sources  map[string]bool
}

func New(basePath, sources string) *App {
	if basePath == "" {
		basePath = "/"
	}
	srcMap := make(map[string]bool)
	for _, s := range strings.Split(sources, ",") {
		s = strings.TrimSpace(s)
		if s != "" {
			srcMap[s] = true
		}
	}
	if len(srcMap) == 0 {
		srcMap["dictionary"] = true
		srcMap["urban"] = true
	}
	return &App{BasePath: basePath, Sources: srcMap}
}

func (a *App) handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFS(TemplatesFS, "templates/index.html")
	if err != nil {
		log.Printf("dictionary index template error: %v", err)
		http.Error(w, "template error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	tmpl.ExecuteTemplate(w, "index.html", map[string]string{"BasePath": a.BasePath})
}

func (a *App) RegisterRoutes(mux *http.ServeMux) {
	h := &handler{sources: a.Sources}
	mux.HandleFunc("GET /{$}", a.handleIndex)
	mux.HandleFunc("GET /swagger.json", HandleSwagger)
	mux.HandleFunc("GET /help.md", handleHelpMarkdown)
	mux.HandleFunc("GET /api/lookup", h.handleLookup)
	mux.HandleFunc("GET /api/lookup/{word}", h.handleLookup)
}
