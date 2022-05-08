package viewport

import (
	"github.com/blackmann/go-gurl/lib"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestResponseHeadersModel_Update_Headers(t *testing.T) {
	instance := newResponseHeadersModel()

	headers := http.Header{}
	headers.Set("content-type", "application/json")

	instance, _ = instance.Update(lib.Response{Headers: headers})

	assert.Contains(t, instance.View(), "Content-Type")
	assert.Contains(t, instance.View(), "application/json")
}

func TestResponseHeadersModel_Update_Reset(t *testing.T) {
	instance := newResponseHeadersModel()

	headers := http.Header{}
	headers.Set("content-type", "application/json")

	instance, _ = instance.Update(lib.Response{Headers: headers})

	// confirm it exists
	assert.Contains(t, instance.View(), "Content-Type")

	instance, _ = instance.Update(lib.Reset)
	assert.NotContains(t, instance.View(), "Content-Type")
}
