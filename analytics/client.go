package analytics

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client interface {
	SendSearchEvent(SearchEvent) (*SearchEventResponse, error)
	SendSearchesEvent([]SearchEvent) (*SearchEventsResponse, error)
	SendClickEvent(ClickEvent) (*ClickEventResponse, error)
	SendCustomEvent(CustomEvent) (*CustomEventResponse, error)
	GetVisit() (*VisitResponse, error)
	GetStatus() (*StatusResponse, error)
	DeleteVisit() (bool, error)
}

type Config struct {
	Token     string
	UserAgent string
}

func NewClient(c Config) (Client, error) {
	return &client{
		token:      c.Token,
		endpoint:   "https://usageanalytics.coveo.com/v14/",
		httpClient: http.DefaultClient,
		useragent:  c.UserAgent,
	}, nil
}

type client struct {
	httpClient *http.Client
	token      string
	endpoint   string
	useragent  string
}

type SearchEvent struct {
	QueryText       string `json:"queryText"`
	AdvancedQuery   string `json:"advancedQuery"`
	NumberOfResults int    `json:"numberOfResults"`
	Contextual      bool   `json:"contextual"`
	ResponseTime    int    `json:"responseTime"`
	Results         []struct {
		documentUri     string `json:"documentUri"`
		documentUriHash string `json:"documentUriHash"`
	} `json:"results"`
	UserGroups          []string               `json:"userGroups"`
	DocumentURI         string                 `json:"documentUri"`
	DocumentURIHash     string                 `json:"DocumentUriHash"`
	SearchQueryUID      string                 `json:"searchQueryUid"`
	CollectionName      string                 `json:"collectionName"`
	SourceName          string                 `json:"sourceName"`
	DocumentPosition    string                 `json:"documentPosition"`
	ActionCause         string                 `json:"actionCause"`
	DocumentTitle       string                 `json:"documentTitle"`
	DocumentURL         string                 `json:"documentUrl"`
	QueryPipeline       string                 `json:"queryPipeline"`
	RankingModifier     string                 `json:"rankingModifier"`
	EventValue          string                 `json:"eventValue"`
	EventType           string                 `json:"eventType"`
	LastSearchQueryUid  string                 `json:"lastSearchQueryUid"`
	UserDisplayName     string                 `json:"userDisplayName"`
	UserAgent           string                 `json:"userAgent"`
	Anonymous           bool                   `json:"anonymous"`
	CustomData          map[string]interface{} `json:"customData,omitempty"`
	Device              string                 `json:"device"`
	Mobile              bool                   `json:"mobile"`
	SplitTestRunName    string                 `json:"splitTestRunName"`
	SplitTestRunVersion string                 `json:"splitTestRunVersion"`
	OriginLevel1        string                 `json:"originLevel1"`
	OriginLevel2        string                 `json:"originLevel2"`
	OriginLevel3        string                 `json:"originLevel3"`
	Username            string                 `json:"username"`
	Language            string                 `json:"language"`
}
type ClickEvent struct {
	DocumentURI         string                 `json:"documentUri"`
	DocumentURIHash     string                 `json:"DocumentUriHash"`
	SearchQueryUID      string                 `json:"searchQueryUid"`
	CollectionName      string                 `json:"collectionName"`
	SourceName          string                 `json:"sourceName"`
	DocumentPosition    string                 `json:"documentPosition"`
	ActionCause         string                 `json:"actionCause"`
	DocumentTitle       string                 `json:"documentTitle"`
	DocumentURL         string                 `json:"documentUrl"`
	QueryPipeline       string                 `json:"queryPipeline"`
	RankingModifier     string                 `json:"rankingModifier"`
	EventValue          string                 `json:"eventValue"`
	EventType           string                 `json:"eventType"`
	LastSearchQueryUid  string                 `json:"lastSearchQueryUid"`
	UserDisplayName     string                 `json:"userDisplayName"`
	UserAgent           string                 `json:"userAgent"`
	Anonymous           bool                   `json:"anonymous"`
	CustomData          map[string]interface{} `json:"customData,omitempty"`
	Device              string                 `json:"device"`
	Mobile              bool                   `json:"mobile"`
	SplitTestRunName    string                 `json:"splitTestRunName"`
	SplitTestRunVersion string                 `json:"splitTestRunVersion"`
	OriginLevel1        string                 `json:"originLevel1"`
	OriginLevel2        string                 `json:"originLevel2"`
	OriginLevel3        string                 `json:"originLevel3"`
	Username            string                 `json:"username"`
	Language            string                 `json:"language"`
}

type CustomEvent struct {
	EventValue          string                 `json:"eventValue"`
	EventType           string                 `json:"eventType"`
	LastSearchQueryUid  string                 `json:"lastSearchQueryUid"`
	UserDisplayName     string                 `json:"userDisplayName"`
	UserAgent           string                 `json:"userAgent"`
	Anonymous           bool                   `json:"anonymous"`
	CustomData          map[string]interface{} `json:"customData,omitempty"`
	Device              string                 `json:"device"`
	Mobile              bool                   `json:"mobile"`
	SplitTestRunName    string                 `json:"splitTestRunName"`
	SplitTestRunVersion string                 `json:"splitTestRunVersion"`
	OriginLevel1        string                 `json:"originLevel1"`
	OriginLevel2        string                 `json:"originLevel2"`
	OriginLevel3        string                 `json:"originLevel3"`
	Username            string                 `json:"username"`
	Language            string                 `json:"language"`
}

type SearchEventResponse struct{}
type SearchEventsResponse struct{}
type ClickEventResponse struct{}
type CustomEventResponse struct{}
type VisitResponse struct{}

func (c *client) SendSearchEvent(event SearchEvent) (*SearchEventResponse, error) {
	return nil, nil
}
func (c *client) SendSearchesEvent(event []SearchEvent) (*SearchEventsResponse, error) {
	return nil, nil
}
func (c *client) SendClickEvent(event ClickEvent) (*ClickEventResponse, error) {
	return nil, nil
}
func (c *client) SendCustomEvent(event CustomEvent) (*CustomEventResponse, error) {
	return nil, nil
}
func (c *client) GetVisit() (*VisitResponse, error) {
	return nil, nil
}
func (c *client) DeleteVisit() (bool, error) {
	return false, nil
}
