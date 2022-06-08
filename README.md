# Panacea Oracle

An oracle which validates off-chain data to be transacted in the data exchange protocol of the Panacea chain while preserving privacy

## Features

- Validating that data meets the requirements of a specific deal/pool
  - with utilizing TEE (Trusted Execution Environment) for preserving privacy
- Providing encrypted data to buyers

## Building and Testing

It's recommended to run the `oracled` in the secure enclave.
So, the following commands build the `oracled` using EGo which requires Intel SGX. For more details, please see [EGo installation guide](https://docs.edgeless.systems/ego/#/getting-started/install).

```bash
make build    # generates a binary: ./build/oracled
make test

# https://docs.edgeless.systems/ego/#/workflows/build?id=sign-and-run
EXE="./build/oracled" make ego-sign
```

If you build the `oracled` without using EGo, please run make commands with the explicit `GO` environment variable.
Then, enclave-related features will not work.

```bash
GO=go make build
GO=go make test
```

## Running

```bash
oracled init  # initialize configs in the app home dir (~/.oracle)
oracled start
```
