package oracle

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/blcvn/corev4-explorer/entities"
	"github.com/blcvn/corev4-explorer/types"
	"github.com/godror/godror"
)

const (
	sql_INIT_TABLE_TASK = `CREATE TABLE %s (
		"ID" NUMBER GENERATED by default on null as IDENTITY,
		"TYPE" NUMBER,
		"STATUS" NUMBER,
		"BLOCK_NUM" NUMBER,
		"NETWORK_NAME" VARCHAR(255),
		"CHANNEL_NAME" VARCHAR(255),
		"TASK_TRANSFORM_ID" NUMBER,
		"CREATE_AT" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	sql_INSERT_TASK = `INSERT INTO %s (TYPE, STATUS, BLOCK_NUM, NETWORK_NAME, CHANNEL_NAME)
	VALUES(:1, :2, :3, :4, :5) RETURNING ID INTO :id `

	sql_UPDATE_TASK_TRANSFORM_ID = `UPDATE %s SET "TASK_TRANSFORM_ID" = :1 WHERE "TYPE" = :2 AND "TASK_TRANSFORM_ID" IS NULL`

	sql_UPDATE_TASK_STATUS = `UPDATE %s SET "STATUS" = :1 WHERE "ID" = :2`
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

func (odb *OracleDB) InsertTask(ctx context.Context, dbTransaction *sql.Tx, task *types.TaskOracle) (int64, error) {

	// Prepare the SQL statement
	tableName := fmt.Sprintf("%s.%s", odb.Schema, appconfig.OnchainTableNameMapping["task"])
	sqlExec := fmt.Sprintf(sql_INSERT_TASK, tableName)
	var id int64
	if dbTransaction != nil {
		_, err := dbTransaction.ExecContext(ctx, sqlExec,
			task.Type,
			task.Status,
			task.BlockNumber,
			task.NetworkName,
			task.ChannelName,
			sql.Named("id", sql.Out{Dest: &id}),
		)

		if err != nil {
			return 0, err
		}
		return id, nil
	} else {
		_, err := odb.Client.ExecContext(ctx,
			sqlExec,
			task.Type,
			task.Status,
			task.BlockNumber,
			task.NetworkName,
			task.ChannelName,
			sql.Named("id", sql.Out{Dest: &id}),
		)

		if err != nil {
			return 0, err
		}
		return id, nil
	}
}

func (odb *OracleDB) CreateNewTaskTransform(taskStatus int32) (*entities.Task, error) {
	var err error
	var taskId int64
	tfTask := &entities.Task{
		Type:   int32(entities.TaskTransform),
		Status: taskStatus,
	}
	ctx, cancel := context.WithTimeout(odb.ctx, appconfig.TransactionDBTimeout)
	defer cancel()
	sqlTx, _err := odb.Client.BeginTx(ctx, &sql.TxOptions{})
	if _err != nil {
		return nil, _err
	}

	defer func() {
		if err != nil {
			if rollbackErr := sqlTx.Rollback(); rollbackErr != nil {
				odb.logger.Errorf("[%s] Failed to RollBackTransaction : %s", rollbackErr.Error())
			}
		}
	}()
	taskTfOracle := odb.transform.TasksToTaskOracle([]*entities.Task{tfTask})
	taskId, err = odb.InsertTask(ctx, sqlTx, taskTfOracle)
	tfTask.Id = taskId

	tableName := fmt.Sprintf("%s.%s", odb.Schema, appconfig.OnchainTableNameMapping["task"])
	sql := fmt.Sprintf(sql_UPDATE_TASK_TRANSFORM_ID, tableName)
	res, err := sqlTx.ExecContext(ctx, sql, taskId, entities.TaskSync)
	if err != nil {
		return nil, err
	}

	if count, _ := res.RowsAffected(); count == 0 {
		err = errors.New("no task need transform")
		return nil, err
	}

	err = sqlTx.Commit()
	if err != nil {
		return nil, err
	}

	return tfTask, nil
}
func (odb *OracleDB) CreateNewTaskDelta(taskStatus int32) (*entities.Task, error) {
	return nil, nil
}
func (odb *OracleDB) CreateNewTaskBalance(taskStatus int32) (*entities.Task, error) {
	return nil, nil
}
func (odb *OracleDB) UpdateTasksStatus(taskId int64, taskStatus int32) error {
	return nil
}
