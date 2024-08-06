FROM golang:latest

WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/items cmd/items/main.go
EXPOSE 8080
CMD ["./bin/items"]