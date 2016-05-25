package search

type Query struct {
	Q                     string            `json:"q,omitempty"`
	AQ                    string            `json:"aq,omitempty"`
	CQ                    string            `json:"cq,omitempty"`
	DQ                    string            `json:"dq,omitempty"`
	NumberOfResults       int               `json:"numberOfResults,omitempty"`
	FirstResult           int               `json:"firstResult,omitempty"`
	GroupByRequests       []*GroupByRequest `json:"groupBy,omitempty"`
	Tab                   string            `json:"tab,omitempty"`
	PartialMatch          bool              `json:"partialMatch,omitempty"`
	PartialMatchKeywords  int               `json:"partialMatchKeywords,omitempty"`
	PartialMatchThreshold string            `json:"partialMatchTreshold"`
	Pipeline              string            `json:"pipeline,omitempty"`
}

type GroupByRequest struct {
	Field                 string `json:"field"`
	MaximumNumberOfValues int    `json:"maximumNumberOfValues,omitempty"`
	SortCriteria          string `json:"sortCriteria,omitempty"`
	InjectionDepth        int    `json:"injectionDepth,omitempty"`
}
