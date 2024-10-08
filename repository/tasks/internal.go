package tasks

import "github.com/blcvn/corev4-explorer/entities"

type dbRepo interface {
	CreateTask(task *entities.Tasks) error
	GetAllTaskByStatus(taskType int32, taskStatus int32) ([]*entities.Tasks, error)
	UpdateTasksStatus(taskIds []int64, taskStatus int32) error
}
