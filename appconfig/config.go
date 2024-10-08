package appconfig

import (
	"time"

	"github.com/blcvn/corev3-libs/ledger/hyperledger/config"
)

var (
	NetWorkType string
	FabricCfg   config.FabricConfig

	TransformInterval    time.Duration
	DeltaTaskScheduler   string
	BalanceTaskScheduler string
)

func InitConfig() {}
