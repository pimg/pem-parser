@PHONY:run
run:
	go run main.go

@PHONY:test
test:
	go test ./...

@PHONY: lint
lint:
	golangci-lint run --out-format=colored-line-number

@PHONY: build
build:
	go build -o pem-parser
