package usecases_test

import (
	"errors"
	"testing"

	"github.com/blcvn/corev4-explorer/common"
	"github.com/blcvn/corev4-explorer/entities"
	"github.com/blcvn/corev4-explorer/usecases"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// PerformTransformTask successfully loads tasks and creates a new task
func TestPerformTransformTask_Success(t *testing.T) {
    // Arrange
    mockTaskRep := new(usecases.MockTaskRepo)
    mockTaskHandler := new(usecases.MockTasksHandler)
    uc := usecases.NewTaskUsecase(mockTaskRep, mockTaskHandler)
    tasks := []*entities.Tasks{{}, {}}
    mockTaskRep.On("LoadTasks", int32(entities.TaskSync), int32(entities.TaskOpen)).Return(tasks, nil)
    mockTaskRep.On("CreateTask", int32(entities.TaskTransform), int32(entities.TaskOpen), mock.Anything).Return(&entities.Tasks{}, nil)
    mockTaskRep.On("UpdateTasks", mock.Anything, mock.Anything, int32(entities.TaskProcessing), mock.Anything).Return(nil)
    mockTaskHandler.On("PerformTasks", mock.Anything, mock.Anything).Return(nil)
    mockTaskRep.On("UpdateTasks", mock.Anything, mock.Anything, int32(entities.TaskDone), mock.Anything).Return(nil)

    // Act
    err := uc.PerformTransformTask()

    // Assert
    assert.Nil(t, err)
    mockTaskRep.AssertExpectations(t)
    mockTaskHandler.AssertExpectations(t)
}

// PerformTransformTask handles error when LoadTasks fails
func TestPerformTransformTask_LoadTasksError(t *testing.T) {
    // Arrange
    mockTaskRep := new(usecases.MockTaskRepo)
    mockTaskHandler := new(usecases.MockTasksHandler)
    uc := usecases.NewTaskUsecase(mockTaskRep, mockTaskHandler)
    expectedError := common.NewUnknownError(errors.New("LoadTasks error"))
    mockTaskRep.On("LoadTasks", int32(entities.TaskSync), int32(entities.TaskOpen)).Return(nil, expectedError)

    // Act
    err := uc.PerformTransformTask()

    // Assert
    assert.Equal(t, expectedError, err)
    mockTaskRep.AssertExpectations(t)
}
