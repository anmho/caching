package api

import (
	"github.com/google/uuid"
	"net/url"
)

type TodoFilters struct {
	UserID *uuid.UUID
}

type TodoFilterFunc func(f *TodoFilters)

func WithUserID(userID uuid.UUID) TodoFilterFunc {
	return func(f *TodoFilters) {
		f.UserID = &userID
	}
}

const (
	UserIDKey string = "user-id"
)

func parseFilters(queryParams url.Values) ([]TodoFilterFunc, error) {
	filters := make([]TodoFilterFunc, 0)
	for key, values := range queryParams {
		if len(values) == 0 {
			continue
		}
		// We will ignore the first value
		value := values[0]
		switch key {
		case UserIDKey:
			userID, err := uuid.Parse(value)
			if err != nil {
				return nil, err
			}
			filters = append(filters, WithUserID(userID))
		}
	}
	return filters, nil
}
