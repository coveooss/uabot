package search

// Response A collection of results from the Coveo index following a query.
// It contains the results that were returned from the query and some metadata.
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

// Result A single result returned from a query to the Coveo index.
type Result struct {
	Title          string                 `json:"title"`
	URI            string                 `json:"uri"`
	Excerpt        string                 `json:"excerpt"`
	FirstSentences string                 `json:"firstSentences"`
	Score          int                    `json:"score"`
	PercentScore   float32                `json:"percentScore"`
	ClickURI       string                 `json:"clickUri"`
	Raw            map[string]interface{} `json:"raw"`
}

// GroupByResult The result of a group by request to the index.
type GroupByResult struct {
	Field  string `json:"field"`
	Values []struct {
		Value           string `json:"value"`
		NumberOfResults int    `json:"numberOfResults"`
		Score           int    `json:"score"`
		ValueType       string `json:"valueType"`
	} `json:"values"`
}
