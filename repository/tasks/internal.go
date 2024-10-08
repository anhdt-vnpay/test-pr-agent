package tasks

import "github.com/blcvn/corev4-explorer/entities"

type dbRepo interface {
	ProcessTasksByStatus(taskType int32, taskStatus int32) ([]entities.Tasks, error)
}
