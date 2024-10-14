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
	sql_INIT_TABLE_ONCHAIN_TRANSACTION = `CREATE TABLE %s (
		TX_HASH VARCHAR2(255),
		TRACE_NO VARCHAR2(255),
		ORIGIN_TRACE_NO VARCHAR2(255),
		SENDER_ID VARCHAR2(255),
		RECEIVER_ID VARCHAR2(255),
		FEE_ACCOUNT_ID VARCHAR2(255),
		SENDER_AMOUNT NUMBER,
		RECEIVER_AMOUNT NUMBER,
		FEE_AMOUNT NUMBER,
		BLOCK_NUM NUMBER,
		NETWORK_NAME VARCHAR2(255),
		CHANNEL_NAME VARCHAR2(255),
		TX_TYPE VARCHAR2(255),
		TIME_ONCHAIN TIMESTAMP,
		TIME_OFFCHAIN TIMESTAMP NOT NULL,
		PRE_DAY VARCHAR2(255),
		CURRENT_DAY VARCHAR2(255),
		NEXT_DAY VARCHAR2(255),
		END_OF_MONTH VARCHAR2(255),
		TASK_TRANSFORM_ID VARCHAR2(255)
	)
	PARTITION BY RANGE (TIME_OFFCHAIN)
	INTERVAL(NUMTODSINTERVAL(1, 'DAY')) 
	(
		PARTITION p_0 VALUES LESS THAN (TIMESTAMP '2023-01-01 00:00:00')
	)`

	sql_INIT_INDEX_1_ONCHAIN_TRANSACTION = `CREATE UNIQUE INDEX IDX_%s_TRACENO ON %s.%s (TRACE_NO)`

	sql_INIT_INDEX_2_ONCHAIN_TRANSACTION = `CREATE INDEX IDX_%s_TIME_OFFCHAIN ON %s.%s (TIME_OFFCHAIN)`

	sql_INSERT_ONCHAIN_TRANSACTION = `INSERT INTO %s (
		TX_HASH, TRACE_NO, ORIGIN_TRACE_NO,
		SENDER_ID, RECEIVER_ID, FEE_ACCOUNT_ID,
		SENDER_AMOUNT, RECEIVER_AMOUNT, FEE_AMOUNT,
		BLOCK_NUM, NETWORK_NAME, CHANNEL_NAME, TX_TYPE,
		TIME_ONCHAIN, TIME_OFFCHAIN, TASK_TRANSFORM_ID
	) VALUES (
		:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14, :15, :16)`
)

func (odb *OracleDB) InitTableOnchainTransaction(ctx context.Context) error {
	tableName := appconfig.OnchainTableNameMapping["onchain_transaction"]

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

	sql2 := fmt.Sprintf(sql_INIT_TABLE_ONCHAIN_TRANSACTION, tableName)
	_, err = odb.Client.ExecContext(ctx, sql2)
	if err != nil {
		odb.logger.Errorf("Failed to create table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("table %s is created successfully", tableName)

	sql3 := fmt.Sprintf(sql_INIT_INDEX_1_ONCHAIN_TRANSACTION, tableName, odb.Schema, tableName)
	_, err = odb.Client.ExecContext(ctx, sql3)
	if err != nil {
		odb.logger.Errorf("Failed to index of table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("index of %s is created successfully", tableName)

	sql4 := fmt.Sprintf(sql_INIT_INDEX_2_ONCHAIN_TRANSACTION, tableName, odb.Schema, tableName)
	_, err = odb.Client.ExecContext(ctx, sql4)
	if err != nil {
		odb.logger.Errorf("Failed to index of table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("index of %s is created successfully", tableName)

	return nil
}

func (odb *OracleDB) SaveOnchainTransaction(ctx context.Context, dbTransaction *sql.Tx, tx *types.OnchainTransactionsOracle) (int64, error) {
	// Prepare the SQL statement
	tableName := fmt.Sprintf("%s.%s", odb.Schema, appconfig.OnchainTableNameMapping["onchain_transaction"])
	sql := fmt.Sprintf(sql_INSERT_ONCHAIN_TRANSACTION, tableName)

	var rowsAffected int64

	if dbTransaction != nil {
		result, err := dbTransaction.ExecContext(ctx,
			sql,
			tx.Txhash,
			tx.TraceNo,
			tx.OriginTraceNo,
			tx.SenderId,
			tx.ReceiverId,
			tx.FeeAccountId,
			tx.SenderAmount,
			tx.ReceiverAmount,
			tx.FeeAmount,
			tx.BlockNum,
			tx.NetworkName,
			tx.ChannelName,
			tx.TxType,
			tx.TimeOnchain,
			tx.TimeOffchain,
			tx.TransformTaskId,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	} else {

		result, err := odb.Client.ExecContext(ctx,
			sql,
			tx.Txhash,
			tx.TraceNo,
			tx.OriginTraceNo,
			tx.SenderId,
			tx.ReceiverId,
			tx.FeeAccountId,
			tx.SenderAmount,
			tx.ReceiverAmount,
			tx.FeeAmount,
			tx.BlockNum,
			tx.NetworkName,
			tx.ChannelName,
			tx.TxType,
			tx.TimeOnchain,
			tx.TimeOffchain,
			tx.TransformTaskId,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	}

	return rowsAffected, nil
}

func (odb *OracleDB) SaveOnchainTransactionOracle(ctx context.Context, dbTransaction *sql.Tx, tx *types.OnchainTransactionsOracle) (int64, error) {
	// Prepare the SQL statement
	tableName := fmt.Sprintf("%s.%s", odb.Schema, appconfig.OnchainTableNameMapping["onchain_transaction"])
	sql := fmt.Sprintf(sql_INSERT_ONCHAIN_TRANSACTION, tableName)

	var rowsAffected int64

	if dbTransaction != nil {
		result, err := dbTransaction.ExecContext(ctx,
			sql,
			tx.Txhash,
			tx.TraceNo,
			tx.OriginTraceNo,
			tx.SenderId,
			tx.ReceiverId,
			tx.FeeAccountId,
			tx.SenderAmount,
			tx.ReceiverAmount,
			tx.FeeAmount,
			tx.BlockNum,
			tx.NetworkName,
			tx.ChannelName,
			tx.TxType,
			tx.TimeOnchain,
			tx.TimeOffchain,
			tx.TransformTaskId,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	} else {

		result, err := odb.Client.ExecContext(ctx,
			sql,
			tx.Txhash,
			tx.TraceNo,
			tx.OriginTraceNo,
			tx.SenderId,
			tx.ReceiverId,
			tx.FeeAccountId,
			tx.SenderAmount,
			tx.ReceiverAmount,
			tx.FeeAmount,
			tx.BlockNum,
			tx.NetworkName,
			tx.ChannelName,
			tx.TxType,
			tx.TimeOnchain,
			tx.TimeOffchain,
			tx.TransformTaskId,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	}

	return rowsAffected, nil
}
