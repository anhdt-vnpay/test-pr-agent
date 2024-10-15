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
	sql_INIT_TABLE_ACCOUNT_BALANCE = `CREATE TABLE %s (
		"ACCOUNT_ID" VARCHAR2(255) PRIMARY KEY,
		"BALANCE" NUMBER,
		"PROCESS_TIME" TIMESTAMP,
		"NETWORK_NAME" VARCHAR(255),
		"CHANNEL_NAME" VARCHAR(255),
		"TASK_BALANCE_ID" VARCHAR2(255)
		)`

	sql_INSERT_ACCOUNT_BALANCE = `INSERT INTO %s 
		(ACCOUNT_ID, BALANCE, PROCESS_TIME,NETWORK_NAME,CHANNEL_NAME,TASK_BALANCE_ID) VALUES (:1, :2, :3, :4, :5, :6)`
)

func (odb *OracleDB) InitTableAccountBalance(ctx context.Context, networkChannelPairs [][]string) error {
	tableName := appconfig.OnchainTableNameMapping["account_balance"]

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

	// var partitionDefs string
	// for _, pair := range networkChannelPairs {
	// 	networkName := strings.ReplaceAll(strings.ToUpper(pair[0]), "-", "")
	// 	channelName := strings.ReplaceAll(strings.ToUpper(pair[1]), "-", "")
	// 	partitionDefs += fmt.Sprintf("PARTITION BALANCE_%s_%s VALUES ('%s', '%s'),", networkName, channelName, pair[0], pair[1])
	// }

	sql2 := fmt.Sprintf(sql_INIT_TABLE_ACCOUNT_BALANCE, tableName)
	_, err = odb.Client.ExecContext(ctx, sql2)
	if err != nil {
		odb.logger.Errorf("Failed to create table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("table %s is created successfully", tableName)

	return nil
}

func (odb *OracleDB) SaveAccountBalanceOracle(ctx context.Context, dbTransaction *sql.Tx, acc types.AccountsBalanceOracle) (int64, error) {
	// Prepare the SQL statement
	tableName := fmt.Sprintf("%s.%s", odb.Schema, appconfig.OnchainTableNameMapping["account_balance"])
	sql := fmt.Sprintf(sql_INSERT_ACCOUNT_BALANCE, tableName)

	var rowsAffected int64

	if dbTransaction != nil {
		result, err := dbTransaction.ExecContext(ctx,
			sql,
			acc.AccountId,
			acc.Balance,
			acc.ProcessTime,
			acc.NetworkName,
			acc.ChannelName,
			acc.TaskId,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	} else {

		result, err := odb.Client.ExecContext(ctx,
			sql,
			acc.AccountId,
			acc.Balance,
			acc.ProcessTime,
			acc.NetworkName,
			acc.ChannelName,
			acc.TaskId,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	}

	return rowsAffected, nil
}

func (odb *OracleDB) CalculateAccountBalance(taskId int64) error {
	return nil
}
