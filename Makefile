export GO111MODULE = on

build_tags := $(strip $(BUILD_TAGS))
BUILD_FLAGS := -tags "$(build_tags)"

OUT_DIR = ./build

.PHONY: all build test install clean

all: build test install

build: go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o $(OUT_DIR)/datavald ./cmd/datavald

test:
	go test -v ./...

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/datavald

clean:
	go clean
	rm -rf $(OUT_DIR)
