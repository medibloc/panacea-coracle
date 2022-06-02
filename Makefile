export GO111MODULE = on

GOBIN ?= ego-go

build_tags := $(strip $(BUILD_TAGS))
BUILD_FLAGS := -tags "$(build_tags)"

OUT_DIR = ./build

.PHONY: all build test install ego-sign clean

all: build test install

build: go.sum
	$(GOBIN) build -mod=readonly $(BUILD_FLAGS) -o $(OUT_DIR)/oracled ./cmd/oracled

test:
	GOBIN=$(GOBIN) ./run_gotest.sh

# Set env vars used by ./e2e/docker-compose.yml before running this target (recommended to use .env file).
e2e-test:
	docker-compose -f ./e2e/docker-compose.yml pull
	docker-compose -f ./e2e/docker-compose.yml up --build --force-recreate --abort-on-container-exit --exit-code-from e2e-test

install: go.sum
	$(GOBIN) install -mod=readonly $(BUILD_FLAGS) ./cmd/oracled

# TODO: more args for private.pem
ego-sign:
	ego sign $(EXE)

clean:
	$(GOBIN) clean
	rm -rf $(OUT_DIR)
