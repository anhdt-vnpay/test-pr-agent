package services

import (
	"context"
)

type ledgerRepo interface {
	Start(ctx context.Context)
}

type syncService struct {
	ledgerRepo ledgerRepo
}

func NewSyncService(ledgerRepo ledgerRepo) *syncService {
	return &syncService{
		ledgerRepo: ledgerRepo,
	}
}

func (s *syncService) Start(ctx context.Context) {
	s.ledgerRepo.Start(ctx)
}
