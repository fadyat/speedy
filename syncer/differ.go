package syncer

type nodeDiff struct {
	id    string
	host  string
	port  int
	state nodeState
}

func newDiff(n *NodeConfig, s nodeState) *nodeDiff {
	return &nodeDiff{
		id:    n.ID,
		host:  n.Host,
		port:  n.Port,
		state: s,
	}
}

type nodeState int

const (
	// nodeStateSynced client have connection to the node and it's up-to-date.
	nodeStateSynced nodeState = iota

	// nodeStateAdded client needs to set up a connection to the node.
	nodeStateAdded

	// nodeStateRemoved client needs to close the connection to the node.
	nodeStateRemoved
)

type cacheConfigDiff struct {
	nodes      map[string]*nodeDiff
	masterInfo *NodeConfig
}

func diff(current, desired *CacheConfig) *cacheConfigDiff {
	var nodes = make(map[string]*nodeDiff)
	for _, n := range current.Nodes {
		nodes[n.ID] = newDiff(n, nodeStateRemoved)
	}

	for _, n := range desired.Nodes {
		if _, ok := nodes[n.ID]; !ok {
			nodes[n.ID] = newDiff(n, nodeStateAdded)
			continue
		}

		nodes[n.ID].state = nodeStateSynced
	}

	return &cacheConfigDiff{
		nodes:      nodes,
		masterInfo: desired.MasterInfo,
	}
}
