#!/bin/sh

SCRIPT_DIR=$(cd `dirname $0` && pwd)

# Init configs the app home dir
rm -rf ~/.oracle
oracled init

# Modify the config.toml template
cat ${SCRIPT_DIR}/config.toml | \
    sed "s|__ORACLE_MNEMONIC__|${E2E_ORACLE_MNEMONIC}|g" | \
    sed "s|__AWS_S3_ACCESS_KEY_ID__|${E2E_AWS_S3_ACCESS_KEY_ID}|g" | \
    sed "s|__AWS_S3_SECRET_ACCESS_KEY__|${E2E_AWS_S3_SECRET_ACCESS_KEY}|g" \
    > ~/.oracle/config.toml

# Start the process
oracled start
