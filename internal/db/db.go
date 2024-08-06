package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/SevvyP/items/pkg"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var ErrNotFound = "item not found"

type DBInterface interface {
	CreateItem(pkg.Item, context.Context) error
	GetItem(string, context.Context) (*pkg.Item, error)
	UpdateItem(pkg.Item, context.Context) error
	DeleteItem(string, context.Context) (*pkg.Item, error)
}

type DB struct {
	client *dynamodb.Client
}

func NewDB() *DB {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	client := dynamodb.NewFromConfig(cfg)
	return &DB{client: client}
}

func (db *DB) CreateItem(item pkg.Item, context context.Context) error {
	if item.Name == "" || item.Description == "" {
		return errors.New("missing name or description")
	}
	newItem, err := attributevalue.MarshalMap(item)
	if err != nil {
		return err
	}
	output, err := db.client.PutItem(context, &dynamodb.PutItemInput{
		TableName: aws.String("items_v1"), Item: newItem,
	})
	if err != nil {
		return err
	}
	err = attributevalue.UnmarshalMap(output.Attributes, &item)
	if err != nil {
		return err
	}
	if item.Name == "" {
		return errors.New("unmarshal failed")
	}
	return err
}

func (db *DB) GetItem(name string, context context.Context) (*pkg.Item, error) {
	if name == "" {
		return nil, errors.New("missing name or description")
	}
	nameattribute, err := attributevalue.Marshal(name)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.GetItemInput{Key: map[string]types.AttributeValue{"name": nameattribute}, TableName: aws.String("items_v1")}
	output, err := db.client.GetItem(context, input)
	if err != nil {
		return nil, err
	}
	var item *pkg.Item
	err = attributevalue.UnmarshalMap(output.Item, &item)
	fmt.Println(output.Item)
	if item.Name == "" {
		return nil, errors.New(ErrNotFound)
	}
	return item, nil
}

// UpdateItem will create the item if it does not exist in the db
// TODO: fix this
func (db *DB) UpdateItem(item pkg.Item, context context.Context) error {
	if item.Name == "" || item.Description == "" {
		return errors.New("missing name or description")
	}
	nameattribute, err := attributevalue.Marshal(item.Name)
	if err != nil {
		return err
	}
	_, err = db.client.UpdateItem(context, &dynamodb.UpdateItemInput{
		TableName: aws.String("items_v1"), Key: map[string]types.AttributeValue{"name": nameattribute}, UpdateExpression: aws.String("set description = :description"), ExpressionAttributeValues: map[string]types.AttributeValue{
			":description": &types.AttributeValueMemberS{Value: item.Description},
		},
	})
	if err != nil {
		return err
	}
	return err
}

func (db *DB) DeleteItem(name string, context context.Context) (*pkg.Item, error) {
	if name == "" {
		return nil, errors.New("missing name or description")
	}
	nameattribute, err := attributevalue.Marshal(name)
	if err != nil {
		return nil, err
	}
	output, err := db.client.DeleteItem(context, &dynamodb.DeleteItemInput{
		TableName: aws.String("items_v1"), Key: map[string]types.AttributeValue{"name": nameattribute}, ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		return nil, err
	}
	item := &pkg.Item{}
	err = attributevalue.UnmarshalMap(output.Attributes, &item)
	if err != nil {
		return nil, err
	}
	if item.Name == "" {
		return nil, errors.New(ErrNotFound)
	}
	return item, nil
}
