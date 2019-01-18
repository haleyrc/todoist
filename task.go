package todoist

import "time"

type Due struct {
	Date      string    `json:"date"`
	Recurring bool      `json:"recurring"`
	Timestamp time.Time `json:"datetime"`
	String    string    `json:"string"`
	Timezone  string    `json:"timezone"`
}

type Tasks []Task

type Task struct {
	CommentCount int64   `json:"comment_count"`
	Completed    bool    `json:"completed"`
	Content      string  `json:"content"`
	Due          Due     `json:"due"`
	ID           int64   `json:"id"`
	Indent       int64   `json:"indent"`
	LabelIDs     []int64 `json:"label_ids"`
	Order        int64   `json:"order"`
	Priority     int64   `json:"priority"`
	ProjectID    int64   `json:"project_id"`
	URL          string  `json:"url"`
}
