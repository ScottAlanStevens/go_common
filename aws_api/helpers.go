package aws_api

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/go-playground/validator"
)

func UnmarshalRequest(body string, result any) (bool, error) {
	reqBytes := []byte(body)

	if len(reqBytes) == 0 {
		return false, nil
	} else {
		err := json.Unmarshal(reqBytes, &result)
		if err != nil {
			return false, err
		}
	}

	return true, nil
}

type Request interface{}

func ValidateRequest[T Request](body string) (*T, *events.APIGatewayProxyResponse) {

	var req T
	_, err := UnmarshalRequest(body, &req)
	if err != nil {
		return nil, BuildErrorResponse(ErrorResponse{
			Message: err.Error(),
		}, 400)
	}
	// if !ok {
	// 	t := new(T)
	// 	req = *t
	// }

	v := validator.New()

	err = v.Struct(&req)
	if err != nil {
		// global.Logger.Debug().Msgf("validation errors: %s", err)
		errors := toConventionalErrors(err.(validator.ValidationErrors), reflect.TypeOf(req))
		errres := BuildErrorResponse(ErrorResponse{
			Message:          ErrorBadRequestMessage,
			ValidationErrors: errors,
		}, 400)
		return nil, errres
	}

	return &req, nil
}

func toConventionalErrors(errs validator.ValidationErrors, t reflect.Type) []string {
	errors := []string{}

	for _, err := range errs {
		// Attempt to find field by name and get json tag name
		field, _ := t.FieldByName(err.StructField())
		var name string

		// If json tag doesn't exist, use lower case of name
		if name = field.Tag.Get("json"); name == "" {
			name = strings.ToLower(err.StructField())
		} else {
			// handles json tags like 'email,omitempty'... we only want the property name.
			name = strings.Split(name, ",")[0]
		}

		switch err.Tag() {
		case "required_if":
			fallthrough
		case "required_unless":
			fallthrough
		case "required":
			errors = append(errors, "'"+name+"' property is required.")
		case "email":
			errors = append(errors, "'"+name+"' should be a valid email.")
		case "min":
			errors = append(errors, "'"+name+"' must be at least "+err.Param()+".")
		case "max":
			errors = append(errors, "'"+name+"' must be at most "+err.Param()+".")
		default:
			errors = append(errors, "'"+name+"' property is invalid.")
		}
	}

	return errors
}

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
