package tasks

import "github.com/blcvn/corev4-explorer/entities"

type dbRepo interface {
	CreateTask(task *entities.Task) error
	GetAllTaskByStatus(taskType int32, taskStatus int32) ([]*entities.Task, error)
	UpdateTasksStatus(taskIds []int64, taskStatus int32) error
}
