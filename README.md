# Panacea Data Market Validator

A sensitive data validator for Panacea Data Marketplace

## Features

- Validating that data can be sold for a specific deal
- Encrypting and storing data for buyers

## Building and Testing

It's recommended to run the data validator in the secure enclave.
So, the following commands build the `datavald` using EGo which requires Intel SGX. For more details, please see [EGo installation guide](https://docs.edgeless.systems/ego/#/getting-started/install).

```bash
make build    # generates a binary: ./build/datavald
make test

# https://docs.edgeless.systems/ego/#/workflows/build?id=sign-and-run
EXE="./build/datavald" make ego-sign
```

If you build the `datavald` without using EGo, please run make commands with the explicit `GOBIN` environment variable.
Then, enclave-related features will not work.

```bash
GOBIN=go make build
GOBIN=go make test
```

## Running

```bash
datavald init  # initialize configs in the app home dir (~/.dataval)
datavald start
```
