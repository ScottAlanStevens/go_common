package aws_api

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/events"
	"github.com/rs/zerolog"
)

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

	if api.options.LogEvents {
		api.logger.Info().
			Interface("event", event).
			Msg("request")
	}

	api.logger.Info().
		Str("path", event.Path).
		Str("resource", event.Resource).
		Str("method", event.HTTPMethod).
		Msg("Routing request...")

	requestOrigin := getRequestOrigin(event.Headers, "origin", "Origin", "Referer", "referer")
	api.requestOrigin = api.getCORSAllowOrigin(requestOrigin)

	if event.HTTPMethod == OPTIONS {
		response := api.HandleOptions(event)
		api.AddCORSResponse(&response)
		if api.options.LogEvents {
			api.logger.Debug().Interface("response", response).Msg("response")
		}
		return response, nil
	}

	handler, err := api.getHandler(event.HTTPMethod, event.Resource)
	if err != nil {
		api.logger.Error().Msgf("Error message: %s, Http status code: %d", err.Error(), 500)
		return *BuildErrorResponse(ErrorResponse{Message: err.Error()}, 500), nil
	}

	response := handler.Handle(event)

	api.AddCORSResponse(&response)

	return response, nil
}

func (api *Api) getHandler(httpMethod string, resource string) (Handler, error) {

	for _, h := range api.options.Handlers {
		if h.HttpMethod == httpMethod && h.Resource == resource {
			return h.HandlerConstructor(api.ctx), nil
		}
	}

	return nil, fmt.Errorf("no handler registered for %s %s", httpMethod, resource)
}
