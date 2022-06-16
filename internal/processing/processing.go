package processing

import (
	"GoBeLvl2/internal/database/dbport"
	"GoBeLvl2/internal/entities"
	"context"
	"fmt"
)

type Processing struct {
	ds *dbport.DbStorage
}

func NewProcessing(ds *dbport.DbStorage) *Processing {
	return &Processing{
		ds: ds,
	}
}

func (p *Processing) AddUser(ctx context.Context, user entities.User) error {
	err := p.ds.AddUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (p *Processing) AddEnv(ctx context.Context, env interface{}) error {

	switch vv := env.(type) {
	case entities.Community:
		err := p.ds.AddCommunity(ctx, vv)
		if err != nil {
			return err
		}
		return nil

	case entities.CorpGroup:
		err := p.ds.AddCorpGroup(ctx, vv)
		if err != nil {
			return err
		}
		return nil

	case entities.Organization:
		err := p.ds.AddOrganization(ctx, vv)
		if err != nil {
			return err
		}
		return nil

	case entities.Project:
		err := p.ds.AddProject(ctx, vv)
		if err != nil {
			return err
		}
		return nil

	default:
		fmt.Printf("I don't know about env %T!\n", vv)
	}

	return nil
}

func (p *Processing) SearchEnv(ctx context.Context, params []string) (interface{}, error) {
	i, err := p.ds.SearchEnv(ctx, params)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (p *Processing) SearchUser(ctx context.Context, params []string) (*entities.User, error) {
	user, err := p.ds.SearchUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (p *Processing) SetUserInEnv(ctx context.Context, user entities.User, env ...interface{}) error {
	for _, v := range env {
		switch vv := v.(type) {
		case entities.Community:
			err := p.ds.SetUserInCommunity(ctx, user, vv)
			if err != nil {
				return err
			}
			return nil

		case entities.CorpGroup:
			err := p.ds.SetUserInCorpGroup(ctx, user, vv)
			if err != nil {
				return err
			}
			return nil

		case entities.Organization:
			err := p.ds.SetUserInOrganization(ctx, user, vv)
			if err != nil {
				return err
			}
			return nil

		case entities.Project:
			err := p.ds.SetUserInProject(ctx, user, vv)
			if err != nil {
				return err
			}
			return nil

		default:
			fmt.Printf("I don't know about env %T!\n", vv)
		}

	}
	return nil
}

func (p *Processing) DeleteUserFromEnv(ctx context.Context, user entities.User, env ...interface{}) error {
	for _, v := range env {
		switch vv := v.(type) {
		case entities.Community:
			err := p.ds.DeleteUserFromCommunity(ctx, user, vv)
			if err != nil {
				return err
			}
			return nil

		case entities.CorpGroup:
			err := p.ds.DeleteUserFromCorpGroup(ctx, user, vv)
			if err != nil {
				return err
			}
			return nil

		case entities.Organization:
			err := p.ds.DeleteUserFromOrganization(ctx, user, vv)
			if err != nil {
				return err
			}
			return nil

		case entities.Project:
			err := p.ds.DeleteUserFromProject(ctx, user, vv)
			if err != nil {
				return err
			}
			return nil

		default:
			fmt.Printf("I don't know about env %T!\n", vv)
		}

	}
	return nil
}
