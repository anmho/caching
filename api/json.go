package api

import (
	"encoding/json"
	"io"
	"net/http"
)



func JSON[T any](data T, w http.ResponseWriter) error {
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