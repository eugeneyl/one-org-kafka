# A Fabric Network deployed on 3 Nodes
A Fabric Network of 1 Orderer with Kafka, 1 Organization with 3 peers, deployed on 3 nodes. This network is tested on Digital Ocean droplets with Ubuntu 18.04 with Hyperledger Fabric 1.4.4. Bare minimum of 2GB RAM for around 60% CPU (high) utilisation.

This setup uses docker swarm. If you want to use extra_hosts, call `git checkout extra_hosts`

### Nodes Configuration

The set up of the nodes are as followed: 

| Node | Zookeeper | Kafka | Orderer | Peer | CLI |
| --- | --- | --- | --- | --- | --- |
| 1 | zookeeper0 | kafka0, kafka1 | orderer0.example.com | peer0.org1.example.com|cli |
| 2 ||  | | peer1.org1.example.com|cli |
| 3 | | | | peer2.org1.example.com|cli |

### Setup Instructions

1. Set up your new droplet. On all 3 droplets, run the following command to start it up.

```bash
sudo apt-get update && sudo apt-get upgrade
```

2. Create a new user and switch into the user.

```bash
sudo adduser frog
sudo usermod -aG sudo frog
su - frog
```

3. Download dependecies for Hyperledger Fabric.

```bash
curl -O https://hyperledger.github.io/composer/latest/prereqs-ubuntu.sh
chmod u+x prereqs-ubuntu.sh
./prereqs-ubuntu.sh

wget https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz
tar -xzvf go1.11.2.linux-amd64.tar.gz
sudo mv go/ /usr/local
nano ~/.bashrc

#(add these 2 lines to end of file)
export GOPATH=/usr/local/go
export PATH=$PATH:$GOPATH/bin

#Log out and log in again for the chances to happen
exit
su - frog
```

4. Fetch fabric image and other tools required to generate channel artefacts and certificates.

```bash
curl -sSL http://bit.ly/2ysbOFE | bash -s
```

5. Set up the overlay network using docker swarm (our example IP address 167.71.121.213)

   - (Node 1) Initialise the swarm network.

   ```bash
   docker swarm init --advertise-addr 167.71.121.213
   ```

   - (Node 2 and Node 3) You will be given a command that will look something like the following. Run this on Node 2 and Node 3.

   ```bash
   docker swarm join --token SWMTKN-1-58xwguh8oa3jj6rcbcm4cyxg9lmitxyv1fs1sn1d5xy51e9arv-1hts5vhxebpjf6fz3kjskpbub 167.71.121.213:2377
   ```

   - (Node 1) Create the overlay network. You will be able to see the swarm network if you call `docker network ls`.

   ```bash
   docker network create --attachable --driver overlay fabric
   ```

   - (Node 2 & 3) Create a busybox container to be linked to the overlay network

   ```bash
   docker run -itd --name mybusybox --network fabric busybox
   ```

6. (Node 1) Download files required to set up the fabric.

```bash
git clone https://github.com/eugeneyl/one-org-kafka.git
```

7. (Node 1) Generate the channel artefacts and certificates required.

```bash
cd fabric-samples
export PATH=$PATH:$PWD/bin
cd one-org-kafka

#Change the ip address of the nodes to the ip address of your droplets. 
nano .env

./generate.sh
```

> <u>Additional Step</u>
>
> Change the FABRIC_CA_SERVER_CA_KEYFILE and FABRIC_CA_SERVER_TLS_KEYFILE of the CA in node1.yaml to reflect the actual key that is generated. You can find the key in one-org-kafka/crypto-config/peerOrganizations/[org1.example.com/ca](http://org1.example.com/ca)

8. (Node1) Zip the file and transfer the file to the other 2 nodes. You can use [filezilla](https://filezilla-project.org/ ) for this transfer.

```bash
cd ..
tar -czvf one-org-kafka.tar.gz one-org-kafka
```

9. (Node 2 and 3) Unzip the folder in the node.

```bash
tar -xzvf one-org-kafka.tar.gz one-org-kafka
```

10. Set up the docker containers for the different components.

```bash
cd one-org-kafka
docker-compose -f node1.yaml up -d
docker-compose -f node2.yaml up -d
docker-compose -f node3.yaml up -d
```

11. (Node 1)  Create the channel block and transfer it to the other nodes using [filezilla](https://filezilla-project.org/ ).

```bash
docker exec cli peer channel create -o orderer0.example.com:7050 -c mychannel -f ./channel-artifacts/channel.tx --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

docker cp cli:/opt/gopath/src/github.com/hyperledger/fabric/peer/mychannel.block .
```

12. (Node 2 and 3) Transfer channel block into the cli container.

```bash
docker cp mychannel.block cli:/opt/gopath/src/github.com/hyperledger/fabric/peer/
```

13. Join all the peers to the channel.

```bash
docker exec cli peer channel join -b mychannel.block
```

14. Install the chain code on all peers, instantiate only on node 1

```bash
docker exec cli peer chaincode install -n orders -v 1.0 -p github.com/chaincode/orders/
```

(Only on Node 1)

```bash
docker exec cli peer chaincode instantiate -o orderer0.example.com:7050 -C mychannel -n orders -v 1.0 -c '{"Args":[]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

> You can also set up other chaincodes in the `./chaincode` directory.

15.  After you installed and instantiated the chaincodes, you can run the following commands on different peers to try out if the network is set up properly.

```bash
docker exec cli peer chaincode invoke -o orderer0.example.com:7050 -C mychannel -n orders -c '{"Args":["initLedger"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

docker exec cli peer chaincode query -C mychannel -n orders -c '{"Args":["queryAllOrders"]}'

docker exec cli peer chaincode invoke -o orderer0.example.com:7050 -C mychannel -n orders -c '{"Args":["createOrder","ORDER14", "23459348", "5493058", "Pending"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem

docker exec cli peer chaincode query -C mychannel -n orders -c '{"Args":["queryOrder", "ORDER2"]}'
```

### Tearing down of network

In order for you to tear down the entire Hyperledger Fabric network, you can call the following command on each droplet.

```bash
docker rm -f $(docker ps -aq) && docker volume prune
```

### Nth peer

In order for you to have more peer nodes, you need to make the following edits:

1. Add the IP address of the new node to the `.env`.
2. Change the count number of peers in `crypto-config.yaml` to desired number.
3. Create a docker compose file for each of the new nodes using the `node_example.yaml` file.

### Next step forward

After you have set up your network, you can add the following feature to familiarise youself with the HLF network.

- Setting up [Hyperledger Explorer](https://github.com/eugeneyl/one-org-kafka/tree/master/explorer) to visualise the network.
- Install other chaincodes. Currently implemented chaincodes:
  1. Sample orders chaincode (Set up during this tutorial)
  2. [Magento Orders](https://github.com/eugeneyl/one-org-kafka/tree/master/chaincode/magento_order)
  3. [Magento Products](https://github.com/eugeneyl/one-org-kafka/tree/master/chaincode/magento_product)
  4. [Basic Tokens](https://github.com/eugeneyl/one-org-kafka/tree/master/chaincode/basic_tokens)
  5. (Coming Soon) ERC-like tokens on Hyperledger Fabric
- Setting up of API Servers to call chaincode functions through API calls. Currently implemented API servers:
  - [Sample API Server](https://github.com/eugeneyl/one-org-kafka/tree/master/sample_apiserver) using sample orders chaincode. 
  - [Magento API Server](https://github.com/eugeneyl/one-org-kafka/tree/master/magento_apiserver) for products and orders chaincode

