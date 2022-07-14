package repository

import (
	"context"
	"sync"

	"github.com/albertogviana/port-service/internal/port"
)

type InMemRepository struct {
	data map[string]*port.Port
	mu   sync.RWMutex
}

func NewInMemRepository() *InMemRepository {
	return &InMemRepository{
		data: make(map[string]*port.Port),
	}
}

func (r *InMemRepository) Create(ctx context.Context, p *port.Port) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := p.Unlocs[0]

	r.data[key] = p

	return nil
}

func (r *InMemRepository) Update(ctx context.Context, p *port.Port) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := p.Unlocs[0]

	r.data[key] = p

	return nil
}

func (r *InMemRepository) FindByUnLoc(ctx context.Context, unloc string) (*port.Port, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	p, ok := r.data[unloc]
	if !ok {
		return nil, port.ErrPortNotFound
	}

	return p, nil
}
