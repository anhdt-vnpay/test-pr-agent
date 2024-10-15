package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/blcvn/corev4-explorer/entities"
	"github.com/blcvn/corev4-explorer/types"
	"github.com/godror/godror"
)

const (
	sql_INIT_TABLE_BLOCK = `CREATE TABLE %s (
		"BLOCK_NUM" NUMBER,
		"DATA_HASH" VARCHAR2(255),
		"PRE_HASH" VARCHAR2(255),
		"TX_COUNT" NUMBER,
		"BLOCK_TIME" TIMESTAMP,
		"PREV_BLOCK_HASH" VARCHAR2(255),
		"BLOCK_HASH" VARCHAR2(255),
		"CHANNEL_GENESIS_HASH" VARCHAR2(255),
		"BLOCKSIZE" NUMBER,
		"NETWORK_NAME" VARCHAR2(255) NOT NULL,
		"CHANNEL_NAME" VARCHAR2(255) NOT NULL
	) PARTITION BY LIST(NETWORK_NAME, CHANNEL_NAME) (%s)`

	sql_INIT_INDEXES_BLOCK = `CREATE UNIQUE INDEX IDX_%s ON %s.%s (BLOCK_NUM, NETWORK_NAME, CHANNEL_NAME)`

	sql_INIT_TABLE_RAW_TRANSACTION = `CREATE TABLE %s (
		"PAYLOAD" BLOB,
		"TX_HASH" VARCHAR2(255),
		"BLOCK_TIME" TIMESTAMP,
		"CHAINCODE_NAME" VARCHAR2(255),
		"VALIDATION_CODE" VARCHAR2(255),
		"CHAINCODE_PROPOSAL_INPUT" VARCHAR2(255),
		"BLOCK_NUM" NUMBER,
		"NETWORK_NAME" VARCHAR2(255),
		"CHANNEL_NAME" VARCHAR2(255)
	) PARTITION BY LIST(NETWORK_NAME, CHANNEL_NAME) (%s)`

	sql_INIT_INDEXES_RAW_TRANSACTION = `CREATE UNIQUE INDEX IDX_%s ON %s.%s (TX_HASH, BLOCK_NUM, NETWORK_NAME, CHANNEL_NAME)`

	sql_INSERT_BLOCK = `INSERT INTO %s (BLOCK_NUM, DATA_HASH, PRE_HASH, TX_COUNT, BLOCK_TIME, BLOCK_HASH, CHANNEL_GENESIS_HASH,
		BLOCKSIZE, NETWORK_NAME, CHANNEL_NAME)
	 VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9, :10)`

	sql_INSERT_RAW_TRANSACTION = `INSERT INTO %s (
			"TX_HASH", "BLOCK_TIME", "CHAINCODE_NAME", "VALIDATION_CODE",
			 "CHAINCODE_PROPOSAL_INPUT", "BLOCK_NUM", "NETWORK_NAME", "CHANNEL_NAME", "PAYLOAD"
		) VALUES (:1, :2, :3, :4, :5, :6, :7, :8, :9)`
)

