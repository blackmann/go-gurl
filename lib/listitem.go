package lib

type ListItem struct {
	Key   string
	Value string
	Ref   interface{}
}

func (h ListItem) FilterValue() string {
	return h.Key
}

func (h ListItem) Title() string {
	return h.Key
}

func (h ListItem) Description() string {
	return h.Value
}
