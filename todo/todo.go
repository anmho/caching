package todo

import (
	"github.com/google/uuid"
	"log/slog"
	"time"
)

type Todo struct {
	ID uuid.UUID `json:"id"`
	// UserID is the ID of the user that created the task.
	UserID      uuid.UUID  `json:"user_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
}

func (t *Todo) IsCompleted() bool {
	return t.CompletedAt != nil
}

func New(userID uuid.UUID, title, description string) *Todo {
	slog.Info(
		"new todo",
		slog.Any("description", description),
	)
	return &Todo{
		ID:          uuid.New(),
		UserID:      userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   nil,
		CompletedAt: nil,
		Title:       title,
		Description: description,
	}
}
