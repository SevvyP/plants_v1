package db

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/SevvyP/items/pkg"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go/middleware"
)

type testStruct struct {
	Test string `dynamodbav:"test"`
}

func TestDB_CreateItem(t *testing.T) {
	type args struct {
		item               pkg.Item
		context            context.Context
		withAPIOptionsFunc func(*middleware.Stack) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errText string
	}{
		{
			name: "create item returns error if no name is provided",
			args: args{
				item: pkg.Item{
					Name:        "",
					Description: "",
				},
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"PutItemMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &dynamodb.PutItemOutput{},
								}, middleware.Metadata{}, nil
							},
						),
						middleware.Before,
					)
				},
			},
			wantErr: true,
			errText: "missing name or description",
		},
		{
			name: "create item returns error if client returns an error",
			args: args{
				item:    pkg.Item{Name: "test", Description: "test"},
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"PutItemMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								attributes, err := attributevalue.MarshalMap(pkg.Item{Name: "test", Description: "test"})
								return middleware.FinalizeOutput{
									Result: &dynamodb.PutItemOutput{Attributes: attributes},
								}, middleware.Metadata{}, err
							},
						),
						middleware.Before,
					)
				},
			},
			wantErr: false,
		},
		{
			name: "create item doesn't return error if client is successful",
			args: args{
				item:    pkg.Item{Name: "test", Description: "test"},
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"PutItemMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: nil,
								}, middleware.Metadata{}, fmt.Errorf("PutItemError")
							},
						),
						middleware.Before,
					)
				},
			},
			wantErr: true,
			errText: "operation error DynamoDB: PutItem, PutItemError",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"), config.WithAPIOptions([]func(*middleware.Stack) error{tt.args.withAPIOptionsFunc}))
			if err != nil {
				t.Fatal(err)
			}
			client := dynamodb.NewFromConfig(cfg)
			db := &DB{client: client}
			err = db.CreateItem(tt.args.item, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.CreateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errText != err.Error() {
				t.Errorf("DB.CreateItem() error = %v, errText = %s", err, tt.errText)
			}
		})
	}
}

func TestDB_GetItem(t *testing.T) {
	type fields struct {
		client *dynamodb.Client
	}
	type args struct {
		name               string
		context            context.Context
		withAPIOptionsFunc func(*middleware.Stack) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pkg.Item
		wantErr bool
		errText string
	}{
		{
			name: "get item returns error if no name is provided",
			args: args{
				name:    "",
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"GetItemMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &dynamodb.GetItemOutput{},
								}, middleware.Metadata{}, nil
							},
						),
						middleware.Before,
					)
				},
			},
			want:    &pkg.Item{},
			wantErr: true,
			errText: "missing name or description",
		},
		{
			name: "get item returns error if client returns error",
			args: args{
				name:    "test",
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"GetItemMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: nil,
								}, middleware.Metadata{}, fmt.Errorf("GetItemError")
							},
						),
						middleware.Before,
					)
				},
			},
			want:    &pkg.Item{},
			wantErr: true,
			errText: "operation error DynamoDB: GetItem, GetItemError",
		},
		{
			name: "get item returns error if returned object is not a item",
			args: args{
				name:    "test",
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"GetItemMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								attributes, err := attributevalue.MarshalMap(testStruct{Test: "test"})
								return middleware.FinalizeOutput{
									Result: &dynamodb.GetItemOutput{Item: attributes},
								}, middleware.Metadata{}, err
							},
						),
						middleware.Before,
					)
				},
			},
			want:    &pkg.Item{},
			wantErr: true,
			errText: "item not found",
		},
		{
			name: "get item doesn't return error if client is successful",
			args: args{
				name:    "test",
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"GetItemMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								attributes, err := attributevalue.MarshalMap(pkg.Item{Name: "test", Description: "test"})
								return middleware.FinalizeOutput{
									Result: &dynamodb.GetItemOutput{Item: attributes},
								}, middleware.Metadata{}, err
							},
						),
						middleware.Before,
					)
				},
			},
			want:    &pkg.Item{Name: "test", Description: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"), config.WithAPIOptions([]func(*middleware.Stack) error{tt.args.withAPIOptionsFunc}))
			if err != nil {
				t.Fatal(err)
			}
			client := dynamodb.NewFromConfig(cfg)
			db := &DB{
				client: client,
			}
			got, err := db.GetItem(tt.args.name, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.GetItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errText != err.Error() {
				t.Errorf("DB.GetItem() error = %v, errText = %s", err, tt.errText)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.GetItem() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_UpdateItem(t *testing.T) {
	type fields struct {
		client *dynamodb.Client
	}
	type args struct {
		item               pkg.Item
		context            context.Context
		withAPIOptionsFunc func(*middleware.Stack) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
		errText string
	}{
		{
			name: "update item returns error if no name is provided",
			args: args{
				item: pkg.Item{
					Name:        "",
					Description: "",
				},
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"MockUpdateItem",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &dynamodb.UpdateItemOutput{},
								}, middleware.Metadata{}, nil
							},
						),
						middleware.Before,
					)
				},
			},
			wantErr: true,
			errText: "missing name or description",
		},
		{
			name: "update item returns error if client returns an error",
			args: args{
				item: pkg.Item{
					Name:        "test",
					Description: "test",
				},
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"MockUpdateItem",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: nil,
								}, middleware.Metadata{}, fmt.Errorf("UpdateItemError")
							},
						),
						middleware.Before,
					)
				},
			},
			wantErr: true,
			errText: "operation error DynamoDB: UpdateItem, UpdateItemError",
		},
		{
			name: "update item doesn't return an error if client is successful",
			args: args{
				item: pkg.Item{
					Name:        "test",
					Description: "test",
				},
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"MockUpdateItem",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								attributes, err := attributevalue.MarshalMap(pkg.Item{Name: "test", Description: "test"})
								return middleware.FinalizeOutput{
									Result: &dynamodb.UpdateItemOutput{Attributes: attributes},
								}, middleware.Metadata{}, err
							},
						),
						middleware.Before,
					)
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"), config.WithAPIOptions([]func(*middleware.Stack) error{tt.args.withAPIOptionsFunc}))
			if err != nil {
				t.Fatal(err)
			}
			client := dynamodb.NewFromConfig(cfg)
			db := &DB{client: client}
			err = db.UpdateItem(tt.args.item, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.UpdateItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errText != err.Error() {
				t.Errorf("DB.UpdateItem() error = %v, errText = %s", err, tt.errText)
			}
		})
	}
}

