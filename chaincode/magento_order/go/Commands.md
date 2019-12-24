Node 1
docker rm -f $(docker ps -aq) && docker volume prune
docker-compose -f docker-compose-simple.yaml up

Node 2
docker exec -it chaincode bash
cd order/go
go build -o order
CORE_CHAINCODE_ID_NAME=mycc:0 CORE_PEER_TLS_ENABLED=false ./order -peer.address peer:7052

Node 3
docker exec -it cli bash
peer chaincode install -p chaincodedev/chaincode/order/go -n mycc -v 0
peer chaincode instantiate -n mycc -v 0 -c '{"Args":[]}' -C myc

//Start Products
peer chaincode invoke -n mycc -c '{"Args":["initOrders"]}' -C myc

/Query Commands
peer chaincode invoke -n mycc -c '{"Args":["queryOrder","28715"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["queryAllOrders"]}' -C myc

//Delete Commands
peer chaincode invoke -n mycc -c '{"Args":["deleteOrder", "28715"]}' -C myc

//Edit Commands
peer chaincode invoke -n mycc -c '{"Args":["editOrder", "28715"]}' -C myc
peer chaincode invoke -n mycc -c '{"Args":["createOrder", "28715"]}' -C myc





