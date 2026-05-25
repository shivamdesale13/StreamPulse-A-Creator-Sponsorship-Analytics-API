APP     := streampulse
BIN_DIR := ./bin

.PHONY: run build clean tidy test vet

run:
	go run ./cmd/server

build:
	go build -o $(BIN_DIR)/$(APP) ./cmd/server

clean:
	rm -rf $(BIN_DIR)

tidy:
	go mod tidy

test:
	go test ./... -v -race

vet:
	go vet ./...
