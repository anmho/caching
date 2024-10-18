package api

import (
	"github.com/anmho/caching/todo"
	"github.com/google/uuid"
	"net/http"
)

func registerRoutes(mux *http.ServeMux, todoService *todo.Service) {
	register(mux, "POST /todos", handleCreateTodo(todoService))
	register(mux, "GET /todos/{id}", handleGetTodoByID(todoService))
	register(mux, "PUT /todos/{id}", handleUpdateTodo(todoService))
	register(mux, "DELETE /todos/{id}", handleDeleteTodo(todoService))
}

type CreateTodoParams struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	UserID      string `json:"user_id"`
}

func handleCreateTodo(todoService *todo.Service) RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		params, err := Read[CreateTodoParams](r.Body)
		if err != nil {
			return NewError(err, WithStatus(http.StatusBadRequest))
		}

		userID, err := uuid.Parse(params.UserID)
		if err != nil {
			return NewError(err, WithStatus(http.StatusBadRequest))
		}

		newTodo, err := todoService.CreateTodo(
			r.Context(),
			userID,
			params.Title,
			params.Description,
		)

		// Something went wrong
		if err != nil {
			return err
		}

		return JSON(http.StatusCreated, newTodo, w)
	}
}

func handleGetTodoByID(todoService *todo.Service) RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		idPath := r.PathValue("idPath")

		id, err := uuid.Parse(idPath)
		if err != nil {
			return NewError(err, WithStatus(http.StatusBadRequest))
		}

		todoItem, err := todoService.FindTodoByID(r.Context(), id)
		if err != nil {
			// handle it centrally
			return err
		}

		return JSON(http.StatusOK, todoItem, w)
	}
}

func handleGetTodos(todoService *todo.Service) RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) error { return nil }
}

func handleUpdateTodo(todoService *todo.Service) RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) error { return nil }
}

func handleDeleteTodo(todoService *todo.Service) RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) error { return nil }
}
