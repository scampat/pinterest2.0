package user

import (
	"context"
	"errors"
	"time"

	"pinterest2.0/users/internal/gateway/boards"
	"pinterest2.0/users/internal/gateway/pins"
	"pinterest2.0/users/pkg/model"
)

var (
	users       = map[int]model.User{}
	nextID      = 1
	ErrNotFound = errors.New("user not found")
)

func CreateUser(username, email string) model.User {
	u := model.User{
		ID:        nextID,
		Username:  username,
		Email:     email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	users[nextID] = u
	nextID++
	return u
}

func GetUser(id int) (model.User, error) {
	u, ok := users[id]
	if !ok {
		return model.User{}, ErrNotFound
	}
	return u, nil
}

func UpdateUser(id int, req model.UpdateUserRequest) (model.User, error) {
	u, ok := users[id]
	if !ok {
		return model.User{}, ErrNotFound
	}

	if req.Username != nil {
		u.Username = *req.Username
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	u.UpdatedAt = time.Now()

	users[id] = u
	return u, nil
}

func DeleteUser(id int) error {
	if _, ok := users[id]; !ok {
		return ErrNotFound
	}
	delete(users, id)
	return nil
}

func GetAllUsers() []model.User {
	userList := make([]model.User, 0, len(users))
	for _, user := range users {
		userList = append(userList, user)
	}
	return userList
}

type boardsGateway interface {
	GetBoardsByUser(ctx context.Context, userID string) ([]boards.Board, error)
}

type pinsGateway interface {
	GetPinsByUser(ctx context.Context, userID string) ([]pins.Pin, error)
}

type Controller struct {
	boardsGateway boardsGateway
	pinsGateway   pinsGateway
}

func New(boardsGateway boardsGateway, pinsGateway pinsGateway) *Controller {
	return &Controller{boardsGateway, pinsGateway}
}

func (c *Controller) CreateUser(ctx context.Context, req model.CreateUserRequest) model.User {
	return CreateUser(req.Username, req.Email)
}

func (c *Controller) GetUser(ctx context.Context, id int) (model.User, error) {
	return GetUser(id)
}

func (c *Controller) UpdateUser(ctx context.Context, id int, req model.UpdateUserRequest) (model.User, error) {
	return UpdateUser(id, req)
}

func (c *Controller) DeleteUser(ctx context.Context, id int) error {
	return DeleteUser(id)
}

func (c *Controller) GetAllUsers(ctx context.Context) []model.User {
	return GetAllUsers()
}

func (c *Controller) GetUserBoards(ctx context.Context, userID string) ([]boards.Board, error) {
	return c.boardsGateway.GetBoardsByUser(ctx, userID)
}

func (c *Controller) GetUserPins(ctx context.Context, userID string) ([]pins.Pin, error) {
	return c.pinsGateway.GetPinsByUser(ctx, userID)
}
