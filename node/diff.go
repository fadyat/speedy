package node

import "github.com/fadyat/speedy/api"

type nodeDiff struct {
	id    string
	host  string
	port  int
	state nodeState
}

func (d *nodeDiff) toNode() *Node {
	return &Node{
		ID:   d.id,
		Host: d.host,
		Port: d.port,
	}
}

func newNodeDiffFromApiNode(
	n *api.Node,
	state nodeState,
) *nodeDiff {
	return &nodeDiff{
		id:    n.Id,
		host:  n.Host,
		port:  int(n.Port),
		state: state,
	}
}

func newNodeDiffFromNode(
	n *Node,
	state nodeState,
) *nodeDiff {
	return &nodeDiff{
		id:    n.ID,
		host:  n.Host,
		port:  n.Port,
		state: state,
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
