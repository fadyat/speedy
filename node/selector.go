package node

import (
	"crypto/rand"
	"math/big"
)

func randomNodeSelector(nc *NodesConfig) *Node {
	if len(nc.keys) == 0 {
		return nil
	}

	var keyIdx, _ = rand.Int(rand.Reader, big.NewInt(int64(len(nc.keys))))
	return nc.nodes[nc.keys[keyIdx.Int64()]]
}
