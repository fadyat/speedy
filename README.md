## speedy

Distributed cache implemented in Go.

Features:

- Client-side consistent hashing
- Distributed leader election via Raft (maybe Paxos/Bully)
- gRPC for communication between nodes

<img src="./docs/images/speedy.png" width="700">

Docs:

- [Sharding algorithms](./docs/sharding.md)
- [Leader election](./docs/leader-election.md)
- [Consensus](./docs/consensus.md)
- [Leader election vs consensus](./docs/leader-election-vs-consensus.md)

Final presentation: [tap](./docs/distributed-cache.pdf)
