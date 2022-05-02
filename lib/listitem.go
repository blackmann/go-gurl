package lib

type Pair struct {
	Key   string
	Value string
	Ref   interface{}
}

func (h Pair) FilterValue() string {
	return h.Key
}

func (h Pair) Title() string {
	return h.Key
}

func (h Pair) Description() string {
	return h.Value
}
