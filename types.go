package dictionary

type DictionaryEntry struct {
	Word     string    `json:"word"`
	Phonetic string    `json:"phonetic,omitempty"`
	Audio    string    `json:"audio,omitempty"`
	Source   string    `json:"source"`
	Meanings []Meaning `json:"meanings"`
}

type LookupResult struct {
	Word    string             `json:"word"`
	Entries []*DictionaryEntry `json:"entries"`
}

type UrbanAPIResponse struct {
	List []UrbanDefinition `json:"list"`
}

type UrbanDefinition struct {
	Definition string `json:"definition"`
	Word       string `json:"word"`
	Example    string `json:"example"`
	ThumbsUp   int    `json:"thumbs_up"`
	ThumbsDown int    `json:"thumbs_down"`
}

type Meaning struct {
	PartOfSpeech string       `json:"part_of_speech"`
	Definitions  []Definition `json:"definitions"`
}

type Definition struct {
	Text     string   `json:"definition"`
	Example  string   `json:"example,omitempty"`
	Synonyms []string `json:"synonyms,omitempty"`
	Antonyms []string `json:"antonyms,omitempty"`
}

type APIResponse struct {
	Word      string `json:"word"`
	Phonetic  string `json:"phonetic"`
	Phonetics []struct {
		Text  string `json:"text"`
		Audio string `json:"audio"`
	} `json:"phonetics"`
	Meanings []struct {
		PartOfSpeech string `json:"partOfSpeech"`
		Definitions  []struct {
			Definition string   `json:"definition"`
			Example    string   `json:"example"`
			Synonyms   []string `json:"synonyms"`
			Antonyms   []string `json:"antonyms"`
		} `json:"definitions"`
	} `json:"meanings"`
}

func (r *APIResponse) ToEntry() *DictionaryEntry {
	entry := &DictionaryEntry{
		Word:     r.Word,
		Phonetic: r.Phonetic,
		Source:   "dictionary",
	}
	for _, p := range r.Phonetics {
		if p.Audio != "" {
			entry.Audio = p.Audio
			break
		}
	}
	for _, m := range r.Meanings {
		meaning := Meaning{PartOfSpeech: m.PartOfSpeech}
		for _, d := range m.Definitions {
			def := Definition{
				Text:     d.Definition,
				Example:  d.Example,
				Synonyms: d.Synonyms,
				Antonyms: d.Antonyms,
			}
			meaning.Definitions = append(meaning.Definitions, def)
		}
		entry.Meanings = append(entry.Meanings, meaning)
	}
	return entry
}
