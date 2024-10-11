package hyperledger

import "github.com/blcvn/corev3-libs/concurrent-map/cmap"

type stateHelper struct {
	syncHeight cmap.ConcurrentMap[string, uint64]
}

func (s *stateHelper) GetBlockHeightOfShard(shardId string) uint64 {
	v, _ := s.syncHeight.Get(shardId)
	return v
}

func (s *stateHelper) UpdateBlockHeightOfShard(shardID string, blockHeight uint64) error {
	s.syncHeight.Set(shardID, blockHeight)
	return nil
}