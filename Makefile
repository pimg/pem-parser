@PHONY:run
run:
	go run main.go -debug

@PHONY:test
test:
	go test ./...

@PHONY: lint
lint:
	golangci-lint run

@PHONY: build
build:
	go build -o pem-parser
