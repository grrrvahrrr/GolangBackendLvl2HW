package dbport

import (
	"GoBeLvl2/internal/entities"
	"context"
)

type DbPort interface {
	AddUser(ctx context.Context, user entities.User) error
	SearchEnv(ctx context.Context, params []string) (interface{}, error)
	SearchUser(ctx context.Context, params []string) (*entities.User, error)
	AddEnv
	SetUserInEnv
	DeleteUserFromEnv
}

type AddEnv interface {
	AddProject(ctx context.Context, pr entities.Project) error
	AddOrganization(ctx context.Context, org entities.Organization) error
	AddCorpGroup(ctx context.Context, cg entities.CorpGroup) error
	AddCommunity(ctx context.Context, com entities.Community) error
}

type SetUserInEnv interface {
	SetUserInProject(ctx context.Context, user entities.User, pr entities.Project) error
	SetUserInOrganization(ctx context.Context, user entities.User, org entities.Organization) error
	SetUserInCorpGroup(ctx context.Context, user entities.User, cg entities.CorpGroup) error
	SetUserInCommunity(ctx context.Context, user entities.User, com entities.Community) error
}

type DeleteUserFromEnv interface {
	DeleteUserFromProject(ctx context.Context, user entities.User, pr entities.Project) error
	DeleteUserFromOrganization(ctx context.Context, user entities.User, org entities.Organization) error
	DeleteUserFromCorpGroup(ctx context.Context, user entities.User, cg entities.CorpGroup) error
	DeleteUserFromCommunity(ctx context.Context, user entities.User, com entities.Community) error
}

type DbStorage struct {
	dbport DbPort
}

func NewDbStorage(dbport DbPort) *DbStorage {
	return &DbStorage{
		dbport: dbport,
	}
}

//Add User
func (ds *DbStorage) AddUser(ctx context.Context, user entities.User) error {
	err := ds.dbport.AddUser(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

//Add env
func (ds *DbStorage) AddProject(ctx context.Context, pr entities.Project) error {
	err := ds.dbport.AddProject(ctx, pr)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DbStorage) AddOrganization(ctx context.Context, org entities.Organization) error {
	err := ds.dbport.AddOrganization(ctx, org)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DbStorage) AddCorpGroup(ctx context.Context, cg entities.CorpGroup) error {
	err := ds.dbport.AddCorpGroup(ctx, cg)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DbStorage) AddCommunity(ctx context.Context, com entities.Community) error {
	err := ds.dbport.AddCommunity(ctx, com)
	if err != nil {
		return err
	}
	return nil
}

//Search
func (ds *DbStorage) SearchEnv(ctx context.Context, params []string) (interface{}, error) {
	i, err := ds.dbport.SearchEnv(ctx, params)
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (ds *DbStorage) SearchUser(ctx context.Context, params []string) (*entities.User, error) {
	user, err := ds.dbport.SearchUser(ctx, params)
	if err != nil {
		return nil, err
	}
	return user, nil
}

//Set User
func (ds *DbStorage) SetUserInCorpGroup(ctx context.Context, user entities.User, cg entities.CorpGroup) error {
	err := ds.dbport.SetUserInCorpGroup(ctx, user, cg)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DbStorage) SetUserInProject(ctx context.Context, user entities.User, pr entities.Project) error {
	err := ds.dbport.SetUserInProject(ctx, user, pr)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DbStorage) SetUserInOrganization(ctx context.Context, user entities.User, org entities.Organization) error {
	err := ds.dbport.SetUserInOrganization(ctx, user, org)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DbStorage) SetUserInCommunity(ctx context.Context, user entities.User, com entities.Community) error {
	err := ds.dbport.SetUserInCommunity(ctx, user, com)
	if err != nil {
		return err
	}
	return nil
}

//Delete User
func (ds *DbStorage) DeleteUserFromCorpGroup(ctx context.Context, user entities.User, cg entities.CorpGroup) error {
	err := ds.dbport.DeleteUserFromCorpGroup(ctx, user, cg)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DbStorage) DeleteUserFromProject(ctx context.Context, user entities.User, pr entities.Project) error {
	err := ds.dbport.DeleteUserFromProject(ctx, user, pr)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DbStorage) DeleteUserFromOrganization(ctx context.Context, user entities.User, org entities.Organization) error {
	err := ds.dbport.DeleteUserFromOrganization(ctx, user, org)
	if err != nil {
		return err
	}
	return nil
}

func (ds *DbStorage) DeleteUserFromCommunity(ctx context.Context, user entities.User, com entities.Community) error {
	err := ds.dbport.DeleteUserFromCommunity(ctx, user, com)
	if err != nil {
		return err
	}
	return nil
}
