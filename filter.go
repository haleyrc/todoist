package todoist

type Filterable interface {
	Task | Project
}

func ByScheduled(t Task) bool {
	return t.Due != nil && t.Due.Datetime != ""
}

func Filter[S ~[]E, E Filterable](coll S, f func(v E) bool) (S, S) {
	var (
		filtered = S{}
		rest     = S{}
	)
	for _, val := range coll {
		if f(val) {
			filtered = append(filtered, val)
		} else {
			rest = append(rest, val)
		}
	}
	return filtered, rest
}
