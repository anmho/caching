package api

import (
	"log/slog"
	"net/http"
)

type RouteHandler func(w http.ResponseWriter, r *http.Request) error

func handleError(w http.ResponseWriter, r *http.Request, err error) {
	switch e := err.(type) {
	case *APIError:
		slog.Error(
			"API error occurred in route",
			slog.String("method", r.Method),
			slog.String("path", r.URL.Path),
			slog.Int("status", e.Status),
			slog.Any("message", e.Message),
			slog.Any("cause", e.Cause),
		)
		JSON(e.Status, e, w)
	default:
		slog.Error("error occurred in route", slog.Any("error", err))
		JSON(
			http.StatusInternalServerError,
			NewError(err, WithStatus(http.StatusInternalServerError)), w)
	}
}

func createHandler(handler RouteHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := handler(w, r)
		if err != nil {
			handleError(w, r, err)
		}
	}
}

func register(mux *http.ServeMux, pattern string, handler RouteHandler) {
	mux.HandleFunc(pattern, createHandler(handler))
}
