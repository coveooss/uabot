package analytics

import (
	"bytes"
	"encoding/json"
	"net/http"
	"fmt"
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
	SendCustomEvent(CustomEvent) error
	GetVisit() (*VisitResponse, error)
	GetStatus() (*StatusResponse, error)
	DeleteVisit() (bool, error)
	GetCookies() ([]*http.Cookie, error)
}

// Config is the configuration of the usageanalytics client
type Config struct {
	// Token is the token used to log into the service remotly
	Token string
	// User agent is the http user agent sent to the service
	UserAgent string
	// IP is used if you want to specify an origin IP to the client
	IP string
}

// NewClient return a capable Coveo Usage Analytics service client. It currently
// uses V14 of the API.
func NewClient(c Config) (Client, error) {
	return &client{
		token:      c.Token,
		endpoint:   "https://usageanalytics.coveo.com/rest/v14/analytics/",
		httpClient: http.DefaultClient,
		useragent:  c.UserAgent,
		ip:         c.IP,
	}, nil
}

type client struct {
	httpClient *http.Client
	token      string
	endpoint   string
	useragent  string
	ip         string
	cookies    []*http.Cookie
}

func NewSearchEvent() (*SearchEvent, error) {
	return &SearchEvent{
		ActionEvent: &ActionEvent{
			Language:     "en",
			Device:       "Bot",
			OriginLevel1: "default",
			OriginLevel2: "All",
		},
		SearchQueryUid: "",
		QueryText:      "",
		ActionCause:    "interfaceLoad",
		Contextual:     false,
	}, nil
}

func NewClickEvent() (*ClickEvent, error) {
	return &ClickEvent{
		ActionEvent: &ActionEvent{
			Language:     "en",
			Device:       "Bot",
			OriginLevel1: "default",
			OriginLevel2: "All",
		},
		DocumentUri: "",
		DocumentUriHash: "",
		SearchQueryUid: "",
		CollectionName: "",
		SourceName: "",
		DocumentPosition: 0,
		ActionCause: "documentOpen",
	}, nil
}

type StatusResponse struct{}
type SearchEventsResponse struct{}
type ClickEventResponse struct{}
type CustomEventResponse struct{}
type VisitResponse struct{}

func (c *client) GetCookies() ([]*http.Cookie, error) {
	return c.cookies, nil
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

func (c *client) SendCustomEvent(event CustomEvent) error {
	return nil
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

	fmt.Println(resp.StatusCode)

	if c.cookies == nil {
		cookies := resp.Cookies()
		c.cookies = cookies
	}

	return nil
}
