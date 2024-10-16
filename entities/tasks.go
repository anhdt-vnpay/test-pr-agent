package entities

import "time"

type TaskType int32

type TaskStatus int32

const (
	TaskSync      TaskType = 1
	TaskTransform TaskType = 2
	TaskDelta     TaskType = 3
	TaskBalance   TaskType = 4
	TaskOffline   TaskType = 5
)

const (
	TaskOpen       TaskStatus = 1
	TaskProcessing TaskStatus = 2
	TaskDone       TaskStatus = 3
	TaskClosed     TaskStatus = 4
	TaskFailed     TaskStatus = 5
)

type Task struct {
	Id          int64
	Type        int32
	Status      int32
	BlockNumber uint64
	StartBlock  uint64
	EndBlock    uint64
	NetworkName string
	ChannelName string
	CreateAt    time.Time
}
