package syncer

import (
	"context"
	"github.com/fadyat/speedy/pkg"
)

func applyChangesToConfigFile(_ context.Context, diff *cacheConfigDiff, configPath string) (bool, error) {
	var (
		cacheConfig = NewDefaultCacheConfig(
			WithMasterInfo(diff.masterInfo),
		)
		changed = false
	)

	for _, d := range diff.nodes {
		switch d.state {
		case nodeStateAdded:
			cacheConfig.Nodes[d.id] = NewNodeConfig(d.id, d.host, d.port)
			changed = true
		case nodeStateRemoved:
			changed = true
		case nodeStateSynced:
			cacheConfig.Nodes[d.id] = NewNodeConfig(d.id, d.host, d.port)
		}
	}

	if !changed {
		return false, nil
	}

	if err := pkg.ToYaml(configPath, cacheConfig); err != nil {
		return false, err
	}

	return true, nil
}
