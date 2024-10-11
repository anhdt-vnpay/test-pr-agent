package usecases

import "github.com/blcvn/corev4-explorer/entities"

func getMaxBlockNumber(tasks []*entities.Task) uint64 {
	max := uint64(0)
	for _, task := range tasks {
		if max < task.BlockNumber {
			max = task.BlockNumber
		}
	}
	return max
}
