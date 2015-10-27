package search

type Query struct {
	Q               string            `json:"q,omitempty"`
	AQ              string            `json:"aq,omitempty"`
	CQ              string            `json:"cq,omitempty"`
	DQ              string            `json:"dq,omitempty"`
	NumberOfResults int               `json:"numberOfResults,omitempty"`
	FirstResult     int               `json:"firstResult,omitempty"`
	GroupBy         []*GroupByRequest `json:"groupBy,omitempty"`
}

type GroupByRequest struct {
	Field                 string `json:"field"`
	MaximumNumberOfValues int    `json:"maximumNumberOfValues,omitempty"`
	SortCriteria          string `json:"sortCriteria,omitempty"`
	InjectionDepth        int    `json:"injectionDepth,omitempty"`
}
