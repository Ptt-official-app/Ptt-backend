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
// they don't have permission to read boardID
func NewNoPermissionForReadBoardArticlesError(r *http.Request, boardID string) []byte {
	m := map[string]string{
		"error":             "no_permission_for_read_board_articles",
		"error_description": "user don't have permission for read board " + boardID,
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return b
}

// NewPermissionError generates a error payload for generic permission error
func NewPermissionError(r *http.Request, err error) []byte {
	m := map[string]string{
		"error":             "permission_error",
		"error_description": err.Error(),
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return b
}

// NewPermissionError generates a error payload for generic server error
func NewServerError(r *http.Request, err error) []byte {
	m := map[string]string{
		"error":             "server_error",
		"error_description": err.Error(),
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return b
}

// NewNoRequriedParameterError generates a error payload for telling client which
// parameter is required
func NewNoRequiredParameterError(r *http.Request, requireParameter string) []byte {
	m := map[string]string{
		"error":             "no_required_parameter",
		"error_description": "required parameter: " + requireParameter + " not found",
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return b
}

func NewParameterShouldBeIntegerError(r *http.Request, parameter string) []byte {
	m := map[string]string{
		"error":             "parameter_should_be_integer",
		"error_description": parameter + " should be integer",
	}
	b, _ := json.MarshalIndent(m, "", "  ")
	return b
}

// NewBoardError generates a error payload when create board fail
func NewNewBoardError(r *http.Request) []byte {
	boardID := r.PostFormValue("title")

	m := map[string]string{
		"error":             "new_board_error",
		"error_description": "new board " + boardID + " failed",
	}
	b, _ := json.MarshalIndent(m, "", "  ")

	return b
}

// TokenPermission generates a error payload for telling client token is not pass
func TokenPermissionError(r *http.Request, err error) []byte {
	m := map[string]string{
		"error":             "token is not pass",
		"error_description": err.Error(),
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return b
}

// NewNoPermissionForReadBoardArticlesError generates a error payload for telling client
// they don't have permission to read boardID
func NewNoPermissionForCreateBoardArticlesError(r *http.Request, boardID string) []byte {
	m := map[string]string{
		"error":             "no_permission_for_create_board_articles",
		"error_description": "user don't have permission to create an article in board " + boardID,
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return b
}

func NewNoPermissionForForwardArticleError(r *http.Request, filename string) []byte {
	m := map[string]string{
		"error":             "no_permission_for_forward_board_article",
		"error_description": "user don't have permission to forward article " + filename,
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return b
}

func NewNoPermissionForForwardArticleToEmailError(r *http.Request, filename, email string) []byte {
	m := map[string]string{
		"error":             "no_permission_for_forward_board_article_to_email",
		"error_description": "user don't have permission to forward article " + filename + " to email " + email,
	}

	b, _ := json.MarshalIndent(m, "", "  ")
	return b
}
