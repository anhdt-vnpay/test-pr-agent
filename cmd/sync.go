package cmd

import (
	"context"

	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/blcvn/corev4-explorer/repository/ledgers/hyperledger"
	"github.com/blcvn/corev4-explorer/services"
)

type ledgerRepo interface {
	Start(ctx context.Context)
}

func sync() {
	ctx := context.Background()
	appconfig.InitConfig()
	ledger := createLedgerRepo()
	syncService := services.NewSyncService(ledger)
	syncService.Start(ctx)
}

func createLedgerRepo() ledgerRepo {
	switch appconfig.NetWorkType {
	case "hyperledger":
		return hyperledger.NewLedgerRepository(appconfig.FabricCfg)
	default:
		return hyperledger.NewLedgerRepository(appconfig.FabricCfg)
	}
}
