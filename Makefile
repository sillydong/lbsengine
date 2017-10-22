# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=bin/example
BINARY_UNIX=$(BINARY_NAME)_unix

all: build-view build build-linux
build-view:
	cd example/view && npm run build
	cd example && go-bindata-assetfs view/dist/...
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./example/...
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)
	rm -rf example/view/dist/*
run-view:
	cd example/view && npm run dev
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
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) ./example/...
