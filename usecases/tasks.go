package usecases

import (
	"fmt"

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
		t:      t,
		h:      h,
		logger: flogging.MustGetLogger("corev4-explorer.usecase.task"),
	}
}

func (uc *taskUc) PerformTransformTask() common.BaseError {
	if task, err := uc.t.CreateTask(int32(entities.TaskTransform), int32(entities.TaskProcessing)); err != nil {
		return err
	} else {
		if err := uc.h.PerformTasks(task.Id, int32(entities.TaskTransform)); err != nil {
			uc.logger.Errorf("PerformTasks id: %d type: %d error: %s", task.Id, entities.TaskTransform, err.Error())
			if _err := uc.t.UpdateTasks(task.Id, task.Type, int32(entities.TaskFailed), nil); _err != nil {
				uc.logger.Errorf("UpdateTasks id: %d status: %d error: %s", task.Id, entities.TaskFailed, _err.Error())
				return _err
			}
		}
		if err := uc.t.UpdateTasks(task.Id, 0, int32(entities.TaskDone), nil); err != nil {
			uc.logger.Errorf("UpdateTasks id: %d status: %d error: %s", task.Id, entities.TaskDone, err.Error())
			return err
		}
	}
	return nil
}

func (uc *taskUc) PerformDeltaTask() common.BaseError {
	fmt.Printf("PerformDeltaTask UC\n")
	return nil
}

func (uc *taskUc) PerformBalanceTask() common.BaseError {
	fmt.Printf("PerformBalanceTask UC\n")
	return nil
}
