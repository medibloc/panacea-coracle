#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(cd `dirname $0` && pwd)

echo "GOBIN: ${GOBIN}"

# NOTE: Please update this array if necessary
PKG_PREFIX=$(grep '^module ' ./go.mod | awk "{print \$2}")
ENCLAVE_TEST_PKGS=(
    "${PKG_PREFIX}/server/tee"
)

#####################################################################
echo "[Step 1] Running non-enclave tests..."

NON_ENCLAVE_TEST_PKGS=$(go list ./... | grep -v /e2e)  # except e2e/*_test.go
for pkg in ${ENCLAVE_TEST_PKGS[@]}; do
    NON_ENCLAVE_TEST_PKGS=$(echo "${NON_ENCLAVE_TEST_PKGS}" | grep -v ${pkg})
done

${GOBIN} test -v ${NON_ENCLAVE_TEST_PKGS}


#####################################################################
echo "[Step 2] Running enclave tests..."

if [ ${GOBIN} != "ego-go" ]; then
    echo "Skipping: EGo disabled"
    exit 0
fi


# Compile the binary of each test package but do not run it.
# Instead, sign the test binary for EGo and run it with the `ego` command.
# NOTE: We need this 'for' loop because Go doesn't support building a test binary for all packages.
for pkg in ${ENCLAVE_TEST_PKGS}; do
    PKG_DIR=${pkg#$PKG_PREFIX}
    ENCLAVE_JSON_PATH="${SCRIPT_DIR}/${PKG_DIR}/testdata/enclave.json"
    TEST_BIN_PATH="./$(jq -c '.exe' -r ${ENCLAVE_JSON_PATH})"

    ${GOBIN} test -c -o ${TEST_BIN_PATH} ${pkg}   # Compile the test binary but do not run it
    ego sign ${ENCLAVE_JSON_PATH}                 # Sign the test binary for EGo

    # Run the test binary with env vars
    EDG_TEST_ENCLAVE_SIGNER_ID_HEX=$(ego signerid ./public.pem) \
        ego run ${TEST_BIN_PATH} -test.v

    rm -f ${TEST_BIN_PATH}
done



