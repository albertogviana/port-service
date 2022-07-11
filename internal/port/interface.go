package port

import (
	"context"

	"github.com/albertogviana/port-service/internal/entity"
)

// RepositoryReader has a list of methods that needs to be implemented by the
// repository for the reader interface.
type RepositoryReader interface {
	FindByUnLoc(ctx context.Context, unloc string) (*entity.Port, error)
}

// RepositoryWriter has a list of methods that needs to be implemented by the
// repository for the writer interface.
type RepositoryWriter interface {
	Create(ctx context.Context, p *entity.Port) error
	Update(ctx context.Context, p *entity.Port) error
}

// Repository combines the reader and writer repository interfaces.
type Repository interface {
	RepositoryReader
	RepositoryWriter
}

// UseCaseReader has a list of methods that needs to be implemented by the
// service for the reader interface.
type UseCaseReader interface{}

// UseCaseWriter has a list of methods that needs to be implemented by the
// service for the writer interface.
type UseCaseWriter interface {
	SavePort(ctx context.Context, p *entity.Port) error
}

// UseCase combines the reader and writer service interfaces.
type UseCase interface {
	UseCaseReader
	UseCaseWriter
}
