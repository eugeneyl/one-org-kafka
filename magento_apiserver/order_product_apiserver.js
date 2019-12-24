var express = require('express');
var bodyParser = require('body-parser');

var app = express();
app.use(bodyParser.json());

var args = process.argv.slice(1)

console.log(args)

var user = args[1]
if (!user) {
    console.error('Enter valid user name')
}

const { FileSystemWallet, Gateway } = require('fabric-network');
const path = require('path');
const ccpPath = path.resolve(__dirname, '..', 'connection-org1.json');

//Check that wallet exist
const walletPath = path.join(process.cwd(), 'wallet');
const wallet = new FileSystemWallet(walletPath);
console.log(`Wallet path: ${walletPath}`);

checkWallet()

async function checkWallet() {
    const userExists = await wallet.exists(user);
    if (!userExists) {
        console.log(`An identity for the user ${user} does not exist in the wallet`);
        console.log('Run the registerUser.js application before retrying');
        process.exit(1);
    }
}

app.post('/api/createProduct/', async function(req, res) {
    try {
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: user, discovery: { enabled: true, asLocalhost: false }});
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('order');
        await contract.submitTransaction('createOrder', req.body.sku);
        console.log('Transaction has been submitted');
        res.send('Transaction has been submitted');
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({error: error});
    } 
});

app.get('/api/queryProduct/', async function(req, res) {
    try {
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: user, discovery: { enabled: true, asLocalhost: false }});
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('product');
        const result = await contract.evaluateTransaction('queryProduct', req.body.id);
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        res.status(200).json({response: result.toString()});
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({error: error});
    }
});

app.get('/api/queryAllProducts/', async function(req, res) {
    try {
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: user, discovery: { enabled: true, asLocalhost: false } });
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('product');
        const result = await contract.evaluateTransaction('queryAllProducts');
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        res.status(200).json({response: result.toString()});
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({error: error});
    }

});

app.put('/api/editProduct/', async function(req, res) {
    try {
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: user, discovery: { enabled: true, asLocalhost: false } });
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('product');
        await contract.submitTransaction('editProduct', req.body.sku);
        console.log('Transaction has been submitted');
        res.send('Transaction has been submitted');
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({error: error});
    }
});

app.delete('/api/deleteProduct/', async function(req, res) {
    try {
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: user, discovery: { enabled: true, asLocalhost: false } });
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('product');
        await contract.submitTransaction('deleteProduct', req.body.id);
        console.log('Transaction has been submitted');
        res.send('Transaction has been submitted');
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({error: error});
    }
});

app.post('/api/createOrder/', async function(req, res) {
    try {
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: user, discovery: { enabled: true, asLocalhost: false }});
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('order');
        await contract.submitTransaction('createOrder', req.body.entityId);
        console.log('Transaction has been submitted');
        res.send('Transaction has been submitted');
        await gateway.disconnect();
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({error: error});
    } 
});

app.get('/api/queryOrder/', async function(req, res) {
    try {
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: user, discovery: { enabled: true, asLocalhost: false } });
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('order');
        const result = await contract.evaluateTransaction('queryOrder', req.body.entityId);
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        res.status(200).json({response: result.toString()});
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({error: error});
    }
});

app.get('/api/queryAllOrders/', async function(req, res) {
    try {
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: user, discovery: { enabled: true, asLocalhost: false } });
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('order');
        const result = await contract.evaluateTransaction('queryAllOrders');
        console.log(`Transaction has been evaluated, result is: ${result.toString()}`);
        res.status(200).json({response: result.toString()});
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({error: error});
    }

});

app.put('/api/editOrder/', async function(req, res) {
    try {
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: user, discovery: { enabled: true, asLocalhost: false } });
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('order');
        await contract.submitTransaction('editOrder', req.body.entityId);
        console.log('Transaction has been submitted');
        res.send('Transaction has been submitted');
    } catch (error){
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({error: error});
    }
});

app.delete('/api/deleteOrder/', async function(req, res) {
    try {
        const gateway = new Gateway();
        await gateway.connect(ccpPath, { wallet, identity: user, discovery: { enabled: true, asLocalhost: false } });
        const network = await gateway.getNetwork('mychannel');
        const contract = network.getContract('order');
        await contract.submitTransaction('deleteOrder', req.body.entityId);
        console.log('Transaction has been submitted');
        res.send('Transaction has been submitted');
    } catch (error) {
        console.error(`Failed to evaluate transaction: ${error}`);
        res.status(500).json({error: error});
    }
});

app.listen(5000);

