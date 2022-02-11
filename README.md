# Panacea Data Market Validator

A sensitive data validator for Panacea Data Marketplace

## Features

- Validating that data can be sold for a specific deal
- Encrypting and storing data for buyers

## Building and Testing

```bash
make build
make test
make install
```

## Running

```bash
# Supported log levels: panic, fatal, error, warn, info, debug, trace
DATAVAL_LOG_LEVEL=info \
DATAVAL_HTTP_LADDR=0.0.0.0:8080 \
DATAVAL_PANACEA_GRPC_ADDR=0.0.0.0:9090 \
DATAVAL_VALIDATOR_MNEMONIC={Your mnemonic} \
DATAVAL_AWS_S3_BUCKET="my-s3-bucket" \
DATAVAL_AWS_S3_REGION="ap-northeast-2" \
datavald
```
