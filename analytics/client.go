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

type SearchEvent struct{}
type ClickEvent struct{}
type CustomEvent struct{}
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
