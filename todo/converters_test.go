package todo

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"strconv"
	"testing"
	"time"
)

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}
	return t
}

func Test_serializeTodoDynamo(t *testing.T) {
	todoID := uuid.New()
	userID := uuid.New()

	createdAt := time.Now()
	updatedAt := createdAt.Add(time.Hour * 1)
	completedAt := updatedAt.Add(time.Hour * 1)

	tests := []struct {
		desc string
		todo *Todo

		expectedDynamoItem map[string]types.AttributeValue
	}{
		{
			desc: "happy path: valid completed todo item",
			todo: &Todo{
				ID:          todoID,
				UserID:      userID,
				CreatedAt:   createdAt,
				UpdatedAt:   &updatedAt,
				CompletedAt: &completedAt,
				Title:       "my completed todo",
				Description: "completed todo description",
			},
			expectedDynamoItem: map[string]types.AttributeValue{
				"ID":          &types.AttributeValueMemberS{Value: todoID.String()},
				"UserID":      &types.AttributeValueMemberS{Value: userID.String()},
				"CreatedAt":   &types.AttributeValueMemberN{Value: strconv.FormatInt(createdAt.UnixMilli(), 10)},
				"UpdatedAt":   &types.AttributeValueMemberN{Value: strconv.FormatInt(updatedAt.UnixMilli(), 10)},
				"CompletedAt": &types.AttributeValueMemberN{Value: strconv.FormatInt(completedAt.UnixMilli(), 10)},
				"Title":       &types.AttributeValueMemberS{Value: "my completed todo"},
				"Description": &types.AttributeValueMemberS{Value: "completed todo description"},
			},
		},
		{
			desc: "happy path: valid incomplete updated todo item",
			todo: &Todo{
				ID:          todoID,
				UserID:      userID,
				CreatedAt:   createdAt,
				UpdatedAt:   &updatedAt,
				Title:       "my completed todo",
				Description: "completed todo description",
			},
			expectedDynamoItem: map[string]types.AttributeValue{
				"ID":          &types.AttributeValueMemberS{Value: todoID.String()},
				"UserID":      &types.AttributeValueMemberS{Value: userID.String()},
				"CreatedAt":   &types.AttributeValueMemberN{Value: strconv.FormatInt(createdAt.UnixMilli(), 10)},
				"UpdatedAt":   &types.AttributeValueMemberN{Value: strconv.FormatInt(updatedAt.UnixMilli(), 10)},
				"Title":       &types.AttributeValueMemberS{Value: "my completed todo"},
				"Description": &types.AttributeValueMemberS{Value: "completed todo description"},
			},
		},
		{
			desc: "happy path: valid new todo item",
			todo: &Todo{
				ID:          todoID,
				UserID:      userID,
				CreatedAt:   createdAt,
				Title:       "my completed todo",
				Description: "completed todo description",
			},
			expectedDynamoItem: map[string]types.AttributeValue{
				"ID":          &types.AttributeValueMemberS{Value: todoID.String()},
				"UserID":      &types.AttributeValueMemberS{Value: userID.String()},
				"CreatedAt":   &types.AttributeValueMemberN{Value: strconv.FormatInt(createdAt.UnixMilli(), 10)},
				"Title":       &types.AttributeValueMemberS{Value: "my completed todo"},
				"Description": &types.AttributeValueMemberS{Value: "completed todo description"},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.desc, func(t *testing.T) {
			dynamoItem := serializeTodoDynamo(tc.todo)
			assert.Equal(t, len(tc.expectedDynamoItem), len(dynamoItem))
			for key := range tc.expectedDynamoItem {
				assert.Equal(t, tc.expectedDynamoItem[key], dynamoItem[key], key)
			}
			for key := range dynamoItem {
				assert.Equal(t, tc.expectedDynamoItem[key], dynamoItem[key], key)
			}
		})
	}
}

func Test_deserializeTodoDynamo(t *testing.T) {

}
