package port

import (
	"context"
	"fmt"

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
		err := s.repo.Create(ctx, p)
		if err != nil {
			return fmt.Errorf("error during the create new port: %w", err)
		}

		return nil
	}

	err = s.repo.Update(ctx, p)
	if err != nil {
		return fmt.Errorf("error during the update port: %w", err)
	}

	return nil
}
