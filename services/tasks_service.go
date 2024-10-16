package services

import (
	"time"

	"github.com/blcvn/corev3-libs/flogging"
	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/blcvn/corev4-explorer/common"
	cron "github.com/robfig/cron/v3"
)

type usecases interface {
	PerformTransformTask() common.BaseError
	PerformDeltaTask() common.BaseError
	PerformBalanceTask() common.BaseError
}

type taskService struct {
	usecases              usecases
	transformTaskInterval time.Duration
	deltaTaskScheduler    string
	balanceTaskScheduler  string
	scheduler             *cron.Cron
	logger                *flogging.FabricLogger
}

func NewTaskService(usecases usecases) *taskService {
	return &taskService{
		usecases:              usecases,
		scheduler:             cron.New(),
		transformTaskInterval: appconfig.TransformInterval,
		deltaTaskScheduler:    appconfig.DeltaTaskScheduler,
		balanceTaskScheduler:  appconfig.BalanceTaskScheduler,
		logger:                flogging.MustGetLogger("corev4-explorer.services.task"),
	}
}

func (s *taskService) Start() {
	go func() {
		ticker := time.NewTicker(s.transformTaskInterval)
		for {
			err := s.usecases.PerformTransformTask()
			if err != nil {
				s.logger.Errorf("[catch me] PerformTransformTask: %s", err.Error())
			}
			<-ticker.C
		}
	}()
	// s.scheduler.AddFunc(s.deltaTaskScheduler, func() {
	// 	err := s.usecases.PerformDeltaTask()
	// 	if err != nil {
	// 		s.logger.Errorf("[catch me] PerformDeltaTask: %s", err.Error())
	// 	}
	// })

	// s.scheduler.AddFunc(s.balanceTaskScheduler, func() {
	// 	err := s.usecases.PerformBalanceTask()
	// 	if err != nil {
	// 		s.logger.Errorf("[catch me] PerformBalanceTask: %s", err.Error())
	// 	}
	// })
	// s.scheduler.Start()
}
