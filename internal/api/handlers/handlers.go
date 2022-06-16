package handlers

import (
	"GoBeLvl2/internal/entities"
	"GoBeLvl2/internal/processing"
	"context"
	"fmt"
)

type Handlers struct {
	p *processing.Processing
}

func NewHandlers(processing *processing.Processing) *Handlers {
	return &Handlers{
		p: processing,
	}
}

func (h *Handlers) AddEnv(ctx context.Context, env interface{}) error {
	err := h.p.AddEnv(ctx, env)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handlers) AddUser(ctx context.Context, user entities.User) error {
	err := h.p.AddUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handlers) SearchEnv(ctx context.Context, params []string) (interface{}, error) {
	i, err := h.p.SearchEnv(ctx, params)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (h *Handlers) SearchUser(ctx context.Context, params []string) (*entities.User, error) {
	user, err := h.p.SearchUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (h *Handlers) UpdateUserInEnv(ctx context.Context, user entities.User, operation string, env ...interface{}) error {
	switch operation {
	case "set":
		err := h.p.SetUserInEnv(ctx, user, env)
		if err != nil {
			return err
		}
		return nil
	case "delete":
		err := h.p.DeleteUserFromEnv(ctx, user, env)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("invalid operation : %s", operation)
	}

}
