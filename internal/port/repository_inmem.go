package port

import (
	"context"
	"sync"

	"github.com/albertogviana/port-service/internal/entity"
)

type RepositoryInMem struct {
	data map[string]*entity.Port
	mu   sync.RWMutex
}

func NewRepositoryInMem() *RepositoryInMem {
	return &RepositoryInMem{
		data: make(map[string]*entity.Port),
	}
}

func (r *RepositoryInMem) Create(ctx context.Context, p *entity.Port) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[p.Unloc] = p

	return nil
}

func (r *RepositoryInMem) Update(ctx context.Context, p *entity.Port) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data[p.Unloc] = p

	return nil
}

func (r *RepositoryInMem) FindByUnLoc(ctx context.Context, unloc string) (*entity.Port, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.data[unloc]
	if !ok {
		return nil, ErrPortNotFound
	}

	return p, nil
}
