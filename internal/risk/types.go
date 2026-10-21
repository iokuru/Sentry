package risk

type Entity struct {
	Type  string  `json:"type"`
	Text  string  `json:"text"`
	Score float64 `json:"score"`
}

type Input struct {
	Text       string   `json:"text"`
	Context    string   `json:"context"`
	Sentiment  string   `json:"sentiment"`
	KeyPhrases []string `json:"key_phrases"`
	Entities   []Entity `json:"entities"`
}
