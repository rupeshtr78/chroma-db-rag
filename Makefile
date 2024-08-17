# Variables
BINARY_NAME=chroma-db
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get

# Build targets
all: build

build:
	$(GOBUILD) -o $(BINARY_NAME) cmd/main.go

test:
	$(GOTEST) ./...

clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)

run:
	./$(BINARY_NAME)

# Additional targets
deps:
	$(GOGET) -u ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

.PHONY: all build test clean run deps docker-up docker-down
