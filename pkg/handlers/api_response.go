package handlers

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
)

func apiResponse(status int, body any) (*events.APIGatewayProxyResponse, error) {
	stringBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	resp := events.APIGatewayProxyResponse{
		StatusCode: status,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(stringBody),
	}
	return &resp, nil
}
