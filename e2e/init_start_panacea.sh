#!/bin/bash

set -euxo pipefail

SCRIPT_DIR=$(cd $(dirname $0) && pwd)

CHAIN_ID="testing"

# Init the panacea directory
rm -rf ~/.panacea
panacead init node1 --chain-id ${CHAIN_ID}

# Init accounts
panacead keys add validator
panacead add-genesis-account $(panacead keys show validator -a) 1000000000000umed
panacead gentx validator 1000000umed --commission-rate 0.1 --commission-max-rate 0.2 --commission-max-change-rate 0.01 --min-self-delegation 1 --chain-id ${CHAIN_ID}

echo -e "${E2E_DATA_BUYER_MNEMONIC}\n\n" | panacead keys add curator -i
panacead add-genesis-account $(panacead keys show curator -a) 100000000000umed

echo -e "${E2E_DATAVAL_MNEMONIC}\n\n" | panacead keys add dataval -i
panacead add-genesis-account $(panacead keys show dataval -a) 100000000000umed

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

DATAVAL_ADDR=$(panacead keys show dataval -a)
sed 's|"trusted_data_validators": \[\]|"trusted_data_validators": ["'"${DATAVAL_ADDR}"'"]|g' ${SCRIPT_DIR}/create_deal.json >/tmp/create_deal.json
sed 's|"trusted_data_validators": \[\]|"trusted_data_validators": ["'"${DATAVAL_ADDR}"'"]|g' ${SCRIPT_DIR}/create_pool.json >/tmp/create_pool.json

panacead tx datapool register-data-validator "https://my-endpoint.com" \
  --from dataval \
  --chain-id ${CHAIN_ID} \
  -b block \
  --yes

panacead tx datadeal create-deal \
  --deal-file /tmp/create_deal.json \
  --from curator \
  --chain-id ${CHAIN_ID} \
  -b block \
  --yes

# TODO: It will be changed, when Get Module Address is merged in panacea-core
MODULE_ADDR=$(panacead q datapool module-addr -o json | jq -r '.address')

panacead tx gov submit-proposal wasm-store ${SCRIPT_DIR}/cw721_base.wasm \
  --title "store NFT contract wasm code" \
  --description "store wasm code for x/datapool module" \
  --instantiate-only-address $MODULE_ADDR \
  --run-as $MODULE_ADDR \
  --deposit "10000000000umed" \
  --from validator \
  --gas auto --gas-adjustment 1.3 \
  --chain-id ${CHAIN_ID} \
  -b block \
  --yes

panacead tx gov vote {store proposal id} yes --from validator --gas auto --gas-adjustment 1.3 --chain-id ${CHAIN_ID}  -b block --yes

INST_MSG=$(jq -n --arg name "curator" --arg symbol "CUR" --arg minter $MODULE_ADDR '{"name": $name, "symbol": $symbol, "minter": $minter}')

panacead tx gov submit-proposal instantiate-contract {code id} "$INST_MSG" \
  --label "curator NFT" \
  --title "instantiate NFT contract" \
  --description "instantiate NFT contract for x/datapool module" \
  --run-as MODULE_ADDR \
  --admin MODULE_ADDR \
  --deposit "100000000umed" \
  --from validator \
  --gas auto --gas-adjustment 1.3 \
  --chain-id ${CHAIN_ID} \
  -b block \
  --yes

panacead tx gov vote {instantiation proposal id} yes --from validator --gas auto --gas-adjustment 1.3 --chain-id ${CHAIN_ID} -b block --yes

panacead tx gov submit-proposal param-change ${SCRIPT_DIR}/param_change_sample.json --from validator --gas auto --gas-adjustment 1.3 --chain-id ${CHAIN_ID} -b block --yes

panacead tx gov vote {param-change proposal id} yes --from validator --gas auto --gas-adjustment 1.3 --chain-id ${CHAIN_ID} -b block --yes

panacead tx datapool create-pool /tmp/create_pool.json \
  --from curator \
  --chain-id ${CHAIN_ID} \
  -b block \
  --yes

# Kill the background panacead and wait for the process to exit.
kill ${PID_PANACEAD} && wait ${PID_PANACEAD}

panacead start
