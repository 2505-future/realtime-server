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

type ConnectHandler struct {
	dynamodb infrustructure.IDynamoDB
}

func NewConnectHandler(dynamodb infrustructure.IDynamoDB) *ConnectHandler {
	return &ConnectHandler{
		dynamodb: dynamodb,
	}
}

func (ch *ConnectHandler) HandleRequest(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("start connect")

	connectionID := request.RequestContext.ConnectionID

	log.Printf("connectionId : %s ¥n", connectionID)

	err := ch.dynamodb.Put(connectionID)
	if err != nil {
		fmt.Println(err)
		return http.Create500response()
	}

	fmt.Println("end connect")
	return http.Create200response()
}

func main() {
	client := db.NewDynamoDBClient()
	dynamodb := infrustructure.NewDynamoDB(client, "websocket")
	handler := NewConnectHandler(dynamodb)
	lambda.Start(handler.HandleRequest)
}
