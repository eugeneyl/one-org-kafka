Node 1
docker rm -f $(docker ps -aq) && docker volume prune
docker-compose -f docker-compose-simple.yaml up

Node 2
docker exec -it chaincode bash
cd token/go
go build -o token
CORE_CHAINCODE_ID_NAME=mycc:0 CORE_PEER_TLS_ENABLED=false ./token -peer.address peer:7052

Node 3
docker exec -it cli bash
peer chaincode install -p chaincodedev/chaincode/token/go -n mycc -v 0
peer chaincode instantiate -n mycc -v 0 -c '{"Args":[]}' -C myc

//Start token
peer chaincode invoke -n mycc -c '{"Args":["initToken","Test","tst","{value}"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["initToken","Test","tst","10000"]}' -C myc

{"Address":"0x95effD0B1d4db2a6B0E04eC164fDA777D795fBAc","PrivateKey":"e5b1d1be37abb6f98330532a31cf68e1b1b07f3462d05f49d12ba27ed816a434"}

//Edit Token
peer chaincode invoke -n mycc -c '{"Args":["mintToken","{privateKey}","{value}"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["burnToken","{privateKey}","{value}"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["mintToken","e5b1d1be37abb6f98330532a31cf68e1b1b07f3462d05f49d12ba27ed816a434","200"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["burnToken","e5b1d1be37abb6f98330532a31cf68e1b1b07f3462d05f49d12ba27ed816a434","300"]}' -C myc

//Get token details
peer chaincode invoke -n mycc -c '{"Args":["queryTotalAmount"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["queryTokenName"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["queryTokenSymbol"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["queryReserve"]}' -C myc

//Create account
peer chaincode invoke -n mycc -c '{"Args":["createAccount","{value}"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["createAccount","e5b1d1be37abb6f98330532a31cf68e1b1b07f3462d05f49d12ba27ed816a434","250"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["createAccount","e5b1d1be37abb6f98330532a31cf68e1b1b07f3462d05f49d12ba27ed816a434","650"]}' -C myc

Address: 0x02f2FE32dcB37f426DCf65eEA5E3d133b530c57b
Private Key: a5fbb8d694944cdc0402af1f82632be7b848967d774f9ddafccafbaf05f0281b

Address: 0x5FC575a338D19C205F06a8C9162A833942439e7d
Private Key: af18c62017fd54748002af30f275843acc2332f40f6f5799128f19d3f235653c

//Account details
peer chaincode invoke -n mycc -c '{"Args":["balanceOf","{address}"]}' -C myc

//Account transactions
peer chaincode invoke -n mycc -c '{"Args":["transfer","{addressFrom}","{privateKey}","{addressTo}","{value}"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["buyToken","{address}","{privateKey}","{value}"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["sellToken","{address}","{privateKey}","{value}"]}' -C myc

//Approval transactions
peer chaincode invoke -n mycc -c '{"Args":["approve","{ownerAddress}","{privateKey}","{spenderAddress}","{value}"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["increaseAllowance","{ownerAddress}","{privateKey}","{spenderAddress}","{value}"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["decreaseAllowance","{ownerAddress}","{privateKey}","{spenderAddress}","{value}"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["transferFrom","{spenderAddress}","{privateKey}","{ownerAddress}","{toAddress}""{value}"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["approve","0x02f2FE32dcB37f426DCf65eEA5E3d133b530c57b","a5fbb8d694944cdc0402af1f82632be7b848967d774f9ddafccafbaf05f0281b","0x5FC575a338D19C205F06a8C9162A833942439e7d","40"]}' -C myc

//Approval details
peer chaincode invoke -n mycc -c '{"Args":["allowance","{ownerAddress}","{spenderAddress}"]}' -C myc

peer chaincode invoke -n mycc -c '{"Args":["allowance","0x02f2FE32dcB37f426DCf65eEA5E3d133b530c57b","0x5FC575a338D19C205F06a8C9162A833942439e7d"]}' -C myc


