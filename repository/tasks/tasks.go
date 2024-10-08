package tasks

type taskRepo struct {
	dbRepo dbRepo
}

func NewTaskRepository(dbRepo dbRepo) *taskRepo {
	return &taskRepo{
		dbRepo: dbRepo,
	}
}
