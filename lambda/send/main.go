package main

import (
	"context"
	"log"
	"os"

	"websocket/db"
	"websocket/http"
	"websocket/infrustructure"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

type SendHandler struct {
	dynamodb infrustructure.IDynamoDB
}

func NewSendHandler(dynamodb infrustructure.IDynamoDB) *SendHandler {
	return &SendHandler{
		dynamodb: dynamodb,
	}
}

type Event struct {
	RoomID     string `json:"roomID"`
	JsonString string `json:"json"`
}

func (sh *SendHandler) HandleRequest(ctx context.Context, event Event) (events.APIGatewayProxyResponse, error) {
	log.Println("start send message")

	roomID := event.RoomID
	if roomID == "" {
		log.Println("roomID is empty")
		return http.Create400response("roomID is empty")
	}

	var connectionIDs []string
	err := sh.dynamodb.GetConnectionIDs(roomID, &connectionIDs)
	if err != nil {
		log.Println("failed to get connection IDs:", err)
		return http.Create500response()
	}

	endpoint := os.Getenv("APIGATEWAY_WEBSOCKET_ENDPOINT")
	for _, connectionID := range connectionIDs {
		err = sh.sendMessage(ctx, endpoint, connectionID, event.JsonString)
	}
	if err != nil {
		log.Println("failed to send message:", err)
		return http.Create500response()
	}

	log.Println("end send message")
	return http.Create200response()
}

func (sh *SendHandler) sendMessage(ctx context.Context, endpoint, connectionID, message string) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
	input := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         []byte(message),
	}
	_, err = client.PostToConnection(ctx, input)

	log.Println("send message to connectionID:", connectionID)
	log.Println("message:", message)

	return err
}

func main() {
	client := db.NewDynamoDBClient()
	dynamodb := infrustructure.NewDynamoDB(client, "websocket")
	handler := NewSendHandler(dynamodb)
	lambda.Start(handler.HandleRequest)
}
