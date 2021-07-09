package http

import (
	"encoding/json"
	"net/http"
	"strings"
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

// NewMethodNotAllowedError generates a error payload for client
// and shows client's request method with its path, and returns which method this
// path supported
func NewMethodNotAllowedError(r *http.Request, supportedMethods []string) []byte {
	m := map[string]string{
		"error": "method_not_allowed",
		"error_description": "path " + r.URL.Path + " not allows method " +
			r.Method + ", allows: " + strings.Join(supportedMethods, ", "),
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return b

}

// NewNoPermissionForReadBoardArticlesError generates a error payload for telling client
// their don't have permission to read boardID
func NewNoPermissionForReadBoardArticlesError(r *http.Request, boardID string) []byte {
	m := map[string]string{
		"error":             "no_permission_for_read_board_articles",
		"error_description": "user don't have permission for read board " + boardID,
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return b

}
