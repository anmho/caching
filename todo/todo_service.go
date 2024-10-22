package todo

import (
	"context"
	"github.com/anmho/caching/cache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"log/slog"
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

func (s *Service) FindTodoByID(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*Todo, error) {
	params := &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"UserID": &types.AttributeValueMemberS{Value: userID.String()},
			"ID":     &types.AttributeValueMemberS{Value: id.String()},
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

func (s *Service) ListUserTodos(ctx context.Context, userID uuid.UUID) ([]*Todo, error) {
	// add pagination with pagination token?
	input := &dynamodb.QueryInput{
		TableName:              aws.String(TodoItemsTableName),
		ConsistentRead:         aws.Bool(true),
		KeyConditionExpression: aws.String("UserID = :userID"),
		ReturnConsumedCapacity: types.ReturnConsumedCapacityTotal,
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":userID": &types.AttributeValueMemberS{
				Value: userID.String(),
			},
		},
	}
	output, err := s.dynamoClient.Query(ctx, input)
	if err != nil {
		return nil, err
	}
	todos := make([]*Todo, 0)
	for _, item := range output.Items {
		todo, err := deserializeTodoDynamo(item)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}

	return todos, nil
}

func formatDate(t *time.Time) string {
	return t.Format(time.RFC3339)
}

func parseDate(date string) (time.Time, error) {
	return time.Parse(time.RFC3339, date)
}

type UpdateParams struct {
	Completed   bool   `json:"completed" validate:"required"`
	Title       string `json:"title" validate:"required"`
	Description string `json:"description" validate:"required"`
}

func (s *Service) UpdateTodo(
	ctx context.Context,
	userID uuid.UUID,
	id uuid.UUID,
	params *UpdateParams) error {

	expressionAttributeValues := map[string]types.AttributeValue{
		":description": &types.AttributeValueMemberS{
			Value: params.Description,
		},
		":title": &types.AttributeValueMemberS{
			Value: params.Title,
		},
		":updatedAt": &types.AttributeValueMemberS{
			Value: formatDate(aws.Time(time.Now().UTC())),
		},
	}

	if params.Completed {
		expressionAttributeValues[":completedAt"] = &types.AttributeValueMemberS{
			Value: formatDate(aws.Time(time.Now().UTC())),
		}
	} else {
		expressionAttributeValues[":completedAt"] = nil
	}

	slog.Info("UpdateTodo", slog.Any("params", params), slog.Any(":completedAt", expressionAttributeValues[":completedAt"]))

	input := &dynamodb.UpdateItemInput{
		Key: map[string]types.AttributeValue{
			"UserID": &types.AttributeValueMemberS{
				Value: userID.String(),
			},
			"ID": &types.AttributeValueMemberS{
				Value: id.String(),
			},
		},
		TableName: aws.String(TodoItemsTableName),
		UpdateExpression: aws.String(`
			SET 
				Title = :title, 
				Description = :description,
				UpdatedAt = :updatedAt,
				CompletedAt = :completedAt
		`),
		ExpressionAttributeValues: expressionAttributeValues,
		ReturnConsumedCapacity:    types.ReturnConsumedCapacityTotal,
		// does not return values to save rcu
	}

	_, err := s.dynamoClient.UpdateItem(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) DeleteTodo(userID uuid.UUID, todoID uuid.UUID) {

}
