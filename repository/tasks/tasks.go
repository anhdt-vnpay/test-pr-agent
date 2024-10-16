package tasks

import (
	"errors"

	"github.com/blcvn/corev4-explorer/common"
	"github.com/blcvn/corev4-explorer/entities"
)

type taskRepo struct {
	dbRepo dbRepo
}

func NewTaskRepository(dbRepo dbRepo) *taskRepo {
	return &taskRepo{
		dbRepo: dbRepo,
	}
}

func (t *taskRepo) CreateTask(taskType int32, taskStatus int32) (*entities.Task, common.BaseError) {
	switch taskType {
	case int32(entities.TaskTransform):
		task, err := t.dbRepo.CreateNewTaskTransform(taskStatus)
		if err != nil {
			return nil, common.NewError(common.ERROR_INTERNAL, err)
		}
		return task, nil
	case int32(entities.TaskDelta):
		task, err := t.dbRepo.CreateNewTaskDelta(taskStatus)
		if err != nil {
			return nil, common.NewError(common.ERROR_INTERNAL, err)
		}
		return task, nil
	case int32(entities.TaskBalance):
		task, err := t.dbRepo.CreateNewTaskBalance(taskStatus)
		if err != nil {
			return nil, common.NewError(common.ERROR_INTERNAL, err)
		}
		return task, nil
	default:
		return nil, common.NewError(common.ERROR_INTERNAL, errors.New("task type invalid"))
	}
}

func (t *taskRepo) UpdateTasks(tasksId int64, tasksType int32, taskStatus int32, data any) common.BaseError {
	err := t.dbRepo.UpdateTasksStatus(tasksId, taskStatus)
	if err != nil {
		return common.NewError(common.ERROR_INTERNAL, err)
	}
	return nil
}
