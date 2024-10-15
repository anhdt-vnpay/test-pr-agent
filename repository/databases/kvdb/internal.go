package kvdb

import (
	"github.com/blcvn/corev3-libs/storage/kv"
	"github.com/blcvn/corev4-explorer/appconfig"
)

func NewKvDB(path string, opt any) (*kv.KvDb, error) {
	return kv.NewKvDb(path, appconfig.KvDBType, opt)
}
