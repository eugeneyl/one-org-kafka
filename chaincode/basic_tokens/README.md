## Sample Token Chaincode

This is a chaincode for a token on the Hyperledger Fabric with simple functionalities. This includes:

- Create new user accounts.
- Check total value and individual value.
- Transfer tokens from one user to another.
- Add more token externally to any accounts. 

> This is just an experimental token chaincode and it is not meant for deployment.

### Pre-Requisite

- Hyperledger Fabric set-up completed with all dependencies installed.

### Setting up the chaincode

- Install the chain code on all nodes

```bash
docker exec cli peer chaincode install -n basic-tokens -v 1.0 -p github.com/chaincode/basic-tokens/
```

- Instantiate the chain code on any one of the nodes.

```bash
docker exec cli peer chaincode instantiate -o orderer0.example.com:7050 -C mychannel -n basic-tokens -v 1.0 -c '{"Args":[]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem 
```

- Initialise the token. This will create 3 master account, 1 for each of the peers. 

```bash
docker exec cli peer chaincode invoke -o orderer0.frogfrogjump.com:7050 -C mychannel -n basic-tokens -c '{"Args":["initLedger"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/frogfrogjump.com/orderers/orderer0.frogfrogjump.com/msp/tlscacerts/tlsca.frogfrogjump.com-cert.pem
```

- You will be able to see that each of these accounts are successfully initated by the bollowing command. You will also see a string of numbers before each account. This will be the "address" of these accounts used for transactions.

```bash
docker exec cli peer chaincode query -C mychannel -n basic-tokens -c '{"Args":["queryAllAccounts"]}' 
```

### Testing the chaincode

After setting up the chaincode, you can test the different functions available:

##### Query

- queryValue(address)	Returns the value of the specified account based on the address.

```bash
docker exec cli peer chaincode query -C mychannel -n basic-tokens -c '{"Args":["queryValue", "{address}"]}' 
```

- queryTotalAmount()	Returns the total amount of all accounts in the ledger.

```bash
docker exec cli peer chaincode query -C mychannel -n basic-tokens -c '{"Args":["queryTotalAmount"]}' 
```

##### Invoke

- createAcoount(username, amount)	Creates an user account with given username and amount. Returns address of account.

```bash
docker exec cli peer chaincode invoke -o orderer0.frogfrogjump.com:7050 -C mychannel -n tokens -c '{"Args":["createAccount", "{username}", “{amount}"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/frogfrogjump.com/orderers/orderer0.frogfrogjump.com/msp/tlscacerts/tlsca.frogfrogjump.com-cert.pem 
```

- transferFrom(from_address, to_address, amount)	Transfer specified amount from one account to another.

```bash
docker exec cli peer chaincode invoke -o orderer0.frogfrogjump.com:7050 -C mychannel -n tokens -c '{"Args":["transferFrom", "{from}", "{to}", "{amount]"}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/frogfrogjump.com/orderers/orderer0.frogfrogjump.com/msp/tlscacerts/tlsca.frogfrogjump.com-cert.pem
```

- addTokens(address, amount)	Add additional tokens to a user account. 

```bash
docker exec cli peer chaincode invoke -o orderer0.frogfrogjump.com:7050 -C mychannel -n tokens -c '{"Args":["addTokens", "{address}", "{amount]"}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/frogfrogjump.com/orderers/orderer0.frogfrogjump.com/msp/tlscacerts/tlsca.frogfrogjump.com-cert.pem
```

