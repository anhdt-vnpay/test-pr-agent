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
	uc := usecases.NewTaskUsecase(mockTaskRep)
	tasks := []*entities.Tasks{{}, {}}
	mockTaskRep.On("LoadTasks", int32(entities.TaskSync), int32(entities.TaskOpen)).Return(tasks, nil)
	mockTaskRep.On("CreateTask", int32(entities.TaskTransform), int32(entities.TaskOpen), mock.Anything).Return(nil, nil)

	// Act
	err := uc.PerformTransformTask()

	// Assert
	assert.Nil(t, err)
	mockTaskRep.AssertExpectations(t)
}

// PerformTransformTask handles error when LoadTasks fails
func TestPerformTransformTask_LoadTasksError(t *testing.T) {
	// Arrange
	mockTaskRep := new(usecases.MockTaskRepo)
	uc := usecases.NewTaskUsecase(mockTaskRep)
	expectedError := common.NewUnknownError(errors.New("LoadTasks failed"))
	mockTaskRep.On("LoadTasks", int32(entities.TaskSync), int32(entities.TaskOpen)).Return(nil, expectedError)

	// Act
	err := uc.PerformTransformTask()

	// Assert
	assert.Equal(t, expectedError, err)
	mockTaskRep.AssertExpectations(t)
}
