package infrustructure

import (
	"context"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/apigatewaymanagementapi"
)

type IMessageSender interface {
	SendMessage(ctx context.Context, connectionID string, message []byte) error
}

type MessageSender struct{}

func NewMessageSender() *MessageSender {
	return &MessageSender{}
}

var endpoint = os.Getenv("APIGATEWAY_WEBSOCKET_ENDPOINT")

func (sh *MessageSender) SendMessage(ctx context.Context, connectionID string, message []byte) error {
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := apigatewaymanagementapi.NewFromConfig(cfg, func(o *apigatewaymanagementapi.Options) {
		o.BaseEndpoint = aws.String(endpoint)
	})
	input := &apigatewaymanagementapi.PostToConnectionInput{
		ConnectionId: aws.String(connectionID),
		Data:         message,
	}
	_, err = client.PostToConnection(ctx, input)

	log.Println("send message to connectionID:", connectionID)
	log.Println("message:", message)

	return err
}
