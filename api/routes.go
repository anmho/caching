package api

import "net/http"

func registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /", handleCreateTodo)
	mux.HandleFunc("GET /", handleGetTodo)
	mux.HandleFunc("GET /", handleCreateTodo)
	mux.HandleFunc("GET /", handleCreateTodo)
}

type Todo struct {
	Title       string
	Description string
	Completed   bool
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {

}

func handleGetTodo(w http.ResponseWriter, r *http.Request) {

}

func handleUpdateTodo(w http.ResponseWriter, r *http.Request) {

}

func handleDeleteTodo(w http.ResponseWriter, r *http.Request) {

}
