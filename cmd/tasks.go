package cmd

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/blcvn/corev4-explorer/appconfig"
	"github.com/blcvn/corev4-explorer/repository/databases"
	"github.com/blcvn/corev4-explorer/repository/tasks"
	"github.com/blcvn/corev4-explorer/services"
	"github.com/blcvn/corev4-explorer/usecases"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func performTasks(metricPort int) {
	ctx := context.Background()
	appconfig.InitConfig()

	taskDBRepo, err := databases.NewTaskDBRepository(ctx)
	if err != nil {
		panic(fmt.Sprintf("Can't create taskDBRepo: %s", err.Error()))
	}
	dataDBRepo, err := databases.NewDataDBRepository(ctx)
	if err != nil {
		panic(fmt.Sprintf("Can't create taskDBRepo: %s", err.Error()))
	}
	taskRepo := tasks.NewTaskRepository(taskDBRepo)
	taskHandler := usecases.NewTaskHandler(dataDBRepo)

	uc := usecases.NewTaskUsecase(taskRepo, taskHandler)

	taskService := services.NewTaskService(uc)
	taskService.Start()

	router := http.NewServeMux()

	router.Handle("/metrics", promhttp.Handler())

	err = http.ListenAndServe(fmt.Sprintf(":%d", metricPort), router)
	if err != nil {
		log.Fatalf("Error while starting prometheus metrics at port %d. Error: %s", metricPort, err.Error())
	}
}
