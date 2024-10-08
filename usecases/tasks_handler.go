package usecases

type dbRepo interface {
	CalculateAccountDelta(taskId string) error
	CalculateAccountBalance(taskId string) error

	SaveTransformData(taskId string, data any)
}

type taskHandler struct {
	dbRepo dbRepo
}

func NewTaskHandler(dbRepo dbRepo) *taskHandler {
	return &taskHandler{
		dbRepo: dbRepo,
	}
}
