## Magento Orders Chaincode

This are chaincodes that initialise and manage the ledger based on the magento orders api.

### Requirements

- Hyperledger Fabric set-up completed with all dependencies installed.
- Magento backend set up. You will need to have the access tokens for magento

### Setting up the chaincode

- Edit the chaincode to add in the `apiHost`, `magentoAccessToken` and `verifyCodeToken`

```go
const (
	apiHost            = "Magento API URL"
	magentoAccessToken = "Magento bearer token"
	verifyCodeToken    = "Magento basic token"
)
```



- You will need to install the chain code on all nodes

```bash
docker exec cli peer chaincode install -n magento_order -v 1.0 -p github.com/chaincode/magento_order/go
```

- Instantiate it on any one of the nodes.

```bash
docker exec cli peer chaincode instantiate -o orderer0.example.com:7050 -C mychannel -n magento_order -v 1.0 -c '{"Args":[]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

- Initialise the ledger. (This should only be called once)

```bash
docker exec cli peer chaincode invoke -o orderer0.example.com:7050 -C mychannel -n magento_order -c '{"Args":["initOrders"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

### Testing the chaincode

After setting up the chaincode, you can test the different functions available:

#### Query

- queryOrder(id)

```bash
docker exec cli peer chaincode query -C mychannel -n magento_order -c '{"Args":["queryOrder", "{id}"]}'
```

- queryAllOrders()

```bash
docker exec cli peer chaincode query -C mychannel -n magento_order -c '{"Args":["queryAllOrders"]}'
```

#### Invoke

- createOrder(sku)

```bash
docker exec cli peer chaincode invoke -o orderer0.example.com:7050 -C mychannel -n magento_order -c '{"Args":["createOrder","{sku}"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

- editOrder(sku)

```bash
docker exec cli peer chaincode invoke -o orderer0.example.com:7050 -C mychannel -n magento_order -c '{"Args":["editOrder","{sku}"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

- deleteOrder(id)

```bash
docker exec cli peer chaincode invoke -o orderer0.example.com:7050 -C mychannel -n magento_order -c '{"Args":["deleteOrder","{id}"]}' --tls true --cafile /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer0.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
```

