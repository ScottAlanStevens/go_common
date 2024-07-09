package aws_api

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
)

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
	HandlerConstructor func(api *Api) Handler
}

type CORSOptions struct {
	AllowedOrigins   []string
	AllowedHeaders   []string
	AllowCredentials bool
}

type ApiOptions struct {
	Handlers []HandlerRegistration
	CORS     CORSOptions
}

func NewApi(logger zerolog.Logger, options ApiOptions) IApi {

	if len(options.CORS.AllowedOrigins) == 0 {
		panic("must have at least 1 CORS origin")
	}

	return &Api{
		logger:  logger,
		options: options,
	}
}

func (api *Api) HandleEvent(ctx context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	api.ctx = ctx
	// global.SetRequestId(&event.RequestContext.RequestID)

	api.logger.Info().
		Str("path", event.Path).
		Str("resource", event.Resource).
		Str("method", event.HTTPMethod).
		Msg("Routing request...")

	if event.HTTPMethod == OPTIONS {
		return api.HandleOptions(event), nil
	}

	handler, err := api.getHandler(event.HTTPMethod, event.Resource)
	if err != nil {
		api.logger.Error().Msgf("Error message: %s, Http status code: %d", err.Error(), 500)
		return *BuildErrorResponse(ErrorResponse{Message: err.Error()}, 500), nil
	}

	requestOrigin := getRequestOrigin(event.Headers, "origin", "Origin", "Referer", "referer")
	api.requestOrigin = api.getCORSAllowOrigin(requestOrigin)

	response := handler.Handle(event)

	api.AddCORSResponse(&response)

	return response, nil
}

func (api *Api) getHandler(httpMethod string, resource string) (Handler, error) {

	for _, h := range api.options.Handlers {
		if h.HttpMethod == httpMethod && h.Resource == resource {
			return h.HandlerConstructor(api), nil
		}
	}

	return nil, fmt.Errorf("no handler registered for %s %s", httpMethod, resource)
}
