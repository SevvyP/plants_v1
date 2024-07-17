package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/SevvyP/plants/pkg"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

var ErrNotFound = "item not found"

type DBInterface interface {
	CreatePlant(pkg.Plant, context.Context) error
	GetPlant(string, context.Context) (*pkg.Plant, error)
	UpdatePlant(pkg.Plant, context.Context) error
	DeletePlant(string, context.Context) (*pkg.Plant, error)
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

func (db *DB) CreatePlant(plant pkg.Plant, context context.Context) error {
	if plant.Name == "" || plant.Description == "" {
		return errors.New("missing name or description")
	}
	item, err := attributevalue.MarshalMap(plant)
	if err != nil {
		return err
	}
	output, err := db.client.PutItem(context, &dynamodb.PutItemInput{
		TableName: aws.String("plants_v1"), Item: item,
	})
	if err != nil {
		return err
	}
	err = attributevalue.UnmarshalMap(output.Attributes, &plant)
	if err != nil {
		return err
	}
	if plant.Name == "" {
		return errors.New("unmarshal failed")
	}
	return err
}

func (db *DB) GetPlant(name string, context context.Context) (*pkg.Plant, error){
	if name == "" {
		return nil, errors.New("missing name or description")
	}
	nameattribute, err := attributevalue.Marshal(name)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.GetItemInput{Key: map[string]types.AttributeValue{"name": nameattribute}, TableName: aws.String("plants_v1")}
	output, err := db.client.GetItem(context, input)
	if err !=nil {
		return nil, err
	}
	var plant *pkg.Plant
	err = attributevalue.UnmarshalMap(output.Item, &plant)
	fmt.Println(output.Item)
	if plant.Name == "" {
		return nil, errors.New(ErrNotFound)
	}
	return plant, nil
}

// UpdatePlant will create the item if it does not exist in the db
// TODO: fix this
func (db *DB) UpdatePlant(plant pkg.Plant, context context.Context) error {
	if plant.Name == "" || plant.Description == "" {
		return errors.New("missing name or description")
	}
	nameattribute, err := attributevalue.Marshal(plant.Name)
	if err != nil {
		return err
	}
	_, err = db.client.UpdateItem(context, &dynamodb.UpdateItemInput{
		TableName: aws.String("plants_v1"), Key: map[string]types.AttributeValue{"name": nameattribute}, UpdateExpression: aws.String("set description = :description"), ExpressionAttributeValues: map[string]types.AttributeValue{
            ":description": &types.AttributeValueMemberS{Value: plant.Description}, 
        },
	})
	if err != nil {
		return err
	}
	return err
}

func (db *DB) DeletePlant(name string, context context.Context) (*pkg.Plant, error) {
	if name == "" {
		return nil, errors.New("missing name or description")
	}
	nameattribute, err := attributevalue.Marshal(name)
	if err != nil {
		return nil, err
	}
	output, err := db.client.DeleteItem(context, &dynamodb.DeleteItemInput{
		TableName: aws.String("plants_v1"), Key: map[string]types.AttributeValue{"name": nameattribute}, ReturnValues: types.ReturnValueAllOld,
	})
	if err != nil {
		return nil, err
	}
	plant := &pkg.Plant{}
	err = attributevalue.UnmarshalMap(output.Attributes, &plant)
	if err != nil {
		return nil, err
	}
	if plant.Name == "" {
		return nil, errors.New(ErrNotFound)
	}
	return plant, nil
}