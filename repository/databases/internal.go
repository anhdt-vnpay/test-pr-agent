package databases

import (
	"github.com/blcvn/corev4-explorer/entities"
)

type taskDB interface {
	CreateTask(task *entities.Task) error
	GetAllTaskByStatus(taskType int32, Tasktatus int32) ([]*entities.Task, error)
	UpdateTaskStatus(taskIds []int64, Tasktatus int32) error
}

type dataDB interface {
	CalculateAccountDelta(taskId int64) error
	CalculateAccountBalance(taskId int64) error

	SaveTransformData(taskId int64, data any) error
	SaveBlockAndRawTxs(block *entities.Block, rawTxs []*entities.RawTransaction) error
}
