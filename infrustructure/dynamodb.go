package infrustructure

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type IDynamoDB interface {
	Put(connectionID, roomID, userID, iconUrl, power, weight, volume, cd string) error
	Delete(connectionID string) error
	GetConnectionIDs(roomID string, connectionIDs *[]string) error
	Get(connectionId string) (string, string, error)
}

type DynamoDB struct {
	client    *dynamodb.Client
	tableName string
}

func NewDynamoDB(client *dynamodb.Client, tableName string) *DynamoDB {
	return &DynamoDB{
		client:    client,
		tableName: tableName,
	}
}

func (d *DynamoDB) Put(connectionID, roomID, userID, iconUrl, power, weight, volume, cd string) error {
	_, err := d.client.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item: map[string]types.AttributeValue{
			"connectionId": &types.AttributeValueMemberS{Value: connectionID},
			"roomId":       &types.AttributeValueMemberS{Value: roomID},
			"userId":       &types.AttributeValueMemberS{Value: userID},
			"iconUrl":      &types.AttributeValueMemberS{Value: iconUrl},
			"power":        &types.AttributeValueMemberN{Value: power},
			"weight":       &types.AttributeValueMemberN{Value: weight},
			"volume":       &types.AttributeValueMemberN{Value: volume},
			"cd":           &types.AttributeValueMemberN{Value: cd},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to put item: %w", err)
	}
	return nil
}

func (d *DynamoDB) Delete(connectionID string) error {
	_, err := d.client.DeleteItem(context.TODO(), &dynamodb.DeleteItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"connectionId": &types.AttributeValueMemberS{Value: connectionID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}
	return nil
}

func (d *DynamoDB) GetConnectionIDs(roomID string, connectionIDs *[]string) error {
	output, err := d.client.Query(context.TODO(), &dynamodb.QueryInput{
		TableName:              aws.String(d.tableName),
		IndexName:              aws.String("roomId-index"),
		KeyConditionExpression: aws.String("roomId = :roomId"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":roomId": &types.AttributeValueMemberS{Value: roomID},
		},
	})
	if err != nil {
		return fmt.Errorf("failed to query by roomId: %w", err)
	}

	for _, item := range output.Items {
		if val, ok := item["connectionId"].(*types.AttributeValueMemberS); ok {
			*connectionIDs = append(*connectionIDs, val.Value)
		}
	}

	return nil
}

func (d *DynamoDB) Get(connectionId string) (string, string, error) {
	output, err := d.client.GetItem(context.TODO(), &dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]types.AttributeValue{
			"connectionId": &types.AttributeValueMemberS{Value: connectionId},
		},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to get item: %w", err)
	}

	if output.Item != nil {
		if val, ok := output.Item["roomId"].(*types.AttributeValueMemberS); ok {
			roomId := val.Value
			if val, ok := output.Item["userId"].(*types.AttributeValueMemberS); ok {
				userId := val.Value
				return roomId, userId, nil
			}
		}
	}

	return "", "", fmt.Errorf("roomId not found")
}
