## speedy

Distributed cache implemented in Go.

Features:

- Client-side consistent hashing
- Distributed leader election via Bully algorithm
- LRU eviction policy
- gRPC for communication between nodes

<img src="./docs/images/speedy.png" width="700">

Design docs with research of potential candidates:

- [Sharding algorithms](./docs/sharding.md)
- [Leader election](./docs/leader-election.md)
- [Consensus](./docs/consensus.md)
- [Leader election vs consensus](./docs/leader-election-vs-consensus.md)

Final presentation: [tap](./docs/distributed-cache.pdf)
