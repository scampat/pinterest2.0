package pins

import (
	"context"

	"pinterest2.0/pins/internal/repository/memory"
	"pinterest2.0/pins/pkg"
)

type Controller struct {
	repo *memory.Repository
}

func NewController(repo *memory.Repository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) GetPinsByUser(ctx context.Context, userID int) ([]pkg.Pin, error) {
	return c.repo.GetPinsByUser(ctx, userID)
}

func (c *Controller) CreatePin(ctx context.Context, pin pkg.Pin) (pkg.Pin, error) {
	return c.repo.CreatePin(ctx, pin)
}
