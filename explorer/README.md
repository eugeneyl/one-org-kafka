## Hypeledger Explorer Setup

Hyperledger Explorer is a browser that can be used to view the activities on a block chain network. You can follow this guide to set up the Hyperledger Explorer for the multi-node Hyperledger Fabric Network.

### Setting up dependencies

Apart from the dependencies that you already have installed in order to set up Hyperledger fabric, you need to install JQ and PostgreSQL for Hyperledger Explorer.

```bash
sudo apt-get update
sudo apt-get install jq postgresql postgresql-contrib
service postgresql restart
```

### Setting up Explorer

1. Download required files.

```bash
git clone https://github.com/hyperledger/blockchain-explorer.git
```

2. Set up database.

```bash
cd blockchain-explorer/app/persistence/fabric/postgreSQL
chmod -R 775 db/
cd db
./createdb.sh
```

3. Set up the explorer configuration.

There are 2 files that you will be required to change to set up the explorer.

- /one-org-kafka/connection-org1.json
  - Replace {ADMINKEY} to the actual private key file name of the admin. You can get this by running `ls /home/frog/one-org-kafka/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/`. The file name will be a string of numbers followed by `_sk`. Eg: ca4d946a621c0f528af9a3a6d05083d16afa7763da38b7e6c8d48c346b1d3034_sk
- `/home/frog/blockchain-explorer/app/platform/fabric/config.json`
  - Change `profile` to `/home/frog/one-org-kafka/connection-org1.json`
- `/home/frog/blockchain-explorer/start.sh`
  - Change `DISCOVERY_AS_LOCALHOST` to false.

4. Add entries in `/etc/hosts` such that they point to the Fabric Nodes. (This will not be necessary if the domain names are mapped to the right IP addresses.)

```bash
sudo nano /etc/hosts

#Add the following entries in the hosts file
[Fabric-Node-IP] orderer0.example.com
[Fabric-Node-IP] peer0.org1.example.com
[Fabric-Node-IP] peer1.org1.example.com
[Fabric-Node-IP] peer2.org1.example.com
```

5. Build Hyperledger Explorer

```bash
cd blockchain-explorer
./main.sh install
./start.sh
```

Your Hyperledger Explorer should be properly set up and you can access it at http://<Your-IP-Address>:8080. If it prompts you to log in, use admin:adminpw.

If there are any errors, you can refer to `./blockchain-explorer/logs/console/console.log` to troubleshoot.

> You can also use ./main to clean or test the explorer dependencies as well.
>
> - `./main.sh clean` : clear all the dependencies
> - `./main.sh test` : run tests for the explorer setup

### Tearing down of Hyperledger Explorer

In order for to stop the Hyperledger Explorer, you can call the following command. 

```bash
./stop.sh
```

