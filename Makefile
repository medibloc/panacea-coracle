export GO111MODULE = on

OUT_DIR = ./build

.PHONY: all build test install clean

all: build test install

build: go.sum
	go build -mod=readonly -o $(OUT_DIR)/ ./cmd/datavald

test:
	go test -v ./...

install: go.sum
	go install -mod=readonly ./cmd/datavald

clean:
	go clean
	rm -rf $(OUT_DIR)
