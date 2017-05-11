package analytics

type ActionEvent struct {
	Language            string                 `json:"language"`
	Device              string                 `json:"device"`
	OriginLevel1        string                 `json:"originLevel1"`
	OriginLevel2        string                 `json:"originLevel2"`
	OriginLevel3        string                 `json:"originLevel3,omitempty"`
	UserAgent           string                 `json:"userAgent,omitempty"`
	CustomData          map[string]interface{} `json:"customData,omitempty"`
	Anonymous           bool                   `json:"anonymous,omitempty"`
	Username            string                 `json:"username,omitempty"`
	UserDisplayName     string                 `json:"userDisplayName,omitempty"`
	Mobile              bool                   `json:"mobile,omitempty"`
	SplitTestRunName    string                 `json:"splitTestRunName,omitempty"`
	SplitTestRunVersion string                 `json:"splitTestRunVersion,omitempty"`
}

type SearchEvent struct {
	*ActionEvent
	SearchQueryUID  string       `json:"searchQueryUid"`
	QueryText       string       `json:"queryText"`
	ActionCause     string       `json:"actionCause"`
	ActionType      string       `json:"actionType"`
	AdvancedQuery   string       `json:"advancedQuery,omitempty"`
	NumberOfResults int          `json:"numberOfResults,omitempty"`
	Contextual      bool         `json:"contextual"`
	ResponseTime    int          `json:"responseTime,omitempty"`
	QueryPipeline   string       `json:"queryPipeline,omitempty"`
	UserGroups      []string     `json:"userGroups,omitempty"`
	Results         []ResultHash `json:"results,omitempty"`
}

type ResultHash struct {
	DocumentURI     string `json:"documentUri"`
	DocumentURIHash string `json:"documentUriHash"`
}

type ClickEvent struct {
	*ActionEvent
	DocumentURI      string `json:"documentUri"`
	DocumentURIHash  string `json:"documentUriHash"`
	SearchQueryUID   string `json:"searchQueryUid"`
	CollectionName   string `json:"collectionName"`
	SourceName       string `json:"sourceName"`
	DocumentPosition int    `json:"documentPosition"`
	ActionCause      string `json:"actionCause"`
	ViewMethod       string `json:"viewMethod, omitempty"`
	DocumentTitle    string `json:"documentTitle,omitempty"`
	DocumentURL      string `json:"documentUrl,omitempty"`
	QueryPipeline    string `json:"queryPipeline,omitempty"`
	RankingModifier  string `json:"rankingModifier,omitempty"`
}

type CustomEvent struct {
	*ActionEvent
	ActionCause        string `json:"actionCause"`
	ActionType         string `json:"actionType"`
	EventType          string `json:"eventType"`
	EventValue         string `json:"eventValue"`
	LastSearchQueryUID string `json:"lastSearchQueryUid,omitempty"`
}

type ViewEvent struct {
	*ActionEvent
	PageURI        string `json:"location"`
	PageTitle      string `json:"title"`
	ContentIDKey   string `json:"contentIdKey"`
	ContentIDValue string `json:"contentIdValue"`
	ContentType    string `json:"contentType"`
	PageReferrer   string `json:"referrer,omitempty"`
}
