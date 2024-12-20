package tasks

import "github.com/blcvn/corev4-explorer/entities"

type dbRepo interface {
	CreateNewTaskTransform(taskStatus int32) (*entities.Task, error)
	CreateNewTaskDelta(taskStatus int32) (*entities.Task, error)
	CreateNewTaskBalance(taskStatus int32) (*entities.Task, error)
	UpdateTasksStatus(taskId int64, taskStatus int32) error
}
