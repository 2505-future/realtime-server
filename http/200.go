package http

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func Create200response() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string("ok"),
	}, nil
}
