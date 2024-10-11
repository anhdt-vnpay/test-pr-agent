package databases

import (
	"context"

	"github.com/blcvn/corev4-explorer/entities"
)

type taskDB interface {
	CreateTask(task *entities.Task) error
	GetAllTaskByStatus(taskType int32, Tasktatus int32) ([]*entities.Task, error)
	UpdateTaskStatus(taskIds []int64, Tasktatus int32) error
}

type dataDB interface {
	CalculateAccountDelta(taskId string) error
	CalculateAccountBalance(taskId string) error

	SaveTransformData(taskId string, data any)
	SaveBlockAndRawTxs(ctx context.Context, block any, rawTxs []any) error
}
