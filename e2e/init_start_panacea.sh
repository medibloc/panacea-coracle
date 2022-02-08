#!/bin/bash

set -euxo pipefail

SCRIPT_DIR=$(cd `dirname $0` && pwd)

CHAIN_ID="testing"

# Init the panacea directory
rm -rf ~/.panacea
panacead init node1 --chain-id ${CHAIN_ID}

# Init accounts
panacead keys add validator
panacead add-genesis-account $(panacead keys show validator -a) 100000000umed

echo -e "${DATA_BUYER_MNEMONIC}\n\n" | panacead keys add buyer -i
panacead add-genesis-account $(panacead keys show buyer -a) 100000000umed

echo -e "${DATAVAL_VALIDATOR_MNEMONIC}\n\n" | panacead keys add dataval -i
panacead add-genesis-account $(panacead keys show dataval -a) 100000000umed

# Init validator
panacead gentx validator 1000000umed --commission-rate 0.1 --commission-max-rate 0.2 --commission-max-change-rate 0.01  --min-self-delegation 1 --chain-id ${CHAIN_ID}
panacead collect-gentxs

# Run panacead in background and get its pid.
panacead start --grpc.address "0.0.0.0:9091" &
pid=$!

# Wait for the 26657 to be opened
#${SCRIPT_DIR}/wait-for 127.0.0.1:26657
# Wait for the 1st block to be created
sleep 10

panacead tx market create-deal \
  --deal-file ${SCRIPT_DIR}/create_deal.json \
  --from buyer \
  --chain-id ${CHAIN_ID} \
  -b block \
  --yes

# Kill the background panacead and wait for the process to exit.
kill ${pid} && wait ${pid}

panacead start
