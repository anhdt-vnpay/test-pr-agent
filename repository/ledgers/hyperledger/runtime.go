package hyperledger

import (
	"context"

	"github.com/blcvn/corev3-libs/ledger/hyperledger/config"
	"github.com/blcvn/corev3-libs/ledger/hyperledger/listener"
	cmap "github.com/orcaman/concurrent-map/v2"
)

type ledgerRepo struct {
	listenerOption listener.ListenerOption
	shardingMap    cmap.ConcurrentMap[string, ledgerClient]
	processTx      func(shardId string, item any) error
}

func NewLedgerRepository(fCfg config.FabricConfig) iLedgerRepo {
	return &ledgerRepo{}
}

func (l *ledgerRepo) Start(ctx context.Context) {
	l.SetupProcessTx()
	for _, lc := range l.shardingMap.Items() {
		lc.Start(ctx)
	}
}

func (l *ledgerRepo) SetupProcessTx() {
	for _, lc := range l.shardingMap.Items() {
		lc.SetupProcessTx(l.processTx)
	}
}
