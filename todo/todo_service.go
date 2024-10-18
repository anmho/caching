package todo

import (
	"context"
	"fmt"
	"github.com/anmho/caching/cache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"log/slog"
	"strconv"
	"time"
)

const TodoItemsTableName = "TodoItems"

type Service struct {
	dynamoClient  *dynamodb.Client
	cacheStrategy cache.Strategy
}

func WithCacheStrategy(strategy cache.Strategy) func(s *Service) {
	return func(s *Service) {
		s.cacheStrategy = strategy
	}
}

func MakeService(dynamoClient *dynamodb.Client, opts ...func(o *Service)) *Service {
	s := &Service{
		dynamoClient: dynamoClient,
	}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

func serializeTodoDynamo(todo *Todo) map[string]types.AttributeValue {
	if todo == nil {
		return nil
	}

	values := map[string]types.AttributeValue{
		"ID":          &types.AttributeValueMemberS{Value: todo.ID.String()},
		"UserID":      &types.AttributeValueMemberS{Value: todo.UserID.String()},
		"CompletedAt": nil,
		"CreatedAt":   &types.AttributeValueMemberN{Value: fmt.Sprintf("%d", todo.CreatedAt.UnixMilli())}, // Store as Unix timestamp
		"Title":       &types.AttributeValueMemberS{Value: todo.Title},
		"Description": &types.AttributeValueMemberS{Value: todo.Description},
	}

	if todo.IsCompleted() {
		values["CompletedAt"] = &types.AttributeValueMemberS{Value: todo.CompletedAt.Format(time.RFC3339)}
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

func (s *Service) CreateTodo(
	ctx context.Context,
	userID uuid.UUID,
	title string,
	description string,
) (*Todo, error) {

	todo := New(userID, title, description)

	dynamoItem := serializeTodoDynamo(todo)
	result, err := s.dynamoClient.PutItem(ctx, &dynamodb.PutItemInput{
		Item:      dynamoItem,
		TableName: aws.String(TodoItemsTableName),
	})
	if err != nil {
		return nil, err
	}

	slog.Info("create todo result", slog.Any("result", result))

	return todo, nil
}

func (s *Service) FindTodoByID(ctx context.Context, id uuid.UUID) (*Todo, error) {
	params := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id.String()},
		},
		TableName:              aws.String(TodoItemsTableName),
		ConsistentRead:         aws.Bool(true),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
	}

	result, err := s.dynamoClient.GetItem(ctx, params)
	if err != nil {
		return nil, err
	}

	todo, err := deserializeTodoDynamo(result.Item)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *Service) GetAllTodosForUser(ctx context.Context, userID uuid.UUID) (*Todo, error) {
	//input := &dynamodb.ScanInput{
	//	TableName:                 aws.String(TodoItemsTableName),
	//	ConditionalOperator:       "",
	//	ConsistentRead:            nil,
	//	ExclusiveStartKey:         nil,
	//	ExpressionAttributeNames:  nil,
	//	ExpressionAttributeValues: nil,
	//	FilterExpression:          nil,
	//	IndexName:                 nil,
	//	Limit:                     nil,
	//	ProjectionExpression:      nil,
	//	ReturnConsumedCapacity:    "",
	//	ScanFilter:                nil,
	//	Segment:                   nil,
	//	Select:                    "",
	//	TotalSegments:             nil,
	//}
	//s.dynamoClient.Query()
	return nil, nil
}
func (s *Service) UpdateTodo(ctx context.Context, todo *Todo) error {
	//todoItem := serializeTodoDynamo(todo)
	//item, err := s.dynamoClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
	//	Key: map[string]types.AttributeValue{
	//		"ID": &types.AttributeValueMemberS{Value: todo.ID.String()},
	//	},
	//	TableName:                           aws.String(TodoItemsTableName),
	//	ConditionExpression:                 nil,
	//	ConditionalOperator:                 "",
	//	ExpressionAttributeNames:            nil,
	//	ExpressionAttributeValues:           nil,
	//	ReturnConsumedCapacity:              "",
	//	ReturnItemCollectionMetrics:         "",
	//	ReturnValues:                        "",
	//	ReturnValuesOnConditionCheckFailure: "",
	//	UpdateExpression:                    nil,
	//})
	//if err != nil {
	//	return err
	//}
	return nil
}

func (s *Service) DeleteTodo() {

}
