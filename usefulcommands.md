sudo apt-get update && sudo apt-get upgrade
sudo adduser frog
sudo usermod -aG sudo frog

su - frog

curl -O https://hyperledger.github.io/composer/latest/prereqs-ubuntu.sh
chmod u+x prereqs-ubuntu.sh
./prereqs-ubuntu.sh

wget https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz
tar -xzvf go1.11.2.linux-amd64.tar.gz
sudo mv go/ /usr/local
vim ~/.bashrc
#(add these 2 lines to end of file)
export GOPATH=/usr/local/go
export PATH=$PATH:$GOPATH/bin

sudo apt-get install jq

exit
su - frog

curl -sSL http://bit.ly/2ysbOFE | bash -s

docker swarm init --advertise-addr {ip address of master node}

docker swarm join --token SWMTKN-1-58xwguh8oa3jj6rcbcm4cyxg9lmitxyv1fs1sn1d5xy51e9arv-1hts5vhxebpjf6fz3kjskpbub 167.71.121.213:2377

docker network create --attachable --driver overlay fabric
docker run -itd --name mybusybox --network fabric busybox

tar -czvf one-org-kafka.tar.gz one-org-kafka
tar -xzvf one-org-kafka.tar.gz one-org-kafka

export PATH=$PATH:/home/frog/fabric-samples/bin
export PATH=$PATH:$PWD/bin
./generate.sh

docker-compose -f node1.yaml up -d
docker-compose -f node2.yaml up -d
docker-compose -f node3.yaml up -d
docker exec cli peer channel create -o orderer0.frogfrogjump.com:7050 -c mychannel -f ./channel-artifacts/channel.tx

docker exec cli peer channel create -o orderer0.frogfrogjump.com:7050 -c mychannel -f ./channel-artifacts/channel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/frogfrogjump.com/orderers/orderer0.frogfrogjump.com/msp/tlscacerts/tlsca.frogfrogjump.com-cert.pem

docker exec cli peer channel join -b mychannel.block

docker cp cli:/opt/gopath/src/github.com/hyperledger/fabric/peer/mychannel.block .

docker cp mychannel.block cli:/opt/gopath/src/github.com/hyperledger/fabric/peer/

docker exec cli peer channel join -b mychannel.block

docker exec cli peer chaincode install -n orders -v 1.0 -p github.com/chaincode/

docker exec cli peer chaincode instantiate -o orderer0.frogfrogjump.com:7050 -C mychannel -n orders -v 1.0 -c '{"Args":[]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/frogfrogjump.com/orderers/orderer0.frogfrogjump.com/msp/tlscacerts/tlsca.frogfrogjump.com-cert.pem

docker exec cli peer chaincode invoke -o orderer0.frogfrogjump.com:7050 -C mychannel -n orders -c '{"Args":["initLedger"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/frogfrogjump.com/orderers/orderer0.frogfrogjump.com/msp/tlscacerts/tlsca.frogfrogjump.com-cert.pem


docker exec cli peer chaincode query -C mychannel -n orders -c '{"Args":["queryAllOrders"]}' 

docker exec cli peer chaincode invoke -o orderer0.frogfrogjump.com:7050 -C mychannel -n orders -c '{"Args":["createOrder","ORDER14", "23459348", "5493058", "Pending"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/frogfrogjump.com/orderers/orderer0.frogfrogjump.com/msp/tlscacerts/tlsca.frogfrogjump.com-cert.pem

docker exec cli peer chaincode invoke -o orderer0.frogfrogjump.com:7050 -C mychannel -n orders -c '{"Args":["createOrder","ORDER15", "23459348", "5493058", "Pending"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/frogfrogjump.com/orderers/orderer0.frogfrogjump.com/msp/tlscacerts/tlsca.frogfrogjump.com-cert.pem

docker exec cli peer chaincode invoke -o orderer0.frogfrogjump.com:7050 -C mychannel -n orders -c '{"Args":["createOrder","ORDER16", "23459348", "5493058", "Pending"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/frogfrogjump.com/orderers/orderer0.frogfrogjump.com/msp/tlscacerts/tlsca.frogfrogjump.com-cert.pem

docker exec cli peer chaincode query -C mychannel -n orders -c '{"Args":["queryOrder", "ORDER2]}' 

docker rm -f $(docker ps -aq) && docker rmi -f $(docker images | grep dev | awk '{print $3}') && docker volume prune

curl -d '{"orderno":"ORDER12","orderid":"1111111","customerid":"2222222","status":"pending"}' -H "Content-Type: application/json" -X POST http://peer3.frogfrogjump:5000/api/addorder