# A Fabric Network deployed on 4 Nodes
A Fabric Network of 1 Orderer with Kafka, 1 Organization with 3 peers, deployed on 3 nodes.

This setup uses docker swarm. If you want to use extra_hosts, call `git checkout extra_hosts`

The set up of the nodes are as followed: 

| Node | Zookeeper | Kafka | Orderer | Peer | CLI |
| --- | --- | --- | --- | --- | --- |
| 1 | zookeeper0 | kafka0, kafka1 | orderer0.frogfrogjump.com | peer0.org1.frogfrogjump.com|cli |
| 2 ||  | | peer1.org1.frogfrogjump.com|cli |
| 3 | | | | peer2.org1.frogfrogjump.com|cli |

