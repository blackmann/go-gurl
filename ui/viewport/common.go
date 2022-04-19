package viewport

import "net/http"

type headerItem struct {
	key   string
	value string
}

func (h headerItem) FilterValue() string {
	return h.key
}

func (h headerItem) Title() string {
	return h.key
}

func (h headerItem) Description() string {
	return h.value
}

type requestHeaders http.Header
