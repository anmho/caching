package todo

import (
	"context"
	"github.com/anmho/caching/cache"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"log/slog"
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

func (s *Service) GetAllTodosForUser(ctx context.Context, userID uuid.UUID) ([]*Todo, error) {
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

func (s *Service) DeleteTodo(userID uuid.UUID, todoID uuid.UUID) {

}
