package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func JSON[T any](status int, data T, w http.ResponseWriter) error {
	return json.NewEncoder(w).Encode(data)
}

func Read[T any](body io.ReadCloser) (*T, error) {
	t := new(T)
	err := json.NewDecoder(body).Decode(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}

var _ error = (*APIError)(nil)

type ErrorOpt func(e *APIError)

func WithStatus(status int) func(e *APIError) {
	return func(e *APIError) {
		e.Status = status
	}
}
func WithMessage(message string) func(e *APIError) {
	return func(e *APIError) {
		e.Message = message
	}
}
func NewError(
	cause error,
	opts ...ErrorOpt,
) *APIError {
	e := &APIError{Cause: cause}
	for _, opt := range opts {
		opt(e)
	}
	if e.Status == 0 {
		e.Status = http.StatusInternalServerError
	}
	if e.Message == "" {
		e.Message = http.StatusText(e.Status)
	}
	return e
}

type APIError struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Cause   error  `json:"-"`
}

func (e *APIError) Error() string {
	return fmt.Sprintf("APIError - %s", e.Message)
}
