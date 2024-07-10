package aws_api

import (
	"context"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
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

type HandlerConstructor func(ctx context.Context) Handler

type IApi interface {
	HandleEvent(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)
}

type Api struct {
	logger  zerolog.Logger
	options ApiOptions

	ctx           context.Context
	requestOrigin string
}

type HandlerRegistration struct {
	HttpMethod         string
	Resource           string
	HandlerConstructor HandlerConstructor
}

type CORSOptions struct {
	AllowedOrigins   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

type ApiOptions struct {
	Handlers  []HandlerRegistration
	CORS      CORSOptions
	LogEvents bool
}

type ErrorResponse struct {
	Message          string   `json:"message"`
	ValidationErrors []string `json:"validation_errors"`
}
