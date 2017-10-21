# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=bin/example
BINARY_UNIX=$(BINARY_NAME)_unix

all: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v example/main.go
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v example/main.go
	./$(BINARY_NAME)
analysis:
	@echo "engine"
	@ cloc --exclude-dir=example ./
	@echo "example"
	@ cloc example/main.go
	@echo "view"
	@ cloc example/view/src

# Cross compilation
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) example/main.go
