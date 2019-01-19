package todoist

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pborman/uuid"
)

var (
	ErrProjectNotFound = errors.New("project not found")
	ErrLabelNotFound   = errors.New("label not found")
)

const TodoistUrl = "https://beta.todoist.com/API/v8/"

var DefaultTimeout = 5 * time.Second

type QueryParam func(url.Values)

type Client struct {
	key    string
	client *http.Client

	// TODO (RCH): Some way to signal to refresh cache is called for
	projects Projects
	labels   []Label
}

func NewClient(key string) Client {
	return Client{
		key: key,
		client: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

func (c *Client) findLabelInCache(name string) (Label, bool) {
	for _, p := range c.labels {
		if p.Name == name {
			return p, true
		}
	}
	return Label{}, false
}

func (c *Client) FindLabel(name string) (Label, error) {
	if l, found := c.findLabelInCache(name); found {
		return l, nil
	}
	if _, err := c.AllLabels(); err != nil {
		return Label{}, nil
	}
	if l, found := c.findLabelInCache(name); found {
		return l, nil
	}
	return Label{}, ErrLabelNotFound
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
	c.labels = labels

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

type NewTask struct {
	Content   string
	Project   string
	Labels    []string
	Priority  int
	DueString string
}

func (c *Client) AddTask(nt NewTask) (Task, error) {
	type NewTaskRequest struct {
		Content   string  `json:"content"`
		Project   int64   `json:"project_id,omitempty"`
		Labels    []int64 `json:"label_ids,omitempty"`
		Priority  int     `json:"priority"`
		DueString string  `json:"due_string,omitempty"`
	}
	ntr := NewTaskRequest{
		Content:   nt.Content,
		DueString: nt.DueString,
		Priority:  1,
	}

	if nt.Priority <= 4 && nt.Priority >= 1 {
		ntr.Priority = nt.Priority
	}

	if nt.Project != "" {
		project, err := c.FindProject(nt.Project)
		if err != nil {
			return Task{}, err
		}
		ntr.Project = project.ID
	}

	if nt.Labels != nil && len(nt.Labels) > 0 {
		// TODO (RCH): Parallelize this
		labels := make([]int64, 0, len(nt.Labels))
		for _, lbl := range nt.Labels {
			label, err := c.FindLabel(lbl)
			if err != nil {
				return Task{}, err
			}
			labels = append(labels, label.ID)
		}
		ntr.Labels = labels
	}

	var buff bytes.Buffer
	if err := json.NewEncoder(&buff).Encode(&ntr); err != nil {
		return Task{}, err
	}
	resp, err := c.post("tasks", &buff)
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

func (c *Client) findProjectInCache(name string) (Project, bool) {
	return c.projects.FindName(name)
}

func (c *Client) FindProject(name string) (Project, error) {
	if p, found := c.findProjectInCache(name); found {
		return p, nil
	}
	if _, err := c.AllProjects(); err != nil {
		return Project{}, nil
	}
	if p, found := c.findProjectInCache(name); found {
		return p, nil
	}
	return Project{}, ErrProjectNotFound
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
	c.projects = projects

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

func (c *Client) post(path string, r io.Reader) (*http.Response, error) {
	url := TodoistUrl + strings.TrimPrefix(path, "/")

	req, err := http.NewRequest("POST", url, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+c.key)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Request-ID", uuid.New())

	return c.client.Do(req)
}
