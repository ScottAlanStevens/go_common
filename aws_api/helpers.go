package aws_api

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func GetPathParameter(event events.APIGatewayProxyRequest, name string) (bool, string) {
	value := strings.ToLower(event.PathParameters[name])
	return value != "", value
}

func getRequestOrigin(headers map[string]string, names ...string) string {
	for _, name := range names {
		if headers[name] != "" {
			return headers[name]
		}
	}
	return ""
}

func (api *Api) getCORSAllowOrigin(requestOrigin string) string {
	for _, allowedOrigin := range api.options.CORS.AllowedOrigins {
		if strings.EqualFold(requestOrigin, allowedOrigin) {
			return allowedOrigin
		}
	}

	// defaults to one of the allowed origins
	return api.options.CORS.AllowedOrigins[0]
}

func (api *Api) BuildResponse(httpStatusCode int, responseObject interface{}) events.APIGatewayProxyResponse {

	responseBodyString := ""
	if responseObject != nil {

		responseBody, err := json.Marshal(responseObject)
		if err != nil {
			return *api.BuildErrorResponse(err, 500)
		}

		responseBodyString = string(responseBody)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: httpStatusCode,
		Headers: map[string]string{
			"Content-Type":                     "application/json",
			"Access-Control-Allow-Origin":      api.requestOrigin,
			"Access-Control-Allow-Headers":     strings.Join(api.options.CORS.AllowedHeaders, ","),
			"Access-Control-Allow-Credentials": "true",
		},
		Body: responseBodyString,
	}
}

func (api *Api) BuildErrorResponse(err error, httpStatusCode int) *events.APIGatewayProxyResponse {
	api.logger.Error().Msgf("Error message: %s, Http status code: %d", err.Error(), httpStatusCode)

	responseBytes, _ := json.Marshal(ErrorResponse{
		Message: err.Error(),
	})

	return &events.APIGatewayProxyResponse{
		StatusCode: httpStatusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
			// "Access-Control-Allow-Origin":  CorsAllowOrigin,
			// "Access-Control-Allow-Headers": CorsAllowHeaders,
		},
		Body: string(responseBytes),
	}
}
