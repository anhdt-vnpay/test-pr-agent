package appconfig

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/blcvn/corev3-libs/flogging"
	"github.com/blcvn/corev3-libs/ledger/hyperledger/config"
	"github.com/blcvn/corev3-libs/storage/kv"
	"github.com/spf13/viper"
)

var (
	cfg          *viper.Viper
	configLogger = flogging.MustGetLogger("config.config")

	NetWorkType string
	FabricCfg   *config.FabricConfig

	TransformInterval    time.Duration
	DeltaTaskScheduler   string
	BalanceTaskScheduler string

	KvDBType     kv.KvDbType = kv.GoLevelDb
	KvDBRootPath string      = "./tmp"

	TransactionDBTimeout = 60 * time.Second

	MaxQueryGoroutine          = int(1)
	BlockMainQueueSize         = int(100)
	BlockExceptionQueueSize    = int(100)
	NumberBlockParsingWorker   = int(100)
	BlockLoadExceptionInterval = 500 * time.Millisecond
	BlockMinExceptionLoadTime  = 500 * time.Millisecond
	UpdateMissing              = true
	BlockSyncTime              = 10 * time.Second
	PeerGRPCTimeout            = 3 * time.Second
)

var (
	OnchainDbConfig         = make(map[string]string)
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

func InitConfig() {
	loadConfig()
	loadDBConfig()
	loadFabricConfig()
	loadKvDBConfig()
	loadLedgerConfig()
}

func loadConfig() {
	file := "config.yaml"
	key := "CONFIG_FILE"
	if value := os.Getenv(key); value != "" {
		file = value
	}

	cfg = viper.New()
	cfg.SetConfigType("yaml")
	fileName := path.Base(file)
	cfg.SetConfigName(fileName)
	folder := path.Dir(file)
	cfg.AddConfigPath(folder)
	cfg.AddConfigPath("./config/")
	cfg.AddConfigPath("../config/")
	configLogger.Infof("Config file: %s", file)
	err := cfg.ReadInConfig()
	if err != nil {
		configLogger.Fatalf("error on parsing configuration file: %s", err.Error())
	}
}

func loadKvDBConfig() {
	cfg := getConfig()
	eCfg := cfg.GetStringMap("Storage")
	if val, ok := eCfg[strings.ToLower("Type")]; ok {
		KvDBType = kv.KvDbType(val.(int))
	}
	if val, ok := eCfg[strings.ToLower("RootFolder")]; ok {
		KvDBRootPath = val.(string)
	}
}

func loadDBConfig() {
	var err error
	config := getConfig()
	OnchainDbConfig = config.GetStringMapString("DB.onchain")

	val := config.GetStringMapString("DB.onchain.table-name-mapping")
	if len(val) != 0 {
		OnchainTableNameMapping = val
	}

	if val, ok := OnchainDbConfig["transaction-timeout"]; ok {
		TransactionDBTimeout, err = time.ParseDuration(val)
		if err != nil {
			panic("DB.onchain config: invalid TransactionDBTimeout config: " + val)
		}
		configLogger.Infof("DB.onchain.transaction-timeout: %s", TransactionDBTimeout)
	}
}

func loadFabricConfig() {
	fabricConfigFile := os.Getenv("FABRIC_CONFIG_FILE")
	if fabricConfigFile == "" {
		fmt.Println("FABRIC_CONFIG_FILE are not set")
		return
	}

	fmt.Println("Flow: fabricConfigFile :  ", fabricConfigFile)

	var err error
	FabricCfg, err = getConfigFromFile(fabricConfigFile)
	if err != nil {
		configLogger.Fatalf("Cannot load fabric config")
		return
	}
}

func getConfigFromFile(configPath string) (*config.FabricConfig, error) {
	fabricConfig := &config.FabricConfig{}
	configFile, err := os.Open(configPath)

	configLogger.Infof("Load config from %s", configPath)
	defer configFile.Close()
	if err != nil {
		configLogger.Fatalf("Error while open config file %s: %s", configPath, err.Error())
		return nil, err
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(fabricConfig); err != nil {
		configLogger.Fatalf("Error while parsing config file %s: %s", configPath, err.Error())
		return nil, err
	}

	return fabricConfig, nil
}

func loadLedgerConfig() {
	config := getConfig()
	cfgMap := config.GetStringMap("Ledger")

	if val, ok := cfgMap[strings.ToLower("UpdateMissing")]; ok {
		UpdateMissing = val.(bool)
		configLogger.Infof("Ledger.UpdateMissing: %t", UpdateMissing)
	}

	if val, ok := cfgMap[strings.ToLower("BlockSyncTime")]; ok {
		var err error
		BlockSyncTime, err = time.ParseDuration(val.(string))
		if err != nil {
			panic("Ledger config: invalid BlockSyncTime config: " + val.(string))
		}
		configLogger.Infof("Ledger.BlockSyncTime: %s", BlockSyncTime)
	}
	if val, ok := cfgMap[strings.ToLower("MaxQueryGoroutine")]; ok {
		MaxQueryGoroutine = val.(int)
		configLogger.Infof("Ledger.MaxQueryGoroutine: %d", MaxQueryGoroutine)
	}

	if val, ok := cfgMap[strings.ToLower("PeerGRPCTimeout")]; ok {
		var err error
		PeerGRPCTimeout, err = time.ParseDuration(val.(string))
		if err != nil {
			panic("Ledger config: invalid PeerGRPCTimeout config: " + val.(string))
		}
		configLogger.Infof("Ledger.PeerGRPCTimeout: %s", PeerGRPCTimeout)
	}

	if val, ok := cfgMap[strings.ToLower("BlockMainQueueSize")]; ok {
		BlockMainQueueSize = val.(int)
		configLogger.Infof("Ledger.BlockMainQueueSize: %d", BlockMainQueueSize)
	}

	if val, ok := cfgMap[strings.ToLower("BlockExceptionQueueSize")]; ok {
		BlockExceptionQueueSize = val.(int)
		configLogger.Infof("Ledger.BlockExceptionQueueSize: %d", BlockExceptionQueueSize)
	}

	if val, ok := cfgMap[strings.ToLower("NumberBlockParsingWorker")]; ok {
		NumberBlockParsingWorker = val.(int)
		configLogger.Infof("Ledger.NumberBlockParsingWorker: %d", NumberBlockParsingWorker)
	}

	if val, ok := cfgMap[strings.ToLower("BlockLoadExceptionInterval")]; ok {
		var err error
		BlockLoadExceptionInterval, err = time.ParseDuration(val.(string))
		if err != nil {
			panic("Ledger config: invalid BlockLoadExceptionInterval config: " + val.(string))
		}
		configLogger.Infof("Ledger.BlockLoadExceptionInterval: %s", BlockLoadExceptionInterval)
	}

	if val, ok := cfgMap[strings.ToLower("BlockMinExceptionLoadTime")]; ok {
		var err error
		BlockMinExceptionLoadTime, err = time.ParseDuration(val.(string))
		if err != nil {
			panic("Ledger config: invalid BlockMinExceptionLoadTime config: " + val.(string))
		}
		configLogger.Infof("Ledger.BlockMinExceptionLoadTime: %s", BlockMinExceptionLoadTime)
	}
}

func getConfig() *viper.Viper {
	return cfg
}
