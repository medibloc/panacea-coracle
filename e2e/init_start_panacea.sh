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
panacead gentx validator 1000000umed --commission-rate 0.1 --commission-max-rate 0.2 --commission-max-change-rate 0.01  --min-self-delegation 1 --chain-id ${CHAIN_ID}

echo -e "${E2E_DATA_BUYER_MNEMONIC}\n\n" | panacead keys add buyer -i
panacead add-genesis-account $(panacead keys show buyer -a) 100000000umed

echo -e "${E2E_DATAVAL_MNEMONIC}\n\n" | panacead keys add dataval -i
panacead add-genesis-account $(panacead keys show dataval -a) 100000000umed

panacead collect-gentxs

# Run panacead in background and get its pid.
# gRPC is temporarily disabled, so that other docker containers do not think that panacead is ready to serve.
panacead start --grpc.enable=false &
PID_PANACEAD=$!

# Wait for the 26657 to be opened
#TODO: After modifying Dockerfile of panacea-core to use alpine, ${SCRIPT_DIR}/wait-for 127.0.0.1:26657
# Wait for the 1st block to be created
sleep 10

panacead tx bank send $(panacead keys show dataval -a) $(panacead keys show validator -a) 100umed --chain-id ${CHAIN_ID} -b block --yes
panacead tx bank send $(panacead keys show validator -a) $(panacead keys show dataval -a) 100umed --chain-id ${CHAIN_ID} -b block --yes

DATAVAL_ADDR=$(panacead keys show dataval -a)
sed 's|"trusted_data_validators": \[\]|"trusted_data_validators": ["'"${DATAVAL_ADDR}"'"]|g' ${SCRIPT_DIR}/create_deal.json > /tmp/create_deal.json

panacead tx datadeal create-deal \
  --deal-file /tmp/create_deal.json \
  --from buyer \
  --chain-id ${CHAIN_ID} \
  -b block \
  --yes

# Kill the background panacead and wait for the process to exit.
kill ${PID_PANACEAD} && wait ${PID_PANACEAD}

panacead start
