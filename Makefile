PATH_TO_BINARY=./cmd/
.DEFAULT_GOAL := run

run: tidy
	go run ${PATH_TO_BINARY}$(BIN)/$(BIN).go

clear:
	rm -rf target
	go clean

tidy:
	go mod tidy

