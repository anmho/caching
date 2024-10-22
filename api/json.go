package api

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"io"
	"net/http"
)

var (
	validate = validator.New(validator.WithRequiredStructEnabled())
)

func JSON[T any](status int, data T, w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func Read[T any](body io.ReadCloser) (*T, error) {

	t := new(T)
	dec := json.NewDecoder(body)
	dec.DisallowUnknownFields()

	err := json.NewDecoder(body).Decode(t)

	if err != nil {
		return nil, err
	}
	err = validate.Struct(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
