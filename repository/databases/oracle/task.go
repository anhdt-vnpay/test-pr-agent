package oracle

import (
	"context"
	"fmt"

	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/blcvn/corev4-explorer/types"
	"github.com/godror/godror"
)

const (
	sql_INIT_TABLE_TASK = `CREATE TABLE %s (
		"ID" NUMBER PRIMARY KEY,
		"TYPE" NUMBER,
		"STATUS" NUMBER,
		"BLOCK_NUM" NUMBER,
		"NETWORK_NAME" VARCHAR(255),
		"CHANNEL_NAME" VARCHAR(255),
		"CREATE_AT" TIMESTAMP
	)`
	sql_INSERT_TASK = `INSERT INTO %s (ID, TYPE, STATUS, BLOCK_NUM, NETWORK_NAME, CHANNEL_NAME, CREATE_AT)
	VALUES(:1, :2, :3, :4, :5, :6, :7)`
)

func (odb *OracleDB) InitTableTask(ctx context.Context) error {
	tableName := appconfig.OnchainTableNameMapping["task"]

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

	sql2 := fmt.Sprintf(sql_INIT_TABLE_TASK, tableName)
	_, err = odb.Client.ExecContext(ctx, sql2)
	if err != nil {
		odb.logger.Errorf("Failed to create table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("table %s is created successfully", tableName)

	return nil
}

func (odb *OracleDB) InsertTask(ctx context.Context, task types.TaskOracle) error {

	// Prepare the SQL statement
	tableName := fmt.Sprintf("%s.%s", odb.Schema, appconfig.OnchainTableNameMapping["task_manager"])
	sql := fmt.Sprintf(sql_INSERT_TASK, tableName)

	_, err := odb.Client.ExecContext(ctx,
		sql,
		task.Id,
		task.Type,
		task.Status,
		task.BlockNumber,
		task.NetworkName,
		task.ChannelName,
		task.CreateAt,
	)

	if err != nil {
		return err
	}

	return nil
}
