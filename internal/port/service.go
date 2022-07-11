package port

import (
	"context"
	"github.com/albertogviana/port-service/internal/entity"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{
		repo: repo,
	}
}

// SavePort saves a port in the database or returns an error.
func (s *Service) SavePort(ctx context.Context, p *entity.Port) error {
	currentPort, err := s.repo.FindByUnLoc(ctx, p.Unloc)

	if err != nil && err != ErrPortNotFound {
		return err
	}

	if currentPort == nil {
		return s.repo.Create(ctx, p)
	}

	return s.repo.Update(ctx, p)
}
