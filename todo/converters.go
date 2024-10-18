package todo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"strconv"
	"time"
)

func serializeTodoDynamo(todo *Todo) map[string]types.AttributeValue {
	if todo == nil {
		return nil
	}

	values := map[string]types.AttributeValue{
		"ID":          &types.AttributeValueMemberS{Value: todo.ID.String()},
		"UserID":      &types.AttributeValueMemberS{Value: todo.UserID.String()},
		"CreatedAt":   &types.AttributeValueMemberN{Value: strconv.FormatInt(todo.CreatedAt.UnixMilli(), 10)}, // Store as Unix timestamp
		"Title":       &types.AttributeValueMemberS{Value: todo.Title},
		"Description": &types.AttributeValueMemberS{Value: todo.Description},
	}

	if todo.UpdatedAt != nil {
		values["UpdatedAt"] = &types.AttributeValueMemberN{Value: strconv.FormatInt(todo.UpdatedAt.UnixMilli(), 10)}
	}

	if todo.IsCompleted() {
		values["CompletedAt"] = &types.AttributeValueMemberN{Value: strconv.FormatInt(todo.CompletedAt.UnixMilli(), 10)}
	}

	return values
}

func deserializeTodoDynamo(item map[string]types.AttributeValue) (*Todo, error) {
	todo := new(Todo)
	if idField, ok := item["ID"].(*types.AttributeValueMemberS); ok {
		id, err := uuid.Parse(idField.Value)
		if err != nil {
			return nil, err
		}

		todo.ID = id
	}

	if createdAtField, ok := item["CreatedAt"].(*types.AttributeValueMemberN); ok {
		createdAtUnixMillis, err := strconv.ParseInt(createdAtField.Value, 10, 64)
		if err != nil {
			return nil, err
		}

		createdAt := time.UnixMilli(createdAtUnixMillis)
		todo.CreatedAt = createdAt
	}

	if updatedAtField, ok := item["UpdatedAt"].(*types.AttributeValueMemberN); ok {
		updatedAtUnixMillis, err := strconv.ParseInt(updatedAtField.Value, 10, 64)
		if err != nil {
			return nil, err
		}

		updatedAt := time.UnixMilli(updatedAtUnixMillis)
		todo.UpdatedAt = &updatedAt
	}

	if completedAtField, ok := item["CompletedAt"].(*types.AttributeValueMemberN); ok {
		completedAtUnixMillis, err := strconv.ParseInt(completedAtField.Value, 10, 64)
		if err != nil {
			return nil, err
		}

		completedAt := time.UnixMilli(completedAtUnixMillis)
		todo.CompletedAt = &completedAt
	}

	if titleField, ok := item["Title"].(*types.AttributeValueMemberS); ok {
		title := titleField.Value
		todo.Title = title
	}

	if descriptionField, ok := item["Description"].(*types.AttributeValueMemberS); ok {
		description := descriptionField.Value
		todo.Description = description
	}

	return todo, nil
}
