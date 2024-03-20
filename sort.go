package todoist

import (
	"cmp"
	"slices"
	"time"
)

type Sortable interface {
	Task | Project
}

func ByDatetime(a, b Task) int {
	ti, _ := time.Parse(time.RFC3339, a.Due.Datetime)
	tj, _ := time.Parse(time.RFC3339, b.Due.Datetime)
	switch {
	case ti.After(tj):
		return 1
	case ti.Equal(tj):
		return 0
	default:
		return -1
	}
}

func ByPriority(a, b Task) int {
	return cmp.Compare(b.Priority, a.Priority)
}

func Sort[S ~[]E, E Sortable](coll S, f func(a, b E) int) {
	slices.SortFunc(coll, f)
}
