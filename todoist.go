package todoist

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"time"
)

const baseURL = "https://api.todoist.com/rest/v2"

type Client struct {
	APIToken   string
	HTTPClient *http.Client
	Verbose    bool
}

func NewClient(apiToken string) (*Client, error) {
	if apiToken == "" {
		return nil, fmt.Errorf("new client: api token is requried")
	}

	c := &Client{
		APIToken:   apiToken,
		HTTPClient: &http.Client{Timeout: 5 * time.Second},
		Verbose:    false,
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

func responseError(resp *http.Response) error {
	re := Error{Code: resp.StatusCode}
	if responseBytes, err := httputil.DumpResponse(resp, true); err != nil {
		re.Response = string(responseBytes)
	}
	if bodyBytes, err := io.ReadAll(resp.Body); err != nil {
		re.Body = string(bodyBytes)
	}
	return re
}

func unmarshal(r io.Reader, dst any) error {
	bytes, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	if err := json.Unmarshal(bytes, dst); err != nil {
		return fmt.Errorf("unmarshal: %w", err)
	}

	return nil
}
