# plants_v1
Plant info API

# Dynamo DB
To run the project you must first configure the awscli with `aws configure` and a access key id, secret access key, and region. Additonally you will need a Dynamo table named `plants`, and the above access key will need to belong to a user with read and write access to the table.

# Auth0
To run this project you will need auth0 set up for api access. Any request to the running application witll require an auth0 bearer token from the correct domain and audience.

# Running 
Required environment variables: 
```
// AWS - can be set through aws configure
AWS_ACCESS_KEY_ID='{dynamo db access key}'
AWS_SECRET_ACCESS_KEY='{dynamo db secret key}'
AWS_DEFAULT_REGION='{region of dynamo db table}'

// AUTH0
AUTH0_DOMAIN='{your auth0 domain}'
AUTH0_AUDIENCE='{your auth0 api id}'
```
Environment variables can be set manually or through a `.env` file within the root directory

To run:
```
git clone git@github.com:SevvyP/plants_v1.git
cd plants_v1
go mod download
go run cmd/plants/main.go
```

# Notes
Using aws sdk2 mocking technique from k.goto: https://dev.to/aws-builders/testing-with-aws-sdk-for-go-v2-without-interface-mocks-55de