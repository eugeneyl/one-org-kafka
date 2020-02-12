## Sample API Server

This guide documents the set up of an API server for Hyperledger Fabric using ExpressJS. This is a sample API server based on the `orders.go` chaincode. 

### Pre-Requisite

- Hyperledger Fabric set-up completed with all depencies installed.
- Connection configuration file for Hyperledger Fabric (eg. `connection-org1.json`)
- Orders Chaincode installed and initialised on Hyperledger Fabric.

### Setting up API Server

- Clone the API Server in the folder with the connection configuration 
- Create a wallet for the API server.

```bash
npm install
node enrollAdmin.js
node registerUser.js
```

> If you are running this API on a seperate server, you will need the following files:
>
> - `connection-org1.json`
> - `wallet/user2`
>
> Also, in the `apiserver,js` file, change the path of the wallet and connection profile such that it points to the actual path of the files above. 

- Add entries in `/etc/hosts` such that they point to the Fabric Nodes. (This will not be necessary if the domain names are mapped to the right IP addresses.)

```bash
sudo nano /etc/hosts

#Add the following entries in the hosts file
[Fabric-Node-IP] orderer0.example.com
[Fabric-Node-IP] peer0.org1.example.com
[Fabric-Node-IP] peer1.org1.example.com
[Fabric-Node-IP] peer2.org1.example.com
```

- Install the required modules.

```bash
npm install
npm install express body-parser --save
```

- Run the API server. (You can use pm2 or docker to run this in the background.)

```bash
node apiserver.js
```
