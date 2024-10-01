package usecases

import (
	"github.com/blcvn/corev4-explorer/common"
	"github.com/blcvn/corev4-explorer/entities"
)

type taskUc struct {
	t taskRepo
}

func NewTaskUsecase(t taskRepo) *taskUc {
	return &taskUc{
		t: t,
	}
}

func (uc *taskUc) PerformTransformTask() common.BaseError {
	if tasks, err := uc.t.LoadTasks(int32(entities.TaskSync), int32(entities.TaskOpen)); err == nil {
		blockNumber := getMaxBlockNumber(tasks)
		if _, err := uc.t.CreateTask(int32(entities.TaskTransform), int32(entities.TaskOpen), blockNumber); err == nil {
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
