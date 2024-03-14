package todoist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"
	"os"
	"time"

	"github.com/haleyrc/uuid"
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

func (c *Client) delete(path string) (*http.Request, error) {
	return c.makeRequest(http.MethodDelete, path, nil)
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}

	if c.Verbose {
		if b, err := httputil.DumpResponse(resp, true); err == nil {
			fmt.Fprintln(os.Stderr, string(b))
		}
	}

	return resp, nil
}

func (c *Client) get(path string) (*http.Request, error) {
	return c.makeRequest(http.MethodGet, path, nil)
}

func (c *Client) makeRequest(method string, path string, body io.Reader) (*http.Request, error) {
	url := baseURL + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Request-Id", uuid.NewString())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))
	return req, nil
}

func (c *Client) post(path string, body any) (*http.Request, error) {
	var r io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		r = bytes.NewBuffer(bodyBytes)
	}
	req, err := c.makeRequest(http.MethodPost, path, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
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
