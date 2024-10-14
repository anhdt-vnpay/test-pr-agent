package oracle

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/blcvn/corev4-explorer/types"
	"github.com/godror/godror"
)

const (
	sql_INIT_TABLE_ACCOUNT_TRANSACTIONS = `CREATE TABLE %s (
		TRACE_NO      VARCHAR(255) ,
		TX_HASH       VARCHAR(255) ,
		ACCOUNT_ID    VARCHAR(255) NOT NULL,
		AMOUNT       NUMBER NOT NULL,
		TIME_ONCHAIN  TIMESTAMP NOT NULL,
		TIME_OFFCHAIN TIMESTAMP NOT NULL,
		NETWORK_NAME  VARCHAR(255) ,
		CHANNEL_NAME  VARCHAR(255),
		BLOCK_NUM     NUMBER NOT NULL,
		PRE_DAY TIMESTAMP,
		CURRENT_DAY TIMESTAMP, 
		NEXT_DAY TIMESTAMP ,
		END_OF_MONTH VARCHAR2(255),
		CREATED_AT TIMESTAMP,
		TASK_TRANSFORM_ID VARCHAR2(255)
	)
	PARTITION BY RANGE (CREATED_AT)
	INTERVAL(NUMTODSINTERVAL(1, 'HOUR')) 
	(
	   PARTITION ACC_TX_0 VALUES LESS THAN (TIMESTAMP '2023-01-01 00:00:00')
	)`

	sql_INIT_INDEXES_ACCOUNT_TRANSACTIONS = `CREATE INDEX IDX_%s_CREATED_AT ON %s.%s (CREATED_AT)`

	sql_INSERT_ACCOUNT_TRANSACTIONS = `INSERT INTO %s (TRACE_NO, TX_HASH, ACCOUNT_ID, AMOUNT, TIME_OFFCHAIN,
		TIME_ONCHAIN, CREATED_AT ,NETWORK_NAME, CHANNEL_NAME, BLOCK_NUM, PRE_DAY, CURRENT_DAY, NEXT_DAY, END_OF_MONTH, TASK_TRANSFORM_ID)
		 VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15)`
)

func (odb *OracleDB) InitTableAccountTransactions(ctx context.Context) error {
	tableName := appconfig.OnchainTableNameMapping["account_transactions"]

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

	sql2 := fmt.Sprintf(sql_INIT_TABLE_ACCOUNT_TRANSACTIONS, tableName)
	_, err = odb.Client.ExecContext(ctx, sql2)
	if err != nil {
		odb.logger.Error("Failed to create table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("table %s is created successfully", tableName)

	sql3 := fmt.Sprintf(sql_INIT_INDEXES_ACCOUNT_TRANSACTIONS, tableName, odb.Schema, tableName)
	_, err = odb.Client.ExecContext(ctx, sql3)
	if err != nil {
		odb.logger.Errorf("Failed to index of table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("indexes of %s are created successfully", tableName)

	return nil
}

func (odb *OracleDB) SaveAccountTxOracle(ctx context.Context, dbTransaction *sql.Tx, accList *types.AccountTxsOracle) (int64, error) {

	// Prepare the SQL statement
	tableName := fmt.Sprintf("%s.%s", odb.Schema, appconfig.OnchainTableNameMapping["account_transactions"])
	sql := fmt.Sprintf(sql_INSERT_ACCOUNT_TRANSACTIONS, tableName)

	var rowsAffected int64

	if dbTransaction != nil {
		result, err := dbTransaction.ExecContext(ctx,
			sql,
			accList.TraceNo,
			accList.Txhash,
			accList.AccountId,
			accList.Amount,
			accList.TimeOffchain,
			accList.TimeOnchain,
			accList.CreatedAt,
			accList.NetworkName,
			accList.ChannelName,
			accList.BlockNum,
			accList.PreDay,
			accList.CurrentDay,
			accList.NextDay,
			accList.EndOfMonth,
			accList.TransformTaskId,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	} else {
		result, err := odb.Client.ExecContext(ctx,
			sql,
			accList.TraceNo,
			accList.Txhash,
			accList.AccountId,
			accList.Amount,
			accList.TimeOffchain,
			accList.TimeOnchain,
			accList.CreatedAt,
			accList.NetworkName,
			accList.ChannelName,
			accList.BlockNum,
			accList.PreDay,
			accList.CurrentDay,
			accList.NextDay,
			accList.EndOfMonth,
			accList.TransformTaskId,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	}

	return rowsAffected, nil
}
