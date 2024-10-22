package todo

import (
	"github.com/google/uuid"
	"net/url"
)

type Filters struct {
	UserID *uuid.UUID
}

type FilterFunc func(f *Filters)

func WithUserID(userID uuid.UUID) FilterFunc {
	return func(f *Filters) {
		f.UserID = &userID
	}
}

const (
	UserIDKey string = "user-id"
)

func parseFilters(queryParams url.Values) ([]FilterFunc, error) {
	filters := make([]FilterFunc, 0)
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
