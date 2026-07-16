package zazu

import (
	"fmt"
	"net/http"
	"strconv"
)

// Error is the API error envelope, mirroring the other Zazu SDKs' hierarchy:
// { "error": { "type": ..., "message": ..., "param": ... } }. Match on Kind
// (or use errors.As with the sentinel helpers) instead of subclassing.
type Error struct {
	Status     int
	Kind       string // authentication, forbidden, not_found, validation, rate_limit, server, api
	Type       string // the API's error.type field
	Message    string
	Param      string
	RequestID  string
	RetryAfter int // seconds; only set for rate_limit
	Body       map[string]any
}

func (e *Error) Error() string {
	if e.Param != "" {
		return fmt.Sprintf("zazu: %s (%d %s, param %s)", e.Message, e.Status, e.Kind, e.Param)
	}
	return fmt.Sprintf("zazu: %s (%d %s)", e.Message, e.Status, e.Kind)
}

// ConfigurationError is returned by New when the client can't be built.
type ConfigurationError struct{ Message string }

func (e *ConfigurationError) Error() string { return "zazu: " + e.Message }

// ConnectionError wraps transport-level failures (timeouts, DNS, refused).
type ConnectionError struct{ Message string }

func (e *ConnectionError) Error() string { return "zazu: connection error: " + e.Message }

func newError(status int, header http.Header, body map[string]any) *Error {
	e := &Error{
		Status:    status,
		Body:      body,
		RequestID: header.Get("X-Request-Id"),
	}

	if payload, ok := body["error"].(map[string]any); ok {
		e.Type, _ = payload["type"].(string)
		e.Message, _ = payload["message"].(string)
		e.Param, _ = payload["param"].(string)
	}

	switch {
	case status == 401:
		e.Kind = "authentication"
	case status == 403:
		e.Kind = "forbidden"
	case status == 404:
		e.Kind = "not_found"
	case status == 422:
		e.Kind = "validation"
	case status == 429:
		e.Kind = "rate_limit"
		if retry := header.Get("Retry-After"); retry != "" {
			e.RetryAfter, _ = strconv.Atoi(retry)
		}
	case status >= 500:
		e.Kind = "server"
	default:
		e.Kind = "api"
	}

	if e.Message == "" {
		e.Message = http.StatusText(status)
	}
	return e
}
