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
DATAVAL_LOG_LEVEL=info \  # panic, fatal, error, warn, info, debug, trace
DATAVAL_HTTP_LADDR=0.0.0.0:8080 \
datavald
```
