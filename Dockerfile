FROM golang:latest

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/plants cmd/plants/main.go
EXPOSE 8080
CMD ["./bin/plants"]