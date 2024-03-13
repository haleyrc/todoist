package todoist

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/haleyrc/uuid"
)

type Project struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Color          string `json:"color"`
	ParentID       string `json:"parent_id"`
	Order          int64  `json:"order"`
	CommentCount   int64  `json:"comment_count"`
	IsShared       bool   `json:"is_shared"`
	IsFavorite     bool   `json:"is_favorite"`
	IsInboxProject bool   `json:"is_inbox_project"`
	IsTeamInbox    bool   `json:"is_team_inbox"`
	ViewStyle      string `json:"view_style"`
	URL            string `json:"url"`
}

func (c *Client) GetProjects(ctx context.Context) ([]Project, error) {
	url := baseURL + "/projects"

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("client: get projects: %w", err)
	}

	req.Header.Set("X-Request-Id", uuid.NewString())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: get projects: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client: get projects: %w", responseError(resp))
	}

	var projects []Project
	if err := unmarshal(resp.Body, &projects); err != nil {
		return nil, fmt.Errorf("client: get projects: %w", err)
	}

	return projects, nil
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
