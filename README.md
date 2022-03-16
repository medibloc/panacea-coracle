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
EDG_DATAVAL_LOG_LEVEL=info \
EDG_DATAVAL_HTTP_LADDR=0.0.0.0:8080 \
EDG_DATAVAL_PANACEA_GRPC_ADDR=0.0.0.0:9090 \
EDG_DATAVAL_VALIDATOR_MNEMONIC="category quick episode sugar argue napkin clap imitate adult square cake oven best village unlock river pilot symbol sick october cart sheriff cream valid" \
EDG_DATAVAL_AWS_S3_BUCKET="data-market" \
EDG_DATAVAL_AWS_S3_REGION="ap-northeast-2" \
datavald
```
