# plants_v1
Plant info API

# Dynamo DB
To run the project you must first configure the awscli with `aws configure` and a access key id, secret access key, and region. Additonally you will need a Dynamo table named `plants`, and the above access key will need to belong to a user with read and write access to the table.

# Running 
```
cd plants_v1
go run cmd/plants/main.go
```

# Notes
Using aws sdk2 mocking technique from k.goto: https://dev.to/aws-builders/testing-with-aws-sdk-for-go-v2-without-interface-mocks-55de