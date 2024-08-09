package pkg

type Plant struct {
	Name        string `json:"name" dynamodbav:"name"`
	Description string `json:"description" dynamodbav:"description"`
}
