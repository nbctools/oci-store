.PHONY: lint-go
lint-go:
	golangci-lint run ./...

run:
	go run ./main.go

build:
	@go build -o bin/oci-store ./...

quick-check:
	./demo/quick-test.sh	

.PHONY : demo
demo:
	./demo/multi-backend-demo.sh

