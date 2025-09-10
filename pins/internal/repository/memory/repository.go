package memory

import (
	"context"
	"sync"
	"time"

	"pinterest2.0/pins/pkg"
)

type Repository struct {
	mu     sync.RWMutex
	pins   map[int]pkg.Pin
	nextID int
}

func NewRepository() *Repository {
	return &Repository{
		pins:   make(map[int]pkg.Pin),
		nextID: 1,
	}
}

func (r *Repository) GetPinsByUser(ctx context.Context, userID int) ([]pkg.Pin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []pkg.Pin
	for _, p := range r.pins {
		if p.UserID == userID {
			result = append(result, p)
		}
	}
	return result, nil
}

func (r *Repository) CreatePin(ctx context.Context, pin pkg.Pin) (pkg.Pin, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	pin.ID = r.nextID
	pin.CreatedAt = time.Now()
	r.pins[pin.ID] = pin
	r.nextID++

	return pin, nil
}