func (odb *OracleDB) InitTableBlock(ctx context.Context, networkChannelPairs [][]string) error {
	tableName := appconfig.OnchainTableNameMapping["block"]

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

	var partitionDefs string
	for _, pair := range networkChannelPairs {
		networkName := strings.ReplaceAll(strings.ToUpper(pair[0]), "-", "")
		channelName := strings.ReplaceAll(strings.ToUpper(pair[1]), "-", "")
		partitionDefs += fmt.Sprintf("PARTITION BLOCK_%s_%s VALUES ('%s', '%s'),", networkName, channelName, pair[0], pair[1])
	}

	sql2 := fmt.Sprintf(sql_INIT_TABLE_BLOCK, tableName, strings.TrimSuffix(partitionDefs, ","))
	_, err = odb.Client.ExecContext(ctx, sql2)
	if err != nil {
		odb.logger.Error("Failed to create table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("table %s is created successfully", tableName)

	sql3 := fmt.Sprintf(sql_INIT_INDEXES_BLOCK, tableName, odb.Schema, tableName)
	_, err = odb.Client.ExecContext(ctx, sql3)
	if err != nil {
		odb.logger.Error(err)
		return err
	}

	odb.logger.Infof("indexes of %s are created successfully", tableName)

	return nil
}

func (odb *OracleDB) InitTableRawTransaction(ctx context.Context, networkChannelPairs [][]string) error {
	tableName := appconfig.OnchainTableNameMapping["raw_transaction"]

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

	var partitionDefs string
	for _, pair := range networkChannelPairs {
		networkName := strings.ReplaceAll(strings.ToUpper(pair[0]), "-", "")
		channelName := strings.ReplaceAll(strings.ToUpper(pair[1]), "-", "")
		partitionDefs += fmt.Sprintf("PARTITION RAW_TX_%s_%s VALUES ('%s', '%s'),", networkName, channelName, pair[0], pair[1])
	}

	sql2 := fmt.Sprintf(sql_INIT_TABLE_RAW_TRANSACTION, tableName, strings.TrimSuffix(partitionDefs, ","))
	_, err = odb.Client.ExecContext(ctx, sql2)
	if err != nil {
		odb.logger.Errorf("Failed to create table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("table %s is created successfully", tableName)

	sql3 := fmt.Sprintf(sql_INIT_INDEXES_RAW_TRANSACTION, tableName, odb.Schema, tableName)
	_, err = odb.Client.ExecContext(ctx, sql3)
	if err != nil {
		odb.logger.Errorf("Failed to index of table %s : %v", tableName, err)
		return err
	}

	odb.logger.Infof("indexes of %s are created successfully", tableName)

	return nil
}

func (odb *OracleDB) SaveBlock(ctx context.Context, dbTransaction *sql.Tx, block *types.BlocksOracle) (int64, error) {

	tableName := fmt.Sprintf("%s.%s", odb.Schema, appconfig.OnchainTableNameMapping["block"])
	sql := fmt.Sprintf(sql_INSERT_BLOCK, tableName)

	var rowsAffected int64

	if dbTransaction != nil {
		result, err := dbTransaction.ExecContext(ctx, sql,
			block.Blocknum,
			block.Datahash,
			block.Prehash,
			block.Txcount,
			block.BlockTime,
			block.Blockhash,
			block.ChannelGenesisHash,
			block.Blksize,
			block.NetworkName,
			block.ChannelName,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	} else {

		result, err := odb.Client.ExecContext(ctx,
			sql,
			block.Blocknum,
			block.Datahash,
			block.Prehash,
			block.Txcount,
			block.BlockTime,
			block.Blockhash,
			block.ChannelGenesisHash,
			block.Blksize,
			block.NetworkName,
			block.ChannelName,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	}

	return rowsAffected, nil
}

func (odb *OracleDB) SaveRawTransaction(ctx context.Context, dbTransaction *sql.Tx, rawTransaction *types.RawTransactionsOracle) (int64, error) {

	// Prepare the SQL statement
	tableName := fmt.Sprintf("%s.%s", odb.Schema, appconfig.OnchainTableNameMapping["raw_transaction"])
	sql := fmt.Sprintf(sql_INSERT_RAW_TRANSACTION, tableName)

	var rowsAffected int64

	if dbTransaction != nil {
		result, err := dbTransaction.ExecContext(ctx,
			sql,
			rawTransaction.Txhash,
			rawTransaction.BlockTime,
			rawTransaction.ChaincodeName,
			rawTransaction.ValidationCode,
			rawTransaction.ChaincodeProposalInput,
			rawTransaction.BlockNum,
			rawTransaction.NetworkName,
			rawTransaction.ChannelName,
			rawTransaction.Payload,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	} else {

		result, err := odb.Client.ExecContext(ctx,
			sql,
			rawTransaction.Txhash,
			rawTransaction.BlockTime,
			rawTransaction.ChaincodeName,
			rawTransaction.ValidationCode,
			rawTransaction.ChaincodeProposalInput,
			rawTransaction.BlockNum,
			rawTransaction.NetworkName,
			rawTransaction.ChannelName,
			rawTransaction.Payload,
		)

		if err != nil {
			return 0, err
		}

		rowsAffected, _ = result.RowsAffected()

	}

	return rowsAffected, nil
}

func (odb *OracleDB) SaveTransformData(taskId int64, data any) error {
	return nil
}
func (odb *OracleDB) SaveBlockAndRawTxs(block *entities.Block, rawTxs []*entities.RawTransaction) error {
	var err error
	blockTrace := fmt.Sprintf("%s_%s_%d", block.NetworkName, block.ChannelName, block.Blocknum)
	ctx, cancel := context.WithTimeout(odb.ctx, appconfig.TransactionDBTimeout)
	defer cancel()
	sqlTx, _err := odb.Client.BeginTx(ctx, &sql.TxOptions{})
	if _err != nil {
		return _err
	}

	defer func() {
		if err != nil {
			odb.logger.Errorf("[%s] SaveBlockAndRawTxs error: %s", err.Error())
			odb.logger.Warnf("[%s] Rolling back DB Transaction !", blockTrace)
			if rollbackErr := sqlTx.Rollback(); rollbackErr != nil {
				odb.logger.Errorf("[%s] Failed to RollBackTransaction : %s", blockTrace, rollbackErr.Error())
			}
		}
	}()

	var errChan = make(chan error, 3)
	var wg sync.WaitGroup

	blocksOracle := odb.transform.BlockToBlockOracle([]*entities.Block{block})
	if len(blocksOracle.Blocknum) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _err := odb.SaveBlock(ctx, sqlTx, blocksOracle)
			if _err != nil {
				errChan <- _err
			}
		}()
	}

	rawTxsOracle := odb.transform.RawTxsToRawTxsOracle(rawTxs)
	if len(rawTxsOracle.BlockNum) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _err := odb.SaveRawTransaction(ctx, sqlTx, rawTxsOracle)
			if _err != nil {
				errChan <- _err
			}
		}()
	}

	taskSyncOracle := odb.transform.TasksToTaskOracle([]*entities.Task{
		{
			Type:        int32(entities.TaskSync),
			Status:      int32(entities.TaskDone),
			BlockNumber: block.Blocknum,
			NetworkName: block.NetworkName,
			ChannelName: block.ChannelName,
		},
	})
	if len(taskSyncOracle.BlockNumber) > 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_err := odb.InsertTask(ctx, sqlTx, taskSyncOracle)
			if _err != nil {
				errChan <- _err
			}
		}()
	}

	wg.Wait()
	close(errChan)
	for err := range errChan {
		return err
	}
	err = sqlTx.Commit()
	return err
}
