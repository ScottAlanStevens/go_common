package aws_api

import "github.com/aws/aws-lambda-go/events"

const (
	GET     string = "GET"
	POST    string = "POST"
	PATCH   string = "PATCH"
	DELETE  string = "DELETE"
	OPTIONS string = "OPTIONS"
)

type Handler interface {
	handle(event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
}

type ErrorResponse struct {
	Message          string   `json:"message"`
	ValidationErrors []string `json:"validation_errors"`
}
