package databases

import (
	"github.com/blcvn/corev4-explorer/entities"
)

type taskDB interface {
	CreateNewTaskTransform(taskStatus int32) (*entities.Task, error)
	CreateNewTaskDelta(taskStatus int32) (*entities.Task, error)
	CreateNewTaskBalance(taskStatus int32) (*entities.Task, error)
	UpdateTasksStatus(taskId int64, taskStatus int32) error
}

type dataDB interface {
	CalculateAccountDelta(taskId int64) error
	CalculateAccountBalance(taskId int64) error

	QueryRawTxsByTaskTransformID(taskId int64) ([]*entities.RawTransaction, error)
	SaveTransformData(taskId int64, onchainTransactions []*entities.OnchainTransaction, accounts []*entities.Account, accountTxsList []*entities.AccountTx) error
	SaveBlockAndRawTxs(block *entities.Block, rawTxs []*entities.RawTransaction) error
}
