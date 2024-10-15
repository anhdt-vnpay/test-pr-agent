package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/blcvn/corev4-explorer/repository/databases"
	"github.com/blcvn/corev4-explorer/repository/ledgers/hyperledger"
	"github.com/blcvn/corev4-explorer/services"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type ledgerRepo interface {
	Start(ctx context.Context)
}

func sync(metricPort int) {
	ctx := context.Background()
	appconfig.InitConfig()
	ledger := createLedgerRepo(ctx)
	syncService := services.NewSyncService(ledger)
	syncService.Start(ctx)

	router := http.NewServeMux()

	router.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(fmt.Sprintf(":%d", metricPort), router)
	if err != nil {
		log.Fatalf("Error while starting prometheus metrics at port %d. Error: %s", metricPort, err.Error())
	}
}

func createLedgerRepo(ctx context.Context) ledgerRepo {
	dbRepo, err := databases.NewDataDBRepository(ctx)
	if err != nil {
		panic(fmt.Sprintf("Can't create connect to db: %s", err.Error()))
	}
	switch appconfig.NetWorkType {
	case "hyperledger":
		return hyperledger.NewLedgerRepository(appconfig.FabricCfg, dbRepo)
	default:
		return hyperledger.NewLedgerRepository(appconfig.FabricCfg, dbRepo)
	}
}
