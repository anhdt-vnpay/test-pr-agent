package hyperledger

import (
	"context"

	"github.com/blcvn/corev3-libs/ledger/hyperledger/helper"
)

type iLedgerRepo interface {
	Start(ctx context.Context)
	SetupProcessTx()
}

type ledgerClient interface {
	Start(ctx context.Context)
	SetupProcessTx(func(shardId string, item any) error)
}

type stateHelper interface {
	GetBlockHeightOfShard(shardId string) uint64
	UpdateBlockHeightOfShard(shardID string, blockHeight uint64) error
}

type transform interface {
	TransformBlockStateProto(ret *helper.TxState, shardId string, blockNumber uint64, txIndex, txCount int) (any, error)
}

type dbRepo interface {
	SaveBlockAndRawTxs(ctx context.Context, block any, rawTxs []any) error
}
