package todo

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"log/slog"
	"time"
)

var (
	FieldNotFoundError = errors.New("field not found")
	DynamoDBTypeError  = errors.New("type of field in dynamo did not match expected type")
	DateParsingError   = errors.New("date string did not match expected format")
)

func NewDynamoDBTypeError(field string) error {
	return fmt.Errorf("%s: %w", field, DynamoDBTypeError)
}

func NewDateParsingError(field string, cause error) error {
	return fmt.Errorf("field %s did not match expected format: %w", field, cause)
}

func serializeTodoDynamo(todo *Todo) map[string]types.AttributeValue {
	if todo == nil {
		return nil
	}

	values := map[string]types.AttributeValue{
		"ID":          &types.AttributeValueMemberS{Value: todo.ID.String()},
		"UserID":      &types.AttributeValueMemberS{Value: todo.UserID.String()},
		"CreatedAt":   &types.AttributeValueMemberS{Value: formatDate(&todo.CreatedAt)}, // Store as Unix timestamp
		"Title":       &types.AttributeValueMemberS{Value: todo.Title},
		"Description": &types.AttributeValueMemberS{Value: todo.Description},
	}

	if todo.UpdatedAt != nil {
		values["UpdatedAt"] = &types.AttributeValueMemberS{Value: formatDate(todo.UpdatedAt)}
	}

	if todo.IsCompleted() {
		values["CompletedAt"] = &types.AttributeValueMemberS{Value: formatDate(todo.CompletedAt)}
	}

	return values
}

func parseDynamoField[T types.AttributeValue](fieldName string, item map[string]types.AttributeValue) (T, error) {
	if field, ok := item[fieldName].(T); ok {
		return field, nil
	} else {
		return *new(T), fmt.Errorf("%s: %w", fieldName, FieldNotFoundError) // could also be incorrect type
	}
}

func parseStringFromDynamo(fieldName string, item map[string]types.AttributeValue) (string, error) {
	field, ok := item[fieldName].(*types.AttributeValueMemberS)
	if !ok {
		return "", fmt.Errorf("%s: %w", fieldName, FieldNotFoundError)
	}
	return field.Value, nil
}

func parseUUIDFromDynamo(fieldName string, item map[string]types.AttributeValue) (uuid.UUID, error) {
	uuidField, err := parseDynamoField[*types.AttributeValueMemberS](fieldName, item)
	if err != nil {
		return uuid.Nil, err
	}

	id, err := uuid.Parse(uuidField.Value)
	if err != nil {
		return uuid.Nil, nil
	}
	return id, nil
}

func parseDateFromDynamo(fieldName string, item map[string]types.AttributeValue) (time.Time, error) {
	dateField, err := parseDynamoField[*types.AttributeValueMemberS](fieldName, item)
	if err != nil {
		return time.Time{}, err
	}

	date, err := parseDate(dateField.Value)
	if err != nil {
		return time.Time{}, err
	}

	return date, nil
}

func deserializeTodoDynamo(item map[string]types.AttributeValue) (*Todo, error) {
	todo := new(Todo)
	var err error

	todo.UserID, err = parseUUIDFromDynamo("UserID", item)
	if err != nil {
		return nil, err
	}

	todo.ID, err = parseUUIDFromDynamo("ID", item)
	if err != nil {
		return nil, err
	}

	// mandatory field
	todo.CreatedAt, err = parseDateFromDynamo("CreatedAt", item)
	if err != nil {
		return nil, err
	}

	updatedAtField, hasField := item["UpdatedAt"]
	if hasField {
		updatedAtField, ok := updatedAtField.(*types.AttributeValueMemberS)
		if !ok {
			return nil, NewDynamoDBTypeError("UpdatedAt")
		}

		updatedAt, err := parseDate(updatedAtField.Value)
		if err != nil {
			return nil, NewDateParsingError("UpdatedAt", err)
		}
		todo.UpdatedAt = &updatedAt
	}
	completedAtField, hasField := item["CompletedAt"]
	if hasField {
		completedAtField, ok := completedAtField.(*types.AttributeValueMemberS)
		if !ok {
			return nil, NewDynamoDBTypeError("CompletedAt")
		}

		completedAt, err := parseDate(completedAtField.Value)
		if err != nil {
			return nil, NewDateParsingError("CompletedAt", err)
		}

		todo.CompletedAt = &completedAt
	}

	todo.Title, err = parseStringFromDynamo("Title", item)
	if err != nil {
		return nil, err
	}

	todo.Description, err = parseStringFromDynamo("Description", item)
	if err != nil {
		return nil, err
	}

	slog.Info("deserializeTodoDynamo",
		slog.Any("item", item),
		slog.Any("todo", todo))

	return todo, nil
}
