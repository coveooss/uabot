package analytics

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

const (
	// EndpointProduction is the Analytics production endpoint
	EndpointProduction = "https://usageanalytics.coveo.com/rest/v15/analytics/"
	// EndpointStaging is the Analytics staging endpoint
	EndpointStaging = "https://usageanalyticsstaging.coveo.com/rest/v15/analytics/"
	// EndpointDevelopment is the Analytics development endpoint
	EndpointDevelopment = "https://usageanalyticsdev.coveo.com/rest/v15/analytics/"
)

// Client is the basic element of the usage analytics service, it wraps a http
// client. with the appropriate calls to the usage analytics service.
type Client interface {

	// SendSearchEvent sends a searchEvent to the analytics service, as the
	// response is not important it only returns an error
	SendSearchEvent(*SearchEvent) error

	// SendSearchesEvent sends multiple searchEvent to the analytics service,
	// using the batch call, as the response is not important it only
	// returns an error
	SendSearchesEvent([]SearchEvent) error

	// SendClickEvent sends a click to the analytics service, as the
	// response is not important it only returns an error
	SendClickEvent(*ClickEvent) error

	// SendCustomEvent sends a custom event to the analytics service, as the
	// response is not important it only returns an error
	SendCustomEvent(CustomEvent) error

	// SendViewEvent sends a view event to the analytics service, as the
	// response is not important it only returns an error
	SendViewEvent(*ViewEvent) error

	GetVisit() (*VisitResponse, error)
	GetStatus() (*StatusResponse, error)
	DeleteVisit() (bool, error)
	GetCookies() []*http.Cookie
}

// Config is the configuration of the usageanalytics client
type Config struct {
	// Token is the token used to log into the service remotly
	Token string
	// User agent is the http user agent sent to the service
	UserAgent string
	// IP is used if you want to specify an origin IP to the client
	IP string
	// Endpoint is used if you want to use custom endpoints (dev,staging,testing)
	Endpoint string
}

// NewClient return a capable Coveo Usage Analytics service client. It currently
// uses V15 of the API.
func NewClient(c Config) Client {
	if len(c.Endpoint) == 0 {
		c.Endpoint = EndpointProduction
	}
	return &client{
		token:      c.Token,
		endpoint:   c.Endpoint,
		httpClient: http.DefaultClient,
		useragent:  c.UserAgent,
		ip:         c.IP}
}

type client struct {
	httpClient *http.Client
	token      string
	endpoint   string
	useragent  string
	ip         string
	cookies    []*http.Cookie
}

// NewSearchEvent creates a new SearchEvent which can then be altered
func NewSearchEvent() *SearchEvent {
	return &SearchEvent{
		ActionEvent: &ActionEvent{
			Language:     "en",
			Device:       "Bot",
			OriginLevel1: "default",
			OriginLevel2: "All",
		},
		SearchQueryUID: "",
		QueryText:      "",
		ActionCause:    "interfaceLoad",
		Contextual:     false,
	}
}

// NewClickEvent creates a new ClickEvent which can then be altered
func NewClickEvent() *ClickEvent {
	return &ClickEvent{
		ActionEvent: &ActionEvent{
			Language:     "en",
			Device:       "Bot",
			OriginLevel1: "default",
			OriginLevel2: "All",
		},
		DocumentURI:      "",
		DocumentURIHash:  "",
		SearchQueryUID:   "",
		CollectionName:   "",
		SourceName:       "",
		DocumentPosition: 0,
		ActionCause:      "documentOpen",
	}
}

// NewCustomEvent creates a new SearchEvent which can then be altered
func NewCustomEvent() *CustomEvent {
	return &CustomEvent{
		ActionEvent: &ActionEvent{
			Language:     "en",
			Device:       "Bot",
			OriginLevel1: "default",
			OriginLevel2: "All",
		},
		EventType:          "",
		EventValue:         "",
		LastSearchQueryUID: "",
	}
}

// NewViewEvent creates a ViewEvent which can then be changed
func NewViewEvent() *ViewEvent {
	return &ViewEvent{
		ActionEvent: &ActionEvent{
			Language:     "en",
			Device:       "Bot",
			OriginLevel1: "default",
			OriginLevel2: "All",
		},
		PageURI:      "",
		PageReferrer: "",
		PageTitle:    "",
	}
}

// StatusResponse is the response to a Status service call
type StatusResponse struct{}

// SearchEventsResponse is the response to a SearchEvent call
type SearchEventsResponse struct{}

// ClickEventResponse is the response to a ClickEvent call
type ClickEventResponse struct{}

// CustomEventResponse is the response to a CustomEvent call
type CustomEventResponse struct{}

// VisitResponse is the response to a Visit call
type VisitResponse struct{}

func (c *client) GetCookies() []*http.Cookie {
	return c.cookies
}

func (c *client) SendSearchEvent(event *SearchEvent) error {
	err := c.sendEventRequest("search/", event)
	return err
}

func (c *client) SendSearchesEvent(event []SearchEvent) error {
	return nil
}

func (c *client) SendClickEvent(event *ClickEvent) error {
	err := c.sendEventRequest("click/", event)
	return err
}

// SendCustomEvent Send a request to usage analytics to create a new custom event.
func (c *client) SendCustomEvent(event CustomEvent) error {
	err := c.sendEventRequest("custom/", event)
	return err
}

// SendViewEvent Send a request to usage analytics to create a new view event.
func (c *client) SendViewEvent(event *ViewEvent) error {
	err := c.sendEventRequest("view/", event)
	return err
}

func (c *client) GetVisit() (*VisitResponse, error) {
	return nil, nil
}

// DeleteVisit forgets the cookie to usageanalytics, the call to the server
// currently does the same thing. This will probably change in the future
func (c *client) DeleteVisit() (bool, error) {
	c.cookies = nil
	return true, nil
}

func (c *client) GetStatus() (*StatusResponse, error) {
	return nil, nil
}

func (c *client) sendEventRequest(path string, event interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.endpoint+path, &buf)
	if err != nil {
		return err
	}

	if c.cookies != nil {
		for _, cookie := range c.cookies {
			req.AddCookie(cookie)
		}
	}
	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accepts", "application/json")
	req.Header.Set("User-Agent", c.useragent)
	if c.ip != "" {
		req.Header.Add("X-Forwarded-For", c.ip)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.cookies == nil {
		cookies := resp.Cookies()
		c.cookies = cookies
	}

	return nil
}

func (c *client) sendRawEventRequest(method string, path string, body string) error {
	req, err := http.NewRequest(method, c.endpoint+path, strings.NewReader(body))
	if err != nil {
		return err
	}

	if c.cookies != nil {
		for _, cookie := range c.cookies {
			req.AddCookie(cookie)
		}
	}
	req.Header.Add("Authorization", "Beared "+c.token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accepts", "application/json")
	req.Header.Set("User-Agent", c.useragent)
	if c.ip != "" {
		req.Header.Add("X-Forwarded-For", c.ip)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if c.cookies == nil {
		cookies := resp.Cookies()
		c.cookies = cookies
	}

	return err
}
