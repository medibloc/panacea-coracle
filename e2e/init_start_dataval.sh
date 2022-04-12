#!/bin/sh

SCRIPT_DIR=$(cd `dirname $0` && pwd)

# Init configs the app home dir
rm -rf ~/.dataval
datavald init

mkdir ${SCRIPT_DIR}/config
touch ${SCRIPT_DIR}/config/data_encryption_key.sealed

# Modify the config.toml template
cat ${SCRIPT_DIR}/config.toml | \
    sed "s|__VALIDATOR_MNEMONIC__|${E2E_DATAVAL_MNEMONIC}|g" | \
    sed "s|__AWS_S3_ACCESS_KEY_ID__|${E2E_AWS_S3_ACCESS_KEY_ID}|g" | \
    sed "s|__AWS_S3_SECRET_ACCESS_KEY__|${E2E_AWS_S3_SECRET_ACCESS_KEY}|g" \
    > ~/.dataval/config.toml

# Start the process
datavald start
