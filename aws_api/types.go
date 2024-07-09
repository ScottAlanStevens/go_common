package aws_api

import (
	"errors"

	"github.com/aws/aws-lambda-go/events"
)

const (
	GET     string = "GET"
	POST    string = "POST"
	PATCH   string = "PATCH"
	DELETE  string = "DELETE"
	OPTIONS string = "OPTIONS"
)

var (
	ErrorUnauthorizedMessage   = "UNAUTHORIZED"
	ErrorInternalServerMessage = "INTERNAL_SERVER_ERROR"
	ErrorNotFoundMessage       = "NOT_FOUND"
	ErrorBadRequestMessage     = "BAD_REQUEST"

	ErrorUnauthorized   = errors.New(ErrorUnauthorizedMessage)
	ErrorInternalServer = errors.New(ErrorInternalServerMessage)
	ErrorNotFound       = errors.New(ErrorNotFoundMessage)
	ErrorBadRequest     = errors.New(ErrorBadRequestMessage)

	ErrorUnauthorizedResponse = ErrorResponse{
		Message: ErrorUnauthorizedMessage,
	}
	ErrorInternalServerResponse = ErrorResponse{
		Message: ErrorInternalServerMessage,
	}
	ErrorNotFoundResponse = ErrorResponse{
		Message: ErrorNotFoundMessage,
	}
	ErrorBadRequestResponse = ErrorResponse{
		Message: ErrorBadRequestMessage,
	}
)

type Handler interface {
	Handle(event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse
}

type ErrorResponse struct {
	Message          string   `json:"message"`
	ValidationErrors []string `json:"validation_errors"`
}
