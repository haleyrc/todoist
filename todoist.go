package todoist

import (
	"fmt"
	"net/http"
	"time"
)

const baseURL = "https://api.todoist.com/rest/v2"

type Client struct {
	APIToken   string
	HTTPClient *http.Client
}

func NewClient(apiToken string) (*Client, error) {
	if apiToken == "" {
		return nil, fmt.Errorf("new client: api token is requried")
	}

	c := &Client{
		APIToken:   apiToken,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
	}

	return c, nil
}

type Error struct {
	Code     int
	Body     string
	Response string
}

func (err Error) Error() string {
	return fmt.Sprintf("todoist: response error: %d", err.Code)
}
