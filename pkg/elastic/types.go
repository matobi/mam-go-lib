package elastic

import "encoding/json"

type SearchResult struct {
	TookMs   int64 `json:"took"`
	TimedOut bool  `json:"_timed_out"`
	Hits     Hits  `json:"hits"`
}

type Hits struct {
	Total    int     `json:"total"`
	MaxScore float64 `json:"max_score"`
	Hits     []Hit   `json:"hits"`
}

type Hit struct {
	Index  string          `json:"_index"`
	Type   string          `json:"_type"`
	ID     string          `json:"_id"`
	Score  float64         `json:"_score"`
	Source json.RawMessage `json:"_source"`
}
