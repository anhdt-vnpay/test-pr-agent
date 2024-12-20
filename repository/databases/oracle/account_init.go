package oracle

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/blcvn/corev4-explorer/types"
	"github.com/godror/godror"
	_ "github.com/godror/godror"
)

const (
	sql_INIT_TABLE_ACCOUNT_INIT = `CREATE TABLE %s (
		"TRACE_NO" VARCHAR2(255),
		"TX_HASH" VARCHAR2(255),
		"ACCOUNT_ID" VARCHAR2(255) PRIMARY KEY,
		"AMOUNT" NUMBER,
		"BLOCKADE_AMOUNT" NUMBER,
		"BALANCE" NUMBER,
		"TYPE" NUMBER,
		"STATE" NUMBER,
		"TIME_ONCHAIN" TIMESTAMP,
		"TIME_OFFCHAIN" TIMESTAMP,
		"NETWORK_NAME" VARCHAR2(255),
		"CHANNEL_NAME" VARCHAR2(255),
		"BLOCK_NUM" NUMBER,
		"TRANSFORM_ID_TASK" NUMBER
	)`
	sql_INSERT_ACCOUNT_INIT = `INSERT INTO %s (
		TRACE_NO, TX_HASH, ACCOUNT_ID,
		AMOUNT,BLOCKADE_AMOUNT,BALANCE, TYPE, STATE,
		TIME_ONCHAIN, TIME_OFFCHAIN, NETWORK_NAME, CHANNEL_NAME, BLOCK_NUM, TRANSFORM_ID_TASK
	) VALUES (
		:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14)`
	SQL_CHECK_TABLE_EXISTED = `SELECT 1 FROM USER_TABLES WHERE TABLE_NAME = '%s'`
)

func (odb *OracleDB) InitTableAccountInit(ctx context.Context) error {
	tableName := appconfig.OnchainTableNameMapping["account_init"]

	sql1 := fmt.Sprintf(SQL_CHECK_TABLE_EXISTED, tableName)
	result, err := odb.Client.QueryContext(ctx, sql1, godror.FetchArraySize(1))
	if err != nil {
		return err
	}

	defer result.Close()

	var existed int
	for result.Next() {
		if err := result.Scan(&existed); err == nil {
			err := fmt.Errorf("table %s existed", tableName)
			odb.logger.Warn(err)
			return err
		}
	}

	sql2 := fmt.Sprintf(sql_INIT_TABLE_ACCOUNT_INIT, tableName)
	_, err = odb.Client.ExecContext(ctx, sql2)
	if err != nil {
		odb.logger.Errorf("Failed to index of table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("table %s is created successfully", tableName)

	// sql3 := fmt.Sprintf(sql_INIT_INDEXES_ACCOUNT_INIT, odb.Schema, tableName)
	// _, err = odb.Client.ExecContext(ctx, sql3)
	// if err != nil {
	// 	odb.logger.Error("Failed to create table %s : %v", tableName, err)
	// 	return err
	// }

	// odb.logger.Infof("indexes of %s are created successfully", tableName)

	return nil
}

func (odb *OracleDB) SaveAccountOracle(ctx context.Context, dbTransaction *sql.Tx, acc *types.AccountsOracle) (int64, error) {
	// Prepare the SQL statement
	tableName := fmt.Sprintf("%s.%s", odb.Schema, appconfig.OnchainTableNameMapping["account_init"])
	sql := fmt.Sprintf(sql_INSERT_ACCOUNT_INIT, tableName)

	var rowsAffected int64

	if dbTransaction != nil {
		result, err := dbTransaction.ExecContext(ctx,
			sql,
			acc.TraceNo,
			acc.Txhash,
			acc.AccountId,
			acc.Amount,
			acc.BlockadeAmount,
			acc.Balance,
			acc.Type,
			acc.State,
			acc.TimeOnchain,
			acc.TimeOffchain,
			acc.NetworkName,
			acc.ChannelName,
			acc.BlockNum,
			acc.TransformTaskId,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	} else {

		result, err := odb.Client.ExecContext(ctx,
			sql,
			acc.TraceNo,
			acc.Txhash,
			acc.AccountId,
			acc.Amount,
			acc.BlockadeAmount,
			acc.Balance,
			acc.Type,
			acc.State,
			acc.TimeOnchain,
			acc.TimeOffchain,
			acc.NetworkName,
			acc.ChannelName,
			acc.BlockNum,
			acc.TransformTaskId,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	}

	return rowsAffected, nil
}
