# Panacea Data Market Validator

A sensitive data validator for Panacea Data Marketplace

## Features

- Validating that data can be sold for a specific deal
- Encrypting and storing data for buyers

## Building and Testing

It's recommended to run the data validator in the secure enclave.
So, the following commands build the `datavald` using EGo which requires Intel SGX. For more details, please refer the [EGo installation guide](https://docs.edgeless.systems/ego/#/getting-started/install).
Build and test without using **EGo**. However, in this case, an error occurs when using the TEE function.
```bash
make build    # generates a binary: ./build/datavald
make test

# https://docs.edgeless.systems/ego/#/workflows/build?id=sign-and-run
EXE="./build/datavald" make ego-sign  
```

If you build the `datavald` without using EGo, please run make commands with the explicit `GOBIN` environment variable.
```bash
GOBIN=go make build
GOBIN=go make test
```

## Running
If you are building using **EGo**, do this
```bash
# Supported log levels: panic, fatal, error, warn, info, debug, trace
EDG_DATAVAL_LOG_LEVEL=info \
EDG_DATAVAL_HTTP_LADDR=0.0.0.0:8080 \
EDG_DATAVAL_PANACEA_GRPC_ADDR=0.0.0.0:9090 \
EDG_DATAVAL_VALIDATOR_MNEMONIC={Your mnemonic} \
EDG_DATAVAL_AWS_S3_BUCKET="my-s3-bucket" \
EDG_DATAVAL_AWS_S3_REGION="ap-northeast-2" \
EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID="my-access-key" \
EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY="my-secret-access-key" \
ego run datavald
```
Run without using **EGo**
```bash
# Supported log levels: panic, fatal, error, warn, info, debug, trace
EDG_DATAVAL_LOG_LEVEL=info \
EDG_DATAVAL_HTTP_LADDR=0.0.0.0:8080 \
EDG_DATAVAL_PANACEA_GRPC_ADDR=0.0.0.0:9090 \
EDG_DATAVAL_VALIDATOR_MNEMONIC={Your mnemonic} \
EDG_DATAVAL_AWS_S3_BUCKET="my-s3-bucket" \
EDG_DATAVAL_AWS_S3_REGION="ap-northeast-2" \
EDG_DATAVAL_AWS_S3_ACCESS_KEY_ID="my-access-key" \
EDG_DATAVAL_AWS_S3_SECRET_ACCESS_KEY="my-secret-access-key" \
datavald
```
