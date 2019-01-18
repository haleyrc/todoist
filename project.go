package todoist

import (
	"errors"
)

type Project struct {
	CommentCount int64  `json:"comment_count"`
	ID           int64  `json:"id"`
	Indent       int64  `json:"indent"`
	Name         string `json:"name"`
	Order        int64  `json:"order"`
}

type Projects []Project

func (ps Projects) FindName(n string) (Project, error) {
	for _, p := range ps {
		if p.Name == n {
			return p, nil
		}
	}
	return Project{}, errors.New("project not found")
}
