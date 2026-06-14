package dictionary

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strings"
)

type APIEndpoint struct {
	Method      string     `json:"method"`
	Path        string     `json:"path"`
	Summary     string     `json:"summary"`
	QueryParams []APIParam `json:"query_params,omitempty"`
	Response    *APIBody   `json:"response,omitempty"`
}

type APIParam struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Required    bool   `json:"required"`
	Description string `json:"description"`
}

type APIBody struct {
	ContentType string `json:"content_type"`
	Example     string `json:"example"`
}

type APIDoc struct {
	AppName     string        `json:"app_name"`
	Version     string        `json:"version"`
	Description string        `json:"description"`
	Keywords    []string      `json:"keywords,omitempty"`
	Endpoints   []APIEndpoint `json:"endpoints"`
}

var DictionaryAPIDoc = APIDoc{
	AppName:     "Dictionary",
	Version:     "0.1.0",
	Description: "English dictionary and slang lookup using Free Dictionary API and Urban Dictionary",
	Keywords:    []string{"dictionary", "define", "definition", "meaning", "word", "synonym", "antonym", "thesaurus", "slang", "vocabulary", "lookup"},
	Endpoints: []APIEndpoint{
		{
			Method:  "GET",
			Path:    "/api/lookup",
			Summary: "Look up a word definition",
			QueryParams: []APIParam{
				{Name: "q", Type: "string", Required: true, Description: "The word to look up"},
			},
			Response: &APIBody{
				ContentType: "application/json",
				Example:     `{"word":"hello","phonetic":"/həˈloʊ/","meanings":[{"part_of_speech":"noun","definitions":[{"definition":"An utterance of 'hello'; a greeting.","example":"she was getting hellos from everyone"}]},{"part_of_speech":"verb","definitions":[{"definition":"Say or shout 'hello'."}]}]}`,
			},
		},
		{
			Method:  "GET",
			Path:    "/api/lookup/{word}",
			Summary: "Look up a word by path parameter",
			Response: &APIBody{
				ContentType: "application/json",
				Example:     `{"word":"serendipity","phonetic":"/ˌsɛɹ.ən.ˈdɪp.ə.ti/","meanings":[{"part_of_speech":"noun","definitions":[{"definition":"The occurrence of events by chance in a happy way."}]}]}`,
			},
		},
	},
}

func HandleSwagger(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(DictionaryAPIDoc)
}

func RenderDocsHTML(doc APIDoc) template.HTML {
	var sb strings.Builder
	sb.WriteString(`<h3 style="font-size: 0.875rem; font-weight: 600; text-transform: uppercase; letter-spacing: 0.05em; color: var(--color-text-secondary); margin-bottom: 1rem;">API Endpoints</h3><div class="accordion-list">`)
	for _, ep := range doc.Endpoints {
		title := fmt.Sprintf("%s %s — %s", ep.Method, ep.Path, ep.Summary)
		sb.WriteString(fmt.Sprintf(`<details class="glass-panel" style="border-radius: 0.5rem; margin-bottom: 0.5rem;"><summary style="padding: 0.75rem 1rem; cursor: pointer; font-weight: 500; color: var(--color-text-primary);">%s</summary><div style="padding: 0 1rem 0.75rem; font-size: 0.875rem; color: var(--color-text-secondary);">`, template.HTMLEscapeString(title)))
		if ep.Response != nil {
			sb.WriteString(fmt.Sprintf(`<pre style="background: var(--color-bg-base); padding: 0.75rem; border-radius: 0.375rem; overflow-x: auto; font-size: 0.8125rem;">%s</pre>`, template.HTMLEscapeString(prettyJSON(ep.Response.Example))))
		}
		sb.WriteString(`</div></details>`)
	}
	sb.WriteString(`</div>`)
	return template.HTML(sb.String())
}

func prettyJSON(s string) string {
	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return s
	}
	b, _ := json.MarshalIndent(v, "", "  ")
	return string(b)
}
