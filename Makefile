export GO111MODULE = on

build_tags := $(strip $(BUILD_TAGS))
BUILD_FLAGS := -tags "$(build_tags)"

OUT_DIR = ./build

.PHONY: all build test install clean

all: build test install

build: go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o $(OUT_DIR)/datavald ./cmd/datavald

UNIT_TESTS=$(shell go list ./... | grep -v /e2e)  # except e2e/*_test.go
test:
	go test -v $(UNIT_TESTS)

# Set env vars used by ./e2e/docker-compose.yml before running this target.
e2e-test:
	docker-compose -f ./e2e/docker-compose.yml pull
	docker-compose -f ./e2e/docker-compose.yml up --build --force-recreate --abort-on-container-exit --exit-code-from e2e-test

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/datavald

clean:
	go clean
	rm -rf $(OUT_DIR)
