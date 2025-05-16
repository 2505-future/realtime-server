package main

import (
	"context"
	"encoding/json"
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
	messager infrustructure.IMessageSender
}

func NewDisconnectHandler(dynamodb infrustructure.IDynamoDB, messager infrustructure.IMessageSender) *DisconnectHandler {
	return &DisconnectHandler{
		dynamodb: dynamodb,
		messager: messager,
	}
}

func (ch *DisconnectHandler) HandleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("start disconnect")

	userConnectionID := request.RequestContext.ConnectionID

	log.Printf("connectionId : %s ¥n", userConnectionID)

	roomID, userID, err := ch.dynamodb.Get(userConnectionID)
	if err != nil {
		fmt.Println(err)
		return http.Create500response()
	}

	var connectionIDs []string
	err = ch.dynamodb.GetConnectionIDs(roomID, &connectionIDs)
	if err != nil {
		fmt.Println(err)
		return http.Create500response()
	}

	// JSON ペイロード構築
	payload := map[string]interface{}{
		"type": "leave",
		"message": map[string]interface{}{
			"id": userID,
		},
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return http.Create500response()
	}

	err = ch.dynamodb.Delete(userConnectionID)
	if err != nil {
		fmt.Println(err)
		return http.Create500response()
	}

	for _, connectionID := range connectionIDs {
		if connectionID == userConnectionID {
			continue
		}
		if err := ch.messager.SendMessage(ctx, connectionID, jsonBytes); err != nil {
			fmt.Println(err)
			return http.Create500response()
		}
	}

	fmt.Println("end disconnect")
	return http.Create200response()
}

func main() {
	client := db.NewDynamoDBClient()
	dynamodb := infrustructure.NewDynamoDB(client, "websocket")
	messager := infrustructure.NewMessageSender()
	handler := NewDisconnectHandler(dynamodb, messager)
	lambda.Start(handler.HandleRequest)
}
