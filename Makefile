BINARY_NAME=crawler
CMD_PATH=cmd/crawler/main.go

build:
	@echo "Building..."
	go build -o out/$(BINARY_NAME) $(CMD_PATH)

run: build
	@echo "Running..."
	./out/$(BINARY_NAME) https://wagslane.dev 10 10

test:
	@echo "Testing..."
	go test ./...

json: build
	@echo "Running with JSON output..."
	./out/$(BINARY_NAME) -url https://wagslane.dev -pages 10 -json

clean:
	@echo "Cleaning..."
	go clean
	rm -rf out

.PHONY: build run test clean
