package todoist

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/haleyrc/uuid"
)

const (
	ColorBerryRed   = "berry_red"
	ColorBlue       = "blue"
	ColorCharcoal   = "charcoal"
	ColorGrape      = "grape"
	ColorGreen      = "green"
	ColorGrey       = "grey"
	ColorLavender   = "lavender"
	ColorLightBlue  = "light_blue"
	ColorLimeGreen  = "lime_green"
	ColorMagenta    = "magenta"
	ColorMintGreen  = "mint_green"
	ColorOliveGreen = "olive_green"
	ColorOrange     = "orange"
	ColorRed        = "red"
	ColorSalmon     = "salmon"
	ColorSkyBlue    = "sky_blue"
	ColorTaupe      = "taupe"
	ColorTeal       = "teal"
	ColorViolet     = "violet"
	ColorYellow     = "yellow"
)

const (
	ViewStyleBoard = "board"
	ViewStyleList  = "list"
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

type ProjectParams struct {
	Name       string `json:"name,omitempty"`
	ParentID   string `json:"parent_id,omitempty"`
	Color      string `json:"color,omitempty"`
	IsFavorite bool   `json:"is_favorite,omitempty"`
	ViewStyle  string `json:"view_style,omitempty"`
}

func (c *Client) CreateProject(ctx context.Context, params ProjectParams) (*Project, error) {
	url := baseURL + "/projects"

	jsonBytes, err := json.Marshal(params)
	if err != nil {
		return nil, fmt.Errorf("client: create project: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, fmt.Errorf("client: create project: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-Id", uuid.NewString())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: create project: %w", err)
	}
	defer resp.Body.Close()

	if c.Verbose {
		if b, err := httputil.DumpResponse(resp, true); err == nil {
			fmt.Fprintln(os.Stderr, string(b))
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client: create project: %w", responseError(resp))
	}

	var project Project
	if err := unmarshal(resp.Body, &project); err != nil {
		return nil, fmt.Errorf("client: create project: %w", err)
	}

	return &project, nil
}

func (c *Client) DeleteProject(ctx context.Context, id string) error {
	url := baseURL + fmt.Sprintf("/projects/%s", id)

	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return fmt.Errorf("client: delete project: %w", err)
	}

	req.Header.Set("X-Request-Id", uuid.NewString())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("client: delete project: %w", err)
	}
	defer resp.Body.Close()

	if c.Verbose {
		if b, err := httputil.DumpResponse(resp, true); err == nil {
			fmt.Fprintln(os.Stderr, string(b))
		}
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("client: delete project: %w", responseError(resp))
	}

	return nil
}

func (c *Client) GetProject(ctx context.Context, projectID string) (*Project, error) {
	url := baseURL + fmt.Sprintf("/projects/%s", projectID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("client: get project: %w", err)
	}

	req.Header.Set("X-Request-Id", uuid.NewString())
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.APIToken))

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("client: get project: %w", err)
	}
	defer resp.Body.Close()

	if c.Verbose {
		if b, err := httputil.DumpResponse(resp, true); err == nil {
			fmt.Fprintln(os.Stderr, string(b))
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client: get project: %w", responseError(resp))
	}

	var project Project
	if err := unmarshal(resp.Body, &project); err != nil {
		return nil, fmt.Errorf("client: get project: %w", err)
	}

	return &project, nil
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

	if c.Verbose {
		if b, err := httputil.DumpResponse(resp, true); err == nil {
			fmt.Fprintln(os.Stderr, string(b))
		}
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client: get projects: %w", responseError(resp))
	}

	var projects []Project
	if err := unmarshal(resp.Body, &projects); err != nil {
		return nil, fmt.Errorf("client: get projects: %w", err)
	}

	return projects, nil
}
