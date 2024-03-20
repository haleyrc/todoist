package todoist

import (
	"context"
	"fmt"
	"net/http"
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
	req, err := c.post("/projects", params)
	if err != nil {
		return nil, fmt.Errorf("todoist: create project: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("todoist: create project: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("todoist: create project: %w", responseError(resp))
	}

	var project Project
	if err := unmarshal(resp.Body, &project); err != nil {
		return nil, fmt.Errorf("todoist: create project: %w", err)
	}

	return &project, nil
}

func (c *Client) DeleteProject(ctx context.Context, id string) error {
	req, err := c.delete(fmt.Sprintf("/projects/%s", id))
	if err != nil {
		return fmt.Errorf("todoist: delete project: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("todoist: delete project: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("todoist: delete project: %w", responseError(resp))
	}

	return nil
}

func (c *Client) GetProject(ctx context.Context, projectID string) (*Project, error) {
	req, err := c.get(fmt.Sprintf("/projects/%s", projectID))
	if err != nil {
		return nil, fmt.Errorf("todoist: get project: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("todoist: get project: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("todoist: get project: %w", responseError(resp))
	}

	var project Project
	if err := unmarshal(resp.Body, &project); err != nil {
		return nil, fmt.Errorf("todoist: get project: %w", err)
	}

	return &project, nil
}

func (c *Client) GetProjects(ctx context.Context) ([]Project, error) {
	req, err := c.get("/projects")
	if err != nil {
		return nil, fmt.Errorf("todoist: get projects: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("todoist: get projects: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("todoist: get projects: %w", responseError(resp))
	}

	var projects []Project
	if err := unmarshal(resp.Body, &projects); err != nil {
		return nil, fmt.Errorf("todoist: get projects: %w", err)
	}

	return projects, nil
}
