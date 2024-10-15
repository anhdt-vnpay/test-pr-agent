package hyperledger

import (
	"context"

	"github.com/blcvn/corev4-explorer/entities"

	"github.com/blcvn/corev3-libs/ledger/hyperledger/listener"
)

type iLedgerRepo interface {
	Start(ctx context.Context)
}

type ledgerClient interface {
	Start(ctx context.Context)
	SetupProcessBlockJob(func(key string, blockEvent *listener.CommonBlockWrapper) (*listener.CommonBlockWrapper, error))
}

type dbRepo interface {
	SaveBlockAndRawTxs(block *entities.Block, rawTxs []*entities.RawTransaction) error
}

type stateHelper interface {
	GetBlockHeightOfShard(shardId string) uint64
	UpdateBlockHeightOfShard(shardID string, blockHeight uint64) error
}
