package analytics

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type Client interface {
	SendSearchEvent(*SearchEvent) error
	SendSearchesEvent([]SearchEvent) error
	SendClickEvent(ClickEvent) error
	SendCustomEvent(CustomEvent) error
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
		endpoint:   "https://usageanalytics.coveo.com/rest/v14/analytics/",
		httpClient: http.DefaultClient,
		useragent:  c.UserAgent,
	}, nil
}

type client struct {
	httpClient *http.Client
	token      string
	endpoint   string
	useragent  string
	cookies    []*http.Cookie
}

func NewInterfaceLoad() (*SearchEvent, error) {
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

type StatusResponse struct{}
type SearchEventsResponse struct{}
type ClickEventResponse struct{}
type CustomEventResponse struct{}
type VisitResponse struct{}

func (c *client) SendSearchEvent(event *SearchEvent) error {
	err := c.sendEventRequest("search/", event)
	return err
}

func (c *client) SendSearchesEvent(event []SearchEvent) error {
	return nil
}

func (c *client) SendClickEvent(event ClickEvent) error {
	return nil
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
	var buf *bytes.Buffer
	err := json.NewEncoder(buf).Encode(event)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", c.endpoint+path, buf)
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

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	c.cookies = cookies

	return nil
}
