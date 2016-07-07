package search

// Query Struct reprensenting a query sent to the index.
type Query struct {
	Q               string            `json:"q,omitempty"`
	AQ              string            `json:"aq,omitempty"`
	CQ              string            `json:"cq,omitempty"`
	DQ              string            `json:"dq,omitempty"`
	NumberOfResults int               `json:"numberOfResults,omitempty"`
	FirstResult     int               `json:"firstResult,omitempty"`
	GroupByRequests []*GroupByRequest `json:"groupBy,omitempty"`
	Tab             string            `json:"tab,omitempty"`
}

// GroupByRequest Struct representing a GroupByRequest send to the index. It is
// used to send data to the facets of a search page.
type GroupByRequest struct {
	Field                 string `json:"field"`
	MaximumNumberOfValues int    `json:"maximumNumberOfValues,omitempty"`
	SortCriteria          string `json:"sortCriteria,omitempty"`
	InjectionDepth        int    `json:"injectionDepth,omitempty"`
}
