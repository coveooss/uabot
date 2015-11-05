package analytics

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Client interface {
	SendSearchEvent(*SearchEvent, *http.Cookie) (string, []*http.Cookie, error)
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

func (c *client) sendEventRequest(path string, buf io.Reader) error {
	req, err := http.NewRequest("POST", c.endpoint+path, buf)
	if err != nil {
		return err
	}

	if c.cookie != nil {
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

func (c *client) SendSearchEvent(event *SearchEvent) error {
	marshalledEvent, err := json.Marshal(event)
	if err != nil {
		return err
	}
	buf := bytes.NewReader(marshalledEvent)

	err = c.sendEventRequest("search/", buf)
	if err != nil {
		return err
	}

	return nil
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
func (c *client) GetStatus() (*StatusResponse, error) {
	return nil, nil
}
