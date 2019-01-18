package todoist

type Labels []Label

type Label struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Order int64  `json:"order"`
}
