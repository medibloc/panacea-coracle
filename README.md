# Panacea Data Market Validator

A sensitive data validator for Panacea Data Marketplace

## Features

- Validating that data can be sold for a specific deal
- Encrypting and storing data for buyers

## Building and Testing
By default, DataValidator uses EGo.<br/>
It also requires an Intel CPU capable of using SGX.<br/>
Please refer [here](https://docs.edgeless.systems/ego/#/getting-started/install) for installation of EGo.

If you are building using **EGo**, do this
```bash
ego-go build ./cmd/datavald
ego sign datavald
```

Build and test without using **EGo**. However, in this case, an error occurs when using the TEE function.
```bash
make build
make test
make install
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
datavald
```
