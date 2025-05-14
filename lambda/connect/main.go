package main

import (
	"fmt"
	"log"

	"websocket/http"
	"websocket/infrustructure"

	"github.com/aws/aws-lambda-go/events"
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
	params := request.QueryStringParameters

	requiredKeys := []string{"roomID", "power", "weight", "volume", "cd", "userID", "iconUrl"}
	for _, key := range requiredKeys {
		if params[key] == "" {
			return http.Create400response(fmt.Sprintf("%s is empty", key))
		}
	}

	// 必要であればこれらの値を DynamoDB に保存
	err := ch.dynamodb.Put(connectionID, params["roomID"], params["userID"], params["iconUrl"], params["power"], params["weight"], params["volume"], params["cd"])
	if err != nil {
		fmt.Println(err)
		return http.Create500response()
	}

	fmt.Println("end connect")
	return http.Create200response()
}
