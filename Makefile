PATH_TO_BINARY=./cmd/
.DEFAULT_GOAL := run

run: tidy
	go run ${PATH_TO_BINARY}$(BIN)/$(BIN).go

build-wasm: tidy
	mkdir -p target
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.js target/
	cp $$(go env GOROOT)/misc/wasm/wasm_exec.html target/index.html
	GOARCH=wasm GOOS=js go build -o target/$(BIN).wasm ${PATH_TO_BINARY}$(BIN)/$(BIN).go
	sed -i 's/test.wasm/$(BIN).wasm/g;' target/index.html
	sed -i '40idocument.getElementById("runButton").remove();' target/index.html

clear:
	rm -rf target
	go clean

tidy:
	go mod tidy

