package entities

type TaskType int32

type TaskStatus int32

const (
	TaskSync      TaskType = 1
	TaskTransform TaskType = 2
)

const (
	TaskOpen       TaskStatus = 1
	TaskProcessing TaskStatus = 2
	TaskDone       TaskStatus = 3
)

type Tasks struct {
	Id          int64
	Type        int32
	Status      int32
	BlockNumber uint64
}
