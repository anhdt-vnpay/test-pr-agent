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
	CalculateAccountDelta(taskId int64) error
	CalculateAccountBalance(taskId int64) error
	SaveTransformData(taskId int64, onchainTransactions []*entities.OnchainTransaction, accounts []*entities.Account, accountTxsList []*entities.AccountTx) error
	QueryRawTxsByTaskTransformID(taskId int64) ([]*entities.RawTransaction, error)
	// CreateTask(task *entities.Task) error
	// GetAllTaskByStatus(taskType int32, taskStatus int32) ([]*entities.Task, error)
	// UpdateTasksStatus(taskIds []int64, taskStatus int32) error
	// UpdateTaskStatus(taskId int64, taskStatus int32) error
	// FindOneById(taskId int64) (*entities.Task, error)
	// GetRawTransactionByBlockNumber(blockNumber uint64) ([]*entities.RawTransaction, error)
}

type taskHandler struct {
	dbRepo dbRepo
	logger *flogging.FabricLogger
}

func NewTaskHandler(dbRepo dbRepo) *taskHandler {
	return &taskHandler{
		dbRepo: dbRepo,
		logger: flogging.MustGetLogger("corev4-explorer.usecase.task.handler"),
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
	t.logger.Infof("Transform task %d start", taskId)
	rawTxs, err := t.dbRepo.QueryRawTxsByTaskTransformID(taskId)
	if err != nil {
		return common.NewError(common.ERROR_INTERNAL, err)
	}

	var onchainTxList []*entities.OnchainTransaction

	var accountList []*entities.Account

	var accountTxs []*entities.AccountTx

	for _, rawTransaction := range rawTxs {

		if strings.HasPrefix(rawTransaction.ChaincodeProposalInput, "Transaction.") {
			newTransactions, newAccountTxs, err := t.rawTxToTransactions(rawTransaction, taskId)
			if err != nil {
				return err
			}
			onchainTxList = append(onchainTxList, newTransactions...)
			accountTxs = append(accountTxs, newAccountTxs...)

		}

		if strings.HasPrefix(rawTransaction.ChaincodeProposalInput, "Account.") {
			newAccounts, err := t.rawTxToAccounts(rawTransaction, taskId)
			if err != nil {
				return err
			}
			accountList = append(accountList, newAccounts...)

		}

	}
	err = t.dbRepo.SaveTransformData(taskId, onchainTxList, accountList, accountTxs)
	if err != nil {
		return common.NewError(common.ERROR_INTERNAL, err)
	}
	return nil
}

func (t *taskHandler) calculateAccountDelta(taskId int64) common.BaseError {
	return nil
}

func (t *taskHandler) rawTxToTransactions(rawTransaction *entities.RawTransaction, taskId int64) ([]*entities.OnchainTransaction, []*entities.AccountTx, common.BaseError) {
	var onchainTxList []*entities.OnchainTransaction
	var accountTxsList []*entities.AccountTx

	response, err := contractpb.NewTransactionResponse(rawTransaction.Payload)
	if err != nil {
		t.logger.Errorf("Create new transaction response error %s", err.Error())
		return nil, nil, common.NewUnknownError(err)
	}

	for _, attribute := range response.Transaction.Attributes {

		preDay, err := time.Parse(common.YYYYMMDD, attribute.OffchainTime.PreDay)
		if err != nil {
			t.logger.Errorf("Error when parse preday %v error %s", attribute.OffchainTime.PreDay, err.Error())
			// return nil, nil, common.NewUnknownError(err)

		}
		currentDay, err := time.Parse(common.YYYYMMDD, attribute.OffchainTime.CurrentDay)
		if err != nil {
			t.logger.Errorf("Error when parse currentDay %s", err.Error())
			// return nil, nil, common.NewUnknownError(err)

		}
		nextDay, err := time.Parse(common.YYYYMMDD, attribute.OffchainTime.NextDay)
		if err != nil {
			t.logger.Errorf("Error when parse nextDay %s", err.Error())
			// return nil, nil, common.NewUnknownError(err)

		}
		var timeOffchain time.Time
		if attribute.CreatedTime == 0 {
			timeOffchain = rawTransaction.BlockTime
		} else {
			timeOffchain = time.Unix(0, (attribute.CreatedTime)*int64(time.Millisecond))
		}
		onchainTransaction := &entities.OnchainTransaction{
			Txhash:       rawTransaction.Txhash,
			BlockNum:     rawTransaction.BlockNum,
			NetworkName:  rawTransaction.NetworkName,
			ChannelName:  rawTransaction.ChannelName,
			TimeOnchain:  rawTransaction.BlockTime,
			TimeOffchain: timeOffchain,
			TraceNo:      attribute.TraceNo,
			TxType:       response.FuncType.Enum().String(),

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
		if attribute.SenderId != "" {
			accountTxsList = append(accountTxsList, &entities.AccountTx{
				Txhash:          rawTransaction.Txhash,
				TraceNo:         attribute.TraceNo,
				AccountId:       attribute.SenderId,
				Amount:          -attribute.SenderAmount,
				TimeOnchain:     rawTransaction.BlockTime.UTC(),
				TimeOffchain:    timeOffchain,
				NetworkName:     rawTransaction.NetworkName,
				ChannelName:     rawTransaction.ChannelName,
				BlockNum:        rawTransaction.BlockNum,
				PreDay:          preDay,
				NextDay:         nextDay,
				CurrentDay:      currentDay,
				EndOfMonth:      "N",
				CreatedAt:       time.Now().Local(),
				TransformTaskId: taskId,
			})
		}
		if attribute.ReceiverId != "" {
			accountTxsList = append(accountTxsList, &entities.AccountTx{
				Txhash:          rawTransaction.Txhash,
				TraceNo:         attribute.TraceNo,
				AccountId:       attribute.ReceiverId,
				Amount:          attribute.ReceiverAmount,
				TimeOnchain:     rawTransaction.BlockTime.UTC(),
				TimeOffchain:    timeOffchain,
				NetworkName:     rawTransaction.NetworkName,
				ChannelName:     rawTransaction.ChannelName,
				BlockNum:        rawTransaction.BlockNum,
				PreDay:          preDay,
				NextDay:         nextDay,
				CurrentDay:      currentDay,
				EndOfMonth:      "N",
				CreatedAt:       time.Now().Local(),
				TransformTaskId: taskId,
			})
		}
	}
	return onchainTxList, accountTxsList, nil
}

func (t *taskHandler) rawTxToAccounts(rawTransaction *entities.RawTransaction, taskId int64) ([]*entities.Account, common.BaseError) {
	var accountList []*entities.Account
	response, err := contractpb.NewAccountResponse(rawTransaction.Payload)
	if err != nil {
		t.logger.Errorf("Create new account response error %s", err.Error())
		return nil, common.NewUnknownError(err)

	}

	var timeOffchain time.Time
	if response.Account.CreatedTime == 0 {
		timeOffchain = rawTransaction.BlockTime
	} else {
		timeOffchain = time.Unix(0, (response.Account.CreatedTime)*int64(time.Millisecond))
	}
	account := &entities.Account{
		Txhash:          rawTransaction.Txhash,
		TraceNo:         response.TraceNo,
		AccountId:       response.Account.Id,
		Amount:          0,
		BlockadeAmount:  response.Account.BlockadeAmount,
		Balance:         response.Account.Balance,
		Type:            int32(response.Account.Type),
		State:           int32(*response.Account.State.Enum()),
		TimeOnchain:     rawTransaction.BlockTime,
		TimeOffchain:    timeOffchain,
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
