package usecases

import (
	"errors"
	"strings"
	"time"

	"github.com/blcvn/corev3-libs/flogging"
	contractpb "github.com/blcvn/corev3-protos/v2/types/smartcontract"
	"github.com/blcvn/corev4-explorer/common"
	"github.com/blcvn/corev4-explorer/entities"
)

type dbRepo interface {
	CalculateAccountDelta(taskId string) error
	CalculateAccountBalance(taskId string) error
	SaveTransformData(taskId int64, onchainTransactions []*entities.OnchainTransaction, accounts []*entities.Account)
	CreateTask(task *entities.Task) error
	GetAllTaskByStatus(taskType int32, taskStatus int32) ([]*entities.Task, error)
	UpdateTasksStatus(taskIds []int64, taskStatus int32) error
	UpdateTaskStatus(taskId int64, taskStatus int32) error
	FindOneById(taskId int64) (*entities.Task, error)
	GetRawTransactionByBlockNumber(blockNumber uint64) ([]*entities.RawTransaction, error)
}

type taskHandler struct {
	dbRepo dbRepo
	logger *flogging.FabricLogger
}

func NewTaskHandler(dbRepo dbRepo) *taskHandler {
	return &taskHandler{
		dbRepo: dbRepo,
	}
}

func (t *taskHandler) PerformTasks(taskId int64, tasksType int32) common.BaseError {
	switch tasksType {
	case int32(entities.TaskTransform):
		return t.transformData(taskId)
	case int32(entities.TaskBalance):
		return t.calculateAccountDelta(taskId)
	default:
		return common.NewUnknownError(errors.New("unknown task"))
	}
}

func (t *taskHandler) transformData(taskId int64) common.BaseError {
	transformTask, err := t.dbRepo.FindOneById(taskId)
	if err != nil {
		t.logger.Errorf("Can not find task by id %d", taskId)
	}
	rawTransactions, err := t.dbRepo.GetRawTransactionByBlockNumber(transformTask.BlockNumber)
	if err != nil {
		t.logger.Errorf("Can not find transaction by blockNumber %d taskId %d", transformTask.BlockNumber, taskId)
	}

	var onchainTxList []*entities.OnchainTransaction

	var accountList []*entities.Account

	for _, rawTransaction := range rawTransactions {

		if strings.HasPrefix(rawTransaction.FunctionName, "Transaction.") {
			newTransactions, err := t.rawTxToTransactions(rawTransaction, transformTask.Id)
			if err != nil {
				return err
			}
			onchainTxList = append(onchainTxList, newTransactions...)

		}

		if strings.HasPrefix(rawTransaction.FunctionName, "Account.") {
			newAccounts, err := t.rawTxToAccounts(rawTransaction, transformTask.Id)
			if err != nil {
				return err
			}
			accountList = append(accountList, newAccounts...)

		}

	}
	t.dbRepo.SaveTransformData(taskId, onchainTxList, accountList)
	return nil
}

func (t *taskHandler) calculateAccountDelta(taskId int64) common.BaseError {
	return nil
}

func (t *taskHandler) rawTxToTransactions(rawTransaction *entities.RawTransaction, taskId int64) ([]*entities.OnchainTransaction, common.BaseError) {
	var onchainTxList []*entities.OnchainTransaction

	response, err := contractpb.NewTransactionResponse([]byte(rawTransaction.TxData))
	if err != nil {
		t.logger.Errorf("Create new transaction response error %s", err.Error())
		return nil, common.NewUnknownError(err)
	}

	for _, attribute := range response.Transaction.Attributes {

		preDay, err := time.Parse(time.RFC3339, attribute.OffchainTime.PreDay)
		if err != nil {
			t.logger.Errorf("Error when parse preday %s", err.Error())
			return nil, common.NewUnknownError(err)

		}
		currentDay, err := time.Parse(time.RFC3339, attribute.OffchainTime.CurrentDay)
		if err != nil {
			t.logger.Errorf("Error when parse currentDay %s", err.Error())
			return nil, common.NewUnknownError(err)

		}
		nextDay, err := time.Parse(time.RFC3339, attribute.OffchainTime.NextDay)
		if err != nil {
			t.logger.Errorf("Error when parse nextDay %s", err.Error())
			return nil, common.NewUnknownError(err)

		}
		onchainTransaction := &entities.OnchainTransaction{
			Txhash:      rawTransaction.Txhash,
			BlockNum:    rawTransaction.BlockNum,
			NetworkName: rawTransaction.NetworkName,
			ChannelName: rawTransaction.ChannelName,
			TimeOnchain: rawTransaction.BlockTime,
			TraceNo:     attribute.TraceNo,
			TxType:      response.FuncType.Enum().String(),

			SenderId:       attribute.SenderId,
			ReceiverId:     attribute.ReceiverId,
			FeeAccountId:   attribute.FeeAccountId,
			SenderAmount:   attribute.SenderAmount,
			ReceiverAmount: attribute.ReceiverAmount,
			FeeAmount:      attribute.FeeAmount,
			EndOfMonth:     attribute.OffchainTime.EndOfMonth,
			OriginTraceNo:  attribute.RevertTraceNo,

			PreDay:          preDay,
			CurrentDay:      currentDay,
			NextDay:         nextDay,
			TransformTaskId: taskId,
		}
		onchainTxList = append(onchainTxList, onchainTransaction)
	}
	return onchainTxList, nil
}

func (t *taskHandler) rawTxToAccounts(rawTransaction *entities.RawTransaction, taskId int64) ([]*entities.Account, common.BaseError) {
	var accountList []*entities.Account

	response, err := contractpb.NewAccountResponse([]byte(rawTransaction.TxData))
	if err != nil {
		t.logger.Errorf("Create new account response error %s", err.Error())
		return nil, common.NewUnknownError(err)

	}
	account := &entities.Account{
		Txhash:          rawTransaction.Txhash,
		TraceNo:         response.TraceNo,
		AccountId:       response.Account.Id,
		BlockadeAmount:  response.Account.BlockadeAmount,
		Balance:         response.Account.Balance,
		Type:            int32(response.Account.Type),
		State:           int32(*response.Account.State.Enum()),
		TimeOnchain:     rawTransaction.BlockTime,
		NetworkName:     rawTransaction.NetworkName,
		ChannelName:     rawTransaction.ChannelName,
		BlockNum:        rawTransaction.BlockNum,
		TransformTaskId: taskId,
	}
	accountList = append(accountList, account)

	for _, accountResponse := range response.Accounts {
		newAccount := &entities.Account{
			Txhash:          rawTransaction.Txhash,
			TraceNo:         response.TraceNo,
			AccountId:       accountResponse.Id,
			BlockadeAmount:  accountResponse.BlockadeAmount,
			Balance:         accountResponse.Balance,
			Type:            int32(accountResponse.Type),
			State:           int32(*accountResponse.State.Enum()),
			TimeOnchain:     rawTransaction.BlockTime,
			NetworkName:     rawTransaction.NetworkName,
			ChannelName:     rawTransaction.ChannelName,
			BlockNum:        rawTransaction.BlockNum,
			TransformTaskId: taskId,
		}
		accountList = append(accountList, newAccount)

	}
	return accountList, nil
}
