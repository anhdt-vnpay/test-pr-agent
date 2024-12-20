package statehelper

import (
	"encoding/binary"
	"fmt"
	"path"

	"github.com/blcvn/corev3-libs/storage/kv/iterator"
	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/blcvn/corev4-explorer/repository/databases/kvdb"
)

type IkvDb interface {
	Has(key string) (bool, error)
	Get(key string) ([]byte, error)
	Put(key string, value []byte) error
	Delete(key string) error
	Iterator(key string) *iterator.Iterator
	Close() error
}

type stateHelper struct {
	kvdb IkvDb
}

func NewStateHelper() (*stateHelper, error) {
	kdb, err := kvdb.NewKvDB(path.Join(appconfig.KvDBRootPath, "ledger"), nil)
	if err != nil {
		return nil, err
	}
	return &stateHelper{
		kvdb: kdb,
	}, nil
}

func (s *stateHelper) GetBlockHeightOfShard(shardId string) uint64 {
	key := getBlockHeightKey(shardId)
	data, err := s.kvdb.Get(key)
	if err != nil {
		return uint64(0)
	}
	return binary.LittleEndian.Uint64(data)
}

func (s *stateHelper) UpdateBlockHeightOfShard(shardID string, blockHeight uint64) error {
	key := getBlockHeightKey(shardID)
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, uint64(blockHeight))
	return s.kvdb.Put(key, b)
}

func getBlockHeightKey(shardID string) string {
	return fmt.Sprintf("LATEST_HEIGHT_%s", shardID)
}
