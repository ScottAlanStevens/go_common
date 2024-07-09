package aws_api

import (
	"github.com/aws/aws-lambda-go/events"
)

func (api *Api) HandleOptions(event events.APIGatewayProxyRequest) events.APIGatewayProxyResponse {
	res := api.BuildResponse(200, nil)
	res.Headers["Access-Control-Allow-Methods"] = "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	return res
}
