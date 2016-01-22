package search

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// DEFAULTENDPOINT The default endpoint to use for the queries
const DEFAULTENDPOINT string = "https://cloudplatform.coveo.com/rest/search/"

type Client interface {
	Query(q Query) (*Response, error)
}

type Config struct {
	Token     string
	UserAgent string
	// Enpoint to use for the queries
	Endpoint string
}

func NewClient(c Config) (Client, error) {
	endpoint := DEFAULTENDPOINT
	if c.Endpoint != "" {
		endpoint = c.Endpoint
	}
	return &client{
		token:      c.Token,
		endpoint:   endpoint,
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

func (c *client) Query(q Query) (*Response, error) {
	marshalledQuery, err := json.Marshal(q)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewReader(marshalledQuery)

	req, err := http.NewRequest("POST", c.endpoint, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.token)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accepts", "application/json")
	req.Header.Set("User-Agent", c.useragent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	queryResponse := &Response{}
	err = json.NewDecoder(resp.Body).Decode(queryResponse)
	return queryResponse, err
}
