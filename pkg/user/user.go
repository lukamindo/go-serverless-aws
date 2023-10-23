package user

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"github.com/lukamindo/go-serverless-aws/pkg/validators"
)

var (
	ErrorBadRequest     = "bad request"
	ErrorInternalDB     = "internal db error"
	ErrorInternalServer = "internal server error"
)

type (
	User struct {
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}

	Users []User
)

func Create(dynaClient dynamodbiface.DynamoDBAPI, tableName string, req events.APIGatewayProxyRequest) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorBadRequest)
	}

	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorBadRequest)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, fmt.Errorf("%s , %s", ErrorInternalServer, err.Error())
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorInternalDB)
	}
	return &u, nil
}

func List(dynaClient dynamodbiface.DynamoDBAPI, tableName string) (*Users, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorInternalDB)
	}

	users := new(Users)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, users)
	if err != nil {
		return nil, errors.New(ErrorInternalServer)
	}
	return users, nil
}

func Get(dynaClient dynamodbiface.DynamoDBAPI, tableName, email string) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	result, err := dynaClient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorInternalDB)
	}

	user := new(User)
	err = dynamodbattribute.UnmarshalMap(result.Item, user)
	if err != nil {
		return nil, fmt.Errorf("%s , %s", ErrorInternalServer, err.Error())
	}
	return user, nil
}

func Update(dynaClient dynamodbiface.DynamoDBAPI, tableName string, req events.APIGatewayProxyRequest) (*User, error) {
	var u User
	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorBadRequest)
	}

	currectUser, err := Get(dynaClient, tableName, u.Email)
	if err != nil {
		return nil, err
	}
	if currectUser != nil && len(currectUser.Email) != 0 {
		return nil, errors.New(ErrorBadRequest)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorInternalServer)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tableName),
	}

	_, err = dynaClient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorInternalDB)
	}
	return &u, nil
}

func Delete(dynaClient dynamodbiface.DynamoDBAPI, tableName string, req events.APIGatewayProxyRequest) error {
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tableName),
	}

	_, err := dynaClient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorInternalDB)
	}
	return nil
}
