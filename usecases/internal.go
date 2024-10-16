package usecases

import (
	"github.com/blcvn/corev4-explorer/common"
	"github.com/blcvn/corev4-explorer/entities"
)

type taskRepo interface {
	// LoadTasks(taskType int32, taskStatus int32) ([]*entities.Task, common.BaseError)
	CreateTask(taskType int32, taskStatus int32) (*entities.Task, common.BaseError)
	UpdateTasks(tasksId int64, tasksType int32, taskStatus int32, data any) common.BaseError
}

type tasksHandler interface {
	PerformTasks(taskId int64, tasksType int32) common.BaseError
}
