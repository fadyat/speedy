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

func diff(current, desired []*NodeConfig) map[string]*nodeDiff {
	var clusterState = make(map[string]*nodeDiff)
	for _, n := range current {
		clusterState[n.ID] = newDiff(n, nodeStateRemoved)
	}

	for _, n := range desired {
		if _, ok := clusterState[n.ID]; !ok {
			clusterState[n.ID] = newDiff(n, nodeStateAdded)
			continue
		}

		clusterState[n.ID].state = nodeStateSynced
	}

	return clusterState
}
