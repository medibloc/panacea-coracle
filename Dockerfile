FROM golang:alpine AS build-env

# Install minimum necessary dependencies,
RUN apk add --no-cache ca-certificates build-base

WORKDIR /src

COPY . .

# Because dataval refers to panacea-core, 'libwasmvm_muslc.a' (or .so) is required.
ADD https://github.com/CosmWasm/wasmvm/releases/download/v0.14.0/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 220b85158d1ae72008f099a7ddafe27f6374518816dd5873fd8be272c5418026

# Because we want to use 'libwasmvm_muslc.a', the 'muslc' build tag must be passed to build CosmWasm/wasmvm.
RUN BUILD_TAGS=muslc make build


FROM alpine:edge

COPY --from=build-env /src/build/datavald /usr/bin/datavald

EXPOSE 8080

CMD ["/usr/bin/datavald"]
