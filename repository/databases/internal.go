package databases

import (
	"context"

	"github.com/blcvn/corev4-explorer/entities"
)

type taskDB interface {
	CreateTask(task *entities.Tasks) error
	GetAllTaskByStatus(taskType int32, taskStatus int32) ([]*entities.Tasks, error)
	UpdateTasksStatus(taskIds []int64, taskStatus int32) error
}

type dataDB interface {
	CalculateAccountDelta(taskId string) error
	CalculateAccountBalance(taskId string) error

	SaveTransformData(taskId string, data any)
	SaveBlockAndRawTxs(ctx context.Context, block any, rawTxs []any) error
}
