# Go parameters
    BINARY_NAME=petstore

# Main package
    MAIN=./cmd

# Commands
    all:
		go mod vendor
		go test -v ./...
		go build -o $(BINARY_NAME) -mod=vendor $(MAIN)
    build:
		go build -o $(BINARY_NAME) -v $(MAIN)
    test:
		go test -v ./...
    clean:
		go clean $(MAIN)
		rm -rf $(BINARY_NAME)
    run:
		go run $(MAIN)
    generate:
		go generate ./...
    lint:
		golangci-lint run ./... -v
    compose_up:
		docker-compose up
    compose_down:
		docker-compose down