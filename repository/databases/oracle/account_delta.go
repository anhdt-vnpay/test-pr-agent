package oracle

import (
	"context"
	"fmt"

	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/godror/godror"
)

const (
	sql_INIT_TABLE_ACCOUNT_DELTA = `CREATE TABLE %s (
		"ACCOUNT_ID" VARCHAR2(255),
		"AMOUNT" NUMBER,
		"PROCESS_HOUR" TIMESTAMP NOT NULL,
		"CREATED_AT" TIMESTAMP NOT NULL,
		"TIME_CURRENT_DAY" TIMESTAMP,
		"DELTA_HOUR" NUMBER,
		"DELTA_DATE" TIMESTAMP,
		"TASK_DELTA_ID" NUMBER
	)
	PARTITION BY RANGE (PROCESS_HOUR)
	INTERVAL(NUMTODSINTERVAL(1, 'DAY')) 
	(
		PARTITION ACC_DELTA_0 VALUES LESS THAN (TIMESTAMP '2023-01-01 00:00:00')
	)`

	sql_INIT_INDEXES_ACCOUNT_DELTA = `CREATE INDEX IDX_%s_ID_PROCESS_HOUR ON %s.%s (ACCOUNT_ID, PROCESS_HOUR)`

	sql_INSERT_ACCOUNT_DELTA = `INSERT INTO %s (
		"ACCOUNT_ID",
		"AMOUNT",
		"PROCESS_HOUR",
		"TIME_CURRENT_DAY",
		"DELTA_HOUR",
		"DELTA_DATE",
		"TASK_DELTA_ID")
	VALUES(:1, :2, :3, :4, :5, :6, :7)`
)

func (odb *OracleDB) InitTableAccountDelta(ctx context.Context) error {

	tableName := fmt.Sprintf("%s.%s", odb.Schema, appconfig.OnchainTableNameMapping["account_delta"])

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

	sql2 := fmt.Sprintf(sql_INIT_TABLE_ACCOUNT_DELTA, tableName)

	_, err = odb.Client.ExecContext(ctx, sql2)
	if err != nil {
		odb.logger.Errorf("Failed to create table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("table %s is created successfully", tableName)

	sql3 := fmt.Sprintf(sql_INIT_INDEXES_ACCOUNT_DELTA, appconfig.OnchainTableNameMapping["account_delta"], odb.Schema, appconfig.OnchainTableNameMapping["account_delta"])
	_, err = odb.Client.ExecContext(ctx, sql3)
	if err != nil {
		odb.logger.Errorf("Failed to index of table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("indexes of %s are created successfully", tableName)

	return nil
}

func (odb *OracleDB) CalculateAccountDelta(taskId int64) error {
	return nil
}
