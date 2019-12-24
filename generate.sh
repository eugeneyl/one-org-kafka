#!/bin/bash
#
SYS_CHANNEL="sys-channel"
CHANNEL_NAME="mychannel"

export $(egrep -v '^#' .env | xargs)

if [ -d "crypto-config" ]; then
rm -Rf crypto-config
fi
set -x
cryptogen generate --config=./crypto-config.yaml
res=$?
set +x
if [ $res -ne 0 ]; then
    echo "Failed to generate certificates..."
    exit 1
fi
echo
echo "Generate CCP files for Org1 and Org2"
./ccp-generate.sh

echo "##########################################################"
echo "#########  Generating Orderer Genesis block ##############"
echo "##########################################################"

configtxgen -profile OneOrgsOrdererGenesis -channelID $SYS_CHANNEL -outputBlock ./channel-artifacts/genesis.block

res=$?
set +x
if [ $res -ne 0 ]; then
    echo "Failed to generate orderer genesis block..."
    exit 1
fi

echo
echo "#################################################################"
echo "### Generating channel configuration transaction 'channel.tx' ###"
echo "#################################################################"
set -x
configtxgen -profile OneOrgsChannel -outputCreateChannelTx ./channel-artifacts/channel.tx -channelID $CHANNEL_NAME
res=$?
set +x
if [ $res -ne 0 ]; then
    echo "Failed to generate channel configuration transaction..."
    exit 1
fi

echo
echo "#################################################################"
echo "#######    Generating anchor peer update for Org1MSP   ##########"
echo "#################################################################"
set -x
configtxgen -profile OneOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID $CHANNEL_NAME -asOrg Org1MSP
res=$?
set +x
if [ $res -ne 0 ]; then
    echo "Failed to generate anchor peer update for Org1MSP..."
    exit 1
fi
