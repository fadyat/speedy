## Leader election

**Election algorithms** are essential in distributed data systems to have a consensus between nodes, of which node is
labeled as the master node. It is responsible for coordination between nodes or centralized lookups.

### Bully algorithm

It follows all the assumptions discussed above in the Election Algorithm.

Bully algorithm is one of the simplest algorithms for election, in which we enumerate all the processes running in the
system and pick the one with the highest ID as the coordinator.

> In this algorithm, each process has a unique ID
> and every process knows the corresponding ID and IP address of every other process.

There are two situations for election:

- The system is newly initialized, so there is no leader
- One of the nodes notices that the leader is down.

Any process in the system can initiate this algorithm for leader election. Thus, we can have concurrent ongoing
elections.

#### Messages in Bully Algorithm

There can be three types of messages that processes exchange with each other in the bully algorithm:

- **election** - initialized by a process to start an election.
- **ok** - liveness message to a process that has initiated an election.
- **coordinator** - sent by a process to declare itself as the coordinator.

#### Bully Algorithm Steps

- A process with ID `i` initiates an election.
- It sends an election message to all the processes with ID greater than `i`.
- Any process upon receiving the election message returns an `ok` to its predecessor and starts an election
  of its own by sending election to higher ID processes.
- If it receives no `ok` messages, it knows that it's the highest ID process. It thus sends a `coordinator` message to
  all other processes.
- If it receives an `ok` message, it knows that there are higher ID processes in the system. It waits for a
  `coordinator` message from one of them.

`O(n^2)` messages are sent in the worst case.

A lower ID process need to always heartbeat to the coordinator to check if it is still alive, if not, it can initiate
the election algorithm.

The main downside of the bully algorithm is that if the highest-ranked node goes down frequently, it will re-claim
leadership every time it comes back online, causing unnecessary reelections. Synchronization of messages can also be
difficult to maintain, especially as the cluster gets larger and physically distributed.

### Raft

todo: add raft

### Resources

- https://itnext.io/lets-implement-a-basic-leader-election-algorithm-using-go-with-rpc-6cd012515358