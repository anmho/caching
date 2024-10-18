package api

import (
	"github.com/anmho/caching/todo"
	"net/http"
)

func New(todoService *todo.Service) *http.ServeMux {
	mux := http.NewServeMux()

	registerRoutes(mux, todoService)

	return mux
}
