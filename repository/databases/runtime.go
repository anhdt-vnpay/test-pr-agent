package databases

import (
	"context"

	"github.com/blcvn/corev4-explorer/appconfig"
	connector "github.com/blcvn/corev4-explorer/connectors/database"
	"github.com/blcvn/corev4-explorer/repository/databases/oracle"
)

func NewDataDBRepository(ctx context.Context) (dataDB, error) {
	switch appconfig.OnchainDbConfig["type"] {
	case "ORACLE":
		client, err := connector.InitOracleClient(appconfig.OnchainDbConfig)
		if err != nil {
			return nil, err
		}
		db := oracle.NewOracleDB(client, appconfig.OnchainDbConfig["schema"])
		db.InitTables(ctx)
		return db, nil
	default:
		client, err := connector.InitOracleClient(appconfig.OnchainDbConfig)
		if err != nil {
			return nil, err
		}
		db := oracle.NewOracleDB(client, appconfig.OnchainDbConfig["schema"])
		db.InitTables(ctx)
		return db, nil
	}
}

func NewTaskDBRepository(ctx context.Context) (taskDB, error) {
	switch appconfig.OnchainDbConfig["type"] {
	case "ORACLE":
		client, err := connector.InitOracleClient(appconfig.OnchainDbConfig)
		if err != nil {
			return nil, err
		}
		db := oracle.NewOracleDB(client, appconfig.OnchainDbConfig["schema"])
		db.InitTables(ctx)
		return db, nil
	default:
		client, err := connector.InitOracleClient(appconfig.OnchainDbConfig)
		if err != nil {
			return nil, err
		}
		db := oracle.NewOracleDB(client, appconfig.OnchainDbConfig["schema"])
		db.InitTables(ctx)
		return db, nil
	}
}
