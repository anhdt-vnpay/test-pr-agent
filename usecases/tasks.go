package usecases

import (
	"github.com/blcvn/corev3-libs/flogging"
	"github.com/blcvn/corev3-libs/gotracing"
	"github.com/blcvn/corev4-explorer/common"
	"github.com/blcvn/corev4-explorer/entities"
)

type taskUc struct {
	t taskRepo
	h tasksHandler

	logger *flogging.FabricLogger
	tracer *gotracing.Tracer
}

func NewTaskUsecase(t taskRepo, h tasksHandler) *taskUc {
	return &taskUc{
		t: t,
		h: h,
	}
}

func (uc *taskUc) PerformTransformTask() common.BaseError {
	if tasks, err := uc.t.LoadTasks(int32(entities.TaskSync), int32(entities.TaskOpen)); err == nil {
		blockNumber := getMaxBlockNumber(tasks)
		if task, err := uc.t.CreateTask(int32(entities.TaskTransform), int32(entities.TaskOpen), blockNumber); err == nil {
			if err := uc.t.UpdateTasks(task.Id, task.Type, int32(entities.TaskProcessing), blockNumber); err == nil {
				if err := uc.h.PerformTasks(task.Id, task.Type); err == nil {
					if err := uc.t.UpdateTasks(task.Id, task.Type, int32(entities.TaskDone), blockNumber); err == nil {
						return nil
					} else {
						return err
					}
				}
			} else {
				return err
			}
			return nil
		} else {
			return err
		}
	} else {
		return err
	}
}
