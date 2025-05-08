package db

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

func NewDynamoDBClient() *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("ap-northeast-1"))
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	client := dynamodb.NewFromConfig(cfg)

	return client
}
