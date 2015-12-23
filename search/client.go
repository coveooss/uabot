package search

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	// EndpointProduction is the Search production endpoint
	EndpointProduction = "https://cloudplatform.coveo.com/rest/search/"
	// EndpointStaging is the Search staging endpoint
	EndpointStaging = "https://cloudplatformstaging.coveo.com/rest/search/"
	// EndpointDevelopment is the Search development endpoint
	EndpointDevelopment = "https://cloudplatformdev.coveo.com/rest/search/"
)

// Client is the search client to make search requests
type Client interface {
	Query(q Query) (*Response, error)
}

// Config is used to configure a new client
type Config struct {
	Token     string
	UserAgent string
	// Endpoint is used if you want to use custom endpoints (dev,staging,testing)
	Endpoint string
}

// NewClient returns a configured http search client using default http client
func NewClient(c Config) (Client, error) {
	if len(c.Endpoint) == 0 {
		c.Endpoint = EndpointProduction
	}

	return &client{
		token:      c.Token,
		endpoint:   c.Endpoint,
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
