package api

import (
	"errors"
	"github.com/anmho/caching/todo"
	"github.com/google/uuid"
	"log"
	"log/slog"
	"net/http"
)

func registerRoutes(mux *http.ServeMux, todoService *todo.Service) {
	register(mux, "POST /todos", handleCreateTodo(todoService))
	register(mux, "GET /todos", handleListTodos(todoService))
	register(mux, "GET /todos/{id}", handleGetTodoByID(todoService))
	register(mux, "PUT /todos/{id}", handleUpdateTodo(todoService))
	register(mux, "DELETE /todos/{id}", handleDeleteTodo(todoService))
}

type CreateTodoParams struct {
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
	UserID      string `json:"user_id" `
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
		idPath := r.PathValue("id")

		id, err := uuid.Parse(idPath)
		if err != nil {
			return NewError(err, WithStatus(http.StatusBadRequest))
		}

		userID, err := uuid.Parse(r.URL.Query().Get("user-id"))

		todoItem, err := todoService.FindTodoByID(r.Context(), userID, id)
		if err != nil {
			return err
		}

		return JSON(http.StatusOK, todoItem, w)
	}
}

func handleListTodos(todoService *todo.Service) RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		//filters, err := parseFilters(r.URL.Query())
		//if err != nil {
		//	return NewError(err, WithStatus(http.StatusBadRequest))
		//}

		userIDParam := r.URL.Query().Get("user-id")
		userID, err := uuid.Parse(userIDParam)
		if err != nil {
			return NewError(errors.New("user-id is required"), WithStatus(http.StatusBadRequest))
		}

		log.Println("user-id", userIDParam)

		todos, err := todoService.ListUserTodos(r.Context(), userID)

		if err != nil {
			return err
		}

		return JSON(http.StatusOK, todos, w)
	}
}

func handleUpdateTodo(todoService *todo.Service) RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) error {
		params, err := Read[todo.UpdateParams](r.Body)
		if err != nil {
			return NewError(err, WithStatus(http.StatusBadRequest))
		}
		//slog.Info("after parsing params", slog.Any("params", params))

		id, err := uuid.Parse(r.PathValue("id"))
		if err != nil {
			return NewError(err, WithStatus(http.StatusBadRequest))
		}

		userID, err := uuid.Parse(r.URL.Query().Get("user-id"))
		if err != nil {
			return NewError(err, WithStatus(http.StatusBadRequest))
		}

		slog.Info("handleUpdateTodo", slog.Any("params", params), slog.Any("userID", userID))
		err = todoService.UpdateTodo(r.Context(),
			userID,
			id,
			params)
		if err != nil {
			return err
		}

		return nil
	}
}

func handleDeleteTodo(todoService *todo.Service) RouteHandler {
	return func(w http.ResponseWriter, r *http.Request) error { return nil }
}