func TestDB_DeleteItem(t *testing.T) {
	type fields struct {
		client *dynamodb.Client
	}
	type args struct {
		name               string
		context            context.Context
		withAPIOptionsFunc func(*middleware.Stack) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pkg.Item
		wantErr bool
		errText string
	}{
		{
			name: "delete item returns error if no name is provided",
			args: args{
				name:    "",
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"DeleteItemMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: &dynamodb.DeleteItemOutput{},
								}, middleware.Metadata{}, nil
							},
						),
						middleware.Before,
					)
				},
			},
			wantErr: true,
			errText: "missing name or description",
		},
		{
			name: "delete item returns error if client returns error",
			args: args{
				name:    "test",
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"DeleteItemMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								return middleware.FinalizeOutput{
									Result: nil,
								}, middleware.Metadata{}, fmt.Errorf("DeleteItemError")
							},
						),
						middleware.Before,
					)
				},
			},
			want:    &pkg.Item{},
			wantErr: true,
			errText: "operation error DynamoDB: DeleteItem, DeleteItemError",
		},
		{
			name: "delete item returns error if returned object is not a item",
			args: args{
				name:    "test",
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"DeleteItemMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								attributes, err := attributevalue.MarshalMap(testStruct{Test: "test"})
								return middleware.FinalizeOutput{
									Result: &dynamodb.DeleteItemOutput{Attributes: attributes},
								}, middleware.Metadata{}, err
							},
						),
						middleware.Before,
					)
				},
			},
			want:    &pkg.Item{},
			wantErr: true,
			errText: "item not found",
		},
		{
			name: "delete item doesn't return error if client is successful",
			args: args{
				name:    "test",
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
					return stack.Finalize.Add(
						middleware.FinalizeMiddlewareFunc(
							"DeleteItemMock",
							func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								attributes, err := attributevalue.MarshalMap(pkg.Item{Name: "test", Description: "test"})
								return middleware.FinalizeOutput{
									Result: &dynamodb.DeleteItemOutput{Attributes: attributes},
								}, middleware.Metadata{}, err
							},
						),
						middleware.Before,
					)
				},
			},
			want:    &pkg.Item{Name: "test", Description: "test"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion("us-east-1"), config.WithAPIOptions([]func(*middleware.Stack) error{tt.args.withAPIOptionsFunc}))
			if err != nil {
				t.Fatal(err)
			}
			client := dynamodb.NewFromConfig(cfg)
			db := &DB{
				client: client,
			}
			got, err := db.DeleteItem(tt.args.name, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.DeleteItem() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errText != err.Error() {
				t.Errorf("DB.DeleteItem() error = %v, errText = %s", err, tt.errText)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.DeleteItem() = %v, want %v", got, tt.want)
			}
		})
	}
}
