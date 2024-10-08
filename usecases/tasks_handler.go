package usecases

import (
	"errors"

	"github.com/blcvn/corev4-explorer/common"
	"github.com/blcvn/corev4-explorer/entities"
)

type dbRepo interface {
	CalculateAccountDelta(taskId string) error
	CalculateAccountBalance(taskId string) error

	SaveTransformData(taskId string, data any)
}

type taskHandler struct {
	dbRepo dbRepo
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
	return nil
}

func (t *taskHandler) calculateAccountDelta(taskId int64) common.BaseError {
	return nil
}
