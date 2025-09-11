package board

import (
	"context"
	"errors"

	"pinterest2.0/boards/internal/gateway"
	"pinterest2.0/boards/internal/gateway/pins"
	"pinterest2.0/boards/pkg/model"
)

type pinsGateway interface {
	GetPinsByBoard(ctx context.Context, boardID string) ([]pins.Pin, error)
}

type Controller struct {
	pinsGateway pinsGateway
}

var boards = map[string]model.Board{}

func New(pinsGateway pinsGateway) *Controller {
	return &Controller{pinsGateway}
}

func (c *Controller) CreateBoard(id, userID, name string) model.Board {
	b := model.Board{ID: id, UserID: userID, Name: name}
	boards[id] = b
	return b
}

func (c *Controller) GetBoardsByUser(userID string) ([]model.Board, error) {
	var userBoards []model.Board
	for _, b := range boards {
		if b.UserID == userID {
			userBoards = append(userBoards, b)
		}
	}
	if len(userBoards) == 0 {
		return nil, gateway.ErrNotFound
	}
	return userBoards, nil
}

func (c *Controller) GetBoardPins(ctx context.Context, boardID string) ([]pins.Pin, error) {
	return c.pinsGateway.GetPinsByBoard(ctx, boardID)
}

func (c *Controller) GetBoard(id string) (model.Board, error) {
	b, ok := boards[id]
	if !ok {
		return model.Board{}, errors.New("board not found")
	}
	return b, nil
}
