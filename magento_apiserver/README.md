## Products and Orders API Server

This guide documents the set up of an API server for Hyperledger Fabric using ExpressJS. It will also cover the different api calls that are available for the products and orders API Server. 

### Pre-Requisite

- Hyperledger Fabric set-up completed with all depencies installed.
- Connection configuration file for Hyperledger Fabric (eg. `connection-org1.json`)
- Order and Product Chaincode installed and initialised on Hyperledger Fabric with access to Magento backend using API calls. Ledgers are initated for both chaincode. 

### Setting up API Server

- Clone the API Server in the folder with the connection configuration 
- Create a wallet for the API server.

```bash
npm install
node enrollAdmin.js
node registerUser.js {username}
```

> If you are running this API on a seperate server, you will need the following files:
>
> - `connection-org1.json`
> - `wallet/{username}`
>
> Also, in the `order_product_apiserver.js` file, change the path of the wallet and connection profile such that it points to the actual path of the files above. 

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
node order_product_apiserver.js {username}
```



### API Server: Design

There are 2 groups of APIs, one for orders and one for products.

#### Product APIs

- `GET /api/queryProduct/` return a product record of the `id` specified.
  - Parameters: `id` 

- `GET /api/queryAllProducts/` return product records of all products in the ledger. 
- `POST /api/createProduct/` add a new product record by retreving it from the Magento backend based on the `sku` specified.
  - Parameters: `sku` 
- `PUT /api/editProduct/` update a product records by retrieving it from the Magento backend based on the `sku` specified.
  - Parameters: `sku` 
- `DELETE /api/deleteProduct/` remove a product record from the ledger based on the `id`specified.
  - Parameters: `id` 

#### Order APIs

- `GET /api/queryOrder/` return a order record of the `entity_id` specified. 
  - Parameters: `entity_id` 
- `GET /api/queryAllOrders/` return order records for all orders in the ledger.
- `POST /api/createOrder/` add a new order record by retreving it from the Magento backend based on the `entity_id` specified.
  - Parameters: `entity_id` 
- `PUT /api/editOrder/` update an order records by retrieving it from the Magento backend based on the `entity_id` specified.
  - Parameters: `entity_id` 
- `DELETE /api/deleteOrder/` remove an order record from the ledger based on the `entity_id`specified.
  - Parameters: `entity_id` 
