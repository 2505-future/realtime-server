package http

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func Create400response(message string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusBadRequest,
		Body:       message,
	}, nil
}
