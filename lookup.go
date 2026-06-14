package dictionary

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	apiBaseURL   = "https://api.dictionaryapi.dev/api/v2/entries/en/"
	urbanBaseURL = "https://api.urbandictionary.com/v0/define?term="
)

func LookupAll(ctx context.Context, word string, sources map[string]bool) *LookupResult {
	result := &LookupResult{Word: word, Entries: make([]*DictionaryEntry, 0)}

	type entryResult struct {
		entry *DictionaryEntry
		err   error
	}

	ch := make(chan entryResult, 2)
	count := 0

	if sources["dictionary"] {
		count++
		go func() {
			e, err := Lookup(ctx, word)
			ch <- entryResult{e, err}
		}()
	}
	if sources["urban"] {
		count++
		go func() {
			e, err := UrbanLookup(ctx, word)
			ch <- entryResult{e, err}
		}()
	}

	for i := 0; i < count; i++ {
		if r := <-ch; r.err == nil {
			result.Entries = append(result.Entries, r.entry)
		}
	}

	return result
}

func Lookup(ctx context.Context, word string) (*DictionaryEntry, error) {
	if word == "" {
		return nil, fmt.Errorf("word is required")
	}

	reqURL := apiBaseURL + url.PathEscape(word)
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "Localitas Dictionary/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("no definition found for '%s'", word)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("dictionary API returned status %d", resp.StatusCode)
	}

	var apiResponses []APIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResponses); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if len(apiResponses) == 0 {
		return nil, fmt.Errorf("no definition found for '%s'", word)
	}

	return apiResponses[0].ToEntry(), nil
}

func UrbanLookup(ctx context.Context, word string) (*DictionaryEntry, error) {
	if word == "" {
		return nil, fmt.Errorf("word is required")
	}

	reqURL := urbanBaseURL + url.QueryEscape(word)
	req, err := http.NewRequestWithContext(ctx, "GET", reqURL, nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "Localitas Dictionary/1.0")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("urban dictionary API returned status %d", resp.StatusCode)
	}

	var urbanResp UrbanAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&urbanResp); err != nil {
		return nil, fmt.Errorf("parse response: %w", err)
	}

	if len(urbanResp.List) == 0 {
		return nil, fmt.Errorf("no urban definition found for '%s'", word)
	}

	entry := &DictionaryEntry{
		Word:   word,
		Source: "urban",
	}
	meaning := Meaning{PartOfSpeech: "slang"}
	for i, d := range urbanResp.List {
		if i >= 5 {
			break
		}
		meaning.Definitions = append(meaning.Definitions, Definition{
			Text:    cleanUrbanText(d.Definition),
			Example: cleanUrbanText(d.Example),
		})
	}
	entry.Meanings = append(entry.Meanings, meaning)
	return entry, nil
}

func cleanUrbanText(s string) string {
	s = strings.ReplaceAll(s, "[", "")
	s = strings.ReplaceAll(s, "]", "")
	return s
}
