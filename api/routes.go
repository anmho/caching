package api

import (
	"database/sql"
	"net/http"
)

func registerRoutes(mux *http.ServeMux, db *sql.DB) {
	mux.HandleFunc("POST /", handleCreateTodo)
	mux.HandleFunc("GET /", handleGetTodo)
	mux.HandleFunc("GET /", handleCreateTodo)
	mux.HandleFunc("GET /", handleCreateTodo)
}

func handleCreateTodo(w http.ResponseWriter, r *http.Request) {

}

func handleGetTodo(w http.ResponseWriter, r *http.Request) {

}

func handleUpdateTodo(w http.ResponseWriter, r *http.Request) {

}

func handleDeleteTodo(w http.ResponseWriter, r *http.Request) {

}
