package pkg

import "github.com/google/uuid"

type Item struct {
	Name        string    `json:"name" dynamodbav:"name"`
	Description string    `json:"description" dynamodbav:"description"`
	Vendor      uuid.UUID `json:"vendor" dynamodbav:"vendor"`
	Price       float32   `json:"price" dynamodbav:"price"`
}
