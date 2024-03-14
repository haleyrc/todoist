package todoist

import (
	"context"
	"fmt"
	"net/http"
)

type Duration struct {
	Amount int64  `json:"amount"`
	Unit   string `json:"unit"`
}

type Due struct {
	Date        string `json:"date"`
	IsRecurring bool   `json:"is_recurring"`
	Datetime    string `json:"datetime"`
	String      string `json:"string"`
	Timezone    string `json:"timezone"`
}

type Task struct {
	CreatorID    string    `json:"creator_id"`
	CreatedAt    string    `json:"created_at"`
	AssigneeID   string    `json:"assignee_id"`
	AssignerID   string    `json:"assigner_id"`
	CommentCount int64     `json:"comment_count"`
	IsCompleted  bool      `json:"is_completed"`
	Content      string    `json:"content"`
	Description  string    `json:"description"`
	Due          *Due      `json:"due"`
	Duration     *Duration `json:"duration"`
	ID           string    `json:"id"`
	Labels       []string  `json:"labels"`
	Order        int64     `json:"order"`
	Priority     int64     `json:"priority"`
	ProjectID    string    `json:"project_id"`
	SectionID    string    `json:"section_id"`
	ParentID     string    `json:"parent_id"`
	URL          string    `json:"url"`
}

type TaskParams struct {
	Content      string   `json:"content,omitempty"`
	Description  string   `json:"description,omitempty"`
	ProjectID    string   `json:"project_id,omitempty"` // Not in update
	SectionID    string   `json:"section_id,omitempty"` // Not in update
	ParentID     string   `json:"parent_id,omitempty"`  // Not in update
	Order        string   `json:"order,omitempty"`      // Not in update
	Labels       []string `json:"labels,omitempty"`
	Priority     int64    `json:"priority,omitempty"`
	DueString    string   `json:"due_string,omitempty"`
	DueDate      string   `json:"due_date,omitempty"`
	DueDatetime  string   `json:"due_datetime,omitempty"`
	DueLang      string   `json:"due_lang,omitempty"`
	AssigneeID   string   `json:"assignee_id,omitempty"`
	Duration     int64    `json:"duration,omitempty"`
	DurationUnit string   `json:"duration_unit,omitempty"`
}

func (c *Client) CloseTask(ctx context.Context, id string) error {
	req, err := c.post(fmt.Sprintf("/tasks/%s/close", id), nil)
	if err != nil {
		return fmt.Errorf("client: close task: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("client: close task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("client: close task: %w", responseError(resp))
	}

	return nil
}

func (c *Client) CreateTask(ctx context.Context, params TaskParams) (*Task, error) {
	req, err := c.post("/tasks", params)
	if err != nil {
		return nil, fmt.Errorf("client: create task: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, fmt.Errorf("client: create task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("client: create task: %w", responseError(resp))
	}

	var task Task
	if err := unmarshal(resp.Body, &task); err != nil {
		return nil, fmt.Errorf("client: create task: %w", err)
	}

	return &task, nil
}

func (c *Client) DeleteTask(ctx context.Context, id string) error {
	req, err := c.delete(fmt.Sprintf("/tasks/%s", id))
	if err != nil {
		return fmt.Errorf("client: delete task: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("client: delete task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("client: delete task: %w", responseError(resp))
	}

	return nil
}

func (c *Client) ReopenTask(ctx context.Context, id string) error {
	req, err := c.post(fmt.Sprintf("/tasks/%s/reopen", id), nil)
	if err != nil {
		return fmt.Errorf("client: reopen task: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return fmt.Errorf("client: reopen task: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("client: reopen task: %w", responseError(resp))
	}

	return nil
}
