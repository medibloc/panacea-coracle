#!/bin/bash

set -euo pipefail

SCRIPT_DIR=$(cd `dirname $0` && pwd)

echo "GOBIN: ${GOBIN}"

# NOTE: Please update this array if necessary
PKG_PREFIX=$(grep '^module ' ./go.mod | awk "{print \$2}")
TEST_PKGS_WITH_EGO=(
    "${PKG_PREFIX}/server/tee"
)

#####################################################################
echo "[Step 1] Running unit tests that don't require EGo..."

TEST_PKGS_WITHOUT_EGO=$(go list ./... | grep -v /e2e)  # except e2e/*_test.go
for pkg in ${TEST_PKGS_WITH_EGO[@]}; do
    TEST_PKGS_WITHOUT_EGO=$(echo "${TEST_PKGS_WITHOUT_EGO}" | grep -v ${pkg})
done

${GOBIN} test -v ${TEST_PKGS_WITHOUT_EGO}


#####################################################################
echo "[Step 2] Running unit tests that require EGo..."

if [ ${GOBIN} != "ego-go" ]; then
    echo "Skipping: EGo disabled"
    exit 0
fi


# Compile the binary of each test package but do not run it.
# Instead, sign the test binary for EGo and run it with the `ego` command.
# NOTE: We need this 'for' loop because Go doesn't support building a test binary for all packages.
for pkg in ${TEST_PKGS_WITH_EGO}; do
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



