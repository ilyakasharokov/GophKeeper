.PHONY: build
build: build-client build-server

BIN_DIR=./bin/
CLIENT_EXECUTABLE=gophkeeper
CLIENT_WINDOWS=$(BIN_DIR)$(CLIENT_EXECUTABLE)_windows_amd64.exe
CLIENT_LINUX=$(BIN_DIR)$(CLIENT_EXECUTABLE)_linux_amd64
CLIENT_DARWIN=$(BIN_DIR)$(CLIENT_EXECUTABLE)_darwin_amd64

GIT_COMMIT := $(shell git rev-list -1 HEAD)
BUILD_DATE := $(shell date +%FT%T%z)
VERSION := $(shell git describe --tags --abbrev=0 --always)

$(CLIENT_WINDOWS):
	GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o $(CLIENT_WINDOWS) \
		-ldflags="-X 'main.BuildCommit=$(GIT_COMMIT)'\
		 -X 'main.BuildVersion=$(VERSION)'\
		 -X 'main.BuildUser=$(USER)'\
		  -X 'main.BuildDate=$(BUILD_DATE)'" \
		./cmd/client/*.go

$(CLIENT_LINUX):
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $(CLIENT_LINUX) \
		-ldflags="-X main.BuildCommit=$(GIT_COMMIT)'\
		 -X 'main.BuildVersion=$(VERSION)'\
		 -X 'main.BuildUser=$(USER)'\
		  -X 'main.BuildDate=$(BUILD_DATE)'" \
		./cmd/client/*.go

$(CLIENT_DARWIN):
	GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 go build -o $(CLIENT_DARWIN) \
		-ldflags="-X 'main.BuildCommit=$(GIT_COMMIT)'\
		 -X 'main.BuildVersion=$(VERSION)'\
		 -X 'main.BuildUser=$(USER)'\
		  -X 'main.BuildDate=$(BUILD_DATE)'" \
		./cmd/client/*.go

build-client: $(CLIENT_LINUX) $(CLIENT_WINDOWS) $(CLIENT_DARWIN)
	@echo Version: $(VERSION)

build-server:
	@echo "Building the server app to the bin dir"
	CGO_ENABLED=1 go build -o ./bin/gophkeeper_srv \
		-ldflags="-X 'main.BuildCommit=$(GIT_COMMIT)'\
		 -X 'main.BuildVersion=$(VERSION)'\
		 -X 'main.BuildUser=$(USER)'\
		  -X 'main.BuildDate=$(BUILD_DATE)'" \
		./cmd/server/*.go

clean: ## Remove previous build
	rm -f $(CLIENT_LINUX) $(CLIENT_WINDOWS) $(CLIENT_DARWIN)