package db

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/SevvyP/plants/pkg"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/smithy-go/middleware"
)

type testStruct struct {
	Test string `dynamodbav:"test"`
}

func TestDB_CreatePlant(t *testing.T) {
	type args struct {
		plant   pkg.Plant
		context context.Context
		withAPIOptionsFunc func(*middleware.Stack) error
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		errText string
	}{
		{
			name: "create plant returns error if no name is provided",
			args: args{
				plant: pkg.Plant{
					Name: "",
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
			name: "create plant returns error if client returns an error",
            args: args{
				plant: pkg.Plant{Name: "test", Description: "test"},
                context:  context.TODO(),
                withAPIOptionsFunc: func(stack *middleware.Stack) error {
                    return stack.Finalize.Add(
                        middleware.FinalizeMiddlewareFunc(
                            "PutItemMock",
                            func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								attributes, err := attributevalue.MarshalMap(pkg.Plant{Name: "test", Description: "test"})
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
			name: "create plant doesn't return error if client is successful",
            args: args{
				plant: pkg.Plant{Name: "test", Description: "test"},
                context:  context.TODO(),
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
			err = db.CreatePlant(tt.args.plant, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.CreatePlant() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errText != err.Error() {
				t.Errorf("DB.CreatePlant() error = %v, errText = %s", err, tt.errText)
			}
		})
	}
}

func TestDB_GetPlant(t *testing.T) {
	type fields struct {
		client *dynamodb.Client
	}
	type args struct {
		name    string
		context context.Context
		withAPIOptionsFunc func(*middleware.Stack) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pkg.Plant
		wantErr bool
		errText string
	}{
		{
			name: "get plant returns error if no name is provided",
			args: args{
				name: "",
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
			want: &pkg.Plant{},
			wantErr: true,
			errText: "missing name or description",
		},
		{
			name: "get plant returns error if client returns error",
			args: args{
				name: "test",
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
			want: &pkg.Plant{},
			wantErr: true,
			errText: "operation error DynamoDB: GetItem, GetItemError",
		},
		{
			name: "get plant returns error if returned object is not a plant",
			args: args{
				name: "test",
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
			want: &pkg.Plant{},
			wantErr: true,
			errText: "item not found",
		},
		{
			name: "get plant doesn't return error if client is successful",
			args: args{
				name: "test",
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
                    return stack.Finalize.Add(
                        middleware.FinalizeMiddlewareFunc(
                            "GetItemMock",
                            func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								attributes, err := attributevalue.MarshalMap(pkg.Plant{Name: "test", Description: "test"})
                                return middleware.FinalizeOutput{
                                    Result: &dynamodb.GetItemOutput{Item: attributes},
                                }, middleware.Metadata{}, err
                            },
                        ),
                        middleware.Before,
                    )
                },
				
			},
			want: &pkg.Plant{Name: "test", Description: "test"},
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
			got, err := db.GetPlant(tt.args.name, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.GetPlant() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errText != err.Error() {
				t.Errorf("DB.GetPlant() error = %v, errText = %s", err, tt.errText)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.GetPlant() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDB_UpdatePlant(t *testing.T) {
	type fields struct {
		client *dynamodb.Client
	}
	type args struct {
		plant   pkg.Plant
		context context.Context
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
			name: "update plant returns error if no name is provided",
			args: args{
				plant: pkg.Plant{
					Name: "",
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
			name: "update plant returns error if client returns an error",
			args: args{
				plant: pkg.Plant{
					Name: "test",
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
			name: "update plant doesn't return an error if client is successful",
			args: args{
				plant: pkg.Plant{
					Name: "test",
					Description: "test",
				},
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
                    return stack.Finalize.Add(
                        middleware.FinalizeMiddlewareFunc(
                            "MockUpdateItem",
                            func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
                                attributes, err := attributevalue.MarshalMap(pkg.Plant{Name: "test", Description: "test"})
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
			err = db.UpdatePlant(tt.args.plant, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.UpdatePlant() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errText != err.Error() {
				t.Errorf("DB.UpdatePlant() error = %v, errText = %s", err, tt.errText)
			}
		})
	}
}

func TestDB_DeletePlant(t *testing.T) {
	type fields struct {
		client *dynamodb.Client
	}
	type args struct {
		name    string
		context context.Context
		withAPIOptionsFunc func(*middleware.Stack) error
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *pkg.Plant
		wantErr bool
		errText string
	}{
		{
			name: "delete plant returns error if no name is provided",
			args: args{
				name: "",
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
			name: "delete plant returns error if client returns error",
			args: args{
				name: "test",
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
			want: &pkg.Plant{},
			wantErr: true,
			errText: "operation error DynamoDB: DeleteItem, DeleteItemError",
		},
		{
			name: "delete plant returns error if returned object is not a plant",
			args: args{
				name: "test",
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
			want: &pkg.Plant{},
			wantErr: true,
			errText: "item not found",
		},
		{
			name: "delete plant doesn't return error if client is successful",
			args: args{
				name: "test",
				context: context.TODO(),
				withAPIOptionsFunc: func(stack *middleware.Stack) error {
                    return stack.Finalize.Add(
                        middleware.FinalizeMiddlewareFunc(
                            "DeleteItemMock",
                            func(context.Context, middleware.FinalizeInput, middleware.FinalizeHandler) (middleware.FinalizeOutput, middleware.Metadata, error) {
								attributes, err := attributevalue.MarshalMap(pkg.Plant{Name: "test", Description: "test"})
                                return middleware.FinalizeOutput{
                                    Result: &dynamodb.DeleteItemOutput{Attributes: attributes},
                                }, middleware.Metadata{}, err
                            },
                        ),
                        middleware.Before,
                    )
                },
				
			},
			want: &pkg.Plant{Name: "test", Description: "test"},
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
			got, err := db.DeletePlant(tt.args.name, tt.args.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.DeletePlant() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errText != err.Error() {
				t.Errorf("DB.DeletePlant() error = %v, errText = %s", err, tt.errText)
			}
			if err == nil && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DB.DeletePlant() = %v, want %v", got, tt.want)
			}
		})
	}
}
