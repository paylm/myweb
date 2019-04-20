# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
BINARY_NAME=myweb
BINARY_UNIX=$(BINARY_NAME)_unix

all: test build 
build: deps
	$(GOBUILD) -o $(BINARY_NAME) -v
test: testdeps
	export GIN_MODE=release
	$(GOTEST) -v ./...
clean:
	$(GOCLEAN)
		rm -f $(BINARY_NAME)
		rm -f $(BINARY_UNIX)
run:
	$(GOBUILD) -o $(BINARY_NAME) -v .
	#./$(BINARY_NAME)
deps:
	$(GOGET) github.com/gomodule/redigo/redis 
	$(GOGET) github.com/go-sql-driver/mysql
	$(GOGET) github.com/gin-gonic/contrib/sessions 
testdeps:
	$(GOGET) github.com/stretchr/testify/assert

# Cross compilation
build-linux:
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v
docker-build:
		docker run --rm -it -v "$(GOPATH)":/go -w /go/src/bitbucket.org/rsohlich/makepost golang:latest go build -o "$(BINARY_UNIX)" -v

.PHONY: clean
