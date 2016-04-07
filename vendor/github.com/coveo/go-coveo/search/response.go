package search

type Response struct {
	TotalCount         int             `json:"totalCount"`
	TotalCountFiltered int             `json:"totalCountFiltered"`
	Duration           int             `json:"duration"`
	IndexDuration      int             `json:"indexDuration"`
	RequestDuration    int             `json:"requestDuration"`
	SearchUID          string          `json:"searchUid"`
	Pipeline           string          `json:"pipeline"`
	GroupByResults     []GroupByResult `json:"groupByResults,omitempty"`
	Results            []Result        `json:"results,omitempty"`
}

type Result struct {
	Title          string                 `json:"title"`
	URI            string                 `json:"uri"`
	Excerpt        string                 `json:"excerpt"`
	FirstSentences string                 `json:"firstSentences"`
	Score          int                    `json:"score"`
	PercentScore   float32                `json:"percentScore"`
	ClickUri       string                 `json:"clickUri"`
	Raw            map[string]interface{} `json:"raw"`
}

type GroupByResult struct {
	Field  string `json:"field"`
	Values []struct {
		Value           string `json:"value"`
		NumberOfResults int    `json:"numberOfResults"`
		Score           int    `json:"score"`
		ValueType       string `json:"valueType"`
	} `json:"values"`
}
