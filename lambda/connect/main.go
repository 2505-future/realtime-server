package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"websocket/db"
	"websocket/http"
	"websocket/infrustructure"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type ConnectHandler struct {
	dynamodb  infrustructure.IDynamoDB
	messsager infrustructure.IMessageSender
}

func NewConnectHandler(dynamodb infrustructure.IDynamoDB, messager infrustructure.IMessageSender) *ConnectHandler {
	return &ConnectHandler{
		dynamodb:  dynamodb,
		messsager: messager,
	}
}

func (ch *ConnectHandler) HandleRequest(ctx context.Context, request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("start connect")

	connectionID := request.RequestContext.ConnectionID
	params := request.QueryStringParameters

	requiredKeys := []string{"roomID", "power", "weight", "volume", "cd", "userID", "iconUrl"}
	for _, key := range requiredKeys {
		if params[key] == "" {
			return http.Create400response(fmt.Sprintf("%s is empty", key))
		}
	}

	// 文字列 -> 数値へ変換
	power, err := strconv.Atoi(params["power"])
	if err != nil {
		return http.Create400response("invalid power")
	}
	weight, err := strconv.Atoi(params["weight"])
	if err != nil {
		return http.Create400response("invalid weight")
	}
	volume, err := strconv.Atoi(params["volume"])
	if err != nil {
		return http.Create400response("invalid volume")
	}
	cd, err := strconv.Atoi(params["cd"])
	if err != nil {
		return http.Create400response("invalid cd")
	}

	// JSON ペイロード構築
	payload := map[string]interface{}{
		"type": "join",
		"message": map[string]interface{}{
			"id":       params["userID"],
			"icon_url": params["iconUrl"],
			"power":    power,
			"weight":   weight,
			"cd":       cd,
			"volume":   volume,
		},
	}
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return http.Create500response()
	}

	var connectionIDs []string
	if err := ch.dynamodb.GetConnectionIDs(params["roomID"], &connectionIDs); err != nil {
		return http.Create500response()
	}

	for _, connectionID := range connectionIDs {
		if err := ch.messsager.SendMessage(ctx, connectionID, jsonBytes); err != nil {
			fmt.Println(err)
			return http.Create500response()
		}
	}

	err = ch.dynamodb.Put(connectionID, params["roomID"], params["userID"], params["iconUrl"], params["power"], params["weight"], params["volume"], params["cd"])
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
	messager := infrustructure.NewMessageSender()
	handler := NewConnectHandler(dynamodb, messager)
	lambda.Start(handler.HandleRequest)
}
