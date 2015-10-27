package search

type Client struct {
}

type Config struct {
	Token string
}

func NewClient(c Config) (*Client, error) {
	return &client, nil
}

func (c *Client) Query(q Query) (*Response, error) {
}
