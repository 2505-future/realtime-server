package main

import (
	"fmt"
	"log"

	"websocket/db"
	"websocket/http"
	"websocket/infrustructure"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type DisconnectHandler struct {
	dynamodb infrustructure.IDynamoDB
}

func NewDisconnectHandler(dynamodb infrustructure.IDynamoDB) *DisconnectHandler {
	return &DisconnectHandler{
		dynamodb: dynamodb,
	}
}

func (ch *DisconnectHandler) HandleRequest(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("start disconnect")

	connectionID := request.RequestContext.ConnectionID

	log.Printf("connectionId : %s ¥n", connectionID)

	err := ch.dynamodb.Delete(connectionID)
	if err != nil {
		fmt.Println(err)
		return http.Create500response()
	}

	fmt.Println("end disconnect")
	return http.Create200response()
}

func main() {
	client := db.NewDynamoDBClient()
	dynamodb := infrustructure.NewDynamoDB(client, "websocket")
	handler := NewDisconnectHandler(dynamodb)
	lambda.Start(handler.HandleRequest)
}
