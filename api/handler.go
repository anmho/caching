package api

import (
	"log/slog"
	"net/http"
)

type RouteHandler func(w http.ResponseWriter, r *http.Request) error

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	slog.Error(
		"error occurred in route",
		slog.String("path", r.URL.Path),
		slog.Any("error", err),
	)
}

func createHandler(handler RouteHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		handleError(w, r, err)
	}
}

func register(mux *http.ServeMux, pattern string, handler RouteHandler) {
	mux.HandleFunc(pattern, createHandler(handler))
}
