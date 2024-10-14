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

	DBType string
)

var (
	OnchainTableNameMapping = map[string]string{
		"block":                "BLOCK",
		"raw_transaction":      "RAW_TRANSACTION",
		"onchain_transaction":  "ONCHAIN_TRANSACTION",
		"account_init":         "ACCOUNT_INIT",
		"account_transactions": "ACCOUNT_TRANSACTIONS",
		"account_delta":        "ACCOUNT_DELTA",
		"account_balance":      "ACCOUNT_BALANCE",
		"task":                 "TASK",
	}
)

func InitConfig() {}
