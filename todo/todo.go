package todo

import (
	"github.com/google/uuid"
	"time"
)

type Todo struct {
	ID uuid.UUID
	// UserID is the ID of the user that created the task.
	UserID      uuid.UUID
	CreatedAt   time.Time
	CompletedAt *time.Time
	Title       string
	Description string
}

func (t *Todo) IsCompleted() bool {
	return t.CompletedAt != nil
}

func New(userID uuid.UUID, title, description string) *Todo {
	return &Todo{
		ID:          uuid.New(),
		UserID:      userID,
		CreatedAt:   time.Now(),
		CompletedAt: nil,
		Title:       title,
		Description: description,
	}
}
