#!/bin/sh
# echo "Init Genesis Block"
# geth --datadir=/root/.ethereum init /root/genesis.json

echo "Starting Geth Test Environment..."
geth --datadir=/root/.ethereum \
    --rpcapi="db,personal,eth,net,web3" \
    --rpccorsdomain="*" --rpc --rpcaddr="0.0.0.0" \
    --rpcvhosts="*" \
    --dev