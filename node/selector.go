package node

var (
	currentNodeIdx = uint64(0)
)

func oneAfterAnotherNodeSelector(nc *NodesConfig) *Node {
	currentNodeIdx %= uint64(len(nc.keys))
	node := nc.nodes[nc.keys[currentNodeIdx]]
	currentNodeIdx = (currentNodeIdx + 1) % uint64(len(nc.keys))
	return node
}
