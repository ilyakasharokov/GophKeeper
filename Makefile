BINARY_NAME=server
COMPILE_PARAMS:=
COMPILE_ENV:=

ifeq ($(OS),Windows_NT)
	COMPILE_ENV += export GOARCH=amd64 && export GOOS=window &&
	COMPILE_PARAMS += cmd\client\main.go
else
	COMPILE_ENV += export GOARCH=amd64 && export GOOS=linux &&
	COMPILE_PARAMS += ./cmd/client/main.go
endif

build:
	${COMPILE_ENV} go build -o ${BINARY_NAME} ${COMPILE_PARAMS}

run:
	./${BINARY_NAME}

build_and_run: build run

clean:
	go clean
	rm ${BINARY_NAME}

test:
	go vet ./...
	go test -v -cover ./...

.PHONY: test, build