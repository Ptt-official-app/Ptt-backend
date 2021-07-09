package http

import (
	"encoding/json"
	"net/http"
)

// NewPathNotFoundError generates a path not found payload for client
// and shows client's request path in argument r.
func NewPathNotFoundError(r *http.Request) []byte {
	m := map[string]string{
		"error":             "not_found",
		"error_description": "path " + r.URL.Path + " not found",
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return b
}
