package todoist

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const TodoistUrl = "https://beta.todoist.com/API/v8/"

var DefaultTimeout = 5 * time.Second

type QueryParam func(url.Values)

type Client struct {
	key    string
	client *http.Client
}

func NewClient(key string) Client {
	return Client{
		key: key,
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

func (c *Client) GetLabel(id int64) (Label, error) {
	url := fmt.Sprintf("labels/%d", id)
	resp, err := c.get(url)
	if err != nil {
		return Label{}, err
	}
	defer resp.Body.Close()

	var label Label
	if err := json.NewDecoder(resp.Body).Decode(&label); err != nil {
		return Label{}, err
	}

	return label, nil
}

func (c *Client) AllLabels() (Labels, error) {
	resp, err := c.get("/labels")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var labels Labels
	if err := json.NewDecoder(resp.Body).Decode(&labels); err != nil {
		return nil, err
	}

	return labels, nil
}

func WithProjectID(id int64) QueryParam {
	return func(q url.Values) {
		sid := strconv.FormatInt(id, 10)
		q.Add("project_id", sid)
	}
}

func WithLabelID(id int64) QueryParam {
	return func(q url.Values) {
		sid := strconv.FormatInt(id, 10)
		q.Add("label_id", sid)
	}
}

func WithRawFilter(s string) QueryParam {
	return func(q url.Values) {
		q.Add("filter", s)
	}
}

func (c *Client) GetTask(id int64) (Task, error) {
	url := fmt.Sprintf("tasks/%d", id)
	resp, err := c.get(url)
	if err != nil {
		return Task{}, err
	}
	defer resp.Body.Close()

	var task Task
	if err := json.NewDecoder(resp.Body).Decode(&task); err != nil {
		return Task{}, err
	}

	return task, nil
}

func (c *Client) ActiveTasks(params ...QueryParam) (Tasks, error) {
	resp, err := c.get("tasks", params...)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var tasks Tasks
	if err := json.NewDecoder(resp.Body).Decode(&tasks); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *Client) AllProjects() (Projects, error) {
	resp, err := c.get("/projects")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var projects Projects
	if err := json.NewDecoder(resp.Body).Decode(&projects); err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *Client) GetProject(id int64) (Project, error) {
	path := fmt.Sprintf("projects/%d", id)
	resp, err := c.get(path)
	if err != nil {
		return Project{}, err
	}
	defer resp.Body.Close()

	var project Project
	if err := json.NewDecoder(resp.Body).Decode(&project); err != nil {
		return Project{}, err
	}

	return project, nil
}

func (c *Client) get(path string, params ...QueryParam) (*http.Response, error) {
	url := TodoistUrl + strings.TrimPrefix(path, "/")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.key)

	q := req.URL.Query()
	for _, param := range params {
		param(q)
	}
	req.URL.RawQuery = q.Encode()

	return c.client.Do(req)
}
