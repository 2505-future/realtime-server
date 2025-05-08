package http

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

func Create500response() (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "Internal Server Error",
	}, nil
}
