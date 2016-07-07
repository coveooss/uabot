package search

type FacetValue struct {
	Value           string `json:"value"`
	LookupValue     string `json:"lookupValue"`
	NumberOfResults int    `json:"numberOfResults"`
}

type FacetValues struct {
	Values []FacetValue `json:"values"`
}
