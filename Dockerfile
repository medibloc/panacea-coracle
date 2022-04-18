FROM golang:alpine AS build-env

# Install minimum necessary dependencies,
RUN apk add --no-cache ca-certificates build-base git wget

RUN wget https://github.com/stedolan/jq/releases/download/jq-1.6/jq-linux64 && \
    mv jq-linux64 /usr/local/bin/jq && \
    chmod +x /usr/local/bin/jq

WORKDIR /src

COPY . .

# Because dataval refers to panacea-core, 'libwasmvm_muslc.a' (or .so) is required.
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.0.0-beta2/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 3f5de8df9c6b606b4211f90edd681c84b0ecd870fdbf50678b6d9afd783a571c

# Because we want to use 'libwasmvm_muslc.a', the 'muslc' build tag must be passed to build CosmWasm/wasmvm.
RUN BUILD_TAGS=muslc GOBIN=go make build

FROM alpine:edge

COPY --from=build-env /src/build/datavald /usr/bin/datavald

EXPOSE 8080

CMD ["/usr/bin/datavald"]
