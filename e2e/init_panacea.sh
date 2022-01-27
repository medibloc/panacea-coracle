#!/bin/bash

set -eo pipefail

rm -rf ~/.panacea

panacead init node1 --chain-id testing

echo -e "${MNEMONIC}\n\n" | panacead keys add validator -i

panacead add-genesis-account $(panacead keys show validator -a) 100000000000000umed
panacead gentx validator 1000000000000umed --commission-rate 0.1 --commission-max-rate 0.2 --commission-max-change-rate 0.01  --min-self-delegation 1000000 --chain-id testing
panacead collect-gentxs

panacead start