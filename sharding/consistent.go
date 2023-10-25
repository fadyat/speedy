package sharding

const (
	defaultVNodesCount = 100
)

type consistent struct {

	// maxVNodesCount is the maximum number of virtual nodes per shard.
	maxVNodesCount int

	// vnodesKeys is an ordered slice of virtual nodes keys, which is used to
	// determine which shard a key belongs to.
	//
	// binary search is used to find the shard, which is O(log n).
	vnodesKeys []uint32

	// todo: mb changed
	// vnodesShards is a map of virtual nodes keys to shards.
	//
	// used map for faster lookup instead of slice.
	vnodesShards map[uint32]*Shard

	// todo: ....
}

func NewConsistent(
	shards []*Shard,
	hashFn func(key string) uint32,
) Algorithm {
	return &consistent{}
}

func (c *consistent) GetShard(key string) *Shard {
	//TODO implement me
	panic("implement me")
}

func (c *consistent) RegisterShard(shard *Shard) error {
	//TODO implement me
	panic("implement me")
}

func (c *consistent) DeleteShard(shard *Shard) error {
	//TODO implement me
	panic("implement me")
}

func (c *consistent) GetShards() []*Shard {
	//TODO implement me
	panic("implement me")
}
