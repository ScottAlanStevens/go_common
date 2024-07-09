package aws_api

import (
	"encoding/json"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

// func (api *Api) UnmarshalRequest(body string, result any) (bool, error) {
// 	reqBytes := []byte(body)

// 	if len(reqBytes) == 0 {
// 		api.logger.Debug().Msg("empty body")
// 		return false, nil
// 	} else {
// 		err := json.Unmarshal(reqBytes, &result)
// 		if err != nil {
// 			api.logger.Logger.Debug().Msgf("failed to parse json request: %s", err.Error())
// 			return false, err
// 		}
// 	}

// 	return true, nil
// }

// func ValidateRequest[T Request](body string) (*T, *events.APIGatewayProxyResponse) {

// 	var req T
// 	_, err := UnmarshalRequest(body, &req)
// 	if err != nil {
// 		return nil, buildErrorResponse(err, 400)
// 	}
// 	// if !ok {
// 	// 	t := new(T)
// 	// 	req = *t
// 	// }

// 	v := validator.New()

// 	err = v.Struct(&req)
// 	if err != nil {
// 		global.Logger.Debug().Msgf("validation errors: %s", err)
// 		errors := toConventionalErrors(err.(validator.ValidationErrors), reflect.TypeOf(req))
// 		errres := buildValidationErrorResponse(ErrorBadRequestMessage, errors, 400)
// 		return nil, &errres
// 	}

// 	return &req, nil
// }

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

func BuildResponse(httpStatusCode int, responseObject interface{}) events.APIGatewayProxyResponse {
	responseBodyString := ""
	if responseObject != nil {

		responseBody, err := json.Marshal(responseObject)
		if err != nil {
			return *BuildErrorResponse(ErrorResponse{Message: err.Error()}, 500)
		}

		responseBodyString = string(responseBody)
	}

	return events.APIGatewayProxyResponse{
		StatusCode: httpStatusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: responseBodyString,
	}
}

func BuildErrorResponse(errResponse ErrorResponse, httpStatusCode int) *events.APIGatewayProxyResponse {
	responseBytes, _ := json.Marshal(errResponse)

	return &events.APIGatewayProxyResponse{
		StatusCode: httpStatusCode,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
		Body: string(responseBytes),
	}
}

func (api *Api) AddCORSResponse(res *events.APIGatewayProxyResponse) {
	res.Headers["Access-Control-Allow-Origin"] = api.requestOrigin
	res.Headers["Access-Control-Allow-Headers"] = strings.Join(api.options.CORS.AllowedHeaders, ",")
	if api.options.CORS.AllowCredentials {
		res.Headers["Access-Control-Allow-Credentials"] = "true"
	}
}
