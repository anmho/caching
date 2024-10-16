package todo

import "github.com/google/uuid"

type Todo struct {
	ID          uuid.UUID
	Title       string
	Description string
	Completed   bool
}

func NewTodo(title, description string) Todo {
	return Todo{
		Title:       title,
		Description: description,
		Completed:   false,
	}
}
