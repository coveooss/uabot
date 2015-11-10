package analytics

type ActionEvent struct {
	Language            string                 `json:"language"`
	Device              string                 `json:"device"`
	OriginLevel1        string                 `json:"originLevel1"`
	OriginLevel2        string                 `json:"originLevel2"`
	UserAgent           string                 `json:"userAgent",omitempty`
	CustomData          map[string]interface{} `json:"customData,omitempty"`
	Anonymous           bool                   `json:"anonymous",omitempty`
	Username            string                 `json:"username",omitempty`
	UserDisplayName     string                 `json:"userDisplayName",omitempty`
	Mobile              bool                   `json:"mobile",omitempty`
	SplitTestRunName    string                 `json:"splitTestRunName",omitempty`
	SplitTestRunVersion string                 `json:"splitTestRunVersion",omitempty`
	OriginLevel3        string                 `json:"originLevel3",omitempty`
}

type SearchEvent struct {
	*ActionEvent
	SearchQueryUid  string   `json:"searchQueryUid"`
	QueryText       string   `json:"queryText"`
	ActionCause     string   `json:"actionCause"`
	AdvancedQuery   string   `json:"advancedQuery",omitempty`
	NumberOfResults int      `json:"numberOfResults",omitempty`
	Contextual      bool     `json:"contextual"`
	ResponseTime    int      `json:"responseTime",omitempty`
	QueryPipeline   string   `json:"queryPipeline",omitempty`
	UserGroups      []string `json:"userGroups,omitempty"`
	Results         []struct {
		DocumentUri     string `json:"documentUri"`
		DocumentUriHash string `json:"documentUriHash"`
	} `json:"results"`
}

type ClickEvent struct {
	*ActionEvent
	DocumentUri      string `json:"documentUri"`
	DocumentUriHash  string `json:"documentUriHash"`
	SearchQueryUid   string `json:"searchQueryUid"`
	CollectionName   string `json:"collectionName"`
	SourceName       string `json:"sourceName"`
	DocumentPosition int    `json:"documentPosition"`
	ActionCause      string `json:"actionCause"`
	DocumentTitle    string `json:"documentTitle",omitempty`
	DocumentUrl      string `json:"documentUrl",omitempty`
	QueryPipeline    string `json:"queryPipeline",omitempty`
	RankingModifier  string `json:"rankingModifier",omitempty`
}

type CustomEvent struct {
	*ActionEvent
	EventType          string `json:"eventType"`
	EventValue         string `json:"eventValue"`
	LastSearchQueryUid string `json:"lastSearchQueryUid",omitempty`
}
